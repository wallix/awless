package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

var (
	listingFormat string

	listOnlyIDs    bool
	listAllInfra   bool
	listAllAccess  bool
	localResources bool
	sortBy         []string
)

func init() {
	RootCmd.AddCommand(listCmd)
	for _, resource := range []graph.ResourceType{graph.Instance, graph.Vpc, graph.Subnet, graph.SecurityGroup} {
		listCmd.AddCommand(listInfraResourceCmd(resource))
	}
	for _, resource := range []graph.ResourceType{graph.User, graph.Role, graph.Policy, graph.Group} {
		listCmd.AddCommand(listAccessResourceCmd(resource))
	}
	listCmd.AddCommand(listAllCmd)

	listCmd.PersistentFlags().StringVar(&listingFormat, "format", "table", "Format for the display of resources: table or csv")

	listCmd.PersistentFlags().BoolVar(&listOnlyIDs, "ids", false, "List only ids")
	listCmd.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	listCmd.PersistentFlags().StringSliceVar(&sortBy, "sort", []string{"Id"}, "Sort tables by column(s) name(s)")

	listAllCmd.PersistentFlags().BoolVar(&listAllInfra, "infra", false, "List infrastructure resources")
	listAllCmd.PersistentFlags().BoolVar(&listAllAccess, "access", false, "List access resources")
}

var listCmd = &cobra.Command{
	Use:               "list",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,
	Short:             "List various type of items: instances, vpc, subnet ...",
}

var listInfraResourceCmd = func(resourceType graph.ResourceType) *cobra.Command {
	return &cobra.Command{
		Use:   resourceType.PluralString(),
		Short: "List AWS EC2 " + resourceType.PluralString(),

		Run: func(cmd *cobra.Command, args []string) {
			var g *graph.Graph
			var err error
			if localResources {
				g, err = config.LoadInfraGraph()
			} else {
				g, err = aws.InfraService.FetchRDFResources(resourceType)
			}
			exitOn(err)

			printResources(g, resourceType)
		},
	}
}

var listAccessResourceCmd = func(resourceType graph.ResourceType) *cobra.Command {
	return &cobra.Command{
		Use:   resourceType.PluralString(),
		Short: "List AWS IAM " + resourceType.PluralString(),

		Run: func(cmd *cobra.Command, args []string) {
			var g *graph.Graph
			var err error
			if localResources {
				g, err = config.LoadAccessGraph()
			} else {
				g, err = aws.AccessService.FetchRDFResources(resourceType)
			}
			exitOn(err)
			printResources(g, resourceType)
		},
	}
}

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all local resources",

	Run: func(cmd *cobra.Command, args []string) {
		if !listAllInfra && !listAllAccess {
			listAllInfra = true //By default, print only infra
		}
		if listAllInfra {
			if !listOnlyIDs {
				fmt.Println("Infrastructure")
			}
			localInfra, err := config.LoadInfraGraph()
			exitOn(err)
			display.SeveralResourcesOfGraph(localInfra, display.PropertiesDisplayer.Services[aws.InfraServiceName], listOnlyIDs)
		}
		if listAllAccess {
			if !listOnlyIDs {
				fmt.Println("Access")
			}
			localAccess, err := config.LoadAccessGraph()
			exitOn(err)
			display.SeveralResourcesOfGraph(localAccess, display.PropertiesDisplayer.Services[aws.AccessServiceName], listOnlyIDs)
		}
	},
}

func printResources(g *graph.Graph, nodeType graph.ResourceType) {
	displayer := display.BuildOptions(
		display.WithRdfType(nodeType),
		display.WithHeaders(display.DefaultsColumnDefinitions[nodeType]),
		display.WithMaxWidth(shell.GetTerminalWidth()),
		display.WithFormat(listingFormat),
		display.WithIDsOnly(listOnlyIDs),
		display.WithSortBy(sortBy...),
	).SetSource(g).Build()

	exitOn(displayer.Print(os.Stdout))
}

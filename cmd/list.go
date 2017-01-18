package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/rdf"
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
	for _, resource := range []rdf.ResourceType{rdf.Instance, rdf.Vpc, rdf.Subnet} {
		listCmd.AddCommand(listInfraResourceCmd(resource))
	}
	for _, resource := range []rdf.ResourceType{rdf.User, rdf.Role, rdf.Policy, rdf.Group} {
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
	Use:   "list",
	Short: "List various type of items: instances, vpc, subnet ...",
}

var listInfraResourceCmd = func(resourceType rdf.ResourceType) *cobra.Command {
	return &cobra.Command{
		Use:   resourceType.PluralString(),
		Short: "List AWS EC2 " + resourceType.PluralString(),

		Run: func(cmd *cobra.Command, args []string) {
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.InfraFilename))

			} else {
				g, err = aws.InfraService.FetchRDFResources(resourceType)
			}
			exitOn(err)

			printResources(g, resourceType)
		},
	}
}

var listAccessResourceCmd = func(resourceType rdf.ResourceType) *cobra.Command {
	return &cobra.Command{
		Use:   resourceType.PluralString(),
		Short: "List AWS IAM " + resourceType.PluralString(),

		Run: func(cmd *cobra.Command, args []string) {
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.AccessFilename))
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
			localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.InfraFilename))
			exitOn(err)
			display.SeveralResourcesOfGraph(localInfra, display.PropertiesDisplayer.Services[aws.InfraServiceName], listOnlyIDs)
		}
		if listAllAccess {
			if !listOnlyIDs {
				fmt.Println("Access")
			}
			localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.AccessFilename))
			exitOn(err)
			display.SeveralResourcesOfGraph(localAccess, display.PropertiesDisplayer.Services[aws.AccessServiceName], listOnlyIDs)
		}
	},
}

func printResources(g *rdf.Graph, nodeType rdf.ResourceType) {
	var displayer display.GraphDisplayer
	if listOnlyIDs {
		displayer = display.BuildGraphDisplayer(
			[]display.ColumnDefinition{
				display.StringColumnDefinition{Prop: "Id"},
				display.StringColumnDefinition{Prop: "Name"},
			},
			display.Options{RdfType: nodeType, Format: "porcelain", SortBy: sortBy},
		)
	} else {
		maxwidth, err := shell.GetTerminalWidth()
		if err != nil {
			maxwidth = 0
		}
		displayer = display.BuildGraphDisplayer(
			display.DefaultsColumnDefinitions[nodeType],
			display.Options{RdfType: nodeType, Format: listingFormat, SortBy: sortBy, MaxWidth: maxwidth},
		)
	}
	displayer.SetGraph(g)
	err := displayer.Print(os.Stdout)
	exitOn(err)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
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
	for apiName, types := range aws.ResourceTypesPerAPI {
		for _, resType := range types {
			listCmd.AddCommand(listServiceResourceCmd(apiName, resType))
		}
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
	Use:                "list",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,
	Short:              "List various type of items: instances, vpc, subnet ...",
}

var listServiceResourceCmd = func(apiName string, resType string) *cobra.Command {
	return &cobra.Command{
		Use:   cloud.PluralizeResource(resType),
		Short: fmt.Sprintf("List AWS %s %s", apiName, cloud.PluralizeResource(resType)),

		Run: func(cmd *cobra.Command, args []string) {
			var g *graph.Graph
			var err error
			if localResources {
				g, err = config.LoadInfraGraph()
			} else {
				srv, err := cloud.GetServiceForType(resType)
				exitOn(err)
				g, err = srv.FetchByType(resType)
			}
			exitOn(err)

			printResources(g, graph.ResourceType(resType))
		},
	}
}

var listInfraCmd = &cobra.Command{
	Use:   "infra",
	Short: "List ec2 resources",

	Run: listAllInfraResources,
}

var listAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "List iam resources",
	Run:   listAllAccessResources,
}

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all local resources",

	Run: func(cmd *cobra.Command, args []string) {
		listAllInfraResources(cmd, args)
		listAllAccessResources(cmd, args)
	},
}

func printResources(g *graph.Graph, resType graph.ResourceType) {
	displayer := display.BuildOptions(
		display.WithRdfType(resType),
		display.WithHeaders(display.DefaultsColumnDefinitions[resType]),
		display.WithMaxWidth(shell.GetTerminalWidth()),
		display.WithFormat(listingFormat),
		display.WithIDsOnly(listOnlyIDs),
		display.WithSortBy(sortBy...),
	).SetSource(g).Build()

	exitOn(displayer.Print(os.Stdout))
}

func listAllAccessResources(cmd *cobra.Command, args []string) {
	g, err := config.LoadAccessGraph()
	exitOn(err)
	displayer := display.BuildOptions(
		display.WithFormat(listingFormat),
		display.WithIDsOnly(listOnlyIDs),
	).SetSource(g).Build()
	exitOn(displayer.Print(os.Stdout))
}

func listAllInfraResources(cmd *cobra.Command, args []string) {
	g, err := config.LoadInfraGraph()
	exitOn(err)
	displayer := display.BuildOptions(
		display.WithFormat(listingFormat),
		display.WithIDsOnly(listOnlyIDs),
	).SetSource(g).Build()
	exitOn(displayer.Print(os.Stdout))
}

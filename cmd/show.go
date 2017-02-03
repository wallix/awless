package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

func init() {
	for apiName, types := range aws.ResourceTypesPerAPI {
		for _, resType := range types {
			showCmd.AddCommand(showServiceResourceCmd(apiName, resType))
		}
	}

	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:                "show",
	Short:              "Show resource and their relations via a given id: users, groups, instances, vpcs, ...",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,
}

var showServiceResourceCmd = func(apiName, resType string) *cobra.Command {
	command := &cobra.Command{
		Use:   resType + " id",
		Short: fmt.Sprintf("Show properties and relations of an AWS %s %s", apiName, resType),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("id required")
			}
			id := args[0]
			var g *graph.Graph

			srv, err := cloud.GetServiceForType(resType)
			exitOn(err)

			if localResources {
				g = sync.LoadCurrentLocalGraph(srv.Name())
			} else {
				g, err = srv.FetchByType(resType)
			}
			exitOn(err)

			printResource(g, graph.ResourceType(resType), id)

			return nil
		},
	}

	command.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	return command
}

func printResource(g *graph.Graph, resourceType graph.ResourceType, id string) {
	a := graph.Alias(id)
	if aID, ok := a.ResolveToId(g, resourceType); ok {
		id = aID
	}

	resource, err := g.GetResource(resourceType, id)
	if err != nil {
		exitOn(err)
	}

	displayer := display.BuildOptions(
		display.WithHeaders(display.DefaultsColumnDefinitions[resourceType]),
		display.WithFormat(listingFormat),
	).SetSource(resource).Build()

	err = displayer.Print(os.Stderr)
	exitOn(err)
}

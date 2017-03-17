/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/aws"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

var (
	listingFormat      string
	listingFiltersFlag []string
	listOnlyIDs        bool
	sortBy             []string
)

func init() {
	RootCmd.AddCommand(listCmd)

	for _, srvName := range aws.ServiceNames {
		listCmd.AddCommand(listAllResourceInServiceCmd(srvName))
	}

	for _, resType := range aws.ResourceTypes {
		listCmd.AddCommand(listSpecificResourceCmd(resType))
	}

	listCmd.PersistentFlags().StringVar(&listingFormat, "format", "table", "Output format: table, csv, tsv, json (default to table)")
	listCmd.PersistentFlags().StringSliceVar(&listingFiltersFlag, "filter", []string{}, "Filter resources given key/values fields. Ex: --filter type=t2.micro")
	listCmd.PersistentFlags().BoolVar(&listOnlyIDs, "ids", false, "List only ids")
	listCmd.PersistentFlags().StringSliceVar(&sortBy, "sort", []string{"Id"}, "Sort tables by column(s) name(s)")
}

var listCmd = &cobra.Command{
	Use:               "list",
	Aliases:           []string{"ls"},
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook),
	PersistentPostRun: applyHooks(saveHistoryHook, verifyNewVersionHook),
	Short:             "List various type of resources",
}

var listSpecificResourceCmd = func(resType string) *cobra.Command {
	return &cobra.Command{
		Use:   cloud.PluralizeResource(resType),
		Short: fmt.Sprintf("List AWS %s", cloud.PluralizeResource(resType)),

		Run: func(cmd *cobra.Command, args []string) {
			var g *graph.Graph

			if localGlobalFlag {
				if srvName, ok := aws.ServicePerResourceType[resType]; ok {
					g = sync.LoadCurrentLocalGraph(srvName)
				} else {
					exitOn(fmt.Errorf("cannot find service for resource type %s", resType))
				}
			} else {
				srv, err := cloud.GetServiceForType(resType)
				exitOn(err)
				g, err = srv.FetchByType(resType)
				exitOn(err)
			}

			printResources(g, resType)
		},
	}
}

var listAllResourceInServiceCmd = func(srvName string) *cobra.Command {
	return &cobra.Command{
		Use:   srvName,
		Short: fmt.Sprintf("List all %s resources", srvName),

		Run: func(cmd *cobra.Command, args []string) {
			g := sync.LoadCurrentLocalGraph(srvName)
			displayer := console.BuildOptions(
				console.WithFormat(listingFormat),
				console.WithIDsOnly(listOnlyIDs),
			).SetSource(g).Build()
			exitOn(displayer.Print(os.Stdout))
		},
	}
}

func printResources(g *graph.Graph, resType string) {
	displayer := console.BuildOptions(
		console.WithRdfType(resType),
		console.WithHeaders(console.DefaultsColumnDefinitions[resType]),
		console.WithFilters(listingFiltersFlag),
		console.WithMaxWidth(console.GetTerminalWidth()),
		console.WithFormat(listingFormat),
		console.WithIDsOnly(listOnlyIDs),
		console.WithSortBy(sortBy...),
	).SetSource(g).Build()

	exitOn(displayer.Print(os.Stdout))
}

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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

var (
	dryRunSyncFlag        bool
	diffFlag              bool
	syncShowPropetiesFlag bool
	servicesToSyncFlags   map[string]*bool
)

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&dryRunSyncFlag, "dry-run", "d", false, "Display the diff between local and remote cloud, but do not write to disk")
	syncCmd.Flags().BoolVar(&diffFlag, "diff", false, "Display the diff between local and remote cloud, after syncing")
	syncCmd.Flags().BoolVarP(&syncShowPropetiesFlag, "show-properties", "p", false, "Show diff of properties")

	servicesToSyncFlags = make(map[string]*bool)
	for _, service := range aws.ServiceNames {
		servicesToSyncFlags[service] = new(bool)
		syncCmd.Flags().BoolVar(servicesToSyncFlags[service], service, false, fmt.Sprintf("Sync '%s' service only", service))
	}
}

var syncCmd = &cobra.Command{
	Use:                "sync",
	Short:              "Manage your local infrastructure",
	PersistentPreRun:   applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(cmd *cobra.Command, args []string) error {
		var services []cloud.Service
		displayAllServices := true
		for _, srv := range cloud.ServiceRegistry {
			if *servicesToSyncFlags[srv.Name()] {
				displayAllServices = false
			}
		}
		for _, srv := range cloud.ServiceRegistry {
			if displayAllServices || *servicesToSyncFlags[srv.Name()] {
				services = append(services, srv)
			}
		}
		localGraphs := make(map[string]*graph.Graph)
		for _, service := range services {
			localGraphs[service.Name()] = sync.LoadCurrentLocalGraph(service.Name())
		}

		graphs, err := sync.DefaultSyncer.Sync(services...)
		if err != nil {
			return err
		}

		if dryRunSyncFlag && config.AwlessFirstSync {
			exitOn(errors.New("No local data for printing diff. You might want to perfom a full sync first with `awless sync`"))
		}
		if diffFlag || dryRunSyncFlag {
			printFormat := "tree"
			if syncShowPropetiesFlag {
				printFormat = "table"
			}
			region := database.MustGetDefaultRegion()
			root := graph.InitResource(region, graph.Region)
			for _, service := range services {
				diff, err := graph.Differ.Run(root, localGraphs[service.Name()], graphs[service.Name()])
				exitOn(err)
				if diff.HasDiff() {
					fmt.Printf("------%s------\n", strings.ToUpper(service.Name()))
					displayer := console.BuildOptions(
						console.WithFormat(printFormat),
						console.WithRootNode(root),
					).SetSource(diff).Build()
					exitOn(displayer.Print(os.Stdout))
				}
			}
		}

		return nil
	},
}

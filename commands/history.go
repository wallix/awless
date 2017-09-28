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

	"github.com/wallix/awless/aws/services"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/console"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

var (
	showProperties bool
)

func init() {
	RootCmd.AddCommand(historyCmd)

	historyCmd.Flags().BoolVar(&showProperties, "properties", false, "Full diff with resources properties")
}

var historyCmd = &cobra.Command{
	Use:               "history",
	Hidden:            true,
	Short:             "(in progress) Show a infra resource history & changes using your locally sync snapshots",
	PersistentPreRun:  applyHooks(initLoggerHook, initAwlessEnvHook, initCloudServicesHook, initSyncerHook, firstInstallDoneHook),
	PersistentPostRun: applyHooks(verifyNewVersionHook, onVersionUpgrade, networkMonitorHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		region := config.GetAWSRegion()

		root := graph.InitResource(cloud.Region, region)

		var diffs []*sync.Diff

		all, err := sync.DefaultSyncer.List()
		exitOn(err)

		for i := 1; i < len(all); i++ {
			from, err := sync.DefaultSyncer.LoadRev(all[i-1].Id)
			exitOn(err)

			to, err := sync.DefaultSyncer.LoadRev(all[i].Id)
			exitOn(err)

			d, err := sync.BuildDiff(from, to, root.Id())
			exitOn(err)

			diffs = append(diffs, d)
		}

		for _, diff := range diffs {
			displayRevisionDiff(diff, awsservices.InfraService.Name(), root, verboseGlobalFlag)
		}

		return nil
	},
}

func displayRevisionDiff(diff *sync.Diff, cloudService string, root *graph.Resource, verbose bool) {
	fromRevision := "repository creation"
	if diff.From.Id != "" {
		fromRevision = diff.From.Id[:7] + " on " + diff.From.Date.Format("Monday January 2, 15:04")
	}

	var graphdiff *graph.Diff
	if cloudService == awsservices.InfraService.Name() {
		graphdiff = diff.InfraDiff
	}

	if showProperties {
		if graphdiff.HasDiff() {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.Id[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			displayer, err := console.BuildOptions(
				console.WithFormat("table"),
				console.WithRootNode(root),
			).SetSource(graphdiff).Build()
			exitOn(err)
			exitOn(displayer.Print(os.Stdout))
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "properties, from", fromRevision,
				"to", diff.To.Id[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			fmt.Println("No changes.")
		}
	} else {
		if graphdiff.HasDiff() {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.Id[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			displayer, err := console.BuildOptions(
				console.WithFormat("tree"),
				console.WithRootNode(root),
			).SetSource(graphdiff).Build()
			exitOn(err)
			exitOn(displayer.Print(os.Stdout))
			fmt.Println()
		} else if verbose {
			fmt.Println("▶", cloudService, "resources, from", fromRevision,
				"to", diff.To.Id[:7], "on", diff.To.Date.Format("Monday January 2, 15:04"))
			fmt.Println("No resource changes.")
		}
	}
}

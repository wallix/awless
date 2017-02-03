package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	gosync "sync"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync"
)

var diffProperties bool

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().BoolVarP(&diffProperties, "properties", "p", false, "Full diff with resources properties")
}

var diffCmd = &cobra.Command{
	Use:                "diff",
	Short:              "Show diff between your local and remote infra",
	PersistentPreRun:   applyHooks(initAwlessEnvHook, initCloudServicesHook, initSyncerHook, checkStatsHook),
	PersistentPostRunE: saveHistoryHook,

	RunE: func(cmd *cobra.Command, args []string) error {
		if config.AwlessFirstSync {
			exitOn(errors.New("No local data for a diff. You might want to perfom a sync first with `awless sync`"))
		}

		var remoteInfra, remoteAccess *graph.Graph

		var wg gosync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			remoteInfra, err = aws.InfraService.FetchResources()
			exitOn(err)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			remoteAccess, err = aws.AccessService.FetchResources()
			exitOn(err)
		}()

		wg.Wait()

		region := database.MustGetDefaultRegion()

		root := graph.InitResource(region, graph.Region)

		localInfra := sync.LoadCurrentLocalGraph(aws.InfraService.Name())

		infraDiff, err := graph.Differ.Run(root, localInfra, remoteInfra)
		exitOn(err)

		localAccess := sync.LoadCurrentLocalGraph(aws.AccessService.Name())

		accessDiff, err := graph.Differ.Run(root, localAccess, remoteAccess)
		exitOn(err)

		var anyDiffs bool

		if diffProperties {
			if accessDiff.HasDiff() {
				anyDiffs = true
				fmt.Println("------ ACCESS ------")
				displayFullDiff(accessDiff, root)
			}
			if infraDiff.HasDiff() {
				anyDiffs = true
				fmt.Println()
				fmt.Println("------ INFRA ------")
				displayFullDiff(infraDiff, root)
			}
		} else {
			if accessDiff.HasDiff() {
				anyDiffs = true
				fmt.Println("------ ACCESS ------")
				displayer := display.BuildOptions(
					display.WithFormat("tree"),
					display.WithRootNode(root),
				).SetSource(accessDiff).Build()
				exitOn(displayer.Print(os.Stdout))
			}

			if infraDiff.HasDiff() {
				anyDiffs = true
				fmt.Println()
				fmt.Println("------ INFRA ------")
				displayer := display.BuildOptions(
					display.WithFormat("tree"),
					display.WithRootNode(root),
				).SetSource(infraDiff).Build()
				exitOn(displayer.Print(os.Stdout))
			}
		}
		if anyDiffs {
			var yesorno string
			fmt.Print("\nDo you want to perform a sync? (y/n): ")
			fmt.Scanln(&yesorno)
			if strings.TrimSpace(yesorno) == "y" {
				_, err := sync.DefaultSyncer.Sync(aws.InfraService, aws.AccessService)
				exitOn(err)
			}
		}

		return nil
	},
}

func displayFullDiff(diff *graph.Diff, rootNode *graph.Resource) {
	displayer := display.BuildOptions(
		display.WithFormat("table"),
		display.WithRootNode(rootNode),
	).SetSource(diff).Build()
	exitOn(displayer.Print(os.Stdout))
}

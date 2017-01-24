package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/graph"
)

var diffProperties bool

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.PersistentFlags().BoolVarP(&diffProperties, "properties", "p", false, "Full diff with resources properties")
}

var diffCmd = &cobra.Command{
	Use:               "diff",
	Short:             "Show diff between your local and remote infra",
	PersistentPreRun:  initCloudServicesFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(cmd *cobra.Command, args []string) error {
		if config.AwlessFirstSync {
			exitOn(errors.New("No local data for a diff. You might want to perfom a sync first with `awless sync`"))
		}

		var awsInfra *aws.AwsInfra
		var awsAccess *aws.AwsAccess

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			infra, err := aws.InfraService.FetchAwsInfra()
			exitOn(err)
			awsInfra = infra
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			access, err := aws.AccessService.FetchAwsAccess()
			exitOn(err)
			awsAccess = access
		}()

		wg.Wait()

		region := database.MustGetDefaultRegion()

		root, err := node.NewNodeFromStrings(graph.Region.ToRDFString(), region)
		exitOn(err)

		localInfra, err := config.LoadInfraGraph()
		exitOn(err)

		remoteInfra, err := aws.BuildAwsInfraGraph(region, awsInfra)
		exitOn(err)

		infraDiff, err := graph.DefaultDiffer.Run(root, localInfra.Graph, remoteInfra.Graph)
		exitOn(err)

		localAccess, err := config.LoadAccessGraph()
		exitOn(err)

		remoteAccess, err := aws.BuildAwsAccessGraph(region, awsAccess)
		exitOn(err)

		accessDiff, err := graph.DefaultDiffer.Run(root, localAccess.Graph, remoteAccess.Graph)
		exitOn(err)

		var hasDiff bool
		if diffProperties {
			if accessDiff.HasDiff() {
				hasDiff = true
				fmt.Println("------ ACCESS ------")
				display.FullDiff(&graph.Diff{accessDiff}, root, aws.AccessServiceName)
			}
			if infraDiff.HasDiff() {
				hasDiff = true
				fmt.Println()
				fmt.Println("------ INFRA ------")
				display.FullDiff(&graph.Diff{infraDiff}, root, aws.InfraServiceName)
			}
		} else {
			if accessDiff.HasResourceDiff() {
				hasDiff = true
				fmt.Println("------ ACCESS ------")
				display.ResourceDiff(&graph.Diff{accessDiff}, root)
			}

			if infraDiff.HasResourceDiff() {
				hasDiff = true
				fmt.Println()
				fmt.Println("------ INFRA ------")
				display.ResourceDiff(&graph.Diff{infraDiff}, root)
			}
		}
		if hasDiff {
			var yesorno string
			fmt.Print("\nDo you want to perform a sync? (y/n): ")
			fmt.Scanln(&yesorno)
			if strings.TrimSpace(yesorno) == "y" {
				performSync(region)
			}
		}

		return nil
	},
}

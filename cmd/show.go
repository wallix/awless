package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/revision"
)

var (
	numberRevisionsToShow   int
	showRevisionsProperties bool
)

func init() {
	showCmd.AddCommand(showVpcCmd)
	showCmd.AddCommand(showCloudRevisionsCmd)
	showCloudRevisionsCmd.PersistentFlags().IntVarP(&numberRevisionsToShow, "number", "n", 10, "Number of revision to show")
	showCloudRevisionsCmd.PersistentFlags().BoolVarP(&showRevisionsProperties, "properties", "p", false, "Full diff with resources properties")

	RootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show various type of items by id: users, groups, instances, vpcs, ...",
}

var showVpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Show a vpc from a given id",

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("show vpc: id required")
		}
		resp, err := aws.InfraService.Vpc(args[0])
		displayItem(resp, err)
		return nil
	},
}

var showCloudRevisionsCmd = &cobra.Command{
	Use:   "revisions",
	Short: "Show cloud revision history",

	RunE: func(cmd *cobra.Command, args []string) error {
		diffs, err := revision.LastDiffs(config.GitDir, numberRevisionsToShow)
		if err != nil {
			return err
		}
		for _, diff := range diffs {
			displayCommit(diff)
		}
		return nil
	},
}

func displayCommit(diff *revision.CommitDiff) {
	fmt.Println("\tRevision: ", diff.Commit, "- Date: ", diff.Time.Format("Monday January 2, 15:04"))

	root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if showRevisionsProperties {
		if diff.GraphDiff.HasDiff() {
			display.FullDiff(diff.GraphDiff, root)
		} else {
			fmt.Println("No changes.")
		}
	} else {
		if diff.GraphDiff.HasResourceDiff() {
			display.ResourceDiff(diff.GraphDiff, root)
		} else {
			fmt.Println("No resource changes.")
		}
	}
}

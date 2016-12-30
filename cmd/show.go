package cmd

import (
	"errors"
	"fmt"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/revision"
)

var (
	numberRevisionsToShow    int
	showRevisionsProperties  bool
	showRevisionsGroupAll    bool
	showRevisionsGroupByDay  bool
	showRevisionsGroupByWeek bool
)

func init() {
	showCmd.AddCommand(showVpcCmd)
	showCmd.AddCommand(showCloudRevisionsCmd)
	showCloudRevisionsCmd.PersistentFlags().IntVarP(&numberRevisionsToShow, "number", "n", 10, "Number of revision to show")
	showCloudRevisionsCmd.PersistentFlags().BoolVarP(&showRevisionsProperties, "properties", "p", false, "Full diff with resources properties")
	showCloudRevisionsCmd.PersistentFlags().BoolVar(&showRevisionsGroupAll, "group-all", false, "Group all revisions")
	showCloudRevisionsCmd.PersistentFlags().BoolVar(&showRevisionsGroupByWeek, "group-by-week", false, "Group revisions by week")
	showCloudRevisionsCmd.PersistentFlags().BoolVar(&showRevisionsGroupByDay, "group-by-day", false, "Group revisions by day")

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
		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}
		r, err := revision.OpenRepository(config.GitDir)
		if err != nil {
			return err
		}
		param := revision.NoGroup
		if showRevisionsGroupAll {
			param = revision.GroupAll
		}
		if showRevisionsGroupByDay {
			param = revision.GroupByDay
		}
		if showRevisionsGroupByWeek {
			param = revision.GroupByWeek
		}
		accessDiffs, err := r.LastDiffs(numberRevisionsToShow, root, param, config.AccessFilename)
		if err != nil {
			return err
		}
		infraDiffs, err := r.LastDiffs(numberRevisionsToShow, root, param, config.InfraFilename)
		if err != nil {
			return err
		}
		for i := range accessDiffs {
			displayCommit(accessDiffs[i], "Access", root)
			displayCommit(infraDiffs[i], "Infra", root)
		}
		return nil
	},
}

func displayCommit(diff *revision.CommitDiff, commitType string, root *node.Node) {
	parentCommit := diff.ParentCommit
	parentText := "repository creation"
	if parentCommit != "" {
		parentText = parentCommit[:7] + " on " + diff.ParentTime.Format("Monday January 2, 15:04")
	}

	if showRevisionsProperties {
		if diff.GraphDiff.HasDiff() {
			fmt.Println("▶", commitType, "properties, from", parentText,
				"to", diff.ChildCommit[:7], "on", diff.ChildTime.Format("Monday January 2, 15:04"))
			display.FullDiff(diff.GraphDiff, root)
			fmt.Println()
		} else if verboseFlag {
			fmt.Println("▶", commitType, "properties, from", parentText,
				"to", diff.ChildCommit[:7], "on", diff.ChildTime.Format("Monday January 2, 15:04"))
			fmt.Println("No changes.")
		}
	} else {
		if diff.GraphDiff.HasResourceDiff() {
			fmt.Println("▶", commitType, "resources, from", parentText,
				"to", diff.ChildCommit[:7], "on", diff.ChildTime.Format("Monday January 2, 15:04"))
			display.ResourceDiff(diff.GraphDiff, root)
			fmt.Println()
		} else if verboseFlag {
			fmt.Println("▶", commitType, "resources, from", parentText,
				"to", diff.ChildCommit[:7], "on", diff.ChildTime.Format("Monday January 2, 15:04"))
			fmt.Println("No resource changes.")
		}
	}
}

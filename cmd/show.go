package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/display"
	"github.com/wallix/awless/rdf"
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
	//Resources
	for _, resource := range []rdf.ResourceType{rdf.Instance, rdf.Vpc, rdf.Subnet} {
		showCmd.AddCommand(showInfraResourceCmd(resource))
	}
	for _, resource := range []rdf.ResourceType{rdf.User, rdf.Role, rdf.Policy, rdf.Group} {
		showCmd.AddCommand(showAccessResourceCmd(resource))
	}

	//Revisions
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

var showInfraResourceCmd = func(resourceType rdf.ResourceType) *cobra.Command {
	command := &cobra.Command{
		Use:   resourceType.String() + " id",
		Short: "Show the properties of a AWS EC2 " + resourceType.String(),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("id required")
			}
			id := args[0]
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.InfraFilename))

			} else {
				g, err = aws.InfraService.FetchRDFResources(resourceType)
			}
			exitOn(err)
			printResource(g, resourceType, id)
			return nil
		},
	}

	command.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	return command
}

var showAccessResourceCmd = func(resourceType rdf.ResourceType) *cobra.Command {
	command := &cobra.Command{
		Use:   resourceType.String() + " id",
		Short: "Show the properties of a AWS IAM " + resourceType.String(),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("id required")
			}
			id := args[0]
			var g *rdf.Graph
			var err error
			if localResources {
				g, err = rdf.NewGraphFromFile(filepath.Join(config.RepoDir, config.AccessFilename))

			} else {
				g, err = aws.AccessService.FetchRDFResources(resourceType)
			}
			exitOn(err)
			printResource(g, resourceType, id)
			return nil
		},
	}

	command.PersistentFlags().BoolVar(&localResources, "local", false, "List locally sync resources")
	return command
}

func printResource(g *rdf.Graph, resourceType rdf.ResourceType, id string) {
	resource := aws.InitResource(id, resourceType)
	if !resource.ExistsInGraph(g) {
		exitOn(fmt.Errorf("the %s '%s' has not been found", resourceType.String(), id))
	}
	err := resource.UnmarshalFromGraph(g)
	if err != nil {
		exitOn(err)
	}

	var displayer display.ResourceDisplayer
	displayer = display.BuildResourceDisplayer(display.DefaultsColumnDefinitions[resourceType],
		display.Options{Format: listingFormat})
	displayer.SetResource(resource)
	err = displayer.Print(os.Stderr)
	exitOn(err)
}

var showCloudRevisionsCmd = &cobra.Command{
	Use:   "revisions",
	Short: "Show cloud revision history",

	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := node.NewNodeFromStrings("/region", config.GetDefaultRegion())
		if err != nil {
			return err
		}
		r, err := revision.OpenRepository(config.RepoDir)
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
			display.RevisionDiff(accessDiffs[i], aws.AccessServiceName, root, verboseFlag, showRevisionsProperties)
			display.RevisionDiff(infraDiffs[i], aws.InfraServiceName, root, verboseFlag, showRevisionsProperties)
		}
		return nil
	},
}

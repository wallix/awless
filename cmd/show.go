package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/revision"
)

var numberRevisionToShow *int

func init() {
	showCmd.AddCommand(showVpcCmd)
	showCmd.AddCommand(showCloudRevisionsCmd)
	numberRevisionToShow = showCloudRevisionsCmd.Flags().IntP("number", "n", 10, "Number of revision to show")

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
		diffs, err := revision.LastDiffs(config.GitDir, *numberRevisionToShow)
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
	fmt.Println("Id:", diff.Commit, "- Date: ", diff.Time.Format("Monday January 2, 15:04"))

	root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	diff.GraphDiff.FullGraph().VisitDepthFirst(root, func(g *rdf.Graph, n *node.Node, distance int) {
		var nodeBuffer bytes.Buffer
		var lit *literal.Literal
		diffTriples, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
		if len(diffTriples) > 0 && err == nil {
			lit, _ = diffTriples[0].Object().Literal()
		}

		var commonResource bool
		switch lit {
		case rdf.ExtraLiteral:
			color.Set(color.FgGreen)
			fmt.Fprintf(os.Stdout, "+%s, %s\n", n.Type(), n.ID())
			color.Unset()
		case rdf.MissingLiteral:
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s, %s\n", n.Type(), n.ID())
			color.Unset()
		default:
			commonResource = true
			fmt.Fprintf(&nodeBuffer, "%s, %s\n", n.Type(), n.ID())
		}
		if commonResource {
			changedProperties := writeNodeProperties(&nodeBuffer, g, n, diff)
			if changedProperties {
				fmt.Fprint(os.Stdout, nodeBuffer.String())
			}
		}
	})
	if len(diff.GraphDiff.Inserted()) == 0 && len(diff.GraphDiff.Deleted()) == 0 {
		fmt.Println("No changes.")
	}

	fmt.Println("----------------------------------------------------------------------")
}

func writeNodeProperties(writer io.Writer, g *rdf.Graph, n *node.Node, diff *revision.CommitDiff) (hasChanges bool) {
	propertiesT, err := g.TriplesForSubjectPredicate(n, rdf.PropertyPredicate)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}

	for _, t := range propertiesT {
		if diff.GraphDiff.HasInsertedTriple(t) {
			hasChanges = true
			prop, err := aws.NewPropertyFromTriple(t)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return hasChanges
			}
			fmt.Fprint(writer, color.New(color.FgGreen).SprintfFunc()("\t+ %s: %s\n", prop.Key, prop.Value))
		}
		if diff.GraphDiff.HasDeletedTriple(t) {
			hasChanges = true
			prop, err := aws.NewPropertyFromTriple(t)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return hasChanges
			}
			fmt.Fprint(writer, color.New(color.FgRed).SprintfFunc()("\t- %s: %v\n", prop.Key, prop.Value))
		}
	}
	return hasChanges
}

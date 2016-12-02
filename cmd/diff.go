package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
)

func init() {
	RootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff between your local and remote infra",

	RunE: func(cmd *cobra.Command, args []string) error {
		var awsInfra *api.AwsInfra
		var awsAccess *api.AwsAccess

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			infra, err := infraApi.FetchAwsInfra()
			exitOn(err)
			awsInfra = infra
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			access, err := accessApi.FetchAwsAccess()
			exitOn(err)
			awsAccess = access
		}()

		wg.Wait()

		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.InfraFilename))
		if err != nil {
			return err
		}

		remoteInfra, err := rdf.BuildAwsInfraGraph(viper.GetString("region"), awsInfra)
		if err != nil {
			return err
		}

		extras, missings, commons, err := rdf.Compare(viper.GetString("region"), localInfra, remoteInfra)
		if err != nil {
			return err
		}

		rdf.AttachLiteralToAllTriples(extras, rdf.DiffPredicate, rdf.ExtraLiteral)
		rdf.AttachLiteralToAllTriples(missings, rdf.DiffPredicate, rdf.MissingLiteral)

		infraGraph := rdf.NewGraph()
		infraGraph.Merge(extras)
		infraGraph.Merge(missings)
		infraGraph.Merge(commons)

		fmt.Println("------ INFRA ------")
		infraGraph.VisitDepthFirst(root, printWithDiff)

		localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.AccessFilename))
		if err != nil {
			return err
		}

		remoteAccess, err := rdf.BuildAwsAccessGraph(viper.GetString("region"), awsAccess)
		if err != nil {
			return err
		}

		extras, missings, commons, err = rdf.Compare(viper.GetString("region"), localAccess, remoteAccess)
		if err != nil {
			return err
		}

		rdf.AttachLiteralToAllTriples(extras, rdf.DiffPredicate, rdf.ExtraLiteral)
		rdf.AttachLiteralToAllTriples(missings, rdf.DiffPredicate, rdf.MissingLiteral)

		accessGraph := rdf.NewGraph()
		accessGraph.Merge(extras)
		accessGraph.Merge(missings)
		accessGraph.Merge(commons)

		fmt.Println()
		fmt.Println("------ ACCESS ------")
		accessGraph.VisitDepthFirst(root, printWithDiff)

		return nil
	},
}

func printWithDiff(g *rdf.Graph, n *node.Node, distance int) {
	var lit *literal.Literal
	diff, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
	if len(diff) > 0 && err == nil {
		lit, _ = diff[0].Object().Literal()
	}

	var tabs bytes.Buffer
	for i := 0; i < distance; i++ {
		tabs.WriteByte('\t')
	}

	switch lit {
	case rdf.ExtraLiteral:
		color.Set(color.FgGreen)
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		color.Unset()
	case rdf.MissingLiteral:
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		color.Unset()
	default:
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
	}
}

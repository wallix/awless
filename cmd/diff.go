package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud/aws"
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
		if config.AwlessFirstSync {
			return errors.New("No local data for a diff. You might want to perfom a sync first with `awless sync`")
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

		root, err := node.NewNodeFromStrings("/region", viper.GetString("region"))
		if err != nil {
			return err
		}

		localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
		if err != nil {
			return err
		}

		remoteInfra, err := aws.BuildAwsInfraGraph(viper.GetString("region"), awsInfra)
		if err != nil {
			return err
		}

		infraDiffGraph, err := rdf.Diff(root, localInfra, remoteInfra)
		if err != nil {
			return err
		}

		noDiffInfra := graphHasDiff(infraDiffGraph)

		localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
		if err != nil {
			return err
		}

		remoteAccess, err := aws.BuildAwsAccessGraph(viper.GetString("region"), awsAccess)
		if err != nil {
			return err
		}

		accessDiffGraph, err := rdf.Diff(root, localAccess, remoteAccess)
		if err != nil {
			return err
		}

		noDiffAccess := graphHasDiff(accessDiffGraph)

		if !noDiffAccess {
			fmt.Println("------ ACCESS ------")
			accessDiffGraph.VisitDepthFirst(root, printWithDiff)
		}

		if !noDiffInfra {
			fmt.Println()
			fmt.Println("------ INFRA ------")
			infraDiffGraph.VisitDepthFirst(root, printWithDiff)
		}

		if !noDiffInfra || !noDiffAccess {
			var yesorno string
			fmt.Print("\nDo you want to perform a sync? (y/n): ")
			fmt.Scanln(&yesorno)
			if strings.TrimSpace(yesorno) == "y" {
				performSync()
			}
		}

		return nil
	},
}

func graphHasDiff(g *rdf.Graph) bool {
	triples, err := g.TriplesForPredicateName("diff")
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("graph has diff:%s", err.Error()))
		return false
	}
	return len(triples) == 0
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

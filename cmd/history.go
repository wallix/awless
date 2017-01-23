package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/revision/repo"
)

func init() {
	RootCmd.AddCommand(historyCmd)
}

type revPair [2]*repo.Rev

var historyCmd = &cobra.Command{
	Use:               "history",
	Short:             "Show your infrastucture history",
	PersistentPreRun:  initAwlessEnvFn,
	PersistentPostRun: saveHistoryFn,

	RunE: func(cmd *cobra.Command, args []string) error {
		if !repo.IsGitInstalled() {
			fmt.Printf("No history available. You need to install git")
			os.Exit(0)
		}

		rep, err := repo.NewRepo()
		exitOn(err)

		all, err := rep.List()
		exitOn(err)

		root, err := node.NewNodeFromStrings("/region", config.GetDefaultRegion())
		if err != nil {
			return err
		}

		compare := make(chan revPair)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case revPair, ok := <-compare:
					if ok {
						compareRev(root, revPair)
					} else {
						return
					}
				}
			}
		}()

		head, trail := all[0], all[1:]
		head, err = rep.LoadRev(head.Id)
		exitOn(err)

		for _, rev := range trail {
			rev, err := rep.LoadRev(rev.Id)
			exitOn(err)
			compare <- revPair([2]*repo.Rev{head, rev})
			head = rev
		}

		close(compare)

		wg.Wait()

		return nil
	},
}

func compareRev(root *node.Node, revs revPair) {
	rev1, rev2 := revs[0], revs[1]

	infraDiff, err := rdf.NewHierarchicalDiffer().Run(root, rev1.Infra, rev2.Infra)
	exitOn(err)

	accessDiff, err := rdf.NewHierarchicalDiffer().Run(root, rev1.Access, rev2.Access)
	exitOn(err)

	if infraDiff.HasDiff() || accessDiff.HasDiff() {
		fmt.Println(fmt.Sprintf("FROM [%s] TO [%s]", rev1.DateString(), rev2.DateString()))
		if infraDiff.HasDiff() {
			fmt.Println("INFRA:")
			infraDiff.FullGraph().VisitDepthFirst(root, printWithDiff)
		}
		if accessDiff.HasDiff() {
			fmt.Println()
			fmt.Println("ACCESS:")
			accessDiff.FullGraph().VisitDepthFirst(root, printWithDiff)
		}
		fmt.Println()
	}
}

func printWithDiff(g *rdf.Graph, n *node.Node, distance int) {
	var lit *literal.Literal
	diff, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
	if len(diff) > 0 && err == nil {
		lit, _ = diff[0].Object().Literal()
	}

	switch lit {
	case rdf.ExtraLiteral:
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "\t%s, %s\n", n.Type(), n.ID())
		color.Unset()
	case rdf.MissingLiteral:
		color.Set(color.FgGreen)
		fmt.Fprintf(os.Stdout, "\t%s, %s\n", n.Type(), n.ID())
		color.Unset()
	}
}

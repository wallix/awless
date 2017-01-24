package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
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

		root, err := node.NewNodeFromStrings(graph.Region.ToRDFString(), database.MustGetDefaultRegion())
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

	infraDiff, err := graph.HierarchicalDiffer.Run(root, rev1.Infra.Graph, rev2.Infra.Graph)
	exitOn(err)

	accessDiff, err := graph.HierarchicalDiffer.Run(root, rev1.Access.Graph, rev2.Access.Graph)
	exitOn(err)

	if infraDiff.HasDiff() || accessDiff.HasDiff() {
		fmt.Println(fmt.Sprintf("FROM [%s] TO [%s]", rev1.DateString(), rev2.DateString()))
		if infraDiff.HasDiff() {
			fmt.Println("INFRA:")
			g := &graph.Graph{infraDiff.FullGraph()}
			g.Visit(root, printWithDiff)
		}
		if accessDiff.HasDiff() {
			fmt.Println()
			fmt.Println("ACCESS:")
			g := &graph.Graph{accessDiff.FullGraph()}
			g.Visit(root, printWithDiff)
		}
		fmt.Println()
	}
}

func printWithDiff(g *graph.Graph, n *node.Node, distance int) {
	var lit *literal.Literal
	diff, err := g.TriplesInDiff(n)
	if len(diff) > 0 && err == nil {
		lit, _ = diff[0].Object().Literal()
	}

	var litString string
	if lit != nil {
		litString, _ = lit.Text()
	}

	switch litString {
	case "extra":
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "\t%s, %s\n", n.Type(), n.ID())
		color.Unset()
	case "missing":
		color.Set(color.FgGreen)
		fmt.Fprintf(os.Stdout, "\t%s, %s\n", n.Type(), n.ID())
		color.Unset()
	}
}

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

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

var historyCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"who"},
	Short:   "Show your infrastucture history",

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

		compare := make(chan *repo.Rev)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case rev := <-compare:
					select {
					case otherRev, ok := <-compare:
						if ok {
							compareRev(root, rev, otherRev)
						} else {
							return
						}
					}
				case <-time.After(time.Second * 5):
					fmt.Println("done")
					return
				}
			}
		}()

		for _, rev := range all {
			rev, err := rep.LoadRev(rev.Id)
			exitOn(err)
			compare <- rev
		}

		close(compare)

		wg.Wait()

		return nil
	},
}

func compareRev(root *node.Node, rev1, rev2 *repo.Rev) {
	infraDiff, err := rdf.NewHierarchicalDiffer().Run(root, rev1.Infra, rev2.Infra)
	exitOn(err)

	accessDiff, err := rdf.NewHierarchicalDiffer().Run(root, rev1.Access, rev2.Access)
	exitOn(err)

	fmt.Println(fmt.Sprintf("FROM [%s] TO [%s]", rev1.DateString(), rev2.DateString()))
	if !infraDiff.HasDiff() && !accessDiff.HasDiff() {
		fmt.Println("\t\tnone")
	} else {
		if infraDiff.HasDiff() {
			fmt.Println("\t\tINFRA")
			infraDiff.FullGraph().VisitDepthFirst(root, printWithDiff)
		}
		if accessDiff.HasDiff() {
			fmt.Println("\t\tACCESS")
			accessDiff.FullGraph().VisitDepthFirst(root, printWithDiff)
		}
	}
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
	}
}

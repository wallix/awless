package display

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph"
)

// ResourceDiff displays a tree view of a diff with only the changed resources
func ResourceDiff(diff *graph.Diff, rootNode *node.Node) {
	diff.FullGraph().Visit(rootNode, func(g *graph.Graph, n *node.Node, distance int) {
		var lit *literal.Literal
		diff, err := g.TriplesInDiff(n)
		if len(diff) > 0 && err == nil {
			lit, _ = diff[0].Object().Literal()
		}

		var tabs bytes.Buffer
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		var litString string
		if lit != nil {
			litString, _ = lit.Text()
		}

		switch litString {
		case "extra":
			color.Set(color.FgGreen)
			fmt.Fprintf(os.Stdout, "+%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
			color.Unset()
		case "missing":
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stdout, "-%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
			color.Unset()
		default:
			fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), graph.NewResourceType(n.Type()).String(), n.ID())
		}
	})
}

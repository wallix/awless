package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple/node"
)

func TestVisitDepthFirstGraph(t *testing.T) {
	g := NewGraph()

	//       1
	//   2       3       4
	// 5   6   7  8        9
	//                      10
	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")
	three, _ := node.NewNodeFromStrings("/three", "3")
	four, _ := node.NewNodeFromStrings("/four", "4")
	five, _ := node.NewNodeFromStrings("/five", "5")
	six, _ := node.NewNodeFromStrings("/six", "6")
	seven, _ := node.NewNodeFromStrings("/seven", "7")
	eight, _ := node.NewNodeFromStrings("/eight", "8")
	nine, _ := node.NewNodeFromStrings("/nine", "9")
	ten, _ := node.NewNodeFromStrings("/ten", "10")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	g.Add(noErrTriple(one, ParentOfPredicate, three))
	g.Add(noErrTriple(one, ParentOfPredicate, four))
	g.Add(noErrTriple(two, ParentOfPredicate, five))
	g.Add(noErrTriple(two, ParentOfPredicate, six))
	g.Add(noErrTriple(three, ParentOfPredicate, seven))
	g.Add(noErrTriple(three, ParentOfPredicate, eight))
	g.Add(noErrTriple(four, ParentOfPredicate, nine))
	g.Add(noErrTriple(nine, ParentOfPredicate, ten))

	var result bytes.Buffer
	each := func(g *Graph, n *node.Node, distance int) {
		for i := 0; i < distance; i++ {
			result.WriteByte('/')
		}
		result.WriteString(n.ID().String())
	}

	g.VisitDepthFirst(one, each)
	if got, want := result.String(), "1/2//5//6/3//7//8/4//9///10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirst(four, each)
	if got, want := result.String(), "4/9//10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirst(three, each)
	if got, want := result.String(), "3/7/8"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}

func TestGraphSize(t *testing.T) {
	g := NewGraph()
	if got, want := g.IsEmpty(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	g.Add(noErrTriple(two, ParentOfPredicate, one))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

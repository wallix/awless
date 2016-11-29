package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func TestVisitDepthFirstGraph(t *testing.T) {
	g, _ := NewGraph()

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

	noErrTriple := func(s *node.Node, p *predicate.Predicate, o *node.Node) *triple.Triple {
		tri, err := triple.New(s, p, triple.NewNodeObject(o))
		if err != nil {
			t.Fatal(err)
		}
		return tri
	}

	g.Add(noErrTriple(one, parentOf, two))
	g.Add(noErrTriple(one, parentOf, three))
	g.Add(noErrTriple(one, parentOf, four))
	g.Add(noErrTriple(two, parentOf, five))
	g.Add(noErrTriple(two, parentOf, six))
	g.Add(noErrTriple(three, parentOf, seven))
	g.Add(noErrTriple(three, parentOf, eight))
	g.Add(noErrTriple(four, parentOf, nine))
	g.Add(noErrTriple(nine, parentOf, ten))

	var result bytes.Buffer
	each := func(n *node.Node, distance int) {
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

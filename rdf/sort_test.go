package rdf

import (
	"sort"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
)

func TestTripleSorter(t *testing.T) {
	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")
	three, _ := node.NewNodeFromStrings("/three", "3")

	triples := []*triple.Triple{
		noErrTriple(three, ParentOfPredicate, two),
		noErrTriple(one, ParentOfPredicate, two),
		noErrTriple(two, ParentOfPredicate, one),
	}

	sort.Sort(&tripleSorter{triples})

	if got, want := triples[0].Subject().Type().String(), "/one"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := triples[1].Subject().Type().String(), "/three"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := triples[2].Subject().Type().String(), "/two"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestNodeSorter(t *testing.T) {
	nodes := []*node.Node{
		noErrNode("/three", "3"),
		noErrNode("/one", "1"),
		noErrNode("/two", "2"),
	}
	sort.Sort(&nodeSorter{nodes})

	if got, want := nodes[0].ID().String(), "1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := nodes[1].ID().String(), "2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := nodes[2].ID().String(), "3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

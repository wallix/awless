/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple/node"
)

func TestListAttached(t *testing.T) {
	g := NewGraph()
	one, _ := node.NewNodeFromStrings("/any", "1")
	two, _ := node.NewNodeFromStrings("/any", "2")
	three, _ := node.NewNodeFromStrings("/any", "3")
	four, _ := node.NewNodeFromStrings("/any", "4")
	five, _ := node.NewNodeFromStrings("/any", "5")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	g.Add(noErrTriple(one, ParentOfPredicate, three))
	g.Add(noErrTriple(two, ParentOfPredicate, four))
	g.Add(noErrTriple(two, ParentOfPredicate, five))

	verify := func(nodes []*node.Node, expLength int, expected ...*node.Node) {
		if got, want := len(nodes), expLength; got != want {
			t.Fatalf("nodes length: got %d, want %d", got, want)
		}
		for _, ex := range expected {
			var found bool
			for _, n := range nodes {
				if (n.Type().String() == ex.Type().String()) && (n.ID().String() == ex.ID().String()) {
					found = true
				}
			}
			if !found {
				t.Fatalf("%v does not contain %v", nodes, ex)
			}
		}
	}

	nodes, _ := g.ListAttachedTo(one, ParentOfPredicate)
	verify(nodes, 2, two, three)

	nodes, _ = g.ListAttachedTo(two, ParentOfPredicate)
	verify(nodes, 2, four, five)

	nodes, _ = g.ListAttachedFrom(four, ParentOfPredicate)
	verify(nodes, 1, two)

	nodes, _ = g.ListAttachedFrom(two, ParentOfPredicate)
	verify(nodes, 1, one)
}

func TestVisitHierarchically(t *testing.T) {
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

	each := func(g *Graph, n *node.Node, distance int) error {
		for i := 0; i < distance; i++ {
			result.WriteByte('/')
		}
		result.WriteString(n.ID().String())
		return nil
	}

	t.Run("Visit top down", func(t *testing.T) {
		result.Reset()
		g.VisitTopDown(one, each)
		if got, want := result.String(), "1/2//5//6/3//7//8/4//9///10"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}

		result.Reset()
		g.VisitTopDown(four, each)
		if got, want := result.String(), "4/9//10"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}

		result.Reset()
		g.VisitTopDown(three, each)
		if got, want := result.String(), "3/7/8"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}
	})

	t.Run("Visit bottom up", func(t *testing.T) {
		result.Reset()
		g.VisitBottomUp(ten, each)
		if got, want := result.String(), "10/9//4///1"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}

		result.Reset()
		g.VisitBottomUp(eight, each)
		if got, want := result.String(), "8/3//1"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}

		result.Reset()
		g.VisitBottomUp(two, each)
		if got, want := result.String(), "2/1"; got != want {
			t.Fatalf("got '%s', want '%s'", got, want)
		}
	})
}

func TestVisitSiblings(t *testing.T) {
	g := NewGraph()
	//                1
	//     2          3           4
	// 5   6  7    8  9  10    11  12
	//                               13
	one, _ := node.NewNodeFromStrings("/mamal", "1")
	two, _ := node.NewNodeFromStrings("/fish", "2")
	three, _ := node.NewNodeFromStrings("/fish", "3")
	four, _ := node.NewNodeFromStrings("/mamal", "4")
	five, _ := node.NewNodeFromStrings("/nemo", "5")
	six, _ := node.NewNodeFromStrings("/nemo", "6")
	seven, _ := node.NewNodeFromStrings("/nemo", "7")
	eight, _ := node.NewNodeFromStrings("/doris", "8")
	nine, _ := node.NewNodeFromStrings("/doris", "9")
	ten, _ := node.NewNodeFromStrings("/nemo", "10")
	eleven, _ := node.NewNodeFromStrings("/dog", "11")
	twelve, _ := node.NewNodeFromStrings("/cat", "12")
	thirteen, _ := node.NewNodeFromStrings("/any", "13")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	g.Add(noErrTriple(one, ParentOfPredicate, three))
	g.Add(noErrTriple(one, ParentOfPredicate, four))
	g.Add(noErrTriple(two, ParentOfPredicate, five))
	g.Add(noErrTriple(two, ParentOfPredicate, six))
	g.Add(noErrTriple(two, ParentOfPredicate, seven))
	g.Add(noErrTriple(three, ParentOfPredicate, eight))
	g.Add(noErrTriple(three, ParentOfPredicate, nine))
	g.Add(noErrTriple(three, ParentOfPredicate, ten))
	g.Add(noErrTriple(four, ParentOfPredicate, eleven))
	g.Add(noErrTriple(four, ParentOfPredicate, twelve))
	g.Add(noErrTriple(twelve, ParentOfPredicate, thirteen))

	var result bytes.Buffer

	each := func(g *Graph, n *node.Node, distance int) error {
		for i := 0; i < distance; i++ {
			result.WriteByte('/')
		}
		result.WriteString(n.ID().String())
		return nil
	}

	tcases := []struct {
		node *node.Node
		out  string
	}{
		{one, "1"},
		{two, "23"}, {three, "23"}, {four, "4"},
		{five, "567"}, {six, "567"}, {seven, "567"},
		{eight, "89"}, {nine, "89"}, {ten, "10"},
		{eleven, "11"}, {twelve, "12"}, {thirteen, "13"},
	}

	for _, tc := range tcases {
		result.Reset()
		g.VisitSiblings(tc.node, each)
		if got, want := result.String(), tc.out; got != want {
			t.Fatalf("node %s: got '%s', want '%s'", tc.node, got, want)
		}
	}
}

func TestGraphSize(t *testing.T) {
	g := NewGraph()
	if got, want := g.IsEmpty(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), uint32(0); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), uint32(1); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	g.Add(noErrTriple(two, ParentOfPredicate, one))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), uint32(2); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

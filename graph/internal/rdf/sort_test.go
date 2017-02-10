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

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
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func NewTripleFromStrings(sub *node.Node, pred string, obj string) (*triple.Triple, error) {
	p, err := predicate.NewImmutable(pred)
	if err != nil {
		return nil, err
	}
	objL, err := literal.DefaultBuilder().Build(literal.Text, obj)
	if err != nil {
		return nil, err
	}
	return triple.New(sub, p, triple.NewLiteralObject(objL))
}

func attachLiteralToNode(g *Graph, n *node.Node, p *predicate.Predicate, lit *literal.Literal) error {
	tri, err := triple.New(n, p, triple.NewLiteralObject(lit))
	if err != nil {
		return err
	}

	g.Add(tri)
	return nil
}

func intersectTriples(a, b []*triple.Triple) []*triple.Triple {
	var inter []*triple.Triple

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				inter = append(inter, a[i])
			}
		}
	}

	return inter
}

func subtractTriples(a, b []*triple.Triple) []*triple.Triple {
	var sub []*triple.Triple

	for i := 0; i < len(a); i++ {
		var found bool
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				found = true
			}
		}
		if !found {
			sub = append(sub, a[i])
		}
	}

	return sub
}

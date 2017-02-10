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

func noErrNode(typ, id string) *node.Node {
	n, _ := node.NewNodeFromStrings(typ, id)
	return n
}

func noErrTriple(s *node.Node, p *predicate.Predicate, o *node.Node) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewNodeObject(o))
	if err != nil {
		panic(err)
	}
	return tri
}

func noErrLiteralTriple(s *node.Node, p *predicate.Predicate, l *literal.Literal) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewLiteralObject(l))
	if err != nil {
		panic(err)
	}
	return tri
}

func MarshalTriples(triples []*triple.Triple) string {
	g := NewGraphFromTriples(triples)
	return g.MustMarshal()
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}

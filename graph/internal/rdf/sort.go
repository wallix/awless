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
	"github.com/google/badwolf/triple/node"
)

type nodeSorter struct {
	nodes []*node.Node
}

func (s *nodeSorter) Len() int {
	return len(s.nodes)
}
func (s *nodeSorter) Less(i, j int) bool {
	return s.nodes[i].ID().String() < s.nodes[j].ID().String()
}

func (s *nodeSorter) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

type tripleSorter struct {
	triples []*triple.Triple
}

func (s *tripleSorter) Len() int {
	return len(s.triples)
}
func (s *tripleSorter) Less(i, j int) bool {
	return s.triples[i].String() < s.triples[j].String()
}

func (s *tripleSorter) Swap(i, j int) {
	s.triples[i], s.triples[j] = s.triples[j], s.triples[i]
}

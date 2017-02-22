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
	"github.com/google/badwolf/triple/predicate"
)

var DefaultDiffer Differ

type Differ interface {
	Run(*node.Node, *Graph, *Graph) (*Diff, error)
}

type Diff struct {
	fromGraph   *Graph
	toGraph     *Graph
	mergedGraph *Graph
	hasDiffs    bool
}

func NewDiff(fromG, toG *Graph) *Diff {
	return &Diff{fromGraph: fromG, toGraph: toG}
}

func (d *Diff) FromGraph() *Graph {
	return d.fromGraph
}

func (d *Diff) ToGraph() *Graph {
	return d.toGraph
}

func (d *Diff) MergedGraph() *Graph {
	d.mergedGraph = d.toGraph.copy()

	fromTriples, err := d.fromGraph.allTriples()
	if err != nil {
		panic(err)
	}

	for _, fromT := range fromTriples {
		if MetaPredicate.ID() == fromT.Predicate().ID() {
			attachLiteralToNode(d.mergedGraph, fromT.Subject(), MetaPredicate, MissingLiteral)
		} else {
			d.mergedGraph.Add(fromT)
		}
	}

	return d.mergedGraph
}

func (d *Diff) HasDiff() bool {
	return d.hasDiffs
}

type hierarchicDiffer struct {
	predicate *predicate.Predicate
}

func (d *hierarchicDiffer) Run(root *node.Node, from *Graph, to *Graph) (*Diff, error) {
	diff := &Diff{fromGraph: from, toGraph: to}

	maxCount := max(from.size(), to.size())
	processing := make(chan *node.Node, maxCount)

	if maxCount < 1 {
		return diff, nil
	}

	processing <- root

	for len(processing) > 0 {
		select {
		case node := <-processing:
			extras, missings, commons, err := compareChildTriplesOf(d.predicate, node, from, to)
			if err != nil {
				return diff, err
			}

			for _, extra := range extras {
				diff.hasDiffs = true
				node, err := extra.Object().Node()
				if err != nil {
					return diff, err
				}
				attachLiteralToNode(diff.toGraph, node, MetaPredicate, ExtraLiteral)
				processing <- node
			}

			for _, missing := range missings {
				diff.hasDiffs = true
				node, err := missing.Object().Node()
				if err != nil {
					return diff, err
				}
				attachLiteralToNode(diff.fromGraph, node, MetaPredicate, ExtraLiteral)
				processing <- node
			}

			for _, nextNodeToProcess := range commons {
				objectNode, err := nextNodeToProcess.Object().Node()
				if err != nil {
					return diff, err
				}
				processing <- objectNode
			}
		}
	}

	return diff, nil
}

func compareChildTriplesOf(onPredicate *predicate.Predicate, root *node.Node, fromGraph *Graph, toGraph *Graph) ([]*triple.Triple, []*triple.Triple, []*triple.Triple, error) {
	var extras, missings, commons []*triple.Triple

	fromTriples, err := fromGraph.TriplesForSubjectPredicate(root, onPredicate)
	if err != nil {
		return extras, missings, commons, err
	}

	toTriples, err := toGraph.TriplesForSubjectPredicate(root, onPredicate)
	if err != nil {
		return extras, missings, commons, err
	}

	extras = append(extras, subtractTriples(toTriples, fromTriples)...)
	missings = append(missings, subtractTriples(fromTriples, toTriples)...)
	commons = append(commons, intersectTriples(fromTriples, toTriples)...)

	return extras, missings, commons, nil
}

func max(a, b uint32) uint32 {
	if a < b {
		return b
	}

	return a
}

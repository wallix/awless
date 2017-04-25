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

package graph

import (
	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

var (
	DefaultDiffer = &hierarchicDiffer{rdf.ParentOf}
	MetaPredicate = "meta"
)

const (
	extraLit   = "extra"
	missingLit = "missing"
)

type Differ interface {
	Run(string, *Graph, *Graph) (*Diff, error)
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
	d.mergedGraph = NewGraph()
	d.mergedGraph.store.Add(d.toGraph.store.Snapshot().Triples()...)

	fromTriples := d.fromGraph.store.Snapshot().Triples()

	for _, fromT := range fromTriples {
		if MetaPredicate == fromT.Predicate() {
			d.mergedGraph.store.Add(tstore.SubjPred(fromT.Subject(), MetaPredicate).StringLiteral(missingLit))
		} else {
			d.mergedGraph.store.Add(fromT)
		}
	}

	return d.mergedGraph
}

func (d *Diff) HasDiff() bool {
	return d.hasDiffs
}

type hierarchicDiffer struct {
	predicate string
}

func (d *hierarchicDiffer) Run(root string, from *Graph, to *Graph) (*Diff, error) {
	diff := &Diff{fromGraph: from, toGraph: to}

	fromSnap := from.store.Snapshot()
	toSnap := to.store.Snapshot()

	maxCount := max(uint32(fromSnap.Count()), uint32(toSnap.Count()))
	processing := make(chan string, maxCount)

	if maxCount < 1 {
		return diff, nil
	}

	processing <- root

	for len(processing) > 0 {
		select {
		case current := <-processing:
			extras, missings, commons, err := compareChildTriplesOf(d.predicate, current, fromSnap, toSnap)
			if err != nil {
				return diff, err
			}

			for _, extra := range extras {
				res, ok := extra.Object().Resource()
				if ok {
					diff.hasDiffs = true
					diff.toGraph.store.Add(tstore.SubjPred(res, MetaPredicate).StringLiteral(extraLit))
					processing <- res
				}
			}

			for _, missing := range missings {
				res, ok := missing.Object().Resource()
				if ok {
					diff.hasDiffs = true
					diff.fromGraph.store.Add(tstore.SubjPred(res, MetaPredicate).StringLiteral(extraLit))
					processing <- res
				}
			}

			for _, nextNodeToProcess := range commons {
				res, ok := nextNodeToProcess.Object().Resource()
				if ok {
					processing <- res
				}
			}
		}
	}

	return diff, nil
}

func compareChildTriplesOf(onPredicate, root string, fromGraph tstore.RDFGraph, toGraph tstore.RDFGraph) ([]tstore.Triple, []tstore.Triple, []tstore.Triple, error) {
	var extras, missings, commons []tstore.Triple

	fromTriples := fromGraph.WithSubjPred(root, onPredicate)
	toTriples := toGraph.WithSubjPred(root, onPredicate)

	extras = append(extras, subtractTriples(toTriples, fromTriples)...)
	missings = append(missings, subtractTriples(fromTriples, toTriples)...)
	commons = append(commons, intersectTriples(fromTriples, toTriples)...)

	return extras, missings, commons, nil
}

func intersectTriples(a, b []tstore.Triple) []tstore.Triple {
	var inter []tstore.Triple

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].Equal(b[j]) {
				inter = append(inter, a[i])
			}
		}
	}

	return inter
}

func subtractTriples(a, b []tstore.Triple) []tstore.Triple {
	var sub []tstore.Triple

	for i := 0; i < len(a); i++ {
		var found bool
		for j := 0; j < len(b); j++ {
			if a[i].Equal(b[j]) {
				found = true
			}
		}
		if !found {
			sub = append(sub, a[i])
		}
	}

	return sub
}

func max(a, b uint32) uint32 {
	if a < b {
		return b
	}

	return a
}

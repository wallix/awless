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
	d.mergedGraph = d.fromGraph.copy()

	toTriples, err := d.toGraph.allTriples()
	if err != nil {
		panic(err)
	}

	for _, toT := range toTriples {
		if MetaPredicate.ID() == toT.Predicate().ID() {
			attachLiteralToNode(d.mergedGraph, toT.Subject(), MetaPredicate, MissingLiteral)
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
				attachLiteralToNode(diff.fromGraph, node, MetaPredicate, ExtraLiteral)
			}

			for _, missing := range missings {
				diff.hasDiffs = true
				node, err := missing.Object().Node()
				if err != nil {
					return diff, err
				}
				attachLiteralToNode(diff.toGraph, node, MetaPredicate, ExtraLiteral)
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

	extras = append(extras, substractTriples(fromTriples, toTriples)...)
	missings = append(missings, substractTriples(toTriples, fromTriples)...)
	commons = append(commons, intersectTriples(fromTriples, toTriples)...)

	return extras, missings, commons, nil
}

func max(a, b uint32) uint32 {
	if a < b {
		return b
	}

	return a
}

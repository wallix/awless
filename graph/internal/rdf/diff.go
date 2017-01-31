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
	localGraph  *Graph
	remoteGraph *Graph
	mergedGraph *Graph
	hasDiffs    bool
}

func NewDiff(localG, remoteG *Graph) *Diff {
	return &Diff{localGraph: localG, remoteGraph: remoteG}
}

func (d *Diff) LocalGraph() *Graph {
	return d.localGraph
}

func (d *Diff) RemoteGraph() *Graph {
	return d.remoteGraph
}

func (d *Diff) MergedGraph() *Graph {
	d.mergedGraph = d.localGraph.copy()

	remotes, err := d.remoteGraph.allTriples()
	if err != nil {
		panic(err)
	}

	for _, remote := range remotes {
		if MetaPredicate.ID() == remote.Predicate().ID() {
			attachLiteralToNode(d.mergedGraph, remote.Subject(), MetaPredicate, MissingLiteral)
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

func (d *hierarchicDiffer) Run(root *node.Node, local *Graph, remote *Graph) (*Diff, error) {
	diff := &Diff{localGraph: local, remoteGraph: remote}

	maxCount := max(local.size(), remote.size())
	processing := make(chan *node.Node, maxCount)

	if maxCount < 1 {
		return diff, nil
	}

	processing <- root

	for len(processing) > 0 {
		select {
		case node := <-processing:
			extras, missings, commons, err := compareChildTriplesOf(d.predicate, node, local, remote)
			if err != nil {
				return diff, err
			}

			for _, extra := range extras {
				diff.hasDiffs = true
				node, err := extra.Object().Node()
				if err != nil {
					return diff, err
				}
				attachLiteralToNode(diff.localGraph, node, MetaPredicate, ExtraLiteral)
			}

			for _, missing := range missings {
				diff.hasDiffs = true
				node, err := missing.Object().Node()
				if err != nil {
					return diff, err
				}
				attachLiteralToNode(diff.remoteGraph, node, MetaPredicate, ExtraLiteral)
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

func compareChildTriplesOf(onPredicate *predicate.Predicate, root *node.Node, localGraph *Graph, remoteGraph *Graph) ([]*triple.Triple, []*triple.Triple, []*triple.Triple, error) {
	var extras, missings, commons []*triple.Triple

	locals, err := localGraph.TriplesForSubjectPredicate(root, onPredicate)
	if err != nil {
		return extras, missings, commons, err
	}

	remotes, err := remoteGraph.TriplesForSubjectPredicate(root, onPredicate)
	if err != nil {
		return extras, missings, commons, err
	}

	extras = append(extras, substractTriples(locals, remotes)...)
	missings = append(missings, substractTriples(remotes, locals)...)
	commons = append(commons, intersectTriples(locals, remotes)...)

	return extras, missings, commons, nil
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

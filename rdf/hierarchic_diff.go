package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func NewHierarchicalDiffer() Differ {
	return &hierarchicDiffer{ParentOfPredicate}
}

type hierarchicDiffer struct {
	predicate *predicate.Predicate
}

func (d *hierarchicDiffer) Run(root *node.Node, local *Graph, remote *Graph) (*Diff, error) {
	diff := &Diff{graph: NewGraph()}

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
				diff.graph.Add(extra)
				attachLiteralToTriple(diff.graph, extra, DiffPredicate, ExtraLiteral)
			}

			for _, missing := range missings {
				diff.graph.Add(missing)
				attachLiteralToTriple(diff.graph, missing, DiffPredicate, MissingLiteral)
			}

			diff.graph.Add(commons...)

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

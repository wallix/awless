package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
)

func Diff(root *node.Node, local *Graph, remote *Graph) (*Graph, error) {
	diffGraph := NewGraph()

	maxCount := max(local.size(), remote.size())
	processing := make(chan *node.Node, maxCount)

	if maxCount < 1 {
		return diffGraph, nil
	}

	processing <- root

	for len(processing) > 0 {
		select {
		case node := <-processing:
			extras, missings, commons, err := compareChildTriplesOf(node, local, remote)
			if err != nil {
				return diffGraph, err
			}

			for _, extra := range extras {
				diffGraph.Add(extra)
				attachLiteralToTriple(diffGraph, extra, DiffPredicate, ExtraLiteral)
			}

			for _, missing := range missings {
				diffGraph.Add(missing)
				attachLiteralToTriple(diffGraph, missing, DiffPredicate, MissingLiteral)
			}

			diffGraph.Add(commons...)

			for _, nextNodeToProcess := range commons {
				objectNode, err := nextNodeToProcess.Object().Node()
				if err != nil {
					return diffGraph, err
				}
				processing <- objectNode
			}
		}
	}

	return diffGraph, nil
}

func compareChildTriplesOf(root *node.Node, localGraph *Graph, remoteGraph *Graph) ([]*triple.Triple, []*triple.Triple, []*triple.Triple, error) {
	var extras, missings, commons []*triple.Triple

	locals, err := localGraph.TriplesForSubjectPredicate(root, ParentOfPredicate)
	if err != nil {
		return extras, missings, commons, err
	}

	remotes, err := remoteGraph.TriplesForSubjectPredicate(root, ParentOfPredicate)
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

package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
)

func Compare(rootID string, local *Graph, remote *Graph) (*Graph, *Graph, *Graph, error) {
	allextras := NewGraph()
	allmissings := NewGraph()
	allcommons := NewGraph()

	rootNode, err := node.NewNodeFromStrings(REGION, rootID)
	if err != nil {
		return allextras, allmissings, allcommons, err
	}

	maxCount := max(local.size(), remote.size())
	processing := make(chan *node.Node, maxCount)

	processing <- rootNode

	for len(processing) > 0 {
		select {
		case node := <-processing:
			extras, missings, commons, err := compareChildTriplesOf(node, local, remote)
			if err != nil {
				return allextras, allmissings, allcommons, err
			}

			allextras.Add(extras...)
			allmissings.Add(missings...)
			allcommons.Add(commons...)

			for _, nextNodeToProcess := range commons {
				objectNode, err := nextNodeToProcess.Object().Node()
				if err != nil {
					return allextras, allmissings, allcommons, err
				}
				processing <- objectNode
			}
		}
	}

	return allextras, allmissings, allcommons, nil
}

func compareChildTriplesOf(root *node.Node, localGraph *Graph, remoteGraph *Graph) ([]*triple.Triple, []*triple.Triple, []*triple.Triple, error) {
	var extras, missings, commons []*triple.Triple

	locals, err := localGraph.TriplesForSubjectPredicate(root, ParentOf)
	if err != nil {
		return extras, missings, commons, err
	}

	remotes, err := remoteGraph.TriplesForSubjectPredicate(root, ParentOf)
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

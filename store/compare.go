package store

import (
	"context"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/storage/memory"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func Compare(rootID string, local []*triple.Triple, remote []*triple.Triple) ([]*triple.Triple, []*triple.Triple, error) {
	var extras []*triple.Triple
	var missings []*triple.Triple
	var commons []*triple.Triple

	ctx := context.Background()
	localGraph, err := memory.DefaultStore.NewGraph(ctx, "local")
	if err != nil {
		return extras, missings, err
	}
	defer memory.DefaultStore.DeleteGraph(ctx, "local")
	localGraph.AddTriples(ctx, local)

	remoteGraph, err := memory.DefaultStore.NewGraph(ctx, "remote")
	if err != nil {
		return extras, missings, err
	}
	defer memory.DefaultStore.DeleteGraph(ctx, "remote")
	remoteGraph.AddTriples(ctx, remote)

	rootNode, err := node.NewNodeFromStrings("/region", rootID)

	if err != nil {
		return extras, missings, err
	}

	maxCount := max(len(local), len(remote))
	nodesCh := make(chan *node.Node, maxCount)
	nodesCh <- rootNode

	for len(nodesCh) > 0 {
		select {
		case node := <-nodesCh:
			processedExtras, processedMissings, processedCommmons, err := compareTriplesFromNode(node, localGraph, remoteGraph)
			if err != nil {
				return extras, missings, err
			}
			extras = append(extras, processedExtras...)
			missings = append(missings, processedMissings...)
			commons = append(commons, processedCommmons...)
			for _, nextNodeToProcess := range processedCommmons {
				objectNode, err := nextNodeToProcess.Object().Node()
				if err != nil {
					return extras, missings, err
				}
				nodesCh <- objectNode
			}
		}
	}

	return extras, missings, nil
}

func compareTriplesFromNode(rootNode *node.Node, localGraph storage.Graph, remoteGraph storage.Graph) ([]*triple.Triple, []*triple.Triple, []*triple.Triple, error) {
	var extras []*triple.Triple
	var missings []*triple.Triple
	var commons []*triple.Triple

	localTriples, err := triplesForSubjectAndPredicate(localGraph, rootNode, parentOf)
	if err != nil {
		return extras, missings, commons, err
	}

	remoteTriples, err := triplesForSubjectAndPredicate(remoteGraph, rootNode, parentOf)
	if err != nil {
		return extras, missings, commons, err
	}

	extras = append(extras, SubstractTriples(localTriples, remoteTriples)...)
	missings = append(missings, SubstractTriples(remoteTriples, localTriples)...)
	commons = append(commons, IntersectTriples(localTriples, remoteTriples)...)

	return extras, missings, commons, nil
}

func triplesForSubjectAndPredicate(graph storage.Graph, subject *node.Node, predicate *predicate.Predicate) ([]*triple.Triple, error) {
	errc := make(chan error)
	tric := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- graph.TriplesForSubjectAndPredicate(context.Background(), subject, predicate, storage.DefaultLookup, tric)
	}()

	var triples []*triple.Triple

OutLoop:
	for {
		select {
		case e := <-errc:
			if e != nil {
				return triples, e
			}
		case t, ok := <-tric:
			if !ok {
				break OutLoop
			}
			triples = append(triples, t)
		}
	}
	return triples, nil
}

func max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

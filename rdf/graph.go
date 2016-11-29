package rdf

import (
	"context"
	"sort"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/storage/memory"
	"github.com/google/badwolf/triple/node"
)

func newMemGraph(name string) (storage.Graph, error) {
	g, err := memory.DefaultStore.NewGraph(context.Background(), name)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func visitDepthFirst(g storage.Graph, root *node.Node, each func(*node.Node, int), distances ...int) error {
	var dist int
	if len(distances) == 0 {
		dist = 1
	} else {
		dist = distances[0]
	}

	each(root, dist)

	relations, err := triplesForSubjectAndPredicate(g, root, parentOf)
	if err != nil {
		return err
	}

	var childs []*node.Node
	for _, relation := range relations {
		n, err := relation.Object().Node()
		if err != nil {
			return err
		}
		childs = append(childs, n)
	}

	sort.Sort(&nodeSorter{childs})

	for _, child := range childs {
		visitDepthFirst(g, child, each, dist+1)
	}

	return nil
}

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

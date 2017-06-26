package triplestore

import (
	"fmt"
	"sort"
)

// A tree is defined from a RDF Graph
// when given a specific predicate as an edge and
// considering triples pointing to RDF resource Object
//
// The tree defined by the graph/predicate should have no cycles
// and node should have at most one parent
type Tree struct {
	g         RDFGraph
	predicate string
}

func NewTree(g RDFGraph, pred string) *Tree {
	if g == nil {
		panic("given RDF graph is nil")
	}
	return &Tree{g: g, predicate: pred}
}

// Traverse the tree in pre-order depth first search
func (t *Tree) TraverseDFS(node string, each func(RDFGraph, string, int) error, depths ...int) error {
	var depth int
	if len(depths) > 0 {
		depth = depths[0]
	}

	if err := each(t.g, node, depth); err != nil {
		return err
	}

	triples := t.g.WithSubjPred(node, t.predicate)

	var childs []string
	for _, tri := range triples {
		n, ok := tri.Object().Resource()
		if !ok {
			return fmt.Errorf("object is not a resource identifier")
		}
		childs = append(childs, n)
	}

	sort.Strings(childs)

	for _, child := range childs {
		t.TraverseDFS(child, each, depth+1)
	}

	return nil
}

// Traverse all ancestors from the given node
func (t *Tree) TraverseAncestors(node string, each func(RDFGraph, string, int) error, depths ...int) error {
	var depth int
	if len(depths) > 0 {
		depth = depths[0]
	}

	if err := each(t.g, node, depth); err != nil {
		return err
	}

	triples := t.g.WithPredObj(t.predicate, Resource(node))

	var parents []string
	for _, tri := range triples {
		parents = append(parents, tri.Subject())
	}

	sort.Strings(parents)

	for _, parent := range parents {
		t.TraverseAncestors(parent, each, depth+1)
	}

	return nil
}

// Traverse siblings of given node. Passed function allow to output the sibling criteria
func (t *Tree) TraverseSiblings(node string, siblingCriteriaFunc func(RDFGraph, string) (string, error), each func(RDFGraph, string, int) error) error {
	triples := t.g.WithPredObj(t.predicate, Resource(node))

	if len(triples) == 0 {
		return each(t.g, node, 0)
	}

	if len(triples) != 1 {
		return fmt.Errorf("tree[%s]: node %s with more than 1 parent: %v", t.predicate, node, triples)
	}

	otherChildTriples := t.g.WithSubjPred(triples[0].Subject(), t.predicate)

	var childs []string
	for _, c := range otherChildTriples {
		child, ok := c.Object().Resource()
		if !ok {
			return fmt.Errorf("object is not a resource identifier")
		}
		childs = append(childs, child)
	}

	sort.Strings(childs)

	nodeCriteria, err := siblingCriteriaFunc(t.g, node)
	if err != nil {
		return err
	}

	for _, child := range childs {
		childCriteria, err := siblingCriteriaFunc(t.g, child)
		if err != nil {
			return err
		}
		if nodeCriteria == childCriteria {
			if err := each(t.g, child, 0); err != nil {
				return err
			}
		}
	}

	return nil
}

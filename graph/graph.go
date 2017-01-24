package graph

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Graph struct {
	*rdf.Graph
}

func NewGraph() *Graph {
	return &Graph{rdf.NewGraph()}
}

func NewGraphFromFile(filepath string) (*Graph, error) {
	g, err := rdf.NewGraphFromFile(filepath)
	return &Graph{g}, err
}

func NewParentOfTriple(subject, obj *node.Node) (*triple.Triple, error) {
	return triple.New(subject, rdf.ParentOfPredicate, triple.NewNodeObject(obj))
}

func NewRegionTypeTriple(subject *node.Node) (*triple.Triple, error) {
	return triple.New(subject, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
}

func (g *Graph) Visit(root *node.Node, each func(*Graph, *node.Node, int), distances ...int) error {
	foreach := func(g *rdf.Graph, n *node.Node, i int) {
		each(&Graph{g}, n, i)
	}

	return g.VisitDepthFirst(root, foreach, distances...)
}

func (g *Graph) VisitUnique(root *node.Node, each func(*Graph, *node.Node, int) error) error {
	foreach := func(g *rdf.Graph, n *node.Node, i int) error {
		return each(&Graph{g}, n, i)
	}

	return g.VisitDepthFirstUnique(root, foreach)
}

func (g *Graph) CountChildrenOfTypeForNode(node *node.Node, childType ResourceType) (int, error) {
	return g.CountTriplesForSubjectAndPredicateObjectOfType(node, rdf.ParentOfPredicate, childType.ToRDFString())
}

func (g *Graph) CountChildrenForNode(node *node.Node) (int, error) {
	return g.CountTriplesForSubjectAndPredicate(node, rdf.ParentOfPredicate)
}

func (g *Graph) TriplesInDiff(node *node.Node) ([]*triple.Triple, error) {
	return g.TriplesForSubjectPredicate(node, rdf.DiffPredicate)
}

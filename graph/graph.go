package graph

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type CloudGraph interface {
	GetResource(t ResourceType, id string) (*Resource, error)
	GetAllResources(t ResourceType) ([]*Resource, error)
}

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

func (g *Graph) GetResource(t ResourceType, id string) (*Resource, error) {
	resource := InitResource(id, t)

	node, err := resource.BuildRdfSubject()
	if err != nil {
		return resource, err
	}

	propsTriples, err := g.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return resource, err
	}

	for _, t := range propsTriples {
		prop, err := NewPropertyFromTriple(t)
		if err != nil {
			return resource, err
		}
		resource.properties[prop.Key] = prop.Value
	}

	return resource, nil
}

func (g *Graph) GetAllResources(t ResourceType) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.NodesForType(t.ToRDFString())
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r, err := g.GetResource(t, node.ID().String())
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
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

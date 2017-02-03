package graph

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Graph struct {
	rdfG *rdf.Graph
}

func NewGraph() *Graph {
	return &Graph{rdf.NewGraph()}
}

func NewGraphFromFile(filepath string) (*Graph, error) {
	g, err := rdf.NewGraphFromFile(filepath)
	return &Graph{g}, err
}

func (g *Graph) AddResource(resources ...*Resource) error {
	for _, res := range resources {
		triples, err := res.marshalRDF()
		if err != nil {
			return err
		}

		g.rdfG.Add(triples...)
	}
	return nil
}

func (g *Graph) AddGraph(graph *Graph) {
	g.rdfG.AddGraph(graph.rdfG)
}

func (g *Graph) AddParent(parent, child *Resource) error {
	n, err := child.toRDFNode()
	if err != nil {
		return err
	}

	parentN, err := node.NewNodeFromStrings(parent.Type().ToRDFString(), parent.Id())
	if err != nil {
		return err
	}

	t, err := triple.New(parentN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
	if err != nil {
		return err
	}

	g.rdfG.Add(t)

	return nil
}

func (g *Graph) Unmarshal(data []byte) error {
	return g.rdfG.Unmarshal(data)
}

func (g *Graph) MustMarshal() string {
	return g.rdfG.MustMarshal()
}

func (g *Graph) Marshal() ([]byte, error) {
	return g.rdfG.Marshal()
}

func (g *Graph) GetResource(t ResourceType, id string) (*Resource, error) {
	resource := InitResource(id, t)

	node, err := resource.toRDFNode()
	if err != nil {
		return resource, err
	}

	propsTriples, err := g.rdfG.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return resource, err
	}
	if err := resource.Properties.unmarshalRDF(propsTriples); err != nil {
		return resource, err
	}

	metaTriples, err := g.rdfG.TriplesForSubjectPredicate(node, rdf.MetaPredicate)
	if err != nil {
		return resource, err
	}
	if err := resource.Meta.unmarshalRDF(metaTriples); err != nil {
		return resource, err
	}

	return resource, nil
}

func (g *Graph) GetAllResources(t ResourceType) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.rdfG.NodesForType(t.ToRDFString())
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

func (g *Graph) VisitChildren(root *Resource, each func(*Resource, int)) error {
	rootNode, err := root.toRDFNode()
	if err != nil {
		return err
	}

	foreach := func(rdfG *rdf.Graph, n *node.Node, i int) {
		res, err := g.GetResource(newResourceType(n), n.ID().String())
		if err != nil {
			panic(err)
		}
		each(res, i)
	}

	return g.rdfG.VisitDepthFirst(rootNode, foreach)
}

func (g *Graph) CountChildrenOfTypeForNode(res *Resource, childType ResourceType) (int, error) {
	n, err := node.NewNodeFromStrings(res.Type().ToRDFString(), res.Id())
	if err != nil {
		return 0, err
	}
	return g.rdfG.CountTriplesForSubjectAndPredicateObjectOfType(n, rdf.ParentOfPredicate, childType.ToRDFString())
}

func (g *Graph) CountChildrenForNode(res *Resource) (int, error) {
	n, err := node.NewNodeFromStrings(res.Type().ToRDFString(), res.Id())
	if err != nil {
		return 0, err
	}
	return g.rdfG.CountTriplesForSubjectAndPredicate(n, rdf.ParentOfPredicate)
}

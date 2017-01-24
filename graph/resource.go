package graph

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Resource struct {
	kind       ResourceType
	id         string
	properties Properties
}

func InitResource(id string, kind ResourceType) *Resource {
	return &Resource{id: id, kind: kind, properties: make(Properties)}
}

func InitFromRdfNode(n *node.Node) *Resource {
	return InitResource(n.ID().String(), NewResourceType(n.Type()))
}

func (res *Resource) Properties() Properties {
	return res.properties
}

func (res *Resource) Type() ResourceType {
	return res.kind
}

func (res *Resource) Id() string {
	return res.id
}

func (res *Resource) ExistsInGraph(g *Graph) bool {
	r, err := res.BuildRdfSubject()
	if err != nil {
		return false
	}

	nodes, err := g.NodesForType(res.kind.ToRDFString())
	if err != nil {
		return false
	}
	for _, n := range nodes {
		if n.UUID().String() == r.UUID().String() {
			return true
		}
	}
	return false
}

func (res *Resource) BuildRdfSubject() (*node.Node, error) {
	return node.NewNodeFromStrings(res.kind.ToRDFString(), res.id)
}

func (res *Resource) UnmarshalFromGraph(g *Graph) error {
	node, err := res.BuildRdfSubject()
	if err != nil {
		return err
	}

	triples, err := g.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return err
	}

	for _, t := range triples {
		prop, err := NewPropertyFromTriple(t)
		if err != nil {
			return err
		}
		res.properties[prop.Key] = prop.Value
	}

	return nil
}

func (res *Resource) MarshalToTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	n, err := res.BuildRdfSubject()
	if err != nil {
		return triples, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, res.kind.ToRDFString()); err != nil {
		return triples, err
	}
	t, err := triple.New(n, rdf.HasTypePredicate, triple.NewLiteralObject(lit))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	for propKey, propValue := range res.properties {
		prop := Property{Key: propKey, Value: propValue}
		if propT, err := prop.tripleFromNode(n); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	return triples, nil
}

func LoadResourcesFromGraph(g *Graph, t ResourceType) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.NodesForType(t.ToRDFString())
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r := InitResource(node.ID().String(), t)
		if err := r.UnmarshalFromGraph(g); err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
}

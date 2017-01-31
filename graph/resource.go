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
	Properties Properties
}

func InitResource(id string, kind ResourceType) *Resource {
	return &Resource{id: id, kind: kind, Properties: make(Properties)}
}

func (res *Resource) Type() ResourceType {
	return res.kind
}

func (res *Resource) Id() string {
	return res.id
}

func (res *Resource) toRDFNode() (*node.Node, error) {
	return node.NewNodeFromStrings(res.kind.ToRDFString(), res.id)
}

func (res *Resource) marshalToTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	n, err := res.toRDFNode()
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

	for propKey, propValue := range res.Properties {
		prop := Property{Key: propKey, Value: propValue}
		if propT, err := prop.tripleFromNode(n); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	return triples, nil
}

type Properties map[string]interface{}

func (props Properties) unmarshalRDF(triples []*triple.Triple) error {
	for _, tr := range triples {
		prop := &Property{}
		if err := prop.unmarshalRDF(tr); err != nil {
			return err
		}
		props[prop.Key] = prop.Value
	}

	return nil
}
package cloud

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

func PropertyTriple(subject *node.Node, propertyKey string, propertyValue interface{}) (*triple.Triple, error) {
	prop := Property{Key: propertyKey, Value: propertyValue}
	json, err := json.Marshal(prop)
	if err != nil {
		return nil, err
	}
	var propL *literal.Literal
	if propL, err = literal.DefaultBuilder().Build(literal.Text, string(json)); err != nil {
		return nil, err
	}
	if propT, err := triple.New(subject, rdf.PropertyPredicate, triple.NewLiteralObject(propL)); err != nil {
		return nil, err
	} else {
		return propT, nil
	}
}

func LoadPropertiesTriples(g *rdf.Graph, node *node.Node) (Properties, error) {
	triples, err := g.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return nil, err
	}
	properties := make(Properties)
	for _, t := range triples {
		oL, e := t.Object().Literal()
		if e != nil {
			return properties, e
		}
		propStr, e := oL.Text()
		if e != nil {
			return properties, e
		}
		var prop Property
		e = json.Unmarshal([]byte(propStr), &prop)
		if e != nil {
			return properties, e
		}
		properties[prop.Key] = prop.Value
	}
	return properties, nil
}

func AddNodeWithPropertiesToTriples(nodeType, id string, resource interface{}, resourceProperties map[string]map[string]string, triples *[]*triple.Triple) (*node.Node, error) {
	n, err := node.NewNodeFromStrings(nodeType, id)
	if err != nil {
		return nil, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, nodeType); err != nil {
		return nil, err
	}
	t, err := triple.New(n, rdf.HasTypePredicate, triple.NewLiteralObject(lit))
	if err != nil {
		return nil, err
	}
	*triples = append(*triples, t)

	value := reflect.ValueOf(resource)
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, fmt.Errorf("can not fetch cloud resource. %v is not a valid pointer.", value)
	}

	nodeV := value.Elem()

	for propertyId, cloudId := range resourceProperties[nodeType] {
		sourceField := nodeV.FieldByName(cloudId)
		if sourceField.IsValid() && !sourceField.IsNil() {
			if propT, err := PropertyTriple(n, propertyId, sourceField.Interface()); err != nil {
				return nil, err
			} else {
				*triples = append(*triples, propT)
			}
		}
	}

	return n, nil
}

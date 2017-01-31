package graph

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Property struct {
	Key   string
	Value interface{}
}

func (prop *Property) tripleFromNode(subject *node.Node) (*triple.Triple, error) {
	propL, err := prop.marshalRDF()
	if err != nil {
		return nil, err
	}
	if propT, err := triple.New(subject, rdf.PropertyPredicate, propL); err != nil {
		return nil, err
	} else {
		return propT, nil
	}
}

func (prop *Property) marshalRDF() (*triple.Object, error) {
	json, err := json.Marshal(prop)
	if err != nil {
		return nil, err
	}
	var propL *literal.Literal
	if propL, err = literal.DefaultBuilder().Build(literal.Text, string(json)); err != nil {
		return nil, err
	}
	return triple.NewLiteralObject(propL), nil
}

func (prop *Property) unmarshalRDF(t *triple.Triple) error {
	if t.Predicate().String() != rdf.PropertyPredicate.String() {
		return fmt.Errorf("unmarshaling property: triple expected property predicate got '%s'", t.Predicate().String())
	}

	oL, err := t.Object().Literal()
	if err != nil {
		return err
	}
	propStr, err := oL.Text()
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(propStr), prop); err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(strings.ToLower(prop.Key), "time"), strings.HasSuffix(strings.ToLower(prop.Key), "date"):
		t, err := time.Parse(time.RFC3339, fmt.Sprint(prop.Value))
		if err == nil {
			prop.Value = t
		}
	case strings.HasSuffix(strings.ToLower(prop.Key), "rules"):
		var propRules struct {
			Key   string
			Value []*FirewallRule
		}
		err = json.Unmarshal([]byte(propStr), &propRules)
		if err == nil {
			prop.Value = propRules.Value
		}
	case strings.HasSuffix(strings.ToLower(prop.Key), "routes"):
		var propRoutes struct {
			Key   string
			Value []*Route
		}
		err = json.Unmarshal([]byte(propStr), &propRoutes)
		if err == nil {
			prop.Value = propRoutes.Value
		}
	}

	return nil
}

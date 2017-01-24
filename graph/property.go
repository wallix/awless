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

type Properties map[string]interface{}

type Property struct {
	Key   string
	Value interface{}
}

func (prop *Property) tripleFromNode(subject *node.Node) (*triple.Triple, error) {
	propL, err := prop.ToLiteralObject()
	if err != nil {
		return nil, err
	}
	if propT, err := triple.New(subject, rdf.PropertyPredicate, propL); err != nil {
		return nil, err
	} else {
		return propT, nil
	}
}

func (prop *Property) ToLiteralObject() (*triple.Object, error) {
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

func NewPropertyFromTriple(t *triple.Triple) (*Property, error) {
	if t.Predicate().String() != rdf.PropertyPredicate.String() {
		return nil, fmt.Errorf("This triple has not a property prediate: %s", t.Predicate().String())
	}
	oL, err := t.Object().Literal()
	if err != nil {
		return nil, err
	}
	propStr, err := oL.Text()
	if err != nil {
		return nil, err
	}

	var prop Property
	err = json.Unmarshal([]byte(propStr), &prop)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasSuffix(strings.ToLower(prop.Key), "time"), strings.HasSuffix(strings.ToLower(prop.Key), "date"):
		t, err := time.Parse(time.RFC3339, fmt.Sprint(prop.Value))
		if err == nil {
			prop.Value = t
		}
	}

	return &prop, nil
}

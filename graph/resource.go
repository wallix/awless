/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Resource struct {
	kind       ResourceType
	id         string
	Properties Properties
	Meta       Properties
}

func InitResource(id string, kind ResourceType) *Resource {
	return &Resource{id: id, kind: kind, Properties: make(Properties), Meta: make(Properties)}
}

func newResourceType(n *node.Node) ResourceType {
	typ := strings.TrimPrefix(n.Type().String(), "/")
	return ResourceType(typ)
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

func (res *Resource) marshalRDF() ([]*triple.Triple, error) {
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
		propL, err := prop.marshalRDF()
		if err != nil {
			return nil, err
		}
		if propT, err := triple.New(n, rdf.PropertyPredicate, propL); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	for metaKey, metaValue := range res.Meta {
		prop := Property{Key: metaKey, Value: metaValue}
		propL, err := prop.marshalRDF()
		if err != nil {
			return nil, err
		}
		if propT, err := triple.New(n, rdf.MetaPredicate, propL); err != nil {
			return nil, err
		} else {
			triples = append(triples, propT)
		}
	}

	return triples, nil
}

type Properties map[string]interface{}

func (props Properties) Substract(other Properties) Properties {
	sub := make(Properties)

	for propK, propV := range props {
		var found bool
		if otherV, ok := other[propK]; ok {
			if reflect.DeepEqual(propV, otherV) {
				found = true
			}
		}
		if !found {
			sub[propK] = propV
		}
	}

	return sub
}

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

type Property struct {
	Key   string
	Value interface{}
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
	oL, err := t.Object().Literal()
	if err != nil {
		return err
	}
	propStr, err := oL.Text()
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(propStr), prop); err != nil {
		fmt.Println("cannot unmarshal %s: %s\n", propStr, err)
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
	case strings.HasSuffix(strings.ToLower(prop.Key), "grants"):
		var propGrants struct {
			Key   string
			Value []*Grant
		}
		err = json.Unmarshal([]byte(propStr), &propGrants)
		if err == nil {
			prop.Value = propGrants.Value
		}
	}

	return nil
}

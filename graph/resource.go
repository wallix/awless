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
	"strconv"
	"strings"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Resource struct {
	kind, id string

	Properties Properties
	Meta       Properties
}

func InitResource(id string, kind string) *Resource {
	return &Resource{id: id, kind: kind, Properties: make(Properties), Meta: make(Properties)}
}

func newResourceType(n *node.Node) string {
	return strings.TrimPrefix(n.Type().String(), "/")
}

func (res *Resource) String() string {
	var identifier string
	if res == nil || (res.Id() == "" && res.Type() == "") {
		return "[none]"
	}
	if name, ok := res.Properties["Name"]; ok && name != "" {
		identifier = fmt.Sprintf("@%v", name)
	} else {
		identifier = res.Id()
	}
	return fmt.Sprintf("%s[%s]", identifier, res.Type())
}

func (res *Resource) Type() string {
	return res.kind
}

func (res *Resource) Id() string {
	return res.id
}

// Compare only the id and type of the resources (no properties nor meta)
func (res *Resource) Same(other *Resource) bool {
	if res == nil && other == nil {
		return true
	}
	if res == nil || other == nil {
		return false
	}
	return res.Id() == other.Id() && res.Type() == other.Type()
}

func (res *Resource) toRDFNode() (*node.Node, error) {
	return node.NewNodeFromStrings("/"+res.kind, res.id)
}

func (res *Resource) marshalRDF() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	n, err := res.toRDFNode()
	if err != nil {
		return triples, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, "/"+res.kind); err != nil {
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

type Resources []*Resource

func (res Resources) Map(f func(*Resource) string) (out []string) {
	for _, r := range res {
		out = append(out, f(r))
	}
	return
}

type Properties map[string]interface{}

func (props Properties) Subtract(other Properties) Properties {
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
		fmt.Printf("cannot unmarshal %s: %s\n", propStr, err)
	}

	switch {
	case strings.HasSuffix(strings.ToLower(prop.Key), "time"), strings.HasSuffix(strings.ToLower(prop.Key), "date"):
		t, err := time.Parse(time.RFC3339, fmt.Sprint(prop.Value))
		if err == nil {
			prop.Value = t
		}
	case strings.HasSuffix(strings.ToLower(prop.Key), "timestamp"):
		tmstp, err := strconv.Atoi(fmt.Sprint(prop.Value))
		if err == nil {
			prop.Value = time.Unix(int64(tmstp), 0)
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

type ResourceById []*Resource

func (r ResourceById) Len() int           { return len(r) }
func (r ResourceById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResourceById) Less(i, j int) bool { return r[i].Id() < r[j].Id() }

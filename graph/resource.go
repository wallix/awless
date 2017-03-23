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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	cloudrdf "github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Resource struct {
	kind, id string

	Properties Properties
	Meta       Properties
}

func InitResource(kind, id string) *Resource {
	return &Resource{id: id, kind: kind, Properties: make(Properties), Meta: make(Properties)}
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
	return node.NewNodeFromStrings("/node", res.id)
}

func (res *Resource) marshalFullRDF() ([]*triple.Triple, error) {
	var triples []*triple.Triple

	triples = append(triples, rdf.Subject(res.id).Predicate(cloudrdf.RdfType).Object(strings.Title(res.Type()), cloudrdf.CloudOwlNS))

	for key, value := range res.Meta {
		if key == "diff" {
			triples = append(triples, rdf.Subject(res.id).Predicate(string(rdf.MetaPredicate.ID())).Literal(fmt.Sprint(value)))
		}
	}

	for key, value := range res.Properties {
		propId, err := getPropertyRDFId(key)
		if err != nil {
			return triples, fmt.Errorf("marshalling property: %s", err)
		}

		propType, err := getPropertyDefinedBy(propId)
		if err != nil {
			return triples, fmt.Errorf("marshalling property: %s", err)
		}
		dataType, err := getPropertyDataType(propId)
		if err != nil {
			return triples, fmt.Errorf("marshalling property: %s", err)
		}
		switch propType {
		case cloudrdf.RdfsLiteral, cloudrdf.RdfsClass:
			obj, err := marshalToRdfObject(value, propType, dataType)
			if err != nil {
				return triples, fmt.Errorf("marshalling property '%s': %s", key, err)
			}
			triple, err := triple.New(rdf.MustBuildNode(res.id), rdf.MustBuildPredicate(propId), obj)
			if err != nil {
				return triples, err
			}
			triples = append(triples, triple)
		case cloudrdf.RdfsList:
			switch dataType {
			case cloudrdf.XsdString:
				list, ok := value.([]string)
				if !ok {
					return triples, fmt.Errorf("marshalling property '%s': expected a string slice, got a %T", key, value)
				}
				for _, l := range list {
					triples = append(triples, rdf.Subject(res.id).Predicate(propId).Literal(l))
				}
			case cloudrdf.RdfsClass:
				list, ok := value.([]string)
				if !ok {
					return triples, fmt.Errorf("marshalling property '%s': expected a string slice, got a %T", key, value)
				}
				for _, l := range list {
					triples = append(triples, rdf.Subject(res.id).Predicate(propId).Object(l))
				}
			case cloudrdf.NetFirewallRule:
				list, ok := value.([]*FirewallRule)
				if !ok {
					return triples, fmt.Errorf("marshalling property '%s': expected a firewall rule slice, got a %T", key, value)
				}
				for _, r := range list {
					ruleId := randomRdfId()
					triples = append(triples, rdf.Subject(res.id).Predicate(propId).Object(ruleId))
					triples = append(triples, r.marshalToTriples(ruleId)...)
				}
			case cloudrdf.NetRoute:
				list, ok := value.([]*Route)
				if !ok {
					return triples, fmt.Errorf("marshalling property '%s': expected a route slice, got a %T", key, value)
				}
				for _, r := range list {
					routeId := randomRdfId()
					triples = append(triples, rdf.Subject(res.id).Predicate(propId).Object(routeId))
					triples = append(triples, r.marshalToTriples(routeId)...)
				}
			case cloudrdf.Grant:
				list, ok := value.([]*Grant)
				if !ok {
					return triples, fmt.Errorf("marshalling property '%s': expected a grant slice, got a %T", key, value)
				}
				for _, g := range list {
					grantId := randomRdfId()
					triples = append(triples, rdf.Subject(res.id).Predicate(propId).Object(grantId))
					triples = append(triples, g.marshalToTriples(grantId)...)
				}
			default:
				return triples, fmt.Errorf("marshalling property '%s': unexpected rdfs:DataType: %s", key, dataType)
			}

		default:
			return triples, fmt.Errorf("marshalling property '%s': unexpected rdfs:isDefinedBy: %s", key, propType)
		}

	}
	return triples, nil
}

func marshalToRdfObject(i interface{}, definedBy, dataType string) (*triple.Object, error) {
	switch definedBy {
	case cloudrdf.RdfsLiteral:
		switch dataType {
		case cloudrdf.XsdDateTime:
			datetime, ok := i.(time.Time)
			if !ok {
				return nil, fmt.Errorf("expected a time, got a %T", i)
			}
			txt, _ := datetime.MarshalText()
			return triple.NewLiteralObject(rdf.MustBuildLiteral(string(txt))), nil
		default:
			return triple.NewLiteralObject(rdf.MustBuildLiteral((fmt.Sprint(i)))), nil
		}
	case cloudrdf.RdfsClass:
		return triple.NewNodeObject(rdf.MustBuildNode(fmt.Sprint(i))), nil
	default:
		return nil, fmt.Errorf("unexpected rdfs:isDefinedBy: %s", definedBy)
	}
}

func (res *Resource) unmarshalFullRdf(gph *rdf.Graph) error {
	triples, err := gph.TriplesForSubjectOnly(rdf.MustBuildNode(res.Id()))
	if err != nil {
		return err
	}
	rTobj, err := marshalResourceType(res.Type())
	if err != nil {
		return err
	}
	if !gph.HasTriple(rdf.Subject(res.Id()).Predicate(cloudrdf.RdfType).ObjectNode(rTobj)) {
		return fmt.Errorf("resource with id '%s' has not type '%s'", res.Id(), res.Type())
	}
	for _, t := range triples {
		pred := string(t.Predicate().ID())

		if !isRDFProperty(pred) || isRDFSubProperty(pred) {
			continue
		}
		propKey, err := getPropertyLabel(pred)
		if err != nil {
			return fmt.Errorf("unmarshalling property: label: %s", err)
		}
		propVal, err := getPropertyValue(gph, t.Object(), pred)
		if err != nil {
			return fmt.Errorf("unmarshalling property %s: val: %s", propKey, err)
		}
		if isRDFList(pred) {
			dataType, err := getPropertyDataType(pred)
			if err != nil {
				return fmt.Errorf("unmarshalling property: datatype: %s", err)
			}
			switch dataType {
			case cloudrdf.RdfsClass, cloudrdf.XsdString:
				list, ok := res.Properties[propKey].([]string)
				if !ok {
					list = []string{}
				}
				list = append(list, propVal.(string))
				res.Properties[propKey] = list
			case cloudrdf.NetFirewallRule:
				list, ok := res.Properties[propKey].([]*FirewallRule)
				if !ok {
					list = []*FirewallRule{}
				}
				list = append(list, propVal.(*FirewallRule))
				res.Properties[propKey] = list
			case cloudrdf.NetRoute:
				list, ok := res.Properties[propKey].([]*Route)
				if !ok {
					list = []*Route{}
				}
				list = append(list, propVal.(*Route))
				res.Properties[propKey] = list
			case cloudrdf.Grant:
				list, ok := res.Properties[propKey].([]*Grant)
				if !ok {
					list = []*Grant{}
				}
				list = append(list, propVal.(*Grant))
				res.Properties[propKey] = list
			default:
				return fmt.Errorf("unmarshalling property: unexpected datatype %s", dataType)
			}
		} else {
			res.Properties[propKey] = propVal
		}
	}
	return nil
}

func (r *Resource) unmarshalMeta(gph *rdf.Graph) error {
	triples, err := gph.TriplesForSubjectPredicate(rdf.MustBuildNode(r.Id()), rdf.MetaPredicate)
	if err != nil {
		return err
	}
	for _, t := range triples {
		lit, err := t.Object().Literal()
		if err != nil {
			return err
		}
		text, err := lit.Text()
		if err != nil {
			return err
		}
		r.Meta["diff"] = text
	}
	return nil
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

func marshalResourceType(rT string) (*triple.Object, error) {
	n, err := node.NewNodeFromStrings("/node", fmt.Sprintf("%s:%s", cloudrdf.CloudOwlNS, strings.Title(rT)))
	if err != nil {
		return nil, err
	}
	return triple.NewNodeObject(n), nil
}

func resolveResourceType(g *rdf.Graph, id string) (string, error) {
	typeTs, err := g.TriplesForSubjectPredicate(rdf.MustBuildNode(id), rdf.MustBuildPredicate(cloudrdf.RdfType))
	if err != nil {
		return "", err
	}
	if len(typeTs) != 1 {
		return "", fmt.Errorf("cannot resolve unique type for resource '%s', got: %v", id, typeTs)
	}
	return unmarshalResourceType(typeTs[0].Object())
}

func unmarshalResourceType(obj *triple.Object) (string, error) {
	node, err := obj.Node()
	if err != nil {
		return "", err
	}
	return strings.ToLower(rdf.TrimNS(node.ID().String())), nil
}

type Property struct {
	Key   string
	Value interface{}
}

type ResourceById []*Resource

func (r ResourceById) Len() int           { return len(r) }
func (r ResourceById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResourceById) Less(i, j int) bool { return r[i].Id() < r[j].Id() }

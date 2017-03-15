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

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
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

func (g *Graph) AddGraph(gph *Graph) {
	g.rdfG.AddGraph(gph.rdfG)
}

func (g *Graph) AddParentRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.ParentOfPredicate)
}

func (g *Graph) AddAppliesOnRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.AppliesOnPredicate)
}

func (g *Graph) GetResource(t string, id string) (*Resource, error) {
	resource := InitResource(id, t)

	node, err := resource.toRDFNode()
	if err != nil {
		return resource, err
	}

	propsTriples, err := g.rdfG.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return resource, err
	}
	if er := resource.Properties.unmarshalRDF(propsTriples); er != nil {
		return resource, er
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

func (g *Graph) FindResource(id string) (*Resource, error) {
	byId := &ById{id}
	resources, err := byId.Resolve(g)
	if err != nil {
		return nil, err
	}
	if len(resources) == 1 {
		return resources[0], nil
	} else if len(resources) > 1 {
		return nil, fmt.Errorf("multiple resources with id '%s' found", id)
	}

	return nil, nil
}

func (g *Graph) FindResourcesByProperty(key string, value interface{}) ([]*Resource, error) {
	byProperty := ByProperty{key, value}
	return byProperty.Resolve(g)
}

func (g *Graph) GetAllResources(t string) ([]*Resource, error) {
	byType := &ByType{t}
	return byType.Resolve(g)
}

func (g *Graph) ResolveResources(resolvers ...Resolver) ([]*Resource, error) {
	var resources []*Resource
	for _, resolv := range resolvers {
		rs, err := resolv.Resolve(g)
		if err != nil {
			return resources, err
		}
		resources = append(resources, rs...)
	}

	return resources, nil
}

type Resolver interface {
	Resolve(g *Graph) ([]*Resource, error)
}

type ById struct {
	Id string
}

func (r *ById) Resolve(g *Graph) ([]*Resource, error) {
	var resources []*Resource

	triples, err := g.rdfG.TriplesForGivenPredicate(rdf.HasTypePredicate)
	if err != nil {
		return resources, err
	}

	for _, triple := range triples {
		sub := triple.Subject()
		if sub.ID().String() == r.Id {
			res, err := g.GetResource(newResourceType(sub), sub.ID().String())
			if err != nil {
				return resources, nil
			}
			resources = append(resources, res)
		}
	}

	return resources, nil
}

type ByProperty struct {
	Name string
	Val  interface{}
}

func (r *ByProperty) Resolve(g *Graph) ([]*Resource, error) {
	var resources []*Resource

	prop := Property{Key: r.Name, Value: r.Val}
	propL, err := prop.marshalRDF()
	if err != nil {
		return resources, err
	}
	triples, err := g.rdfG.TriplesForPredicateObject(rdf.PropertyPredicate, propL)
	if err != nil {
		return resources, err
	}
	for _, t := range triples {
		s := t.Subject()
		r, err := g.GetResource(newResourceType(s), s.ID().String())
		if err != nil {
			return resources, err
		}
		resources = append(resources, r)
	}
	return resources, nil
}

type ByType struct {
	typ string
}

func (r *ByType) Resolve(g *Graph) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.rdfG.NodesForType("/" + r.typ)
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r, err := g.GetResource(r.typ, node.ID().String())
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (g *Graph) ListResourcesAppliedOn(start *Resource) ([]*Resource, error) {
	var resources []*Resource

	node, err := start.toRDFNode()
	if err != nil {
		return resources, err
	}

	relations, err := g.rdfG.ListAttachedTo(node, rdf.AppliesOnPredicate)
	if err != nil {
		return resources, err
	}
	for _, node := range relations {
		res, err := g.GetResource(newResourceType(node), node.ID().String())
		if err != nil {
			return resources, err
		}
		resources = append(resources, res)
	}

	return resources, nil
}

func (g *Graph) ListResourcesDependingOn(start *Resource) ([]*Resource, error) {
	var resources []*Resource

	node, err := start.toRDFNode()
	if err != nil {
		return resources, err
	}

	relations, err := g.rdfG.ListAttachedFrom(node, rdf.AppliesOnPredicate)
	if err != nil {
		return resources, err
	}
	for _, node := range relations {
		res, err := g.GetResource(newResourceType(node), node.ID().String())
		if err != nil {
			return resources, err
		}
		resources = append(resources, res)
	}

	return resources, nil
}

func (g *Graph) Accept(v Visitor) error {
	return v.Visit(g)
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

func (g *Graph) addRelation(one, other *Resource, pred *predicate.Predicate) error {
	n, err := other.toRDFNode()
	if err != nil {
		return err
	}

	oneN, err := node.NewNodeFromStrings("/"+one.Type(), one.Id())
	if err != nil {
		return err
	}

	t, err := triple.New(oneN, pred, triple.NewNodeObject(n))
	if err != nil {
		return err
	}

	g.rdfG.Add(t)

	return nil
}

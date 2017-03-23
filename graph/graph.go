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
		triples, err := res.marshalFullRDF()
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
	resource := InitResource(t, id)

	if err := resource.unmarshalFullRdf(g.rdfG); err != nil {
		return resource, err
	}

	if err := resource.unmarshalMeta(g.rdfG); err != nil {
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

func (g *Graph) GetAllResources(typs ...string) ([]*Resource, error) {
	byTypes := &ByTypes{typs}
	return byTypes.Resolve(g)
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
		id := node.ID().String()
		rT, err := resolveResourceType(g.rdfG, id)
		if err != nil {
			return resources, err
		}
		res, err := g.GetResource(rT, id)
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

	oneN, err := node.NewNodeFromStrings("/node", one.Id())
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

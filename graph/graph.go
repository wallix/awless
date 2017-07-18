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
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Graph struct {
	store tstore.Source
}

func NewGraph() *Graph {
	return &Graph{tstore.NewSource()}
}

func NewGraphFromFile(filepath string) (*Graph, error) {
	g := NewGraph()
	f, err := os.Open(filepath)
	if err != nil {
		return g, err
	}
	ts, err := tstore.NewBinaryDecoder(f).Decode()
	if err != nil {
		return g, err
	}
	g.store.Add(ts...)
	return g, nil
}

func (g *Graph) AddResource(resources ...*Resource) error {
	for _, res := range resources {
		triples, err := res.marshalFullRDF()
		if err != nil {
			return err
		}

		for relType, attachedRes := range res.Relations {
			switch relType {
			case rdf.ChildrenOfRel:
				for _, attached := range attachedRes {
					if err := g.AddParentRelation(attached, res); err != nil {
						return err
					}
				}
			case rdf.DependingOnRel:
				for _, attached := range attachedRes {
					if err := g.AddAppliesOnRelation(attached, res); err != nil {
						return err
					}
				}
			}
		}

		g.store.Add(triples...)
	}
	return nil
}

func (g *Graph) AddGraph(other *Graph) {
	g.store.Add(other.store.CopyTriples()...)
}

func (g *Graph) AddParentRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.ParentOf)
}

func (g *Graph) AddAppliesOnRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.ApplyOn)
}

func (g *Graph) GetResource(t string, id string) (*Resource, error) {
	resource := InitResource(t, id)
	snap := g.store.Snapshot()
	if err := resource.unmarshalFullRdf(snap); err != nil {
		return resource, err
	}

	if err := resource.unmarshalMeta(snap); err != nil {
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

	snap := g.store.Snapshot()
	for _, tri := range snap.WithPredObj(rdf.ApplyOn, tstore.Resource(start.Id())) {
		id := tri.Subject()
		rT, err := resolveResourceType(snap, id)
		if err != nil {
			if err == errTypeNotFound {
				resources = append(resources, NotFoundResource(id))
				continue
			}
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

func (g *Graph) ListResourcesAppliedOn(start *Resource) ([]*Resource, error) {
	var resources []*Resource

	snap := g.store.Snapshot()

	for _, tri := range snap.WithSubjPred(start.Id(), rdf.ApplyOn) {
		id, ok := tri.Object().Resource()
		if !ok {
			return resources, fmt.Errorf("triple %s %s: object is not a resource identifier", start.Id(), rdf.ApplyOn)
		}
		rT, err := resolveResourceType(snap, id)
		if err != nil {
			if err == errTypeNotFound {
				resources = append(resources, NotFoundResource(id))
				continue
			}
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
	ts, err := tstore.NewBinaryDecoder(bytes.NewReader(data)).Decode()
	if err != nil {
		return err
	}
	g.store.Add(ts...)
	return nil
}

func (g *Graph) UnmarshalMultiple(readers ...io.Reader) error {
	dec := tstore.NewDatasetDecoder(tstore.NewBinaryDecoder, readers...)
	ts, err := dec.Decode()
	if err != nil {
		return err
	}
	g.store.Add(ts...)
	return nil
}

func (g *Graph) MustMarshal() string {
	var buff bytes.Buffer
	if err := tstore.NewBinaryEncoder(&buff).Encode(g.store.CopyTriples()...); err != nil {
		panic(err)
	}
	return string(buff.Bytes())
}

func (g *Graph) MarshalTo(w io.Writer) error {
	return tstore.NewBinaryEncoder(w).Encode(g.store.CopyTriples()...)
}

func (g *Graph) addRelation(one, other *Resource, pred string) error {
	g.store.Add(tstore.SubjPred(one.Id(), pred).Resource(other.Id()))
	return nil
}

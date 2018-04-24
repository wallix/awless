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

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Graph struct {
	store tstore.Source
}

func NewGraph() *Graph {
	return &Graph{tstore.NewSource()}
}

func NewGraphFromFiles(files ...string) (cloud.GraphAPI, error) {
	g := NewGraph()

	var readers []io.Reader
	for _, f := range files {
		if reader, err := os.Open(f); err != nil {
			return g, err
		} else {
			readers = append(readers, reader)
		}
	}

	err := g.UnmarshalFromReaders(readers...)
	return g, err
}

func (g *Graph) AsRDFGraphSnaphot() tstore.RDFGraph {
	return g.store.Snapshot()
}

func (g *Graph) AddResource(resources ...*Resource) error {
	for _, res := range resources {
		triples, err := res.marshalFullRDF()
		if err != nil {
			return err
		}

		for relType, attachedRes := range res.relations {
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
	resources, err := byId.Resolve(g.store.Snapshot())
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
	return byProperty.Resolve(g.store.Snapshot())
}

func (g *Graph) FindAncestor(res *Resource, resourceType string) *Resource {
	var found *Resource
	find := func(res *Resource, depth int) error {
		if res.Type() == resourceType {
			found = res
			return nil
		}
		return nil
	}
	g.Accept(&ParentsVisitor{From: res, Each: find})
	return found
}

func (g *Graph) GetAllResources(typs ...string) ([]*Resource, error) {
	byTypes := &ByTypes{typs}
	return byTypes.Resolve(g.store.Snapshot())
}

func (g *Graph) ResolveResources(resolvers ...Resolver) ([]*Resource, error) {
	var resources []*Resource
	snap := g.store.Snapshot()
	for _, resolv := range resolvers {
		rs, err := resolv.Resolve(snap)
		if err != nil {
			return resources, err
		}
		resources = append(resources, rs...)
	}

	return resources, nil
}

func (g *Graph) FilterGraph(q cloud.Query) (cloud.GraphAPI, error) {
	if len(q.ResourceType) != 1 {
		return nil, fmt.Errorf("invalid query: must have exactly one resource type, got %d", len(q.ResourceType))
	}
	resourceType := q.ResourceType[0]
	return g.Filter(resourceType, func(r *Resource) bool {
		if q.Matcher == nil {
			return true
		}
		return q.Matcher.Match(r)
	})
}

func (g *Graph) Find(q cloud.Query) ([]cloud.Resource, error) {
	var resources []*Resource
	var err error
	switch len(q.ResourceType) {
	case 0:
		return nil, fmt.Errorf("invalid query: need at least one resource type")
	case 1:
		var filtered cloud.GraphAPI
		filtered, err = g.FilterGraph(q)
		if err != nil {
			return nil, err
		}
		resources, err = filtered.(*Graph).GetAllResources(q.ResourceType[0])
		if err != nil {
			return nil, err
		}
	default:
		if q.Matcher != nil {
			return nil, fmt.Errorf("invalid query: can not filter whith multiple resource types")
		}
		resources, err = g.GetAllResources(q.ResourceType...)
		if err != nil {
			return nil, err
		}
	}
	var res []cloud.Resource
	for _, r := range resources {
		res = append(res, r)
	}
	return res, nil
}

func (g *Graph) FindWithProperties(props map[string]interface{}) ([]cloud.Resource, error) {
	var resolvers []Resolver
	for k, v := range props {
		resolvers = append(resolvers, &ByProperty{Key: k, Value: v})
	}
	resources, err := g.ResolveResources(&And{Resolvers: resolvers})
	if err != nil {
		return nil, err
	}

	var res []cloud.Resource
	for _, r := range resources {
		res = append(res, r)
	}
	return res, nil
}

func (g *Graph) FindOne(q cloud.Query) (cloud.Resource, error) {
	filtered, err := g.FilterGraph(q)
	if err != nil {
		return nil, err
	}
	resources, err := filtered.(*Graph).GetAllResources(q.ResourceType[0])
	if err != nil {
		return nil, err
	}
	switch len(resources) {
	case 0:
		return nil, fmt.Errorf("resource not found")
	case 1:
		return resources[0], nil
	default:
		return nil, fmt.Errorf("multiple resources found")
	}
}

func (g *Graph) Merge(mg cloud.GraphAPI) error {
	toMerge, ok := mg.(*Graph)
	if !ok {
		return fmt.Errorf("can not merge graphs, graph to merge is not a *graph.Graph, but a %T", mg)
	}
	g.AddGraph(toMerge)
	return nil
}

func (g *Graph) ResourceRelations(from cloud.Resource, relation string, recursive bool) (collect []cloud.Resource, err error) {
	collectFunc := func(r *Resource, depth int) error {
		if depth == 1 || recursive {
			collect = append(collect, r)
		}
		return nil
	}
	switch relation {
	case rdf.ChildrenOfRel:
		err = g.Accept(&ChildrenVisitor{From: from.(*Resource), IncludeFrom: false, Relation: rdf.ParentOf, Each: collectFunc})
	case rdf.DependingOnRel:
		err = g.Accept(&ParentsVisitor{From: from.(*Resource), IncludeFrom: false, Relation: rdf.ApplyOn, Each: collectFunc})
	case rdf.ApplyOn:
		err = g.Accept(&ChildrenVisitor{From: from.(*Resource), IncludeFrom: false, Relation: rdf.ApplyOn, Each: collectFunc})
	default:
		err = g.Accept(&ParentsVisitor{From: from.(*Resource), IncludeFrom: false, Relation: relation, Each: collectFunc})
	}
	return
}

func (g *Graph) VisitRelations(from cloud.Resource, relation string, includeFrom bool, each func(cloud.Resource, int) error) (err error) {
	eachFunc := func(r *Resource, depth int) error {
		return each(r, depth)
	}
	switch relation {
	case rdf.ChildrenOfRel:
		err = g.Accept(&ChildrenVisitor{From: from.(*Resource), IncludeFrom: includeFrom, Relation: rdf.ParentOf, Each: eachFunc})
	case rdf.DependingOnRel:
		err = g.Accept(&ParentsVisitor{From: from.(*Resource), IncludeFrom: includeFrom, Relation: rdf.ApplyOn, Each: eachFunc})
	case rdf.ApplyOn:
		err = g.Accept(&ChildrenVisitor{From: from.(*Resource), IncludeFrom: includeFrom, Relation: rdf.ApplyOn, Each: eachFunc})
	default:
		err = g.Accept(&ParentsVisitor{From: from.(*Resource), IncludeFrom: includeFrom, Relation: relation, Each: eachFunc})
	}
	return
}

func (g *Graph) ResourceSiblings(res cloud.Resource) (collect []cloud.Resource, err error) {
	err = g.Accept(&SiblingsVisitor{From: res.(*Resource), IncludeFrom: false, Each: func(r *Resource, depth int) error {
		collect = append(collect, r)
		return nil
	}})
	return collect, err
}

func ResolveResourcesWithProp(snap tstore.RDFGraph, resType, propName, propVal string) ([]*Resource, error) {
	resolv := ByTypeAndProperty{
		Type:  resType,
		Key:   propName,
		Value: propVal,
	}
	return resolv.Resolve(snap)
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
	ts, err := tstore.NewAutoDecoder(bytes.NewReader(data)).Decode()
	if err != nil {
		return err
	}
	g.store.Add(ts...)
	return nil
}

func (g *Graph) UnmarshalFromReaders(readers ...io.Reader) error {
	dec := tstore.NewDatasetDecoder(tstore.NewLenientNTDecoder, readers...)
	ts, err := dec.Decode()
	if err != nil {
		return err
	}
	g.store.Add(ts...)
	return nil
}

func (g *Graph) MustMarshal() string {
	var buff bytes.Buffer
	if err := tstore.NewLenientNTEncoder(&buff).Encode(g.store.CopyTriples()...); err != nil {
		panic(err)
	}
	return string(buff.Bytes())
}

func (g *Graph) MarshalTo(w io.Writer) error {
	return tstore.NewLenientNTEncoder(w).Encode(g.store.CopyTriples()...)
}

func (g *Graph) addRelation(one, other *Resource, pred string) error {
	g.store.Add(tstore.SubjPred(one.Id(), pred).Resource(other.Id()))
	return nil
}

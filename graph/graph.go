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

func (g *Graph) AddGraph(graph *Graph) {
	g.rdfG.AddGraph(graph.rdfG)
}

func (g *Graph) AddParentRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.ParentOfPredicate)
}

func (g *Graph) AddAppliesOnRelation(parent, child *Resource) error {
	return g.addRelation(parent, child, rdf.AppliesOnPredicate)
}

func (g *Graph) GetResource(t ResourceType, id string) (*Resource, error) {
	resource := InitResource(id, t)

	node, err := resource.toRDFNode()
	if err != nil {
		return resource, err
	}

	propsTriples, err := g.rdfG.TriplesForSubjectPredicate(node, rdf.PropertyPredicate)
	if err != nil {
		return resource, err
	}
	if err := resource.Properties.unmarshalRDF(propsTriples); err != nil {
		return resource, err
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
	triples, err := g.rdfG.TriplesForGivenPredicate(rdf.HasTypePredicate)
	if err != nil {
		return nil, err
	}

	for _, triple := range triples {
		sub := triple.Subject()
		if sub.ID().String() == id {
			return g.GetResource(newResourceType(sub), sub.ID().String())
		}
	}

	return nil, nil
}

func (g *Graph) FindResourcesByProperty(key string, value interface{}) ([]*Resource, error) {
	var res []*Resource
	prop := Property{Key: key, Value: value}
	propL, err := prop.marshalRDF()
	if err != nil {
		return res, err
	}
	triples, err := g.rdfG.TriplesForPredicateObject(rdf.PropertyPredicate, propL)
	if err != nil {
		return res, err
	}
	for _, t := range triples {
		s := t.Subject()
		r, err := g.GetResource(newResourceType(s), s.ID().String())
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (g *Graph) GetAllResources(t ResourceType) ([]*Resource, error) {
	var res []*Resource
	nodes, err := g.rdfG.NodesForType(t.ToRDFString())
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r, err := g.GetResource(t, node.ID().String())
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
	for _, node := range relations {
		res, err := g.GetResource(newResourceType(node), node.ID().String())
		if err != nil {
			return resources, err
		}
		resources = append(resources, res)
	}

	return resources, nil
}

func (g *Graph) VisitChildren(start *Resource, each func(*Resource, int)) error {
	startNode, err := start.toRDFNode()
	if err != nil {
		return err
	}

	foreach := func(rdfG *rdf.Graph, n *node.Node, i int) {
		res, err := g.GetResource(newResourceType(n), n.ID().String())
		if err != nil {
			panic(err)
		}
		each(res, i)
	}

	return g.rdfG.VisitTopDown(startNode, foreach)
}

func (g *Graph) VisitParents(start *Resource, each func(*Resource, int)) error {
	startNode, err := start.toRDFNode()
	if err != nil {
		return err
	}

	foreach := func(rdfG *rdf.Graph, n *node.Node, i int) {
		res, err := g.GetResource(newResourceType(n), n.ID().String())
		if err != nil {
			panic(err)
		}
		each(res, i)
	}

	return g.rdfG.VisitBottomUp(startNode, foreach)
}

func (g *Graph) VisitSiblings(res *Resource, each func(*Resource, int)) error {
	resNode, err := res.toRDFNode()
	if err != nil {
		return err
	}

	foreach := func(rdfG *rdf.Graph, n *node.Node, i int) {
		res, err := g.GetResource(newResourceType(n), n.ID().String())
		if err != nil {
			panic(err)
		}
		each(res, i)
	}

	return g.rdfG.VisitSiblings(resNode, foreach)
}

func (g *Graph) CountChildrenOfTypeForNode(res *Resource, childType ResourceType) (int, error) {
	n, err := node.NewNodeFromStrings(res.Type().ToRDFString(), res.Id())
	if err != nil {
		return 0, err
	}
	return g.rdfG.CountTriplesForSubjectAndPredicateObjectOfType(n, rdf.ParentOfPredicate, childType.ToRDFString())
}

func (g *Graph) CountChildrenForNode(res *Resource) (int, error) {
	n, err := node.NewNodeFromStrings(res.Type().ToRDFString(), res.Id())
	if err != nil {
		return 0, err
	}
	return g.rdfG.CountTriplesForSubjectAndPredicate(n, rdf.ParentOfPredicate)
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

	oneN, err := node.NewNodeFromStrings(one.Type().ToRDFString(), one.Id())
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

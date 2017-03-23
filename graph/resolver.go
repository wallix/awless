package graph

import (
	"fmt"

	"github.com/wallix/awless/cloud/properties"
	cloudrdf "github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Resolver interface {
	Resolve(g *Graph) ([]*Resource, error)
}

type ById struct {
	Id string
}

func (r *ById) Resolve(g *Graph) ([]*Resource, error) {
	resolver := &ByProperty{Name: properties.ID, Val: r.Id}
	return resolver.Resolve(g)
}

type ByProperty struct {
	Name string
	Val  interface{}
}

func (r *ByProperty) Resolve(g *Graph) ([]*Resource, error) {
	var resources []*Resource
	rdfpropLabel, ok := cloudrdf.Labels[r.Name]
	if !ok {
		return resources, fmt.Errorf("resolve resources by property: undefined property label '%s'", r.Name)
	}
	rdfProp, ok := cloudrdf.RdfProperties[rdfpropLabel]
	if !ok {
		return resources, fmt.Errorf("resolve resources by property: undefined property definition '%s'", rdfpropLabel)
	}
	obj, err := marshalToRdfObject(r.Val, rdfProp.RdfsDefinedBy, rdfProp.RdfsDataType)
	if err != nil {
		return resources, fmt.Errorf("resolve resources: unmarshaling property '%s': '%s'", r.Name, err)
	}
	triples, err := g.rdfG.TriplesForPredicateObject(rdf.MustBuildPredicate(rdfpropLabel), obj)
	if err != nil {
		return resources, err
	}
	for _, t := range triples {
		s := t.Subject()
		rt, err := resolveResourceType(g.rdfG, s.ID().String())
		if err != nil {
			return resources, err
		}
		r := InitResource(rt, s.ID().String())

		if err := r.unmarshalFullRdf(g.rdfG); err != nil {
			return resources, err
		}
		resources = append(resources, r)
	}
	return resources, nil
}

type And struct {
	Resolvers []Resolver
}

func (r *And) Resolve(g *Graph) (result []*Resource, err error) {
	if len(r.Resolvers) == 0 {
		return
	}
	result, err = r.Resolvers[0].Resolve(g)
	if err != nil {
		return
	}
	gg := NewGraph()
	err = gg.AddResource(result...)
	if err != nil {
		return
	}
	for _, resolv := range r.Resolvers {
		result, err = resolv.Resolve(gg)
		if err != nil {
			return
		}
		gg = NewGraph()
		err = gg.AddResource(result...)
		if err != nil {
			return
		}
	}
	return
}

type ByType struct {
	Typ string
}

func (r *ByType) Resolve(g *Graph) ([]*Resource, error) {
	var resources []*Resource

	typObj, err := marshalResourceType(r.Typ)
	if err != nil {
		return resources, err
	}
	triples, err := g.rdfG.TriplesForPredicateObject(rdf.MustBuildPredicate(cloudrdf.RdfType), typObj)
	if err != nil {
		return resources, err
	}
	for _, t := range triples {
		s := t.Subject()
		r := InitResource(r.Typ, s.ID().String())
		err := r.unmarshalFullRdf(g.rdfG)
		if err != nil {
			return resources, err
		}
		resources = append(resources, r)
	}
	return resources, nil
}

type ByTypes struct {
	Typs []string
}

func (r *ByTypes) Resolve(g *Graph) ([]*Resource, error) {
	var res []*Resource
	for _, t := range r.Typs {
		bt := &ByType{t}
		all, err := bt.Resolve(g)
		if err != nil {
			return res, err
		}
		res = append(res, all...)
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

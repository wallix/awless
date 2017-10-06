package graph

import (
	"fmt"

	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Resolver interface {
	Resolve(snap tstore.RDFGraph) ([]*Resource, error)
}

type ById struct {
	Id string
}

func (r *ById) Resolve(snap tstore.RDFGraph) ([]*Resource, error) {
	resolver := &ByProperty{Key: properties.ID, Value: r.Id}
	return resolver.Resolve(snap)
}

type ByTypeAndProperty struct {
	Type  string
	Key   string
	Value interface{}
}

func (r *ByTypeAndProperty) Resolve(snap tstore.RDFGraph) ([]*Resource, error) {
	var resources []*Resource

	if r.Value == nil {
		return resources, nil
	}
	rdfpropLabel, ok := rdf.Labels[r.Key]
	if !ok {
		return resources, fmt.Errorf("resolve by property: undefined property label '%s'", r.Key)
	}
	rdfProp, err := rdf.Properties.Get(rdfpropLabel)
	if err != nil {
		return resources, fmt.Errorf("resolve by property: %s", err)
	}
	obj, err := marshalToRdfObject(r.Value, rdfProp.RdfsDefinedBy, rdfProp.RdfsDataType)
	if err != nil {
		return resources, fmt.Errorf("resolve by property: unmarshaling property '%s': %s", r.Key, err)
	}
	for _, t := range snap.WithPredObj(rdfpropLabel, obj) {
		rt, err := resolveResourceType(snap, t.Subject())
		if err != nil {
			return resources, err
		}
		if rt == r.Type {
			r := InitResource(rt, t.Subject())

			if err := r.unmarshalFullRdf(snap); err != nil {
				return resources, err
			}
			resources = append(resources, r)
		}
	}
	return resources, nil
}

type ByProperty struct {
	Key   string
	Value interface{}
}

func (r *ByProperty) Resolve(snap tstore.RDFGraph) ([]*Resource, error) {
	var resources []*Resource
	if r.Value == nil {
		return resources, nil
	}
	rdfpropLabel, ok := rdf.Labels[r.Key]
	if !ok {
		return resources, fmt.Errorf("resolve by property: undefined property label '%s'", r.Key)
	}
	rdfProp, err := rdf.Properties.Get(rdfpropLabel)
	if err != nil {
		return resources, fmt.Errorf("resolve by property: %s", err)
	}
	obj, err := marshalToRdfObject(r.Value, rdfProp.RdfsDefinedBy, rdfProp.RdfsDataType)
	if err != nil {
		return resources, fmt.Errorf("resolve by property: unmarshaling property '%s': %s", r.Key, err)
	}
	for _, t := range snap.WithPredObj(rdfpropLabel, obj) {
		rt, err := resolveResourceType(snap, t.Subject())
		if err != nil {
			return resources, err
		}
		r := InitResource(rt, t.Subject())

		if err := r.unmarshalFullRdf(snap); err != nil {
			return resources, err
		}
		resources = append(resources, r)
	}
	return resources, nil
}

type And struct {
	Resolvers []Resolver
}

func (r *And) Resolve(snap tstore.RDFGraph) (result []*Resource, err error) {
	if len(r.Resolvers) == 0 {
		return
	}
	result, err = r.Resolvers[0].Resolve(snap)
	if err != nil {
		return
	}
	gg := NewGraph()
	if err = gg.AddResource(result...); err != nil {
		return
	}
	for _, resolv := range r.Resolvers {
		result, err = resolv.Resolve(gg.store.Snapshot())
		if err != nil {
			return
		}
		gg = NewGraph()
		if err = gg.AddResource(result...); err != nil {
			return
		}
	}
	return
}

type Or struct {
	Resolvers []Resolver
}

func (r *Or) Resolve(snap tstore.RDFGraph) (result []*Resource, err error) {
	for _, resolv := range r.Resolvers {
		result, err = resolv.Resolve(snap)
		if err != nil {
			return
		}
		if len(result) > 0 {
			return
		}
	}
	return
}

type ByType struct {
	Typ string
}

func (r *ByType) Resolve(snap tstore.RDFGraph) ([]*Resource, error) {
	var resources []*Resource
	typ := namespacedResourceType(r.Typ)
	for _, t := range snap.WithPredObj(rdf.RdfType, tstore.Resource(typ)) {
		r := InitResource(r.Typ, t.Subject())
		err := r.unmarshalFullRdf(snap)
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

func (r *ByTypes) Resolve(snap tstore.RDFGraph) ([]*Resource, error) {
	var res []*Resource
	for _, t := range r.Typs {
		bt := &ByType{t}
		all, err := bt.Resolve(snap)
		if err != nil {
			return res, err
		}
		res = append(res, all...)
	}

	return res, nil
}

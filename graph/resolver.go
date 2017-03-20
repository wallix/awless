package graph

import "github.com/wallix/awless/graph/internal/rdf"

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
	var res []*Resource
	nodes, err := g.rdfG.NodesForType("/" + r.Typ)
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		r, err := g.GetResource(r.Typ, node.ID().String())
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
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
		res, err := g.GetResource(newResourceType(node), node.ID().String())
		if err != nil {
			return resources, err
		}
		resources = append(resources, res)
	}

	return resources, nil
}

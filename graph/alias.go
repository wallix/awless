package graph

import (
	"github.com/wallix/awless/graph/internal/rdf"
)

type Alias string

func (a Alias) ResolveToId(g *Graph, resT ResourceType) (string, bool) {
	prop := Property{Key: "Name", Value: a}
	propL, err := prop.ToLiteralObject()
	if err != nil {
		return "", false
	}
	triples, err := g.TriplesForPredicateObject(rdf.PropertyPredicate, propL)
	if err != nil {
		return "", false
	}
	for _, t := range triples {
		s := t.Subject()
		if s.Type().String() == resT.ToRDFString() {
			return s.ID().String(), true
		}
	}

	return "", false
}

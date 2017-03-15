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
	"github.com/wallix/awless/graph/internal/rdf"
)

type Alias string

func (a Alias) ResolveToId(g *Graph, resT string) (string, bool) {
	prop := Property{Key: "Name", Value: a}
	propL, err := prop.marshalRDF()
	if err != nil {
		return "", false
	}
	triples, err := g.rdfG.TriplesForPredicateObject(rdf.PropertyPredicate, propL)
	if err != nil {
		return "", false
	}
	for _, t := range triples {
		s := t.Subject()
		if s.Type().String() == "/"+resT {
			return s.ID().String(), true
		}
	}

	return "", false
}

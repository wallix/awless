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

package rdf

import (
	"context"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

type QueryType int

const (
	SUBJECT_PREDICATE QueryType = iota
	PREDICATE_OBJECT
	PREDICATE_ONLY
	SUBJECT_ONLY
)

func (g *Graph) TriplesForSubjectPredicate(subject *node.Node, predicate *predicate.Predicate) ([]*triple.Triple, error) {
	return g.returnTriples(SUBJECT_PREDICATE, predicate, subject)
}

func (g *Graph) TriplesForPredicateObject(predicate *predicate.Predicate, object *triple.Object) ([]*triple.Triple, error) {
	return g.returnTriples(PREDICATE_OBJECT, predicate, object)
}

func (g *Graph) TriplesForSubjectOnly(subject *node.Node) ([]*triple.Triple, error) {
	return g.returnTriples(SUBJECT_ONLY, subject)
}

func (g *Graph) CountTriplesForSubjectAndPredicate(subject *node.Node, predicate *predicate.Predicate) (int, error) {
	all, err := g.returnTriples(SUBJECT_PREDICATE, predicate, subject)
	return len(all), err
}

func (g *Graph) TriplesForType(t string) ([]*triple.Triple, error) {
	literal, err := literal.DefaultBuilder().Build(literal.Text, t)
	if err != nil {
		return []*triple.Triple{}, err
	}

	return g.returnTriples(PREDICATE_OBJECT, HasTypePredicate, triple.NewLiteralObject(literal))
}

func (g *Graph) TriplesForGivenPredicate(pred *predicate.Predicate) ([]*triple.Triple, error) {
	return g.returnTriples(PREDICATE_ONLY, pred)
}

func (g *Graph) CountTriplesForSubjectAndPredicateObjectOfType(subject *node.Node, predicate *predicate.Predicate, objectType string) (int, error) {
	all, err := g.returnTriples(SUBJECT_PREDICATE, predicate, subject)
	if err != nil {
		return 0, err
	}

	var count int

	for _, t := range all {
		n, e := t.Object().Node()
		if e != nil {
			return 0, e
		}
		triples, e := g.TriplesForSubjectPredicate(n, HasTypePredicate)
		if e != nil {
			return 0, e
		}
		if len(triples) == 1 {
			hasTypeTriple := triples[0]
			childTypeL, e := hasTypeTriple.Object().Literal()
			if e != nil {
				return 0, e
			}
			childType, e := childTypeL.Text()
			if e != nil {
				return 0, e
			} else if childType == objectType {
				count++
			}
		}
	}

	return count, err
}

func (g *Graph) NodesForType(t string) ([]*node.Node, error) {
	var nodes []*node.Node
	errc := make(chan error)
	nodec := make(chan *node.Node)
	literal, err := literal.DefaultBuilder().Build(literal.Text, t)
	if err != nil {
		return nodes, err
	}

	go func() {
		defer close(errc)
		errc <- g.Subjects(context.Background(), HasTypePredicate, triple.NewLiteralObject(literal), storage.DefaultLookup, nodec)
	}()

	for n := range nodec {
		nodes = append(nodes, n)
	}

	return nodes, <-errc

}

func (g *Graph) returnTriples(kind QueryType, objects ...interface{}) ([]*triple.Triple, error) {
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)

		switch kind {
		case SUBJECT_PREDICATE:
			predicate, subject := objects[0].(*predicate.Predicate), objects[1].(*node.Node)
			errc <- g.TriplesForSubjectAndPredicate(context.Background(), subject, predicate, storage.DefaultLookup, triplec)
		case PREDICATE_OBJECT:
			predicate, object := objects[0].(*predicate.Predicate), objects[1].(*triple.Object)
			errc <- g.TriplesForPredicateAndObject(context.Background(), predicate, object, storage.DefaultLookup, triplec)
		case PREDICATE_ONLY:
			predicate := objects[0].(*predicate.Predicate)
			errc <- g.TriplesForPredicate(context.Background(), predicate, storage.DefaultLookup, triplec)
		case SUBJECT_ONLY:
			subject := objects[0].(*node.Node)
			errc <- g.TriplesForSubject(context.Background(), subject, storage.DefaultLookup, triplec)
		}
	}()

	var triples []*triple.Triple

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
}

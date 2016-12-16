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
)

func (g *Graph) TriplesForSubjectPredicate(subject *node.Node, predicate *predicate.Predicate) ([]*triple.Triple, error) {
	return g.returnTriples(SUBJECT_PREDICATE, predicate, subject)
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

func (g *Graph) TriplesForPredicateName(name string) ([]*triple.Triple, error) {
	predicate, err := predicate.NewImmutable(name)
	if err != nil {
		return []*triple.Triple{}, err
	}

	return g.returnTriples(PREDICATE_ONLY, predicate)
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
		}
	}()

	var triples []*triple.Triple

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
}

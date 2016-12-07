package rdf

import (
	"context"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func (g *Graph) TriplesForSubjectPredicate(subject *node.Node, predicate *predicate.Predicate) ([]*triple.Triple, error) {
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- g.TriplesForSubjectAndPredicate(context.Background(), subject, predicate, storage.DefaultLookup, triplec)
	}()

	var triples []*triple.Triple

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
}

func (g *Graph) TriplesForType(t string) ([]*triple.Triple, error) {
	var triples []*triple.Triple
	errc := make(chan error)
	triplec := make(chan *triple.Triple)
	literal, err := literal.DefaultBuilder().Build(literal.Text, t)
	if err != nil {
		return triples, err
	}

	go func() {
		defer close(errc)
		errc <- g.TriplesForPredicateAndObject(context.Background(), HasTypePredicate, triple.NewLiteralObject(literal), storage.DefaultLookup, triplec)
	}()

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
}

func (g *Graph) TriplesForPredicateName(name string) ([]*triple.Triple, error) {
	var triples []*triple.Triple
	errc := make(chan error)
	triplec := make(chan *triple.Triple)
	p, err := predicate.NewImmutable(name)
	if err != nil {
		return triples, err
	}

	go func() {
		defer close(errc)
		errc <- g.TriplesForPredicate(context.Background(), p, storage.DefaultLookup, triplec)
	}()

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
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

func (g *Graph) CountTriplesForSubjectAndPredicate(subject *node.Node, predicate *predicate.Predicate) (int, error) {
	count := 0
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- g.TriplesForSubjectAndPredicate(context.Background(), subject, predicate, storage.DefaultLookup, triplec)
	}()

	for range triplec {
		count++
	}

	return count, <-errc
}

func (g *Graph) CountTriplesForSubjectAndPredicateObjectOfType(subject *node.Node, predicate *predicate.Predicate, objectType string) (int, error) {
	count := 0
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- g.TriplesForSubjectAndPredicate(context.Background(), subject, predicate, storage.DefaultLookup, triplec)
	}()

	for t := range triplec {
		n, err := t.Object().Node()
		if err != nil {
			return 0, err
		}
		triples, err := g.TriplesForSubjectPredicate(n, HasTypePredicate)
		if err != nil {
			return 0, err
		}
		if len(triples) == 1 {
			hasTypeTriple := triples[0]
			childTypeL, err := hasTypeTriple.Object().Literal()
			if err != nil {
				return 0, err
			}
			childType, err := childTypeL.Text()
			if err != nil {
				return 0, err
			} else if childType == objectType {
				count++
			}
		}
	}

	return count, <-errc
}

package rdf

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"time"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/storage/memory"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

var (
	ParentOf      *predicate.Predicate
	HasType       *predicate.Predicate
	DiffPredicate *predicate.Predicate

	ExtraLiteral   *literal.Literal
	MissingLiteral *literal.Literal
)

func init() {
	var err error
	if ParentOf, err = predicate.NewImmutable("parent_of"); err != nil {
		panic(err)
	}
	if HasType, err = predicate.NewImmutable("has_type"); err != nil {
		panic(err)
	}
	if DiffPredicate, err = predicate.NewImmutable("diff"); err != nil {
		panic(err)
	}
	if ExtraLiteral, err = literal.DefaultBuilder().Build(literal.Text, "extra"); err != nil {
		panic(err)
	}
	if MissingLiteral, err = literal.DefaultBuilder().Build(literal.Text, "missing"); err != nil {
		panic(err)
	}
}

type Graph struct {
	storage.Graph
	triplesCount int
}

func NewGraph() *Graph {
	g, err := memory.DefaultStore.NewGraph(context.Background(), randString())
	if err != nil {
		panic(err) // badwolf implementation: only happens on duplicates names of graph
	}
	return &Graph{Graph: g}
}

func NewGraphFromTriples(triples []*triple.Triple) *Graph {
	g := NewGraph()
	g.triplesCount = len(triples)
	g.Add(triples...)
	return g
}

func NewGraphFromFile(filepath string) (*Graph, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	g := NewGraph()

	if err := g.Unmarshal(data); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Graph) Add(triples ...*triple.Triple) {
	g.triplesCount += len(triples)
	_ = g.AddTriples(context.Background(), triples) // badwolf mem store implementation always returns nil error
}

func (g *Graph) VisitDepthFirst(root *node.Node, each func(*Graph, *node.Node, int), distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	each(g, root, dist)

	relations, err := g.TriplesForSubjectPredicate(root, ParentOf)
	if err != nil {
		return err
	}

	var childs []*node.Node
	for _, relation := range relations {
		n, err := relation.Object().Node()
		if err != nil {
			return err
		}
		childs = append(childs, n)
	}

	sort.Sort(&nodeSorter{childs})

	for _, child := range childs {
		g.VisitDepthFirst(child, each, dist+1)
	}

	return nil
}

func (g *Graph) copy() *Graph {
	newg := NewGraph()

	all, _ := g.allTriples()
	newg.Add(all...)

	return newg
}

func (g *Graph) Merge(other *Graph) *Graph {
	all, _ := other.allTriples()
	g.Add(all...)

	return g
}

func (g *Graph) Substract(other *Graph) *Graph {
	sub := g.copy()

	others, _ := other.allTriples()
	sub.RemoveTriples(context.Background(), others)

	return sub
}

func (g *Graph) Intersect(other *Graph) *Graph {
	inter := NewGraph()

	all, err := g.allTriples()
	if err != nil {
		return nil
	}

	for _, tri := range all {
		exists, err := other.Exist(context.Background(), tri)
		if exists && err == nil {
			inter.Add(tri)
		}
	}

	return inter
}

func (g *Graph) size() int {
	return g.triplesCount
}

func (g *Graph) IsEmpty() bool {
	return g.size() == 0
}

func (g *Graph) allTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- g.Triples(context.Background(), triplec)
	}()

	for t := range triplec {
		triples = append(triples, t)
	}

	return triples, <-errc
}

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
		errc <- g.TriplesForPredicateAndObject(context.Background(), HasType, triple.NewLiteralObject(literal), storage.DefaultLookup, triplec)
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
		errc <- g.Subjects(context.Background(), HasType, triple.NewLiteralObject(literal), storage.DefaultLookup, nodec)
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
		triples, err := g.TriplesForSubjectPredicate(n, HasType)
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

func (g *Graph) Unmarshal(data []byte) error {
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if bytes.Equal(bytes.TrimSpace(line), []byte("")) {
			continue
		}
		triple, err := triple.Parse(string(line), literal.DefaultBuilder())
		if err != nil {
			return err
		}
		g.Add(triple)
	}

	return nil
}

func (g *Graph) Marshal() ([]byte, error) {
	triples, err := g.allTriples()
	if err != nil {
		return nil, err
	}

	sort.Sort(&tripleSorter{triples})

	var out [][]byte
	for _, triple := range triples {
		out = append(out, []byte(triple.String()))
	}

	return bytes.Join(out, []byte("\n")), nil
}

func (g *Graph) MustMarshal() string {
	b, err := g.Marshal()
	if err != nil {
		panic(err)
	}
	return string(b)
}

var rando = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func randString() string {
	return fmt.Sprintf("%d", rando.Intn(1e16))
}

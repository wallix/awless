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

var parentOf *predicate.Predicate

func init() {
	var err error
	if parentOf, err = predicate.NewImmutable("parent_of"); err != nil {
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

func (g *Graph) VisitDepthFirst(root *node.Node, each func(*node.Node, int), distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	each(root, dist)

	relations, err := triplesForSubjectAndPredicate(g, root, parentOf)
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

func (g *Graph) TriplesCount() int {
	return g.triplesCount
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

type nodeSorter struct {
	nodes []*node.Node
}

func (s *nodeSorter) Len() int {
	return len(s.nodes)
}
func (s *nodeSorter) Less(i, j int) bool {
	return s.nodes[i].ID().String() < s.nodes[j].ID().String()
}

func (s *nodeSorter) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

type tripleSorter struct {
	triples []*triple.Triple
}

func (s *tripleSorter) Len() int {
	return len(s.triples)
}
func (s *tripleSorter) Less(i, j int) bool {
	return s.triples[i].String() < s.triples[j].String()
}

func (s *tripleSorter) Swap(i, j int) {
	s.triples[i], s.triples[j] = s.triples[j], s.triples[i]
}

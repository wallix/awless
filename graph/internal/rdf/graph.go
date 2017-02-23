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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/storage/memory"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func init() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

type Graph struct {
	storage.Graph
	triplesCount uint32 // atomic
}

func NewGraph() *Graph {
	g, err := memory.DefaultStore.NewGraph(context.Background(), randString())
	if err != nil {
		panic(err) // badwoclf implementation: only happens on duplicates names of graph
	}
	return &Graph{Graph: g}
}

func NewGraphFromTriples(triples []*triple.Triple) *Graph {
	g := NewGraph()
	atomic.StoreUint32(&g.triplesCount, uint32(len(triples)))
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
	atomic.AddUint32(&g.triplesCount, uint32(len(triples)))
	_ = g.AddTriples(context.Background(), triples) // badwolf mem store implementation always returns nil error
}

func (g *Graph) AddGraph(graph *Graph) {
	all, _ := graph.allTriples()
	g.Add(all...)
}

func (g *Graph) ListAttachedTo(n *node.Node, pred *predicate.Predicate) ([]*node.Node, error) {
	return g.listAttached(n, pred, DOWN)
}

func (g *Graph) ListAttachedFrom(n *node.Node, pred *predicate.Predicate) ([]*node.Node, error) {
	return g.listAttached(n, pred, UP)
}

func (g *Graph) VisitTopDown(root *node.Node, each func(*Graph, *node.Node, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(g, root, dist); err != nil {
		return err
	}

	relations, err := g.TriplesForSubjectPredicate(root, ParentOfPredicate)
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
		g.VisitTopDown(child, each, dist+1)
	}

	return nil
}

func (g *Graph) VisitBottomUp(startNode *node.Node, each func(*Graph, *node.Node, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(g, startNode, dist); err != nil {
		return err
	}

	relations, err := g.TriplesForPredicateObject(ParentOfPredicate, triple.NewNodeObject(startNode))
	if err != nil {
		return err
	}

	var parents []*node.Node
	for _, relation := range relations {
		parents = append(parents, relation.Subject())
	}

	sort.Sort(&nodeSorter{parents})

	for _, child := range parents {
		g.VisitBottomUp(child, each, dist+1)
	}

	return nil
}

func (g *Graph) VisitSiblings(start *node.Node, each func(*Graph, *node.Node, int) error, distances ...int) error {
	relations, err := g.TriplesForPredicateObject(ParentOfPredicate, triple.NewNodeObject(start))
	if err != nil {
		return err
	}

	var parents []*node.Node
	for _, relation := range relations {
		parents = append(parents, relation.Subject())
	}

	if len(parents) == 0 {
		return each(g, start, 0)
	}

	sort.Sort(&nodeSorter{parents})

	for _, parent := range parents {
		triples, err := g.TriplesForSubjectPredicate(parent, ParentOfPredicate)
		if err != nil {
			return nil
		}

		var childs []*node.Node
		for _, triple := range triples {
			child, err := triple.Object().Node()
			if err != nil {
				return err
			}
			childs = append(childs, child)
		}

		sort.Sort(&nodeSorter{childs})

		for _, child := range childs {
			sameType := child.Type().String() == start.Type().String()
			if sameType {
				if err := each(g, child, 0); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *Graph) IsEmpty() bool {
	return g.size() == 0
}

func (g *Graph) HasTriple(t *triple.Triple) bool {
	ok, err := g.Exist(context.Background(), t)
	return ok && err == nil
}

func (g *Graph) copy() *Graph {
	newg := NewGraph()

	all, _ := g.allTriples()
	newg.Add(all...)

	return newg
}

func (g *Graph) size() uint32 {
	return atomic.LoadUint32(&g.triplesCount)
}

func (g *Graph) allTriples() ([]*triple.Triple, error) {
	var triples []*triple.Triple
	errc := make(chan error)
	triplec := make(chan *triple.Triple)

	go func() {
		defer close(errc)
		errc <- g.Triples(context.Background(), storage.DefaultLookup, triplec)
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
	toSort := strings.Split(string(b), "\n")
	sort.Strings(toSort)
	return strings.Join(toSort, "\n")
}

type direction int

const (
	UP direction = iota
	DOWN
)

func (g *Graph) listAttached(n *node.Node, pred *predicate.Predicate, dir direction) ([]*node.Node, error) {
	var nodes []*node.Node
	var attached []*triple.Triple
	var err error

	switch dir {
	case DOWN:
		attached, err = g.TriplesForSubjectPredicate(n, pred)
	case UP:
		attached, err = g.TriplesForPredicateObject(pred, triple.NewNodeObject(n))
	default:
		panic("undefined direction")
	}

	if err != nil {
		return nodes, err
	}

	for _, triple := range attached {
		var nod *node.Node
		switch dir {
		case UP:
			nod = triple.Subject()
		case DOWN:
			nod, err = triple.Object().Node()
		default:
			panic("undefined direction")
		}

		if err != nil {
			return nodes, err
		}

		if nod != nil {
			nodes = append(nodes, nod)
		}
	}

	return nodes, nil
}

func randString() string {
	return fmt.Sprintf("%d", rand.Intn(math.MaxInt32))
}

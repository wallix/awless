// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package memory provides a volatile memory-based implementation of the
// storage.Store and storage.Graph interfaces.
package memory

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/google/badwolf/storage"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

// DefaultStore provides a volatile in memory store.
var DefaultStore storage.Store

func init() {
	DefaultStore = NewStore()
}

type memoryStore struct {
	graphs map[string]storage.Graph
	rwmu   sync.RWMutex
}

// NewStore creates a new memory store.
func NewStore() storage.Store {
	return &memoryStore{
		graphs: make(map[string]storage.Graph),
	}
}

// Name returns the ID of the backend being used.
func (s *memoryStore) Name(ctx context.Context) string {
	return "VOLATILE"
}

// Version returns the version of the driver implementation.
func (s *memoryStore) Version(ctx context.Context) string {
	return "0.1.vcli"
}

// NewGraph creates a new graph.
func (s *memoryStore) NewGraph(ctx context.Context, id string) (storage.Graph, error) {
	g := &memory{
		id:    id,
		idx:   make(map[string]*triple.Triple),
		idxS:  make(map[string]map[string]*triple.Triple),
		idxP:  make(map[string]map[string]*triple.Triple),
		idxO:  make(map[string]map[string]*triple.Triple),
		idxSP: make(map[string]map[string]*triple.Triple),
		idxPO: make(map[string]map[string]*triple.Triple),
		idxSO: make(map[string]map[string]*triple.Triple),
	}

	s.rwmu.Lock()
	defer s.rwmu.Unlock()
	if _, ok := s.graphs[id]; ok {
		return nil, fmt.Errorf("memory.NewGraph(%q): graph already exists", id)
	}
	s.graphs[id] = g
	return g, nil
}

// Graph returns an existing graph if available. Getting a non existing
// graph should return an error.
func (s *memoryStore) Graph(ctx context.Context, id string) (storage.Graph, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	if g, ok := s.graphs[id]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("memory.Graph(%q): graph does not exist", id)
}

// DeleteGraph deletes an existing graph. Deleting a non existing graph
// should return an error.
func (s *memoryStore) DeleteGraph(ctx context.Context, id string) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()
	if _, ok := s.graphs[id]; ok {
		delete(s.graphs, id)
		return nil
	}
	return fmt.Errorf("memory.DeleteGraph(%q): graph does not exist", id)
}

// GraphNames returns the current available graph names in the store.
func (s *memoryStore) GraphNames(ctx context.Context, names chan<- string) error {
	if names == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()
	for k := range s.graphs {
		names <- k
	}
	close(names)
	return nil
}

// memory provides an memory-based volatile implementation of the graph API.
type memory struct {
	id    string
	rwmu  sync.RWMutex
	idx   map[string]*triple.Triple
	idxS  map[string]map[string]*triple.Triple
	idxP  map[string]map[string]*triple.Triple
	idxO  map[string]map[string]*triple.Triple
	idxSP map[string]map[string]*triple.Triple
	idxPO map[string]map[string]*triple.Triple
	idxSO map[string]map[string]*triple.Triple
}

// ID returns the id for this graph.
func (m *memory) ID(ctx context.Context) string {
	return m.id
}

// AddTriples adds the triples to the storage.
func (m *memory) AddTriples(ctx context.Context, ts []*triple.Triple) error {
	m.rwmu.Lock()
	defer m.rwmu.Unlock()
	for _, t := range ts {
		suuid := t.UUID().String()
		sUUID := t.Subject().UUID().String()
		pUUID := t.Predicate().UUID().String()
		oUUID := t.Object().UUID().String()
		// Update master index
		m.idx[suuid] = t

		if _, ok := m.idxS[sUUID]; !ok {
			m.idxS[sUUID] = make(map[string]*triple.Triple)
		}
		m.idxS[sUUID][suuid] = t

		if _, ok := m.idxP[pUUID]; !ok {
			m.idxP[pUUID] = make(map[string]*triple.Triple)
		}
		m.idxP[pUUID][suuid] = t

		if _, ok := m.idxO[oUUID]; !ok {
			m.idxO[oUUID] = make(map[string]*triple.Triple)
		}
		m.idxO[oUUID][suuid] = t

		key := strings.Join([]string{sUUID, pUUID}, ":")
		if _, ok := m.idxSP[key]; !ok {
			m.idxSP[key] = make(map[string]*triple.Triple)
		}
		m.idxSP[key][suuid] = t

		key = strings.Join([]string{pUUID, oUUID}, ":")
		if _, ok := m.idxPO[key]; !ok {
			m.idxPO[key] = make(map[string]*triple.Triple)
		}
		m.idxPO[key][suuid] = t

		key = strings.Join([]string{sUUID, oUUID}, ":")
		if _, ok := m.idxSO[key]; !ok {
			m.idxSO[key] = make(map[string]*triple.Triple)
		}
		m.idxSO[key][suuid] = t
	}
	return nil
}

// RemoveTriples removes the triples from the storage.
func (m *memory) RemoveTriples(ctx context.Context, ts []*triple.Triple) error {
	for _, t := range ts {
		suuid := t.UUID().String()
		sUUID := t.Subject().UUID().String()
		pUUID := t.Predicate().UUID().String()
		oUUID := t.Object().UUID().String()
		// Update master index
		m.rwmu.Lock()
		delete(m.idx, suuid)
		delete(m.idxS[sUUID], suuid)
		delete(m.idxP[pUUID], suuid)
		delete(m.idxO[oUUID], suuid)

		key := strings.Join([]string{sUUID, pUUID}, ":")
		delete(m.idxSP[key], suuid)
		if len(m.idxSP[key]) == 0 {
			delete(m.idxSP, key)
		}

		key = strings.Join([]string{pUUID, oUUID}, ":")
		delete(m.idxPO[key], suuid)
		if len(m.idxPO[key]) == 0 {
			delete(m.idxPO, key)
		}

		key = strings.Join([]string{sUUID, oUUID}, ":")
		delete(m.idxSO[key], suuid)
		if len(m.idxSO[key]) == 0 {
			delete(m.idxSO, key)
		}

		m.rwmu.Unlock()
	}
	return nil
}

// checker provides the mechanics to check if a predicate/triple should be
// considered on a certain operation.
type checker struct {
	max bool
	c   int
	o   *storage.LookupOptions
}

// newChecker creates a new checker for a given LookupOptions configuration.
func newChecker(o *storage.LookupOptions) *checker {
	return &checker{
		max: o.MaxElements > 0,
		c:   o.MaxElements,
		o:   o,
	}
}

// CheckAndUpdate checks if a predicate should be considered and it also updates
// the internal state in case counts are needed.
func (c *checker) CheckAndUpdate(p *predicate.Predicate) bool {
	if c.max && c.c <= 0 {
		return false
	}
	if p.Type() == predicate.Immutable {
		c.c--
		return true
	}
	t, _ := p.TimeAnchor()
	if c.o.LowerAnchor != nil && t.Before(*c.o.LowerAnchor) {
		return false
	}
	if c.o.UpperAnchor != nil && t.After(*c.o.UpperAnchor) {
		return false
	}
	c.c--
	return true
}

// Objects published the objects for the give object and predicate to the
// provided channel.
func (m *memory) Objects(ctx context.Context, s *node.Node, p *predicate.Predicate, lo *storage.LookupOptions, objs chan<- *triple.Object) error {

	sUUID := s.UUID().String()
	pUUID := p.UUID().String()
	spIdx := strings.Join([]string{sUUID, pUUID}, ":")
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(objs)

	ckr := newChecker(lo)
	for _, t := range m.idxSP[spIdx] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			objs <- t.Object()
		}
	}
	return nil
}

// Subject publishes the subjects for the give predicate and object to the
// provided channel.
func (m *memory) Subjects(ctx context.Context, p *predicate.Predicate, o *triple.Object, lo *storage.LookupOptions, subjs chan<- *node.Node) error {
	if subjs == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	pUUID := p.UUID().String()
	oUUID := o.UUID().String()
	poIdx := strings.Join([]string{pUUID, oUUID}, ":")
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(subjs)

	ckr := newChecker(lo)
	for _, t := range m.idxPO[poIdx] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			subjs <- t.Subject()
		}
	}
	return nil
}

// PredicatesForSubjectAndObject publishes all predicates available for the
// given subject and object to the provided channel.
func (m *memory) PredicatesForSubjectAndObject(ctx context.Context, s *node.Node, o *triple.Object, lo *storage.LookupOptions, prds chan<- *predicate.Predicate) error {
	if prds == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	sUUID := s.UUID().String()
	oUUID := o.UUID().String()
	soIdx := strings.Join([]string{sUUID, oUUID}, ":")
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(prds)

	ckr := newChecker(lo)
	for _, t := range m.idxSO[soIdx] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			prds <- t.Predicate()
		}
	}
	return nil
}

// PredicatesForSubject publishes all the predicates known for the given
// subject to the provided channel.
func (m *memory) PredicatesForSubject(ctx context.Context, s *node.Node, lo *storage.LookupOptions, prds chan<- *predicate.Predicate) error {
	if prds == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	sUUID := s.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(prds)
	ckr := newChecker(lo)
	for _, t := range m.idxS[sUUID] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			prds <- t.Predicate()
		}
	}
	return nil
}

// PredicatesForObject publishes all the predicates known for the given object
// to the provided channel.
func (m *memory) PredicatesForObject(ctx context.Context, o *triple.Object, lo *storage.LookupOptions, prds chan<- *predicate.Predicate) error {
	if prds == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	oUUID := o.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(prds)
	ckr := newChecker(lo)
	for _, t := range m.idxO[oUUID] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			prds <- t.Predicate()
		}
	}
	return nil
}

// TriplesForSubject publishes all triples available for the given subject to
// the provided channel.
func (m *memory) TriplesForSubject(ctx context.Context, s *node.Node, lo *storage.LookupOptions, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	sUUID := s.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	ckr := newChecker(lo)
	for _, t := range m.idxS[sUUID] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			trpls <- t
		}
	}
	return nil
}

// TriplesForPredicate publishes all triples available for the given predicate
// to the provided channel.
func (m *memory) TriplesForPredicate(ctx context.Context, p *predicate.Predicate, lo *storage.LookupOptions, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	pUUID := p.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	ckr := newChecker(lo)
	for _, t := range m.idxP[pUUID] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			trpls <- t
		}
	}
	return nil
}

// TriplesForObject publishes all triples available for the given object to the
// provided channel.
func (m *memory) TriplesForObject(ctx context.Context, o *triple.Object, lo *storage.LookupOptions, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	oUUID := o.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	ckr := newChecker(lo)
	for _, t := range m.idxO[oUUID] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			trpls <- t
		}
	}
	return nil
}

// TriplesForSubjectAndPredicate publishes all triples available for the given
// subject and predicate to the provided channel.
func (m *memory) TriplesForSubjectAndPredicate(ctx context.Context, s *node.Node, p *predicate.Predicate, lo *storage.LookupOptions, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	sUUID := s.UUID().String()
	pUUID := p.UUID().String()
	spIdx := strings.Join([]string{sUUID, pUUID}, ":")
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	ckr := newChecker(lo)
	for _, t := range m.idxSP[spIdx] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			trpls <- t
		}
	}
	return nil
}

// TriplesForPredicateAndObject publishes all triples available for the given
// predicate and object to the provided channel.
func (m *memory) TriplesForPredicateAndObject(ctx context.Context, p *predicate.Predicate, o *triple.Object, lo *storage.LookupOptions, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	pUUID := p.UUID().String()
	oUUID := o.UUID().String()
	poIdx := strings.Join([]string{pUUID, oUUID}, ":")
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	ckr := newChecker(lo)
	for _, t := range m.idxPO[poIdx] {
		if ckr.CheckAndUpdate(t.Predicate()) {
			trpls <- t
		}
	}
	return nil
}

// Exist checks if the provided triple exists on the store.
func (m *memory) Exist(ctx context.Context, t *triple.Triple) (bool, error) {
	suuid := t.UUID().String()
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	_, ok := m.idx[suuid]
	return ok, nil
}

// Triples allows to iterate over all available triples by pushing them to the
// provided channel.
func (m *memory) Triples(ctx context.Context, trpls chan<- *triple.Triple) error {
	if trpls == nil {
		return fmt.Errorf("cannot provide an empty channel")
	}
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	defer close(trpls)

	for _, t := range m.idx {
		trpls <- t
	}
	return nil
}

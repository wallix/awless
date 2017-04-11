package triplestore

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

// A source is a persistent yet mutable source or container of triples.
type Source interface {
	Add(...Triple)
	Remove(...Triple)
	Snapshot() RDFGraph
}

// A RDFGraph is an immutable set of triples. It is a snapshot of a source and it is queryable.
type RDFGraph interface {
	Contains(Triple) bool
	Triples() []Triple
	Count() int
	WithSubject(s string) []Triple
	WithPredicate(p string) []Triple
	WithObject(o Object) []Triple
	WithSubjObj(s string, o Object) []Triple
	WithSubjPred(s, p string) []Triple
	WithPredObj(p string, o Object) []Triple
}

type Triples []Triple

func (ts Triples) Equal(others Triples) bool {
	if len(ts) != len(others) {
		return false
	}

	this := make(map[string]struct{})
	for _, tri := range ts {
		this[tri.(*triple).key()] = struct{}{}
	}

	other := make(map[string]struct{})
	for _, tri := range others {
		other[tri.(*triple).key()] = struct{}{}
	}

	return reflect.DeepEqual(this, other)
}

func (ts Triples) Map(fn func(Triple) string) (out []string) {
	for _, t := range ts {
		out = append(out, fn(t))
	}
	return
}

func (ts Triples) String() string {
	joined := strings.Join(ts.Map(
		func(t Triple) string { return fmt.Sprint(t) },
	), "\n")
	return fmt.Sprintf("[%s]", joined)
}

type source struct {
	latestSnap RDFGraph
	updated    uint32 // atomic
	mu         sync.RWMutex
	triples    map[string]Triple
}

// A source is a persistent yet mutable source or container of triples
func NewSource() Source {
	return &source{
		triples:    make(map[string]Triple),
		latestSnap: newGraph(0),
	}
}

func (s *source) isUpdated() bool {
	return atomic.LoadUint32(&s.updated) > 0
}

func (s *source) update() {
	atomic.StoreUint32(&s.updated, uint32(1))
}

func (s *source) reset() {
	atomic.StoreUint32(&s.updated, uint32(0))
}

func (s *source) Add(ts ...Triple) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.update()

	for _, t := range ts {
		tr := t.(*triple)
		s.triples[tr.key()] = t
	}
}
func (s *source) Remove(ts ...Triple) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.update()

	for _, t := range ts {
		tr := t.(*triple)
		delete(s.triples, tr.key())
	}
}

func (s *source) Snapshot() RDFGraph {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isUpdated() {
		return s.latestSnap
	}

	gph := newGraph(len(s.triples))
	defer func() {
		s.latestSnap = gph
		s.reset()
	}()

	for k, t := range s.triples {
		objKey := t.Object().(object).key()
		sub, pred := t.Subject(), t.Predicate()

		gph.s[sub] = append(gph.s[sub], t)
		gph.p[pred] = append(gph.p[pred], t)
		gph.o[objKey] = append(gph.o[objKey], t)

		sp := sub + pred
		gph.sp[sp] = append(gph.sp[sp], t)

		so := sub + objKey
		gph.so[so] = append(gph.so[so], t)

		po := pred + objKey
		gph.po[po] = append(gph.po[po], t)

		gph.spo[k] = t
	}

	for _, t := range gph.spo {
		gph.unique = append(gph.unique, t)
	}

	return gph
}

type graph struct {
	unique     []Triple
	s, p, o    map[string][]Triple
	sp, so, po map[string][]Triple
	spo        map[string]Triple
}

func newGraph(cap int) *graph {
	return &graph{
		s:   make(map[string][]Triple, cap),
		p:   make(map[string][]Triple, cap),
		o:   make(map[string][]Triple, cap),
		sp:  make(map[string][]Triple, cap),
		so:  make(map[string][]Triple, cap),
		po:  make(map[string][]Triple, cap),
		spo: make(map[string]Triple, cap),
	}
}

func (g *graph) Contains(t Triple) bool {
	_, ok := g.spo[t.(*triple).key()]
	return ok
}
func (g *graph) Triples() []Triple {
	return g.unique
}
func (g *graph) Count() int {
	return len(g.spo)
}

func (g *graph) WithSubject(s string) []Triple {
	return g.s[s]
}
func (g *graph) WithPredicate(p string) []Triple {
	return g.p[p]
}
func (g *graph) WithObject(o Object) []Triple {
	return g.o[o.(object).key()]
}
func (g *graph) WithSubjObj(s string, o Object) []Triple {
	return g.so[s+o.(object).key()]
}
func (g *graph) WithSubjPred(s, p string) []Triple {
	return g.sp[s+p]
}
func (g *graph) WithPredObj(p string, o Object) []Triple {
	return g.po[p+o.(object).key()]
}

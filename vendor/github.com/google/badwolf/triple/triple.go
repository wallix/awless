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

// Package triple implements and allows to manipulate BadWolf triples.
package triple

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"github.com/pborman/uuid"
)

// Object is the box that either contains a literal, a predicate or a node.
type Object struct {
	n *node.Node
	p *predicate.Predicate
	l *literal.Literal
}

// String pretty prints the object.
func (o *Object) String() string {
	if o.n != nil {
		return o.n.String()
	}
	if o.l != nil {
		return o.l.String()
	}
	if o.p != nil {
		return o.p.String()
	}
	return "@@@INVALID_OBJECT@@@"
}

// UUID returns a global unique identifier for the given object. It is
// implemented by returning the UIUD of the underlying value stored in the
// object.
func (o *Object) UUID() uuid.UUID {
	switch {
	case o.l != nil:
		return o.l.UUID()
	case o.p != nil:
		return o.p.UUID()
	default:
		return o.n.UUID()
	}
}

// Node attempts to return the boxed node.
func (o *Object) Node() (*node.Node, error) {
	if o.n == nil {
		return nil, fmt.Errorf("triple.Literal does not box a node in %s", o)
	}
	return o.n, nil
}

// Predicate attempts to return the boxed predicate.
func (o *Object) Predicate() (*predicate.Predicate, error) {
	if o.p == nil {
		return nil, fmt.Errorf("triple.Literal does not box a predicate in %s", o)
	}
	return o.p, nil
}

// Literal attempts to return the boxed literal.
func (o *Object) Literal() (*literal.Literal, error) {
	if o.l == nil {
		return nil, fmt.Errorf("triple.Literal does not box a literal in %s", o)
	}
	return o.l, nil
}

// ParseObject attempts to parse an object.
func ParseObject(s string, b literal.Builder) (*Object, error) {
	n, err := node.Parse(s)
	if err == nil {
		return NewNodeObject(n), nil
	}
	l, err := b.Parse(s)
	if err == nil {
		return NewLiteralObject(l), nil
	}
	o, err := predicate.Parse(s)
	if err == nil {

		return NewPredicateObject(o), nil
	}
	return nil, err
}

// NewNodeObject returns a new object that boxes a node.
func NewNodeObject(n *node.Node) *Object {
	return &Object{
		n: n,
	}
}

// NewPredicateObject returns a new object that boxes a predicate.
func NewPredicateObject(p *predicate.Predicate) *Object {
	return &Object{
		p: p,
	}
}

// NewLiteralObject returns a new object that boxes a literal.
func NewLiteralObject(l *literal.Literal) *Object {
	return &Object{
		l: l,
	}
}

// Triple describes a <subject predicate object> used by BadWolf.
type Triple struct {
	s *node.Node
	p *predicate.Predicate
	o *Object
}

// New creates a new triple.
func New(s *node.Node, p *predicate.Predicate, o *Object) (*Triple, error) {
	if s == nil || p == nil || o == nil {
		return nil, fmt.Errorf("triple.New cannot create triples from nil components in <%v %v %v>", s, p, o)
	}
	return &Triple{
		s: s,
		p: p,
		o: o,
	}, nil
}

// Subject returns the subject of the triple.
func (t *Triple) Subject() *node.Node {
	return t.s
}

// Predicate returns the predicate of the triple.
func (t *Triple) Predicate() *predicate.Predicate {
	return t.p
}

// Object returns the object of the triple.
func (t *Triple) Object() *Object {
	return t.o
}

// Equal checks if two triples are identical.
func (t *Triple) Equal(t2 *Triple) bool {
	return uuid.Equal(t.UUID(), t2.UUID())
}

// String marshals the triple into pretty string.
func (t *Triple) String() string {
	return fmt.Sprintf("%s\t%s\t%s", t.s, t.p, t.o)
}

var (
	pSplit *regexp.Regexp
	oSplit *regexp.Regexp
)

func init() {
	pSplit = regexp.MustCompile(">\\s+\"")
	oSplit = regexp.MustCompile("(]\\s+/)|(]\\s+\")")
}

// Parse process the provided text and tries to create a triple. It assumes
// that the provided text contains only one triple.
func Parse(line string, b literal.Builder) (*Triple, error) {
	raw := strings.TrimSpace(line)
	idxp := pSplit.FindIndex([]byte(raw))
	idxo := oSplit.FindIndex([]byte(raw))
	if len(idxp) == 0 || len(idxo) == 0 {
		return nil, fmt.Errorf("triple.Parse could not split s p o  out of %s", raw)
	}
	ss, sp, so := raw[0:idxp[0]+1], raw[idxp[1]-1:idxo[0]+1], raw[idxo[1]-1:]
	s, err := node.Parse(ss)
	if err != nil {
		return nil, fmt.Errorf("triple.Parse failed to parse subject %s with error %v", ss, err)
	}
	p, err := predicate.Parse(sp)
	if err != nil {
		return nil, fmt.Errorf("triple.Parse failed to parse predicate %s with error %v", sp, err)
	}
	o, err := ParseObject(so, b)
	if err != nil {
		return nil, fmt.Errorf("triple.Parse failed to parse object %s with error %v", so, err)
	}
	return New(s, p, o)
}

// Reify given the current triple it returns the original triple and the newly
// reified ones. It also returns the newly created blank node.
func (t *Triple) Reify() ([]*Triple, *node.Node, error) {
	// Function that create the proper reification predicates.
	rp := func(id string, p *predicate.Predicate) (*predicate.Predicate, error) {
		if p.Type() == predicate.Temporal {
			ta, _ := p.TimeAnchor()
			return predicate.NewTemporal(id, *ta)
		}
		return predicate.NewImmutable(id)
	}
	b := node.NewBlankNode()
	s, err := rp("_subject", t.p)
	if err != nil {
		return nil, nil, err
	}
	ts, _ := New(b, s, NewNodeObject(t.s))
	p, err := rp("_predicate", t.p)
	if err != nil {
		return nil, nil, err
	}
	tp, _ := New(b, p, NewPredicateObject(t.p))
	var to *Triple
	if t.o.l != nil {
		o, err := rp("_object", t.p)
		if err != nil {
			return nil, nil, err
		}
		to, _ = New(b, o, NewLiteralObject(t.o.l))
	}
	if t.o.n != nil {
		o, err := rp("_object", t.p)
		if err != nil {
			return nil, nil, err
		}
		to, _ = New(b, o, NewNodeObject(t.o.n))
	}
	if t.o.p != nil {
		o, err := rp("_object", t.p)
		if err != nil {
			return nil, nil, err
		}
		to, _ = New(b, o, NewPredicateObject(t.o.p))
	}

	return []*Triple{t, ts, tp, to}, b, nil
}

// UUID returns a global unique identifier for the given triple. It is
// implemented as the SHA1 UUID of the concatenated UUIDs of the subject,
// predicate, and object.
func (t *Triple) UUID() uuid.UUID {
	var buffer bytes.Buffer

	buffer.Write([]byte(t.s.UUID()))
	buffer.Write([]byte(t.p.UUID()))
	buffer.Write([]byte(t.o.UUID()))

	return uuid.NewSHA1(uuid.NIL, buffer.Bytes())
}

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

// Package storage provides the abstraction to build drivers for BadWolf.
package storage

import (
	"time"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"golang.org/x/net/context"
)

// LookupOptions allows to specify the behavior of the lookup operations.
type LookupOptions struct {
	// MaxElements list the maximum number of elements to return. If not
	// set it returns all the lookup results.
	MaxElements int

	// LowerAnchor, if provided, represents the lower time anchor to be considered.
	LowerAnchor *time.Time

	// UpperAnchor, if provided, represents the upper time anchor to be considered.
	UpperAnchor *time.Time
}

// DefaultLookup provides the default lookup behavior.
var DefaultLookup = &LookupOptions{}

// Store interface describes the low lever API that allows to create new graphs.
type Store interface {
	// Name returns the ID of the backend being used.
	Name(ctx context.Context) string

	// Version returns the version of the driver implementation.
	Version(ctx context.Context) string

	// NewGraph creates a new graph. Creating an already existing graph
	// should return an error.
	NewGraph(ctx context.Context, id string) (Graph, error)

	// Graph returns an existing graph if available. Getting a non existing
	// graph should return an error.
	Graph(ctx context.Context, id string) (Graph, error)

	// DeleteGraph deletes an existing graph. Deleting a non existing graph
	// should return an error.
	DeleteGraph(ctx context.Context, id string) error

	// GraphNames returns the current available graph names in the store.
	GraphNames(ctx context.Context, names chan<- string) error
}

// Graph interface describes the low level API that storage drivers need
// to implement to provide a compliant graph storage that can be used with
// BadWolf.
//
// If you are implementing a driver or just using a low lever driver directly
// it is important for you to keep in mind that you will need to drain the
// provided channel. Otherwise you run the risk of leaking go routines.
type Graph interface {
	// ID returns the id for this graph.
	ID(ctx context.Context) string

	// AddTriples adds the triples to the storage. Adding a triple that already
	// exists should not fail.
	AddTriples(ctx context.Context, ts []*triple.Triple) error

	// RemoveTriples removes the triples from the storage. Removing triples that
	// are not present on the store should not fail.
	RemoveTriples(ctx context.Context, ts []*triple.Triple) error

	// Objects pushes to the provided channel the objects for the given object and
	// predicate. The function does not return immediately but spawns a goroutine
	// to satisfy elements in the channel.
	//
	// Given a subject and a predicate, this method retrieves the objects of
	// triples that match them. By default, if does not limit the maximum number
	// of possible objects returned, unless properly specified by provided lookup
	// options.
	//
	// If the provided predicate is immutable it will return all the possible
	// subject values or the number of max elements specified. There is no
	// requirement on how to sample the returned max elements.
	//
	// If the predicate is an unanchored temporal triple and no time anchors are
	// provided in the lookup options, it will return all the available objects.
	// If time anchors are provided, it will return all the values anchored in the
	// provided time window. If max elements is also provided as part of the
	// lookup options it will return at most max elements. There is no
	// specifications on how that sample should be conducted.
	Objects(ctx context.Context, s *node.Node, p *predicate.Predicate, lo *LookupOptions, objs chan<- *triple.Object) error

	// Subject pushes to the provided channel the subjects for the give predicate
	// and object. The function does not return immediately but spawns a
	// goroutine to satisfy elements in the channel.
	//
	// Given a predicate and an object, this method retrieves the subjects of
	// triples that matches them. By default, it does not limit the maximum number
	// of possible subjects returned, unless properly specified by provided lookup
	// options.
	//
	// If the provided predicate is immutable it will return all the possible
	// subject values or the number of max elements specified. There is no
	// requirement on how to sample the returned max elements.
	//
	// If the predicate is an unanchored temporal triple and no time anchors are
	// provided in the lookup options, it will return all the available subjects.
	// If time anchors are provided, it will return all the values anchored in the
	// provided time window. If max elements is also provided as part of the
	// lookup options it will return the at most max elements. There is no
	// specifications on how that sample should be conducted.
	Subjects(ctx context.Context, p *predicate.Predicate, o *triple.Object, lo *LookupOptions, subs chan<- *node.Node) error

	// PredicatesForSubject pushes to the provided channel all the predicates
	// known for the given subject. The function does not return immediately but
	// spawns a goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available predicates. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided
	// type window would be return. Same sampling consideration apply if max
	// element is provided.
	PredicatesForSubject(ctx context.Context, s *node.Node, lo *LookupOptions, prds chan<- *predicate.Predicate) error

	// PredicatesForObject pushes to the provided channel all the predicates known
	// for the given object. The function returns immediately and spawns a go
	// routine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available predicates. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element
	// is provided.
	PredicatesForObject(ctx context.Context, o *triple.Object, lo *LookupOptions, prds chan<- *predicate.Predicate) error

	// PredicatesForSubjectAndObject pushes to the provided channel all predicates
	// available for the given subject and object. The function does not return
	// immediately but spawns a goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available predicates. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	PredicatesForSubjectAndObject(ctx context.Context, s *node.Node, o *triple.Object, lo *LookupOptions, prds chan<- *predicate.Predicate) error

	// TriplesForSubject pushes to the provided channel all triples available for
	// the given subject. The function does not return immediately but spawns a
	// goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available triples. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	TriplesForSubject(ctx context.Context, s *node.Node, lo *LookupOptions, trpls chan<- *triple.Triple) error

	// TriplesForPredicate pushes to the provided channel all triples available
	// for the given predicate.The function does not return immediately but spawns
	// a goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available triples. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	TriplesForPredicate(ctx context.Context, p *predicate.Predicate, lo *LookupOptions, trpls chan<- *triple.Triple) error

	// TriplesForObject pushes to the provided channel all triples available for
	// the given object. The function does not return immediately but spawns a
	// goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available triples. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	TriplesForObject(ctx context.Context, o *triple.Object, lo *LookupOptions, trpls chan<- *triple.Triple) error

	// TriplesForSubjectAndPredicate pushes to the provided channel all triples
	// available for the given subject and predicate. The function does not return
	// immediately but spawns a goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available triples. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	TriplesForSubjectAndPredicate(ctx context.Context, s *node.Node, p *predicate.Predicate, lo *LookupOptions, trpls chan<- *triple.Triple) error

	// TriplesForPredicateAndObject pushes to the provided channel all triples
	// available for the given predicate and object. The function does not return
	// immediately but spawns a goroutine to satisfy elements in the channel.
	//
	// If the lookup options provide a max number of elements the function will
	// return a sample of the available triples. If time anchor bounds are
	// provided in the lookup options, only predicates matching the provided type
	// window would be return. Same sampling consideration apply if max element is
	// provided.
	TriplesForPredicateAndObject(ctx context.Context, p *predicate.Predicate, o *triple.Object, lo *LookupOptions, trpls chan<- *triple.Triple) error

	// Exist checks if the provided triple exists on the store.
	Exist(ctx context.Context, t *triple.Triple) (bool, error)

	// Triples pushes to the provided channel all available triples in the graph.
	// The function does not return immediately but spawns a goroutine to satisfy
	// elements in the channel.
	Triples(ctx context.Context, trpls chan<- *triple.Triple) error
}

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
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
)

func TestIntersectTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/node<1>  \"to\"@[] /node<1>"))
	a = append(a, parseTriple("/node<2>  \"to\"@[] /node<2>"))
	a = append(a, parseTriple("/node<3>  \"to\"@[] /node<3>"))
	a = append(a, parseTriple("/node<4>  \"to\"@[] /node<4>"))

	b = append(b, parseTriple("/node<0>  \"to\"@[] /node<0>"))
	b = append(b, parseTriple("/node<2>  \"to\"@[] /node<2>"))
	b = append(b, parseTriple("/node<3>  \"to\"@[] /node<3>"))
	b = append(b, parseTriple("/node<5>  \"to\"@[] /node<5>"))
	b = append(b, parseTriple("/node<6>  \"to\"@[] /node<6>"))

	result := intersectTriples(a, b)
	expect = append(expect, parseTriple("/node<2>  \"to\"@[] /node<2>"))
	expect = append(expect, parseTriple("/node<3>  \"to\"@[] /node<3>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestSubtractTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/node<1>  \"to\"@[] /node<1>"))
	a = append(a, parseTriple("/node<2>  \"to\"@[] /node<2>"))
	a = append(a, parseTriple("/node<3>  \"to\"@[] /node<3>"))
	a = append(a, parseTriple("/node<4>  \"to\"@[] /node<4>"))

	b = append(b, parseTriple("/node<0>  \"to\"@[] /node<0>"))
	b = append(b, parseTriple("/node<2>  \"to\"@[] /node<2>"))
	b = append(b, parseTriple("/node<3>  \"to\"@[] /node<3>"))
	b = append(b, parseTriple("/node<5>  \"to\"@[] /node<5>"))
	b = append(b, parseTriple("/node<6>  \"to\"@[] /node<6>"))

	result := subtractTriples(a, b)
	expect = append(expect, parseTriple("/node<1>  \"to\"@[] /node<1>"))
	expect = append(expect, parseTriple("/node<4>  \"to\"@[] /node<4>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = subtractTriples(b, a)
	expect = []*triple.Triple{}
	expect = append(expect, parseTriple("/node<0>  \"to\"@[] /node<0>"))
	expect = append(expect, parseTriple("/node<5>  \"to\"@[] /node<5>"))
	expect = append(expect, parseTriple("/node<6>  \"to\"@[] /node<6>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestAttachTriple(t *testing.T) {
	tri := parseTriple("/node<1>  \"to\"@[] /node<1>")
	objNode, _ := tri.Object().Node()

	g := NewGraphFromTriples([]*triple.Triple{tri})

	l, _ := literal.DefaultBuilder().Build(literal.Text, "trumped")
	attachLiteralToNode(g, objNode, MetaPredicate, l)

	triples, _ := g.TriplesForSubjectPredicate(objNode, MetaPredicate)

	if got, want := len(triples), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	lit, _ := triples[0].Object().Literal()
	if got, want := lit.ToComparableString(), `"trumped"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

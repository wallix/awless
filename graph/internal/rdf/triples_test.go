package rdf

import (
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
)

func TestIntersectTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	a = append(a, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	a = append(a, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	a = append(a, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	b = append(b, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	b = append(b, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	b = append(b, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	b = append(b, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	b = append(b, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	result := intersectTriples(a, b)
	expect = append(expect, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	expect = append(expect, parseTriple("/a<3>  \"to\"@[] /b<3>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestSubstractTriples(t *testing.T) {
	var a, b, expect []*triple.Triple

	a = append(a, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	a = append(a, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	a = append(a, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	a = append(a, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	b = append(b, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	b = append(b, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	b = append(b, parseTriple("/a<3>  \"to\"@[] /b<3>"))
	b = append(b, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	b = append(b, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	result := substractTriples(a, b)
	expect = append(expect, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	expect = append(expect, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = substractTriples(b, a)
	expect = []*triple.Triple{}
	expect = append(expect, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	expect = append(expect, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	expect = append(expect, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	if got, want := MarshalTriples(result), MarshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestAttachTriple(t *testing.T) {
	tri := parseTriple("/a<1>  \"to\"@[] /b<1>")
	objNode, _ := tri.Object().Node()

	g := NewGraphFromTriples([]*triple.Triple{tri})

	l, _ := literal.DefaultBuilder().Build(literal.Text, "trumped")
	attachLiteralToTriple(g, tri, DiffPredicate, l)

	triples, _ := g.TriplesForSubjectPredicate(objNode, DiffPredicate)

	if got, want := len(triples), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	lit, _ := triples[0].Object().Literal()
	if got, want := lit.ToComparableString(), `"trumped"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

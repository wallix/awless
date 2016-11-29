package rdf

import (
	"testing"

	"github.com/google/badwolf/triple"
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

	result := IntersectTriples(a, b)
	expect = append(expect, parseTriple("/a<2>  \"to\"@[] /b<2>"))
	expect = append(expect, parseTriple("/a<3>  \"to\"@[] /b<3>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
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

	result := SubstractTriples(a, b)
	expect = append(expect, parseTriple("/a<1>  \"to\"@[] /b<1>"))
	expect = append(expect, parseTriple("/a<4>  \"to\"@[] /b<4>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = SubstractTriples(b, a)
	expect = []*triple.Triple{}
	expect = append(expect, parseTriple("/a<0>  \"to\"@[] /b<0>"))
	expect = append(expect, parseTriple("/a<5>  \"to\"@[] /b<5>"))
	expect = append(expect, parseTriple("/a<6>  \"to\"@[] /b<6>"))

	if got, want := marshalTriples(result), marshalTriples(expect); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

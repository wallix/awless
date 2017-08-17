package triplestore

import (
	"testing"
	"time"
)

func TestBuildTriple(t *testing.T) {
	tri := SubjPred("subject", "predicate").StringLiteral("any")
	if got, want := tri.Subject(), "subject"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tri.Predicate(), "predicate"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	tri = SubjPredRes("subject", "predicate", "resource")
	if got, want := tri.Subject(), "subject"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tri.Predicate(), "predicate"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	res, _ := tri.Object().Resource()
	if got, want := res, "resource"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	tri, _ = SubjPredLit("subject", "predicate", 3)
	if got, want := tri.Subject(), "subject"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := tri.Predicate(), "predicate"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ := ParseInteger(tri.Object())
	if got, want := lit, 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestBuildObjectFromInterface(t *testing.T) {
	obj, _ := ObjectLiteral(true)
	if got, want := obj, BooleanLiteral(true); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	obj, _ = ObjectLiteral(5)
	if got, want := obj, IntegerLiteral(5); got != want {
		t.Fatalf("got %v, want %v", got, want)

	}
	obj, _ = ObjectLiteral(int64(5))
	if got, want := obj, IntegerLiteral(5); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	obj, _ = ObjectLiteral("any")
	if got, want := obj, StringLiteral("any"); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	now := time.Now()
	obj, _ = ObjectLiteral(now)
	if got, want := obj, DateTimeLiteral(now); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	obj, _ = ObjectLiteral(&now)
	if got, want := obj, DateTimeLiteral(now); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

type stringer struct {
	s string
}

func (s stringer) String() string {
	return s.s
}

func TestBuildAndParseObjectLiteralFromDifferentTypes(t *testing.T) {

	tcases := []struct {
		in  interface{}
		out Object
		exp interface{}
	}{
		{stringer{"stuff"}, StringLiteral("stuff"), "stuff"},

		{float64(2.0), Float64Literal(2.0), float64(2.0)},
		{float32(2.0), Float32Literal(2.0), float32(2.0)},

		{int8(-2), Int8Literal(-2), int8(-2)},
		{int16(-2), Int16Literal(-2), int16(-2)},
		{int32(-2), IntegerLiteral(-2), int(-2)},
		{int64(-2), IntegerLiteral(-2), int(-2)},
		{int(-2), IntegerLiteral(-2), int(-2)},

		{uint8(2), Uint8Literal(2), uint8(2)},
		{uint16(2), Uint16Literal(2), uint16(2)},
		{uint32(2), UintegerLiteral(2), uint(2)},
		{uint64(2), UintegerLiteral(2), uint(2)},
		{uint(2), UintegerLiteral(2), uint(2)},
	}

	for _, tcase := range tcases {
		obj, err := ObjectLiteral(tcase.in)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := obj, tcase.out; !got.Equal(want) {
			t.Fatalf("got %v, want %v", got, want)
		}

		lit, err := ParseLiteral(tcase.out)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := lit, tcase.exp; got != want {
			t.Fatalf("got %v (%T), want %v (%T)", got, got, want, want)
		}
	}
}

func TestUnsupportedLiteralTypesErr(t *testing.T) {
	type any struct{}

	_, err := ObjectLiteral(&any{})
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(UnsupportedLiteralTypeError); !ok {
		t.Fatal("expected error of known type")
	}
}

func TestParseObject(t *testing.T) {
	tri := SubjPred("subject", "predicate").IntegerLiteral(123)
	num, err := ParseInteger(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := num, 123; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	numInt, err := ParseLiteral(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := numInt, 123; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	tri = SubjPred("subject", "predicate").BooleanLiteral(true)
	b, err := ParseBoolean(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := b, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	tri = SubjPred("subject", "predicate").BooleanLiteral(true)
	bInt, err := ParseLiteral(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := bInt, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	now := time.Now()
	tri = SubjPred("subject", "predicate").DateTimeLiteral(now)
	date, err := ParseDateTime(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := date, now.UTC(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	tri = SubjPred("subject", "predicate").DateTimeLiteral(now)
	dateInt, err := ParseLiteral(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := dateInt, now.UTC(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	tri = SubjPred("subject", "predicate").StringLiteral("rdf")
	s, err := ParseString(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := s, "rdf"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	tri = SubjPred("subject", "predicate").StringLiteral("rdf")
	sInt, err := ParseLiteral(tri.Object())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := sInt, "rdf"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	lit, ok := tri.Object().Literal()
	if got, want := ok, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := lit.Value(), "rdf"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := lit.Type(), XsdString; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	_, ok = tri.Object().Resource()
	if got, want := ok, false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}

func TestObjectHasResource(t *testing.T) {
	tri := SubjPred("subject", "predicate").Resource("dbpedia:Bonobo")

	rid, ok := tri.Object().Resource()
	if got, want := ok, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := rid, "dbpedia:Bonobo"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	_, ok = tri.Object().Literal()
	if got, want := ok, false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}

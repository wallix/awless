package utils

import (
	"reflect"
	"testing"
)

type item struct {
	Id string
}

func TestStructWithIdIntersection(t *testing.T) {
	a := []*item{&item{"1"}, &item{"2"}, &item{"3"}, &item{"4"}}
	b := []*item{&item{"0"}, &item{"2"}, &item{"3"}, &item{"5"}, &item{"6"}}
	expect := []string{"2", "3"}

	if r := Intersect(a, b); !reflect.DeepEqual(r, expect) {
		t.Fatalf("got %v, want %v", r, expect)
	}
}

func TestStructWithIdSubstraction(t *testing.T) {
	a := []*item{&item{"1"}, &item{"2"}, &item{"3"}, &item{"4"}}
	b := []*item{&item{"0"}, &item{"2"}, &item{"3"}, &item{"5"}, &item{"6"}}
	expect := []string{"1", "4"}

	if r := Substraction(a, b); !reflect.DeepEqual(r, expect) {
		t.Fatalf("got %v, want %v", r, expect)
	}

	expect = []string{"0", "5", "6"}
	if r := Substraction(b, a); !reflect.DeepEqual(r, expect) {
		t.Fatalf("got %v, want %v", r, expect)
	}
}

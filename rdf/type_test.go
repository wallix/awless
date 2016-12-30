package rdf

import "testing"

func TestResourceTypeToRdfType(t *testing.T) {
	str := "region"
	if got, want := ToRDFType(str), REGION; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := ToResourceType(ToRDFType(str)), str; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := ToResourceType(str), str; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

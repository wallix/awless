package rdf

import "testing"

func TestResourceTypeToRdfType(t *testing.T) {
	if got, want := Region.ToRDFType(), "/region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := Region.String(), "region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	resourceTypes := []ResourceType{Region, Vpc, Subnet, Instance, User, Role, Group, Policy}
	for _, r := range resourceTypes {
		if got, want := NewResourceTypeFromRdfType(r.ToRDFType()), r; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		if got, want := "/"+r.String(), r.ToRDFType(); got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

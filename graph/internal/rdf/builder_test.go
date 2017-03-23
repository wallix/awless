package rdf

import "testing"

func TestBuildTriple(t *testing.T) {
	triple := Subject("id").Predicate("pred").Literal("L33t")
	if got, want := triple.String(), `/node<id>	"pred"@[]	"L33t"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	triple = Subject("id").Predicate("pred").Object("obj")
	if got, want := triple.String(), `/node<id>	"pred"@[]	/node<obj>`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	triple = Subject("id", "ns").Predicate("pred").Object("obj", "ns1", "ns2")
	if got, want := triple.String(), `/node<ns:id>	"pred"@[]	/node<ns1:ns2:obj>`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	triple = Subject("id").Predicate("pred", "ns3", "ns4").Literal("lit")
	if got, want := triple.String(), `/node<id>	"ns3:ns4:pred"@[]	"lit"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTrimNamespaces(t *testing.T) {
	tcases := []struct {
		from, exp string
	}{
		{from: "a:b", exp: "b"},
		{from: "", exp: ""},
		{from: "ns2:a", exp: "a"},
		{from: "c", exp: "c"},
	}
	for _, tcase := range tcases {
		if got, want := TrimNS(tcase.from), tcase.exp; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

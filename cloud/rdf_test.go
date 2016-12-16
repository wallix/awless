package cloud

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"github.com/wallix/awless/rdf"
)

func TestLoadPropertiesTriples(t *testing.T) {
	g := rdf.NewGraph()

	aLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop1", Value: "val1"}))
	if err != nil {
		t.Fatal(err)
	}
	bLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop2", Value: "val2"}))
	if err != nil {
		t.Fatal(err)
	}
	cLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop3", Value: "val3"}))
	if err != nil {
		t.Fatal(err)
	}
	dLiteral, err := literal.DefaultBuilder().Build(literal.Text, mustJsonMarshal(Property{Key: "prop4", Value: "val4"}))
	if err != nil {
		t.Fatal(err)
	}

	one, _ := node.NewNodeFromStrings("/one", "1")
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, aLiteral))
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, bLiteral))
	g.Add(noErrLiteralTriple(one, rdf.PropertyPredicate, cLiteral))
	two, _ := node.NewNodeFromStrings("/two", "2")
	g.Add(noErrLiteralTriple(two, rdf.PropertyPredicate, dLiteral))

	properties, err := LoadPropertiesTriples(g, one)
	if err != nil {
		t.Fatal(err)
	}
	expected := Properties{
		"prop1": "val1",
		"prop2": "val2",
		"prop3": "val3",
	}

	if got, want := properties, expected; !reflect.DeepEqual(properties, expected) {
		t.Fatalf("got %s, want %s", got, want)
	}

	properties, err = LoadPropertiesTriples(g, two)
	expected = Properties{
		"prop4": "val4",
	}

	if got, want := properties, expected; !reflect.DeepEqual(properties, expected) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func mustJsonMarshal(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func noErrLiteralTriple(s *node.Node, p *predicate.Predicate, l *literal.Literal) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewLiteralObject(l))
	if err != nil {
		panic(err)
	}
	return tri
}

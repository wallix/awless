package aws

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

func TestUnmarshalResource(t *testing.T) {
	res := InitResource("inst_1", rdf.Instance)

	g := rdf.NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Tags","Value":[{"Key":"Name","Value":"redis"}]}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":{"Code": 16,"Name":"running"}}"^^type:text`))

	res.UnmarshalFromGraph(g)

	expected := Properties{"Id": "inst_1", "Type": "t2.micro", "PublicIp": "1.2.3.4",
		"State": map[string]interface{}{"Code": float64(16), "Name": "running"},
		"Tags": []interface{}{
			map[string]interface{}{"Key": "Name", "Value": "redis"},
		},
	}

	if got, want := res.properties, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}

	node, err := node.NewNodeFromStrings(rdf.Instance.ToRDFType(), "inst_1")
	if err != nil {
		t.Fatal(err)
	}
	res = InitFromRdfNode(node)
	res.UnmarshalFromGraph(g)
	if got, want := res.properties, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}
}

func TestLoadResources(t *testing.T) {
	g := rdf.NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
	/instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis2"}"^^type:text
	/instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"redis3"}"^^type:text
	/instance<subnet>  "has_type"@[] "/subnet"^^type:text
  /instance<subnet>  "property"@[] "{"Key":"Id","Value":"my subnet"}"^^type:text`))

	expected := []*Resource{
		{kind: rdf.Instance, id: "inst_1", properties: Properties{"Id": "inst_1", "Name": "redis"}},
		{kind: rdf.Instance, id: "inst_2", properties: Properties{"Id": "inst_2", "Name": "redis2"}},
		{kind: rdf.Instance, id: "inst_3", properties: Properties{"Id": "inst_3", "Name": "redis3"}},
	}
	res, err := LoadResourcesFromGraph(g, rdf.Instance)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, r := range expected {
		found := false
		for _, r2 := range res {
			if r2.kind == r.kind && r2.id == r.id && reflect.DeepEqual(r2.properties, r.properties) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%+v not found", r)
		}
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

func newNode(t, id string) *node.Node {
	nodeT, err := node.NewType(t)
	if err != nil {
		panic(err)
	}
	nodeID, err := node.NewID(id)
	if err != nil {
		panic(err)
	}
	return node.NewNode(nodeT, nodeID)
}

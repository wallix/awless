package graph

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/badwolf/triple/node"
)

func TestUnmarshalResource(t *testing.T) {
	res := InitResource("inst_1", Instance)

	g := NewGraph()
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

	node, err := node.NewNodeFromStrings(Instance.ToRDFString(), "inst_1")
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
	g := NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
	/instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis2"}"^^type:text
	/instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"redis3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"CreationDate","Value":"2017-01-10T16:47:18Z"}"^^type:text
	/instance<subnet>  "has_type"@[] "/subnet"^^type:text
  /instance<subnet>  "property"@[] "{"Key":"Id","Value":"my subnet"}"^^type:text`))

	time, _ := time.Parse(time.RFC3339, "2017-01-10T16:47:18Z")

	expected := []*Resource{
		{kind: Instance, id: "inst_1", properties: Properties{"Id": "inst_1", "Name": "redis"}},
		{kind: Instance, id: "inst_2", properties: Properties{"Id": "inst_2", "Name": "redis2"}},
		{kind: Instance, id: "inst_3", properties: Properties{"Id": "inst_3", "Name": "redis3", "CreationDate": time}},
	}
	res, err := LoadResourcesFromGraph(g, Instance)
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

package graph

import (
	"testing"
)

func TestResourceNameToId(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
	/instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis2"}"^^type:text
	/instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"mongo"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"CreationDate","Value":"2017-01-10T16:47:18Z"}"^^type:text
	/subnet<subnet_1>  "has_type"@[] "/subnet"^^type:text
  /subnet<subnet_1>  "property"@[] "{"Key":"Id","Value":"subnet_1"}"^^type:text
  /subnet<subnet_1>  "property"@[] "{"Key":"Name","Value":"mongo"}"^^type:text`))

	tcases := []struct {
		name         string
		resourceType ResourceType
		expectID     string
		ok           bool
	}{
		{name: "redis", resourceType: Instance, expectID: "inst_1", ok: true},
		{name: "redis2", resourceType: Instance, expectID: "inst_2", ok: true},
		{name: "mongo", resourceType: Instance, expectID: "inst_3", ok: true},
		{name: "mongo", resourceType: Subnet, expectID: "subnet_1", ok: true},
		{name: "nothere", expectID: "", ok: false},
	}
	for _, tcase := range tcases {
		a := Alias(tcase.name)
		id, ok := a.ResolveToId(g, tcase.resourceType)
		if got, want := ok, tcase.ok; got != want {
			t.Fatalf("got %t, want %t", got, want)
		}
		if ok {
			if got, want := id, tcase.expectID; got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	}
}

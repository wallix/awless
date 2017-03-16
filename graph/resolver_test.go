package graph

import "testing"

func TestAndResolver(t *testing.T) {
	t.Parallel()
	g := NewGraph()

	g.Unmarshal([]byte(`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /subnet<sub_1>  "has_type"@[] "/subnet"^^type:text
  /subnet<sub_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text`))

	resources, err := g.ResolveResources(&And{[]Resolver{
		&ByType{Typ: "instance"},
		&ByProperty{Name: "Name", Val: "redis"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 1; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	if got, want := resources[0].Id(), "inst_1"; got != want {
		t.Fatalf("got %s want %s", got, want)
	}

	resources, err = g.ResolveResources(&And{[]Resolver{
		&ByType{Typ: "subnet"},
		&ByProperty{Name: "Id", Val: "inst_2"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 0; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

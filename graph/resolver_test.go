package graph_test

import (
	"testing"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestAndResolver(t *testing.T) {
	t.Parallel()
	g := graph.NewGraph()
	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Build(),
		resourcetest.Subnet("sub_1").Prop("Name", "redis").Build(),
	)

	resources, err := g.ResolveResources(&graph.And{[]graph.Resolver{
		&graph.ByType{Typ: "instance"},
		&graph.ByProperty{Name: "Name", Val: "redis"},
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

	resources, err = g.ResolveResources(&graph.And{[]graph.Resolver{
		&graph.ByType{Typ: "subnet"},
		&graph.ByProperty{Name: "ID", Val: "inst_2"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 0; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

package graph_test

import (
	"testing"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestResolvers(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Build(),
		resourcetest.Subnet("sub_1").Prop("Name", "redis").Build(),
	)

	t.Run("and", func(t *testing.T) {
		resources, err := g.ResolveResources(&graph.And{[]graph.Resolver{
			&graph.ByType{Typ: "instance"},
			&graph.ByProperty{Key: "Name", Value: "redis"},
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
			&graph.ByProperty{Key: "ID", Value: "inst_2"},
		},
		})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(resources), 0; got != want {
			t.Fatalf("got %d want %d", got, want)
		}
	})

	t.Run("or", func(t *testing.T) {
		resources, err := g.ResolveResources(&graph.Or{[]graph.Resolver{
			&graph.ByProperty{Key: "Name", Value: "unexisting"},
			&graph.ByProperty{Key: "ID", Value: "sub_1"},
			&graph.ByProperty{Key: "Name", Value: "unexisting"},
		},
		})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(resources), 1; got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if got, want := resources[0].Id(), "sub_1"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	})

	t.Run("by type and property", func(t *testing.T) {
		resources, err := g.ResolveResources(&graph.ByTypeAndProperty{
			Type: "instance", Key: "Name", Value: "redis",
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

		resources, err = g.ResolveResources(&graph.ByTypeAndProperty{
			Type: "subnet", Key: "ID", Value: "inst_2",
		})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(resources), 0; got != want {
			t.Fatalf("got %d want %d", got, want)
		}
	})
}

package graph_test

import (
	"testing"

	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestFilterGraph(t *testing.T) {
	g := graph.NewGraph()

	inst1 := resourcetest.Instance("inst_1").Prop("Tags", []string{"Creator=God", "Env=Cosmos"}).Build()
	inst2 := resourcetest.Instance("inst_2").Prop("Tags", []string{"Creator=God"}).Prop("Name", "redis").Build()
	subnet1 := resourcetest.Subnet("sub_1").Build()

	g.AddResource(inst1, inst2, subnet1)

	filtered, err := g.Filter("subnet")
	if err != nil {
		t.Fatal(err)
	}
	subnets, _ := filtered.GetAllResources("subnet")
	if got, want := len(subnets), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	instances, _ := filtered.GetAllResources("instance")
	if got, want := len(instances), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filterFn := func(r *graph.Resource) bool {
		if r.Properties()[properties.ID] == "inst_1" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterFn)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Properties()[properties.ID], "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filterOne := func(r *graph.Resource) bool {
		if r.Properties()[properties.ID] == "inst_2" {
			return true
		}
		return false
	}
	filterTwo := func(r *graph.Resource) bool {
		if r.Properties()[properties.Name] == "redis" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterOne, filterTwo)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Id(), "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties()["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filtered, _ = g.Filter("instance",
		graph.BuildPropertyFilterFunc(properties.ID, "inst"),
		graph.BuildPropertyFilterFunc(properties.Name, "Redis"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Id(), "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties()["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filtered, _ = g.Filter("instance",
		graph.BuildTagFilterFunc("Creator", "God"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := hasResource(instances, inst1), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := hasResource(instances, inst2), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	filtered, _ = g.Filter("instance",
		graph.BuildTagFilterFunc("Creator", "God"),
		graph.BuildTagFilterFunc("Env", "Cosmos"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := hasResource(instances, inst1), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	filtered, _ = g.Filter("instance",
		graph.BuildTagKeyFilterFunc("Creator"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := hasResource(instances, inst1), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := hasResource(instances, inst2), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	filtered, _ = g.Filter("instance",
		graph.BuildTagValueFilterFunc("Cosmos"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := hasResource(instances, inst1), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}

func hasResource(arr []*graph.Resource, r *graph.Resource) bool {
	for _, a := range arr {
		if a.Id() == r.Id() {
			return true
		}
	}
	return false
}

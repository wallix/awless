package graph

import "testing"

func TestFilterGraph(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(
		`/instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /subnet<sub_1>  "has_type"@[] "/subnet"^^type:text
  /subnet<sub_1>  "property"@[] "{"Key":"Id","Value":"sub_1"}"^^type:text`))
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

	filterFn := func(r *Resource) bool {
		if r.Properties["Id"] == "inst_1" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterFn)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Properties["Id"], "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filterOne := func(r *Resource) bool {
		if r.Properties["Id"] == "inst_2" {
			return true
		}
		return false
	}
	filterTwo := func(r *Resource) bool {
		if r.Properties["Name"] == "redis" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterOne, filterTwo)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Properties["Id"], "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filtered, _ = g.Filter("instance",
		BuildPropertyFilterFunc("Id", "inst"),
		BuildPropertyFilterFunc("Name", "Redis"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Properties["Id"], "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

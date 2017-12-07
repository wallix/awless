/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestGetResource(t *testing.T) {
	g := graph.NewGraph()

	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Prop("Type", "t2.micro").Prop("PublicIP", "1.2.3.4").Prop("State", "running").Build(),
	)

	res, err := g.GetResource("instance", "inst_1")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{properties.ID: "inst_1", properties.Type: "t2.micro", properties.PublicIP: "1.2.3.4",
		properties.State: "running",
		properties.Name:  "redis",
	}

	if got, want := res.Properties(), expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}
}

func TestFindResources(t *testing.T) {
	t.Parallel()
	g := graph.NewGraph()

	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Build(),
		resourcetest.Subnet("sub_1").Prop("Name", "redis").Build(),
	)

	t.Run("FindResource", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResource("inst_1")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := res.Properties()["Name"], "redis"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}

		if res, err = g.FindResource("none"); err != nil {
			t.Fatal(err)
		}
		if res != nil {
			t.Fatalf("expected nil got %v", res)
		}

		if res, err = g.FindResource("sub_1"); err != nil {
			t.Fatal(err)
		}
		if got, want := res.Type(), "subnet"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	})
	t.Run("FindResourcesByProperty", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResourcesByProperty("ID", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		expected := []*graph.Resource{
			resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v want %+v", got, want)
		}
		res, err = g.FindResourcesByProperty("Name", "redis")
		if err != nil {
			t.Fatal(err)
		}
		expected = []*graph.Resource{
			resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
			resourcetest.Subnet("sub_1").Prop("Name", "redis").Build(),
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if res[0].Id() == expected[0].Id() {
			if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		} else {
			if got, want := res[0], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		}
	})
}

func TestGetAllResources(t *testing.T) {
	g := graph.NewGraph()

	time, _ := time.Parse(time.RFC3339, "2017-01-10T16:47:18Z")

	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Prop("Name", "redis2").Build(),
		resourcetest.Instance("inst_3").Prop("Name", "redis3").Prop("Created", time).Build(),
		resourcetest.Subnet("subnet").Prop("Name", "redis").Build(),
	)

	expected := []*graph.Resource{
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Prop("Name", "redis2").Build(),
		resourcetest.Instance("inst_3").Prop("Name", "redis3").Prop("Created", time).Build(),
	}
	res, err := g.GetAllResources("instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, r := range expected {
		found := false
		for _, r2 := range res {
			if r2.Type() == r.Type() && r2.Id() == r.Id() && reflect.DeepEqual(r2.Properties(), r.Properties()) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%+v not found", r)
		}
	}
}

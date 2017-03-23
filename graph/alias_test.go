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
	"testing"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
)

func TestResourceNameToId(t *testing.T) {
	g := graph.NewGraph()
	g.AddResource(
		resourcetest.Instance("inst_1").Prop("Name", "redis").Build(),
		resourcetest.Instance("inst_2").Prop("Name", "redis2").Build(),
		resourcetest.Instance("inst_3").Prop("Name", "mongo").Build(),
		resourcetest.Subnet("subnet_1").Prop("Name", "mongo").Build(),
	)

	tcases := []struct {
		name         string
		resourceType string
		expectID     string
		ok           bool
	}{
		{name: "redis", resourceType: "instance", expectID: "inst_1", ok: true},
		{name: "redis2", resourceType: "instance", expectID: "inst_2", ok: true},
		{name: "mongo", resourceType: "instance", expectID: "inst_3", ok: true},
		{name: "mongo", resourceType: "subnet", expectID: "subnet_1", ok: true},
		{name: "nothere", expectID: "", ok: false},
	}
	for i, tcase := range tcases {
		a := graph.Alias(tcase.name)
		id, ok := a.ResolveToId(g, tcase.resourceType)
		if got, want := ok, tcase.ok; got != want {
			t.Fatalf("%d: got %t, want %t", i, got, want)
		}
		if ok {
			if got, want := id, tcase.expectID; got != want {
				t.Fatalf("%d: got %s, want %s", i, got, want)
			}
		}
	}
}

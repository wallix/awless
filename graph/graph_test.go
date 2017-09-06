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

package graph

import (
	"testing"

	tstore "github.com/wallix/triplestore"
)

func TestFindAncestors(t *testing.T) {
	g := NewGraph()
	inst := InitResource("instance", "inst_1")
	sub := InitResource("subnet", "subnet_1")
	region := InitResource("region", "north-korea")
	g.AddResource(inst, sub, region)
	g.AddParentRelation(sub, inst)
	g.AddParentRelation(region, sub)

	res := g.FindAncestor(inst, "region")
	if got, want := res.Id(), "north-korea"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := res.Type(), "region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	res = g.FindAncestor(inst, "subnet")
	if got, want := res.Id(), "subnet_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := res.Type(), "subnet"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	res = g.FindAncestor(sub, "region")
	if got, want := res.Id(), "north-korea"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := res.Type(), "region"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	res = g.FindAncestor(region, "region")
	if res != nil {
		t.Fatalf("got %v, want nil", res)
	}
}

func TestAddGraphRelation(t *testing.T) {
	t.Run("Add parent", func(t *testing.T) {
		g := NewGraph()
		g.AddResource(InitResource("instance", "inst_1"))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddParentRelation(InitResource("subnet", "subnet_1"), res)

		expTriples := tstore.Triples([]tstore.Triple{
			tstore.SubjPred("inst_1", "rdf:type").Resource("cloud-owl:Instance"),
			tstore.SubjPred("inst_1", "cloud:id").StringLiteral("inst_1"),
			tstore.SubjPred("subnet_1", "cloud-rel:parentOf").Resource("inst_1"),
		})

		if got, want := tstore.Triples(g.store.Snapshot().Triples()), expTriples; !got.Equal(want) {
			t.Fatalf("got\n%v\nwant\n%v\n", got, want)
		}
	})

	t.Run("Add applies on", func(t *testing.T) {
		g := NewGraph()
		g.AddResource(InitResource("instance", "inst_1"))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddAppliesOnRelation(InitResource("subnet", "subnet_1"), res)

		expTriples := tstore.Triples([]tstore.Triple{
			tstore.SubjPred("inst_1", "rdf:type").Resource("cloud-owl:Instance"),
			tstore.SubjPred("inst_1", "cloud:id").StringLiteral("inst_1"),
			tstore.SubjPred("subnet_1", "cloud-rel:applyOn").Resource("inst_1"),
		})

		if got, want := tstore.Triples(g.store.Snapshot().Triples()), expTriples; !got.Equal(want) {
			t.Fatalf("got\n%q\nwant\n%q\n", got, want)
		}
	})
}

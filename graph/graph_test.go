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
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/match"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/cloud/rdf"
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

func TestFind(t *testing.T) {
	g := NewGraph()
	i1 := instResource("i1").prop("Name", "redis").prop("Subnet", "s1").prop(properties.Tags, []string{"TagKey1=TagValue1"}).build()
	i2 := instResource("i2").prop("Subnet", "s1").prop(properties.Tags, []string{"TagKey1=TagValue1"}).build()
	s1 := subResource("s1").prop(properties.Tags, []string{"TagKey2=TagValue2"}).build()
	s2 := subResource("s2").prop("Name", "prod").prop("ActiveServicesCount", 42).build()
	v1 := vpcResource("v1").prop("Name", "prod").prop(properties.Tags, []string{"TagKey2=TagValue2"}).build()
	g.AddResource(i1, i2, s1, s2, v1)
	tcases := []struct {
		query     cloud.Query
		expectRes []cloud.Resource
	}{
		{
			query:     cloud.NewQuery("instance"),
			expectRes: []cloud.Resource{i1, i2},
		},
		{
			query:     cloud.NewQuery("instance").Match(match.Property("Subnet", "s1")),
			expectRes: []cloud.Resource{i1, i2},
		},
		{
			query:     cloud.NewQuery("instance").Match(match.And(match.Property("Subnet", "s1"), match.Property(properties.ID, "i2"))),
			expectRes: []cloud.Resource{i2},
		},
		{
			query:     cloud.NewQuery("subnet"),
			expectRes: []cloud.Resource{s1, s2},
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ID", "s1")),
			expectRes: []cloud.Resource{s1},
		},
		{
			query:     cloud.NewQuery("instance").Match(match.Property("Name", "nothing")),
			expectRes: nil,
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ID", "S1")),
			expectRes: nil,
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ID", "S1").IgnoreCase()),
			expectRes: []cloud.Resource{s1},
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ActiveServicesCount", "42")),
			expectRes: nil,
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ActiveServicesCount", "42").MatchString()),
			expectRes: []cloud.Resource{s2},
		},
		{
			query:     cloud.NewQuery("instance").Match(match.Tag("TagKey1", "TagValue1")),
			expectRes: []cloud.Resource{i1, i2},
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.TagKey("TagKey2")),
			expectRes: []cloud.Resource{s1},
		},
		{
			query:     cloud.NewQuery("vpc").Match(match.TagValue("TagValue2")),
			expectRes: []cloud.Resource{v1},
		},
		{
			query:     cloud.NewQuery("subnet", "vpc"),
			expectRes: []cloud.Resource{s1, s2, v1},
		},
		{
			query:     cloud.NewQuery("instance").Match(match.And(match.Property(properties.ID, "i2"), match.Property(properties.Name, "redis"))),
			expectRes: nil,
		},
		{
			query:     cloud.NewQuery("instance").Match(match.Or(match.Property(properties.ID, "i2"), match.Property(properties.Name, "redis"))),
			expectRes: []cloud.Resource{i1, i2},
		},
	}
	for i, tcase := range tcases {
		res, err := g.Find(tcase.query)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		sort.Slice(res, func(i int, j int) bool {
			if res[i].Type() == res[j].Type() {
				return res[i].Id() <= res[j].Id()
			}
			return res[i].Type() < res[j].Type()
		})
		if got, want := res, tcase.expectRes; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestFindWithProperties(t *testing.T) {
	g := NewGraph()
	i1 := instResource("i1").prop("Name", "redis").prop("Subnet", "s1").prop(properties.Tags, []string{"TagKey1=TagValue1"}).build()
	i2 := instResource("i2").prop("Subnet", "s1").prop(properties.Tags, []string{"TagKey1=TagValue1"}).build()
	s1 := subResource("s1").prop(properties.Tags, []string{"TagKey2=TagValue2"}).build()
	s2 := subResource("s2").prop("Name", "prod").prop("ActiveServicesCount", 42).build()
	v1 := vpcResource("v1").prop("Name", "prod").prop(properties.Tags, []string{"TagKey2=TagValue2"}).build()
	g.AddResource(i1, i2, s1, s2, v1)
	tcases := []struct {
		props     map[string]interface{}
		expectRes []cloud.Resource
	}{
		{
			props:     map[string]interface{}{"Name": "redis"},
			expectRes: []cloud.Resource{i1},
		},
		{
			props:     map[string]interface{}{"Name": "prod"},
			expectRes: []cloud.Resource{s2, v1},
		},
		{
			props:     map[string]interface{}{"Name": "nothing"},
			expectRes: nil,
		},
	}
	for i, tcase := range tcases {
		res, err := g.FindWithProperties(tcase.props)
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		sort.Slice(res, func(i int, j int) bool {
			if res[i].Type() == res[j].Type() {
				return res[i].Id() <= res[j].Id()
			}
			return res[i].Type() < res[j].Type()
		})
		if got, want := res, tcase.expectRes; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestFindOne(t *testing.T) {
	g := NewGraph()
	i1 := instResource("i1").prop("Name", "redis").prop("Subnet", "s1").build()
	i2 := instResource("i2").prop("Subnet", "s1").build()
	s1 := subResource("s1").build()
	v1 := vpcResource("s1").build()
	g.AddResource(i1, i2, s1, v1)
	tcases := []struct {
		query             cloud.Query
		expectRes         cloud.Resource
		expectErrContains string
	}{
		{
			query:     cloud.NewQuery("instance").Match(match.Property("Name", "redis")),
			expectRes: i1,
		},
		{
			query:             cloud.NewQuery("instance").Match(match.Property("Subnet", "s1")),
			expectErrContains: "multiple",
		},
		{
			query:     cloud.NewQuery("subnet"),
			expectRes: s1,
		},
		{
			query:     cloud.NewQuery("subnet").Match(match.Property("ID", "s1")),
			expectRes: s1,
		},
		{
			query:             cloud.NewQuery("instance").Match(match.Property("Name", "nothing")),
			expectErrContains: "not found",
		},
	}
	for i, tcase := range tcases {
		res, err := g.FindOne(tcase.query)
		if tcase.expectErrContains != "" {
			if err == nil {
				t.Fatalf("%d: expect error contains '%s', got nil", i+1, tcase.expectErrContains)
			}
			if !strings.Contains(err.Error(), tcase.expectErrContains) {
				t.Fatalf("%d: expect error contains '%s', got %s", i+1, tcase.expectErrContains, err.Error())
			}
			continue
		}
		if err != nil {
			t.Fatalf("%d: %s", i+1, err)
		}
		if got, want := res, tcase.expectRes; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d: got %v, want %v", i+1, got, want)
		}
	}
}

func TestResourceChildrenAndSiblings(t *testing.T) {
	g := NewGraph()
	v1 := InitResource("vpc", "vpc_1")
	s1 := InitResource("subnet", "sub_1")
	s2 := InitResource("subnet", "sub_2")
	i1 := InitResource("instance", "inst_1")
	i2 := InitResource("instance", "inst_2")
	i3 := InitResource("instance", "inst_3")
	sg1 := InitResource("securitygroup", "secgroup_1")
	g.AddResource(v1, s1, s2, i1, i2, i3, sg1)
	g.AddParentRelation(v1, s1)
	g.AddParentRelation(v1, s2)
	g.AddParentRelation(v1, sg1)
	g.AddParentRelation(s1, i1)
	g.AddParentRelation(s1, i2)
	g.AddParentRelation(s2, i3)
	g.AddAppliesOnRelation(sg1, i1)
	g.AddAppliesOnRelation(sg1, i3)

	t.Run("ResourceChildren", func(t *testing.T) {
		tcases := []struct {
			from         cloud.Resource
			relation     string
			recursive    bool
			expRelations []cloud.Resource
		}{
			{from: v1, relation: rdf.ChildrenOfRel, recursive: false, expRelations: []cloud.Resource{sg1, s1, s2}},
			{from: v1, relation: rdf.ChildrenOfRel, recursive: true, expRelations: []cloud.Resource{i1, i2, i3, sg1, s1, s2}},
			{from: s1, relation: rdf.ChildrenOfRel, recursive: false, expRelations: []cloud.Resource{i1, i2}},
			{from: s1, relation: rdf.ParentOf, recursive: false, expRelations: []cloud.Resource{v1}},
			{from: i1, relation: rdf.ParentOf, recursive: false, expRelations: []cloud.Resource{s1}},
			{from: i1, relation: rdf.ParentOf, recursive: true, expRelations: []cloud.Resource{s1, v1}},
			{from: i1, relation: rdf.DependingOnRel, recursive: false, expRelations: []cloud.Resource{sg1}},
			{from: sg1, relation: rdf.ApplyOn, recursive: false, expRelations: []cloud.Resource{i1, i3}},
		}
		for i, tcase := range tcases {
			res, err := g.ResourceRelations(tcase.from, tcase.relation, tcase.recursive)
			if err != nil {
				t.Fatalf("%d: %s", i+1, err)
			}
			sort.Slice(res, func(i int, j int) bool {
				if res[i].Type() == res[j].Type() {
					return res[i].Id() <= res[j].Id()
				}
				return res[i].Type() < res[j].Type()
			})
			if got, want := res, tcase.expRelations; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: got %v, want %v", i+1, got, want)
			}
		}
	})

	t.Run("ResourceSiblings", func(t *testing.T) {
		tcases := []struct {
			from        cloud.Resource
			expSiblings []cloud.Resource
		}{
			{from: v1, expSiblings: nil},
			{from: s1, expSiblings: []cloud.Resource{s2}},
			{from: s2, expSiblings: []cloud.Resource{s1}},
			{from: i1, expSiblings: []cloud.Resource{i2}},
		}
		for i, tcase := range tcases {
			res, err := g.ResourceSiblings(tcase.from)
			if err != nil {
				t.Fatalf("%d: %s", i+1, err)
			}
			if got, want := res, tcase.expSiblings; !reflect.DeepEqual(got, want) {
				t.Fatalf("%d: got %v, want %v", i+1, got, want)
			}
		}
	})
}

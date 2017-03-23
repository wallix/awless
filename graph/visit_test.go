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

	"github.com/wallix/awless/graph"
)

func TestCollectors(t *testing.T) {
	g := graph.NewGraph()
	i1 := graph.InitResource("instance", "inst_1")
	i2 := graph.InitResource("instance", "inst_2")
	i3 := graph.InitResource("instance", "inst_3")
	s1 := graph.InitResource("subnet", "sub_1")
	s2 := graph.InitResource("subnet", "sub_2")
	v1 := graph.InitResource("vpc", "vpc_1")
	v2 := graph.InitResource("vpc", "vpc_2")
	err := g.AddResource(i1, i2, i3, s1, s2, v1, v2)
	if err != nil {
		t.Fatal(err)
	}
	g.AddParentRelation(s1, i1)
	g.AddParentRelation(s1, i2)
	g.AddParentRelation(s2, i3)
	g.AddParentRelation(v1, s1)
	g.AddParentRelation(v2, s2)

	var collect []*graph.Resource
	tcases := []struct {
		vis graph.Visitor
		exp []*graph.Resource
	}{
		{vis: &graph.ParentsVisitor{From: i1, Each: graph.VisitorCollectFunc(&collect)}, exp: []*graph.Resource{s1, v1}},
		{vis: &graph.ParentsVisitor{From: s2, Each: graph.VisitorCollectFunc(&collect), IncludeFrom: true}, exp: []*graph.Resource{s2, v2}},
		{vis: &graph.ParentsVisitor{From: v2, Each: graph.VisitorCollectFunc(&collect)}, exp: []*graph.Resource{}},
		{vis: &graph.ChildrenVisitor{From: i1, Each: graph.VisitorCollectFunc(&collect)}, exp: []*graph.Resource{}},
		{vis: &graph.ChildrenVisitor{From: s2, Each: graph.VisitorCollectFunc(&collect), IncludeFrom: true}, exp: []*graph.Resource{s2, i3}},
		{vis: &graph.ChildrenVisitor{From: v1, Each: graph.VisitorCollectFunc(&collect)}, exp: []*graph.Resource{s1, i1, i2}},
		{vis: &graph.SiblingsVisitor{From: i1, Each: graph.VisitorCollectFunc(&collect), IncludeFrom: true}, exp: []*graph.Resource{i1, i2}},
		{vis: &graph.SiblingsVisitor{From: s2, Each: graph.VisitorCollectFunc(&collect)}, exp: []*graph.Resource{}},
	}

	for i, tcase := range tcases {
		collect = []*graph.Resource{}

		if err := g.Accept(tcase.vis); err != nil {
			t.Fatal(err)
		}
		if got, want := collect, tcase.exp; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d. got %#v, want %#v", i, got, want)
		}
	}

}

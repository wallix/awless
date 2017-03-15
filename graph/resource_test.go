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
	"testing"
)

func TestSortResource(t *testing.T) {
	resources := []*Resource{{id: "b"}, {id: "c"}, {id: "a"}}
	sort.Sort(ResourceById(resources))

	if got, want := len(resources), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := resources[0].Id(), "a"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := resources[1].Id(), "b"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := resources[2].Id(), "c"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestEqualResources(t *testing.T) {
	i1 := &Resource{id: "inst_1", kind: "instance"}
	i2 := &Resource{id: "inst_2", kind: "instance"}
	i3 := &Resource{id: "toto", kind: "instance"}
	s1 := &Resource{id: "subnet_1", kind: "subnet"}
	s2 := &Resource{id: "subnet_1", kind: "subnet"}
	s3 := &Resource{id: "toto", kind: "subnet"}
	empty := &Resource{}
	tcases := []struct {
		from, to *Resource
		exp      bool
	}{
		{from: i1, to: i1, exp: true},
		{from: i1, to: i2, exp: false},
		{from: i1, to: i3, exp: false},
		{from: i1, to: s1, exp: false},
		{from: s1, to: s2, exp: true},
		{from: s2, to: s1, exp: true},
		{from: s1, to: s3, exp: false},
		{from: i3, to: s3, exp: false},
		{from: empty, to: empty, exp: true},
		{from: empty, to: nil, exp: false},
		{from: nil, to: empty, exp: false},
		{from: nil, to: nil, exp: true},
		{from: empty, to: i1, exp: false},
		{from: i1, to: empty, exp: false},
	}

	for _, tcase := range tcases {
		if tcase.from.Same(tcase.to) != tcase.exp {
			t.Fatalf("expected %t, from %+v, to %+v", tcase.exp, tcase.from, tcase.to)
		}
	}
}

func TestPrintResource(t *testing.T) {
	tcases := []struct {
		res *Resource
		exp string
	}{
		{res: &Resource{id: "inst_1", kind: "instance"}, exp: "inst_1[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", Properties: Properties{"Id": "notthis"}}, exp: "inst_1[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", Properties: Properties{"Id": "notthis", "Name": "to-display"}}, exp: "@to-display[instance]"},
		{res: &Resource{id: "inst_1", kind: "instance", Properties: Properties{"Name": ""}}, exp: "inst_1[instance]"},
		{res: &Resource{kind: "instance", Properties: Properties{"Id": "notthis", "Name": "to-display"}}, exp: "@to-display[instance]"},
		{res: &Resource{}, exp: "[none]"},
		{res: nil, exp: "[none]"},
	}
	for _, tcase := range tcases {
		if got, want := tcase.res.String(), tcase.exp; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestReduceResources(t *testing.T) {
	res := Resources{{id: "1"}, {id: "2"}, {id: "3"}}
	if got, want := res.Map(func(r *Resource) string { return r.String() }), []string{"1[]", "2[]", "3[]"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCompareProperties(t *testing.T) {
	props1 := Properties(map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "3",
		"four":  4,
	})
	props2 := Properties(map[string]interface{}{
		"zero":  0,
		"two":   2,
		"three": "3",
		"four":  "4",
		"five":  "5",
	})

	exp := Properties(map[string]interface{}{"one": 1, "four": 4})
	if got, want := props1.Subtract(props2), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	exp = Properties(map[string]interface{}{"zero": 0, "four": "4", "five": "5"})
	if got, want := props2.Subtract(props1), exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

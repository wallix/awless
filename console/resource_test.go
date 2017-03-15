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

package console

import (
	"bytes"
	"testing"

	"github.com/wallix/awless/graph"
)

func TestResourceDisplay(t *testing.T) {
	g := graph.NewGraph()

	res1 := graph.InitResource("inst_1", "instance")
	res1.Properties = map[string]interface{}{
		"Id":     "inst_1",
		"Name":   "instance 1",
		"Prop 1": "prop 1",
		"Prop 2": "prop 2",
	}
	res2 := graph.InitResource("inst_2", "instance")

	g.AddResource(res1, res2)

	r, err := g.GetResource("instance", "inst_1")
	if err != nil {
		t.Fatal(err)
	}

	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}

	displayer := BuildOptions(
		WithHeaders(headers),
		WithFormat("table"),
	).SetSource(r).Build()

	expected := `+------------+------------+
| PROPERTY ▲ |   VALUE    |
+------------+------------+
| Id         | inst_1     |
| Name       | instance 1 |
| Prop 1     | prop 1     |
| Prop 2     | prop 2     |
+------------+------------+
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

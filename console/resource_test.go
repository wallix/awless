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
	"github.com/wallix/awless/graph/resourcetest"
)

func TestResourceDisplay(t *testing.T) {
	g := graph.NewGraph()

	res1 := resourcetest.Instance("inst_1").Prop("ID", "inst_1").Prop("Name", "instance 1").Build()
	res2 := resourcetest.Instance("inst_2").Prop("ID", "inst_2").Build()

	if err := g.AddResource(res1, res2); err != nil {
		t.Fatal(err)
	}

	r, err := g.GetResource("instance", "inst_1")
	if err != nil {
		t.Fatal(err)
	}

	columns := []string{"ID", "Name"}

	displayer, _ := BuildOptions(
		WithColumns(columns),
		WithFormat("table"),
	).SetSource(r).Build()

	expected := `| PROPERTY ▲ |   VALUE    |
|------------|------------|
| ID         | inst_1     |
| Name       | instance 1 |
`
	var w bytes.Buffer
	if err := displayer.Print(&w); err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

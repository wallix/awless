package display

import (
	"bytes"
	"testing"

	"github.com/wallix/awless/graph"
)

func TestResourceDisplay(t *testing.T) {
	g := graph.NewGraph()

	res1 := graph.InitResource("inst_1", graph.Instance)
	res1.Properties = map[string]interface{}{
		"Id":     "inst_1",
		"Name":   "instance 1",
		"Prop 1": "prop 1",
		"Prop 2": "prop 2",
	}
	res2 := graph.InitResource("inst_2", graph.Instance)

	g.AddResource(res1, res2)

	r, err := g.GetResource(graph.Instance, "inst_1")
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

package display

import (
	"bytes"
	"testing"

	"github.com/wallix/awless/graph"
)

func TestResourceDisplay(t *testing.T) {
	t0 := parseTriple(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`)
	t1 := parseTriple(`/instance<inst_1>	"has_type"@[]	"/instance"^^type:text`)
	t2 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text`)
	t3 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Name","Value":"instance 1"}"^^type:text`)
	t4 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Prop 1","Value":"prop 1"}"^^type:text`)
	t5 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Prop 2","Value":"prop 2"}"^^type:text`)
	t6 := parseTriple(`/instance<inst_2>	"has_type"@[]	"/instance"^^type:text`)
	t7 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_1>`)
	t8 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)
	g := graph.NewGraph()
	g.Add(t0, t1, t2, t3, t4, t5, t6, t7, t8)

	r := graph.InitResource("inst_1", graph.Instance)
	r.UnmarshalFromGraph(g)

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
	err := displayer.Print(&w)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := w.String(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

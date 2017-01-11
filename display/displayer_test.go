package display

import (
	"testing"

	"github.com/fatih/color"
	"github.com/wallix/awless/rdf"
)

func TestTabularDisplays(t *testing.T) {
	instances := []byte(`/region<eu-west-1> "has_type"@[] "/region"^^type:text
  /instance<inst_1>  "has_type"@[] "/instance"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Id","Value":"inst_1"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Name","Value":"redis"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"Type","Value":"t2.micro"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"PublicIp","Value":"1.2.3.4"}"^^type:text
  /instance<inst_1>  "property"@[] "{"Key":"State","Value":"running"}"^^type:text

  /instance<inst_2>  "has_type"@[] "/instance"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Name","Value":"django"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"Type","Value":"t2.medium"}"^^type:text
  /instance<inst_2>  "property"@[] "{"Key":"State","Value":"stopped"}"^^type:text
	
	
  /instance<inst_3>  "has_type"@[] "/instance"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Id","Value":"inst_3"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Name","Value":"apache"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"Type","Value":"t2.xlarge"}"^^type:text
  /instance<inst_3>  "property"@[] "{"Key":"State","Value":"running"}"^^type:text

  /region<eu-west-1>  "parent_of"@[] /instance<inst_1>
  /region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)

	g := rdf.NewGraph()
	g.Unmarshal(instances)

	headers := []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		StringColumnDefinition{Prop: "State"},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}
	displayer := BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance,
		Format:  "csv",
	})
	displayer.SetGraph(g)

	expected := `Id, Name, State, Type, Public IP
inst_1, redis, running, t2.micro, 1.2.3.4
inst_2, django, stopped, t2.medium, 
inst_3, apache, running, t2.xlarge, `

	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance,
		Format:  "csv",
		SortBy:  []string{"Name"},
	})
	displayer.SetGraph(g)

	expected = `Id, Name, State, Type, Public IP
inst_3, apache, running, t2.xlarge, 
inst_2, django, stopped, t2.medium, 
inst_1, redis, running, t2.micro, 1.2.3.4`

	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
		ColoredValueColumnDefinition{
			StringColumnDefinition: StringColumnDefinition{Prop: "State"},
			ColoredValues:          map[string]color.Attribute{"running": color.FgGreen},
		},
		StringColumnDefinition{Prop: "Type"},
		StringColumnDefinition{Prop: "PublicIp", Friendly: "Public IP"},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|  ID ▲  |  NAME  |  STATE  |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
| inst_3 | apache | running | t2.xlarge |           |
+--------+--------+---------+-----------+-----------+
`
	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "id"},
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_3 | apache | running | t2.xlarge |           |
| inst_2 | django | stopped | t2.medium |           |
+--------+--------+---------+-----------+-----------+
`
	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "table",
		SortBy: []string{"state", "name"},
	})
	displayer.SetGraph(g)
	expected = `+--------+--------+---------+-----------+-----------+
|   ID   |  NAME  | STATE ▲ |   TYPE    | PUBLIC IP |
+--------+--------+---------+-----------+-----------+
| inst_3 | apache | running | t2.xlarge |           |
| inst_1 | redis  | running | t2.micro  | 1.2.3.4   |
| inst_2 | django | stopped | t2.medium |           |
+--------+--------+---------+-----------+-----------+
`
	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}

	headers = []ColumnDefinition{
		StringColumnDefinition{Prop: "Id"},
		StringColumnDefinition{Prop: "Name"},
	}
	displayer = BuildGraphDisplayer(headers, Options{
		RdfType: rdf.Instance, Format: "porcelain",
	})
	displayer.SetGraph(g)
	expected = `inst_1
redis
inst_2
django
inst_3
apache`
	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

func TestCompareInterface(t *testing.T) {
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(4)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(1)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1), interface{}(-3)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("bbc")), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}("abc"), interface{}("aac")), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.3)), true; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
	if got, want := valueLowerOrEqual(interface{}(1.2), interface{}(1.1)), false; got != want {
		t.Fatalf("got %t want %t", got, want)
	}
}

package display

import (
	"testing"

	"github.com/wallix/awless/rdf"
)

func TestDisplayAsCSV(t *testing.T) {
	displayer := BuildDisplayer(Options{
		RdfType: rdf.INSTANCE, Format: "csv",
	})

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

  /region<eu-west-1>  "parent_of"@[] /instance<inst_1>
  /region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)

	g := rdf.NewGraph()
	g.Unmarshal(instances)

	displayer.SetGraph(g)
	displayer.SetHeaders([]Header{
		StringHeader{Prop: "Id"},
		StringHeader{Prop: "Name"},
		StringHeader{Prop: "State"},
		StringHeader{Prop: "Type"},
		StringHeader{Prop: "PublicIp", Friendly: "Public IP"},
	})

	expected := `Id, Name, State, Type, Public IP
inst_1, redis, running, t2.micro, 1.2.3.4
inst_2, django, stopped, t2.medium, `

	if got, want := displayer.Print(), expected; got != want {
		t.Fatalf("got \n%s\n\nwant\n\n%s\n", got, want)
	}
}

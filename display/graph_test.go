package display

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

func TestPrintOneResource(t *testing.T) {
	t0 := parseTriple(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`)
	t1 := parseTriple(`/instance<inst_1>	"has_type"@[]	"/instance"^^type:text`)
	t2 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text`)
	t3 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Tags","Value":[{"Key":"Name","Value":"instance 1"}]}"^^type:text`)
	t4 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Prop 1","Value":"prop 1"}"^^type:text`)
	t5 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Prop 2","Value":"prop 2"}"^^type:text`)
	t6 := parseTriple(`/instance<inst_2>	"has_type"@[]	"/instance"^^type:text`)
	t7 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_1>`)
	t8 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)
	g := rdf.NewGraphFromTriples([]*triple.Triple{t0, t1, t2, t3, t4, t5, t6, t7, t8})

	expect := `Instance 'instance 1'
+------------+------------+
| PROPERTY ▲ |   VALUE    |
+------------+------------+
| Id         | inst_1     |
| Name       | instance 1 |
| Prop 1     | prop 1     |
| Prop 2     | prop 2     |
+------------+------------+
`
	var res bytes.Buffer
	err := OneResourceOfGraph(&res, g, "instance", "inst_1", PropertiesDisplayer.Services[aws.InfraServiceName].Resources["instance"])
	if err != nil {
		t.Fatal(err)
	}

	if got, want := res.String(), expect; got != want {
		t.Fatalf("got:\n%s\nwant:\n%s\n", got, want)
	}
}

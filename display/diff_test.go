package display

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

func TestDisplayCommit(t *testing.T) {
	t0 := parseTriple(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`)
	t1 := parseTriple(`/instance<inst_1>	"has_type"@[]	"/instance"^^type:text`)
	t2 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text`)
	t3 := parseTriple(`/instance<inst_2>	"has_type"@[]	"/instance"^^type:text`)
	t4 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_1>`)
	t5 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)
	t6 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /vpc<vpc_1>`)
	t7 := parseTriple(`/vpc<vpc_1>  "parent_of"@[] /instance<inst_1>`)
	graph := rdf.NewGraphFromTriples([]*triple.Triple{t0, t1, t2, t3, t4, t5, t6, t7})
	diff := rdf.NewEmptyDiffFromGraph(graph)
	diff.AddDeleted(t2, rdf.ParentOfPredicate)
	diff.AddDeleted(t3, rdf.ParentOfPredicate)
	diff.AddDeleted(t5, rdf.ParentOfPredicate)
	t8 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"new_id"}"^^type:text`)
	t9 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_3>`)
	t10 := parseTriple(`/instance<inst_3>	"has_type"@[]	"/instance"^^type:text`)
	t11 := parseTriple(`/instance<inst_3>	"property"@[]	"{"Key":"Id","Value":"inst_3"}"^^type:text`)
	t12 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_4>`)
	t13 := parseTriple(`/instance<inst_4>	"has_type"@[]	"/instance"^^type:text`)
	diff.AddInserted(t8, rdf.ParentOfPredicate)
	diff.AddInserted(t9, rdf.ParentOfPredicate)
	diff.AddInserted(t10, rdf.ParentOfPredicate)
	diff.AddInserted(t11, rdf.ParentOfPredicate)
	diff.AddInserted(t12, rdf.ParentOfPredicate)
	diff.AddInserted(t13, rdf.ParentOfPredicate)

	rootNode, err := node.NewNodeFromStrings("/region", "eu-west-1")
	if err != nil {
		t.Fatal(err)
	}
	table, err := tableFromDiff(diff, rootNode)
	if err != nil {
		t.Fatal(err)
	}
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+-----------+---------+----------+--------+
|  TYPE ▲   | NAME/ID | PROPERTY | VALUE  |
+-----------+---------+----------+--------+
| /instance | inst_1  | Id       | inst_1 |
+           +         +          +--------+
|           |         |          | new_id |
+           +---------+----------+--------+
|           | inst_2  |          |        |
+           +---------+----------+--------+
|           | inst_3  | Id       | inst_3 |
+           +---------+----------+--------+
|           | inst_4  |          |        |
+-----------+---------+----------+--------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}

package ast

import (
	"reflect"
	"testing"
)

func TestCloneAST(t *testing.T) {
	tree := &AST{}

	tree.Statements = append(tree.Statements, &Statement{Node: &DeclarationNode{
		Left: &IdentifierNode{Ident: "myvar"},
		Right: &ExpressionNode{
			Action: "create", Entity: "vpc",
			Refs:    map[string]string{"myname": "name"},
			Params:  map[string]interface{}{"count": 1},
			Aliases: map[string]string{"subnet": "my-subnet"},
			Holes:   make(map[string]string),
		}}}, &Statement{Node: &DeclarationNode{
		Left: &IdentifierNode{Ident: "myothervar"},
		Right: &ExpressionNode{
			Action: "create", Entity: "subnet",
			Refs:    make(map[string]string),
			Params:  make(map[string]interface{}),
			Aliases: make(map[string]string),
			Holes:   map[string]string{"vpc": "myvar"},
		}}}, &Statement{Node: &ExpressionNode{
		Action: "create", Entity: "instance",
		Refs:    make(map[string]string),
		Params:  make(map[string]interface{}),
		Aliases: make(map[string]string),
		Holes:   map[string]string{"subnet": "myothervar"},
	}},
	)

	clone := tree.Clone()

	if got, want := clone, tree; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v", got, want)
	}

	clone.Statements[0].Node.(*DeclarationNode).Right.Params["new"] = "trump"

	if got, want := clone.Statements, tree.Statements; reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %s\n\nwant %s", got, want)
	}
}

func TestGetStatementAttributes(t *testing.T) {
	params := map[string]interface{}{"count": 1}
	st := &Statement{Node: &DeclarationNode{
		Left: &IdentifierNode{},
		Right: &ExpressionNode{
			Action: "create", Entity: "vpc", Params: params,
		}}}
	if got, want := st.Action(), "create"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := st.Entity(), "vpc"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := st.Params(), params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}

	st = &Statement{Node: &ExpressionNode{
		Action: "delete", Entity: "subnet", Params: params,
	}}
	if got, want := st.Action(), "delete"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := st.Entity(), "subnet"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := st.Params(), params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

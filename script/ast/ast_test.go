package ast

import (
	"reflect"
	"testing"
)

func TestCloneAST(t *testing.T) {
	tree := &AST{}

	tree.Statements = append(tree.Statements, &DeclarationNode{
		Left: &IdentifierNode{Ident: "myvar"},
		Right: &ExpressionNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 1},
			Holes:  make(map[string]string),
		}}, &DeclarationNode{
		Left: &IdentifierNode{Ident: "myothervar"},
		Right: &ExpressionNode{
			Action: "create", Entity: "subnet",
			Params: make(map[string]interface{}),
			Holes:  map[string]string{"vpc": "myvar"},
		}}, &ExpressionNode{
		Action: "create", Entity: "instance",
		Holes:  map[string]string{"subnet": "myothervar"},
		Params: make(map[string]interface{}),
	},
	)

	clone := tree.Clone()

	if got, want := clone, tree; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v", got, want)
	}

	clone.Statements[0].(*DeclarationNode).Right.Params["new"] = "trump"

	if got, want := clone.Statements, tree.Statements; reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %s\n\nwant %s", got, want)
	}
}

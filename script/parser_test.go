package script

import (
	"reflect"
	"testing"

	"github.com/wallix/awless/script/ast"
)

func TestScriptParsing(t *testing.T) {
	tcases := []struct {
		input    string
		verifyFn func(s *Script)
	}{
		{
			input: `
			myvpc  =   create   vpc  cidr=10.0.0.0/24 num=3
mysubnet = delete subnet vpc={ one } name=bill
create instance count=1 instance.type=t2.micro image=ami-9398d3e0 ip=127.0.0.1
		`,

			verifyFn: func(s *Script) {
				assertStatementsCount(t, s, 3)
				assertDeclarationNode(t, 0, s.Statements, "myvpc", "create", "vpc",
					map[string]interface{}{"cidr": "10.0.0.0/24", "num": 3}, map[string]string{},
				)
				assertDeclarationNode(t, 1, s.Statements, "mysubnet", "delete", "subnet",
					map[string]interface{}{"name": "bill"},
					map[string]string{"vpc": "one"},
				)
				assertExpressionNode(t, 2, s.Statements, "create", "instance",
					map[string]interface{}{"count": 1, "instance.type": "t2.micro", "ip": "127.0.0.1", "image": "ami-9398d3e0"}, map[string]string{},
				)
			},
		},

		{
			input: `create vpc`,
			verifyFn: func(s *Script) {
				assertStatementsCount(t, s, 1)
				assertExpressionNode(t, 0, s.Statements, "create", "vpc", nil, nil)
			},
		},
	}

	for _, tcase := range tcases {
		scrpt, err := Parse(tcase.input)
		if err != nil {
			t.Fatal(err)
		}

		tcase.verifyFn(scrpt)
	}
}

func assertDeclarationNode(t *testing.T, index int, sts []ast.Node, expIdent, expAction, expEntity string, params map[string]interface{}, holes map[string]string) {
	n := sts[index]

	decl, ok := n.(*ast.DeclarationNode)
	if !ok {
		t.Fatalf("statement %d: unexpected type %T", index, n)
	}

	assertIdentifierNode(t, index, decl.Left, expIdent)
	verifyExpressionNode(t, index, decl.Right, expAction, expEntity, params, holes)
}

func assertStatementsCount(t *testing.T, s *Script, count int) {
	if got, want := len(s.Statements), count; got != want {
		t.Fatalf("expected %d statements got %d\n%#v", want, got, s.Statements)
	}
}

func assertIdentifierNode(t *testing.T, index int, n *ast.IdentifierNode, expected string) {
	if got, want := n.Ident, expected; got != want {
		t.Fatalf("statement %d: identifier: got '%s' want '%s'", got, want)
	}
}

func assertExpressionNode(t *testing.T, index int, sts []ast.Node, expAction, expEntity string, params map[string]interface{}, holes map[string]string) {
	n := sts[index]
	verifyExpressionNode(t, index, n, expAction, expEntity, params, holes)
}

func verifyExpressionNode(t *testing.T, index int, n ast.Node, expAction, expEntity string, params map[string]interface{}, holes map[string]string) {
	expr, ok := n.(*ast.ExpressionNode)
	if !ok {
		t.Fatalf("statement %d: unexpected type %T", index, n)
	}

	if got, want := expr.Action, expAction; got != want {
		t.Fatalf("statement %d: action: got '%s' want '%s'", index, got, want)
	}
	if got, want := expr.Entity, expEntity; got != want {
		t.Fatalf("statement %d: entity: got '%s' want '%s'", index, got, want)
	}

	if got, want := expr.Params, params; !reflect.DeepEqual(got, want) {
		t.Fatalf("statement %d: params: got %#v, want %#v", index, got, want)
	}

	if got, want := expr.Holes, holes; !reflect.DeepEqual(got, want) {
		t.Fatalf("statement %d: holes: got %#v, want %#v", index, got, want)
	}
}

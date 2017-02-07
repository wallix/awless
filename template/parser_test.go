package template

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/wallix/awless/template/ast"
)

func TestTemplateParsing(t *testing.T) {

	t.Run("Parse special characters", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(tpl *Template) error
		}{
			{
				input: "attach policy arn=arn:aws:iam::aws:policy/AmazonS3FullAccess",
				verifyFn: func(tpl *Template) error {
					if err := isExpressionNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Allow and ignore comments", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(tpl *Template) error
		}{
			{
				input: "create vpc\n#my comment\ncreate subnet",
				verifyFn: func(tpl *Template) error {
					if got, want := len(tpl.Statements), 2; got != want {
						t.Fatalf("got %d, want %d", got, want)
					}
					if err := isExpressionNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					if err := isExpressionNode(tpl.Statements[1].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
			{
				input: "create vpc \n//my comment\ncreate subnet",
				verifyFn: func(tpl *Template) error {
					if got, want := len(tpl.Statements), 2; got != want {
						t.Fatalf("got %d, want %d", got, want)
					}
					if err := isExpressionNode(tpl.Statements[0].Node); err != nil {
						t.Fatal(err)
					}
					if err := isExpressionNode(tpl.Statements[1].Node); err != nil {
						t.Fatal(err)
					}
					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Onliner statement", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(n ast.Node) error
		}{
			{
				input:    `create vpc`,
				verifyFn: func(n ast.Node) error { return assertParams(n, nil) },
			},
			{
				input:    `create vpc`,
				verifyFn: func(n ast.Node) error { return isExpressionNode(n) },
			},
			{
				input:    `mysubnet = create subnet`,
				verifyFn: func(n ast.Node) error { return isDeclarationNode(n) },
			},
			{
				input: `create vpc cidr=10.0.0.0/24 num=3 ip=127.0.0.1 name=bousin`,
				verifyFn: func(n ast.Node) error {
					return assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/24", "num": 3, "ip": "127.0.0.1", "name": "bousin"})
				},
			},
			{
				input: `create subnet vpc=$myvpc`,
				verifyFn: func(n ast.Node) error {
					return assertRefs(n, map[string]string{"vpc": "myvpc"})
				},
			},
			{
				input: `create instance subnet=@my-subnet`,
				verifyFn: func(n ast.Node) error {
					return assertAliases(n, map[string]string{"subnet": "my-subnet"})
				},
			},
			{
				input: `delete vpc id={my-vpc-id}`,
				verifyFn: func(n ast.Node) error {
					return assertHoles(n, map[string]string{"id": "my-vpc-id"})
				},
			},
			{
				input: `create securitygroup port=20-80`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"port": "20-80"}); err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `myinstance = create instance type={instance.type} cidr=10.0.0.0/25 subnet=@default-subnet vpc=$myvpc`,
				verifyFn: func(n ast.Node) error {
					if err := assertParams(n, map[string]interface{}{"cidr": "10.0.0.0/25"}); err != nil {
						return err
					}
					if err := assertAliases(n, map[string]string{"subnet": "default-subnet"}); err != nil {
						return err
					}
					if err := assertHoles(n, map[string]string{"type": "instance.type"}); err != nil {
						return err
					}
					if err := assertRefs(n, map[string]string{"vpc": "myvpc"}); err != nil {
						return err
					}

					return nil
				},
			},
		}

		for _, tcase := range tcases {
			node, err := ParseStatement(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(node); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})

	t.Run("Multiline parsing", func(t *testing.T) {
		tcases := []struct {
			input    string
			verifyFn func(s *Template) error
		}{
			{
				input: `create vpc
create subnet`,
				verifyFn: func(s *Template) error {
					err := assertExpressionNode(s.Statements[0].Node, "create", "vpc",
						nil,
						nil,
						nil,
						nil,
					)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				input: `
			myvpc  =   create   vpc  cidr=10.0.0.0/24 num=3
mysubnet = delete subnet vpc=$myvpc name={ the_name } cidr=10.0.0.0/25
create instance count=1 instance.type=t2.micro subnet=$mysubnet image=ami-9398d3e0 ip=127.0.0.1
		`,

				verifyFn: func(s *Template) error {
					err := assertDeclarationNode(s.Statements[0].Node, "myvpc", "create", "vpc",
						map[string]string{},
						map[string]interface{}{"cidr": "10.0.0.0/24", "num": 3},
						map[string]string{},
						map[string]string{},
					)
					if err != nil {
						return err
					}

					err = assertDeclarationNode(s.Statements[1].Node, "mysubnet", "delete", "subnet",
						map[string]string{"vpc": "myvpc"},
						map[string]interface{}{"cidr": "10.0.0.0/25"},
						map[string]string{"name": "the_name"},
						map[string]string{},
					)
					if err != nil {
						return err
					}

					err = assertExpressionNode(s.Statements[2].Node, "create", "instance",
						map[string]string{"subnet": "mysubnet"},
						map[string]interface{}{"count": 1, "instance.type": "t2.micro", "ip": "127.0.0.1", "image": "ami-9398d3e0"},
						map[string]string{},
						map[string]string{},
					)

					return err
				},
			},
		}

		for _, tcase := range tcases {
			templ, err := Parse(tcase.input)
			if err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}

			if err := tcase.verifyFn(templ); err != nil {
				t.Fatalf("\ninput: [%s]\nError: %s\n", tcase.input, err)
			}
		}
	})
}

func assertParams(n ast.Node, expected map[string]interface{}) error {
	compare := func(got, want map[string]interface{}) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("params: got %#v, want %#v", got, want)
		}
		return nil
	}

	switch n.(type) {
	case *ast.ExpressionNode:
		params := n.(*ast.ExpressionNode).Params
		return compare(params, expected)
	case *ast.DeclarationNode:
		params := n.(*ast.DeclarationNode).Right.Params
		return compare(params, expected)
	default:
		return errors.New("unexpected node type")
	}

	return nil
}

func assertRefs(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("refs: got %#v, want %#v", got, want)
		}
		return nil
	}

	switch n.(type) {
	case *ast.ExpressionNode:
		params := n.(*ast.ExpressionNode).Refs
		return compare(params, expected)
	case *ast.DeclarationNode:
		params := n.(*ast.DeclarationNode).Right.Refs
		return compare(params, expected)
	default:
		return errors.New("unexpected node type")
	}

	return nil
}

func assertAliases(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("aliases: got %#v, want %#v", got, want)
		}
		return nil
	}

	switch n.(type) {
	case *ast.ExpressionNode:
		params := n.(*ast.ExpressionNode).Aliases
		return compare(params, expected)
	case *ast.DeclarationNode:
		params := n.(*ast.DeclarationNode).Right.Aliases
		return compare(params, expected)
	default:
		return errors.New("unexpected node type")
	}

	return nil
}

func assertHoles(n ast.Node, expected map[string]string) error {
	compare := func(got, want map[string]string) error {
		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("holes: got %#v, want %#v", got, want)
		}
		return nil
	}

	switch n.(type) {
	case *ast.ExpressionNode:
		params := n.(*ast.ExpressionNode).Holes
		return compare(params, expected)
	case *ast.DeclarationNode:
		params := n.(*ast.DeclarationNode).Right.Holes
		return compare(params, expected)
	default:
		return errors.New("unexpected node type")
	}

	return nil
}

func assertDeclarationNode(n ast.Node, expIdent, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	if err := isDeclarationNode(n); err != nil {
		return err
	}

	decl := n.(*ast.DeclarationNode)

	if err := assertIdentifierNode(decl.Left, expIdent); err != nil {
		return err
	}

	if err := verifyExpressionNode(decl.Right, expAction, expEntity, refs, params, holes, aliases); err != nil {
		return err
	}

	return nil
}

func assertIdentifierNode(n *ast.IdentifierNode, expected string) error {
	if got, want := n.Ident, expected; got != want {
		return fmt.Errorf("identifier: got '%s' want '%s'", got, want)
	}
	return nil
}

func assertExpressionNode(n ast.Node, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	return verifyExpressionNode(n, expAction, expEntity, refs, params, holes, aliases)
}

func verifyExpressionNode(n ast.Node, expAction, expEntity string, refs map[string]string, params map[string]interface{}, holes, aliases map[string]string) error {
	if err := isExpressionNode(n); err != nil {
		return err
	}

	expr := n.(*ast.ExpressionNode)

	if got, want := expr.Action, expAction; got != want {
		return fmt.Errorf("action: got '%s' want '%s'", got, want)
	}
	if got, want := expr.Entity, expEntity; got != want {
		return fmt.Errorf("entity: got '%s' want '%s'", got, want)
	}

	if got, want := expr.Params, params; !reflect.DeepEqual(got, want) {
		return fmt.Errorf("params: got %#v, want %#v", got, want)
	}

	if got, want := expr.Refs, refs; !reflect.DeepEqual(got, want) {
		return fmt.Errorf("refs: got %#v, want %#v", got, want)
	}

	if got, want := expr.Holes, holes; !reflect.DeepEqual(got, want) {
		return fmt.Errorf("holes: got %#v, want %#v", got, want)
	}

	if got, want := expr.Aliases, aliases; !reflect.DeepEqual(got, want) {
		return fmt.Errorf("aliases: got %#v, want %#v", got, want)
	}

	return nil
}

func isExpressionNode(n ast.Node) error {
	switch n.(type) {
	case *ast.ExpressionNode:
	default:
		return errors.New("expected expression node")
	}
	return nil
}

func isDeclarationNode(n ast.Node) error {
	switch n.(type) {
	case *ast.DeclarationNode:
	default:
		return errors.New("expected declaration node")
	}
	return nil
}

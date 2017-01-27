package template

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/oklog/ulid"
	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver"
)

type stubDriver struct{}

func (d *stubDriver) Lookup(lookups ...string) driver.DriverFn {
	return func(map[string]interface{}) (interface{}, error) { return nil, nil }
}
func (d *stubDriver) SetLogger(*log.Logger) {}
func (d *stubDriver) SetDryRun(bool)        {}

type errorDriver struct {
	err error
}

func (d *errorDriver) Lookup(lookups ...string) driver.DriverFn {
	return func(map[string]interface{}) (interface{}, error) { return nil, d.err }
}
func (d *errorDriver) SetLogger(*log.Logger) {}
func (d *errorDriver) SetDryRun(bool)        {}

func TestRunDriverOutputOperations(t *testing.T) {
	anErr := errors.New("my error message")

	tcases := []struct {
		input  string
		driver driver.Driver
		expect []*Operation
	}{
		{
			input:  "create vpc cidr=10.0.0.0/25\ndelete subnet id=sub-5f4g3hj",
			driver: &stubDriver{},
			expect: []*Operation{
				&Operation{Line: "create vpc cidr=10.0.0.0/25"},
				&Operation{Line: "delete subnet id=sub-5f4g3hj"},
			},
		},
		{
			input:  "create vpc cidr=10.0.0.0/25",
			driver: &errorDriver{anErr},
			expect: []*Operation{
				&Operation{Line: "create vpc cidr=10.0.0.0/25", Err: anErr},
			},
		},
	}

	for _, tcase := range tcases {
		templ, err := Parse(tcase.input)
		if err != nil {
			t.Fatal(err)
		}
		_, ops, _ := templ.Run(tcase.driver)

		for i, op := range ops {
			if got, want := op.Line, tcase.expect[i].Line; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot '%q'\n\twant '%q'", tcase.input, got, want)
			}
			if got, want := op.Output, tcase.expect[i].Output; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot %s\n\twant %s", tcase.input, got, want)
			}
			if got, want := op.Err, tcase.expect[i].Err; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot %v\n\twant %v", tcase.input, got, want)
			}

			if _, err := ulid.Parse(op.ID); err != nil {
				t.Fatalf("\ninput: '%s'\n cannot parse ulid %s", tcase.input, op.ID)
			}
		}
	}
}

func TestRunDriverOnTemplate(t *testing.T) {
	t.Run("Driver run TWICE multiline statement", func(t *testing.T) {
		s := &Template{&ast.AST{}}

		s.Statements = append(s.Statements, &ast.DeclarationNode{
			Left: &ast.IdentifierNode{Ident: "createdvpc"},
			Right: &ast.ExpressionNode{
				Action: "create", Entity: "vpc",
				Params: map[string]interface{}{"count": 1},
			}}, &ast.DeclarationNode{
			Left: &ast.IdentifierNode{Ident: "createdsubnet"},
			Right: &ast.ExpressionNode{
				Action: "create", Entity: "subnet",
				Refs: map[string]string{"vpc": "createdvpc"},
			}}, &ast.ExpressionNode{
			Action: "create", Entity: "instance",
			Refs: map[string]string{"subnet": "createdsubnet"},
		},
		)

		mDriver := &mockDriver{prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}, {
			action: "create", entity: "subnet",
			expectedParams: map[string]interface{}{"vpc": "mynewvpc"},
		}, {
			action: "create", entity: "instance",
			expectedParams: map[string]interface{}{"subnet": "mynewsubnet"},
		},
		},
		}

		if _, _, err := s.Run(mDriver); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}

		mDriver = &mockDriver{prefix: "myother", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}, {
			action: "create", entity: "subnet",
			expectedParams: map[string]interface{}{"vpc": "myothervpc"},
		}, {
			action: "create", entity: "instance",
			expectedParams: map[string]interface{}{"subnet": "myothersubnet"},
		},
		},
		}

		if _, _, err := s.Run(mDriver); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver visit expression nodes", func(t *testing.T) {
		s := &Template{&ast.AST{}}

		n := &ast.ExpressionNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 1},
		}
		s.Statements = append(s.Statements, n)

		mDriver := &mockDriver{prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		if _, _, err := s.Run(mDriver); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver visit declaration nodes", func(t *testing.T) {
		s := &Template{&ast.AST{}}

		decl := &ast.DeclarationNode{
			Left: &ast.IdentifierNode{Ident: "myvar"},
			Right: &ast.ExpressionNode{
				Action: "create", Entity: "vpc",
				Params: map[string]interface{}{"count": 1},
			},
		}
		s.Statements = append(s.Statements, decl)

		mDriver := &mockDriver{prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		executedTemplate, _, err := s.Run(mDriver)
		if err != nil {
			t.Fatal(err)
		}

		modifiedDecl := executedTemplate.Statements[0].(*ast.DeclarationNode)
		if got, want := modifiedDecl.Left.Val, "mynewvpc"; got != want {
			t.Fatalf("identifier: got %#v, want %#v", got, want)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAliases(t *testing.T) {
	tree := &ast.AST{}

	tree.Statements = append(tree.Statements, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Aliases: map[string]string{"1": "one"},
		}}, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Aliases: map[string]string{"2": "two", "3": "three"},
		}}, &ast.ExpressionNode{
		Aliases: map[string]string{"4": "four"},
	},
	)
	s := &Template{tree}
	expect := map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
		"4": "four",
	}
	if got, want := s.GetAliases(), expect; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestMergeParams(t *testing.T) {
	templ := &Template{&ast.AST{}}

	templ.Statements = append(templ.Statements, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 1},
		}}, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Action: "create", Entity: "subnet",
		}}, &ast.ExpressionNode{
		Action: "create", Entity: "instance",
		Params: map[string]interface{}{"type": "t1", "image": "image1"},
	})
	templ.MergeParams(map[string]interface{}{
		"vpc.count":       10,
		"subnet.cidr":     "10.0.0.0/24",
		"instance.image":  "image2",
		"instance.subnet": "mysubnet",
	})

	var expect []ast.Node
	expect = append(expect, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 10},
		}}, &ast.DeclarationNode{
		Left: &ast.IdentifierNode{},
		Right: &ast.ExpressionNode{
			Action: "create", Entity: "subnet",
			Params: map[string]interface{}{"cidr": "10.0.0.0/24"},
		}}, &ast.ExpressionNode{
		Action: "create", Entity: "instance",
		Params: map[string]interface{}{"type": "t1", "image": "image2", "subnet": "mysubnet"},
	})

	if got, want := templ.Statements, expect; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestResolveTemplate(t *testing.T) {
	t.Run("Holes Resolution", func(t *testing.T) {
		s := &Template{&ast.AST{}}

		expr := &ast.ExpressionNode{
			Holes: map[string]string{"name": "presidentName", "rank": "presidentRank"},
		}
		s.Statements = append(s.Statements, expr)

		decl := &ast.DeclarationNode{
			Right: &ast.ExpressionNode{
				Holes: map[string]string{"age": "presidentAge", "wife": "presidentWife"},
			},
		}
		s.Statements = append(s.Statements, decl)

		fills := map[string]interface{}{
			"presidentName": "trump",
			"presidentRank": 45,
			"presidentAge":  70,
			"presidentWife": "melania",
		}

		s.ResolveTemplate(fills)

		expected := map[string]interface{}{"name": "trump", "rank": 45}
		if got, want := expr.Params, expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := len(expr.Holes), 0; got != want {
			t.Fatalf("length of holes: got %d, want %d", got, want)
		}

		expected = map[string]interface{}{"age": 70, "wife": "melania"}
		if got, want := decl.Right.Params, expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := len(decl.Right.Holes), 0; got != want {
			t.Fatalf("length of holes: got %d, want %d", got, want)
		}
	})

	t.Run("Interactive holes resolution", func(t *testing.T) {
		s := &Template{&ast.AST{}}

		expr := &ast.ExpressionNode{
			Holes: map[string]string{"age": "age_of_president", "name": "name_of_president"},
		}
		s.Statements = append(s.Statements, expr)

		each := func(question string) interface{} {
			if question == "age_of_president" {
				return 70
			}
			if question == "name_of_president" {
				return "trump"
			}

			return nil
		}

		s.InteractiveResolveTemplate(each)

		expected := map[string]interface{}{"age": 70, "name": "trump"}
		if got, want := expr.Params, expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		if got, want := len(expr.Holes), 0; got != want {
			t.Fatalf("length of holes: got %d, want %d", got, want)
		}
	})
}

type expectation struct {
	lookupDone     bool
	action, entity string
	expectedParams map[string]interface{}
}

type mockDriver struct {
	expects []*expectation
	prefix  string
}

func (r *mockDriver) lookupsCalled() error {
	for _, expect := range r.expects {
		if expect.lookupDone == false {
			return fmt.Errorf("lookup for expectation %v not called", expect)
		}
	}

	return nil
}

func (r *mockDriver) Lookup(lookups ...string) driver.DriverFn {
	for _, expect := range r.expects {
		if lookups[0] == expect.action && lookups[1] == expect.entity {
			expect.lookupDone = true

			return func(params map[string]interface{}) (interface{}, error) {
				if got, want := expect.expectedParams, params; !reflect.DeepEqual(got, want) {
					return nil, fmt.Errorf("[%s %s] params mismatch: expected %v, got %v", expect.action, expect.entity, got, want)
				}
				return r.prefix + expect.entity, nil
			}
		}
	}

	return func(params map[string]interface{}) (interface{}, error) {
		return nil, errors.New("Unexpected lookup fallthrough")
	}
}

func (r *mockDriver) SetLogger(*log.Logger) {}
func (r *mockDriver) SetDryRun(bool)        {}

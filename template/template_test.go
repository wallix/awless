/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/oklog/ulid"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver"
)

type noopDriver struct{}

func (d *noopDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	return func(map[string]interface{}) (interface{}, error) { return nil, nil }, nil
}
func (d *noopDriver) SetLogger(*logger.Logger) {}
func (d *noopDriver) SetDryRun(bool)           {}

type errorDriver struct {
	err error
}

func (d *errorDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	return func(map[string]interface{}) (interface{}, error) { return nil, d.err }, nil
}
func (d *errorDriver) SetLogger(*logger.Logger) {}
func (d *errorDriver) SetDryRun(bool)           {}

func TestRunDriverReportsInStatement(t *testing.T) {
	anErr := errors.New("my error message")

	type line struct {
		expString string
		expErr    error
		expResult interface{}
	}

	tcases := []struct {
		input  string
		driver driver.Driver
		lines  []line
	}{
		{
			input:  "create vpc cidr=10.0.0.0/25\ndelete subnet id=sub-5f4g3hj",
			driver: &noopDriver{},
			lines: []line{
				{expString: "create vpc cidr=10.0.0.0/25"},
				{expString: "delete subnet id=sub-5f4g3hj"},
			},
		},
		{
			input:  "create vpc cidr=10.0.0.0/25",
			driver: &errorDriver{anErr},
			lines: []line{
				{expString: "create vpc cidr=10.0.0.0/25", expErr: anErr},
			},
		},
	}

	for _, tcase := range tcases {
		templ, err := Parse(tcase.input)
		if err != nil {
			t.Fatal(err)
		}
		ran, _ := templ.Run(tcase.driver)

		for i, cmd := range ran.CommandNodesIterator() {
			if got, want := cmd.String(), tcase.lines[i].expString; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot '%q'\n\twant '%q'", tcase.input, got, want)
			}
			if got, want := cmd.Result(), tcase.lines[i].expResult; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot %s\n\twant %s", tcase.input, got, want)
			}
			if got, want := cmd.Err(), tcase.lines[i].expErr; got != want {
				t.Fatalf("\ninput: '%s'\n\tgot %v\n\twant %v", tcase.input, got, want)
			}
		}
	}
}

func TestNewTemplateExecutionFromTemplate(t *testing.T) {
	temp, err := Parse("create vpc name=any\ncreate subnet ip=10.0.0.0\ndelete instance id=i-5d678\nstop instance id=i-5d678")
	if err != nil {
		t.Fatal(err)
	}

	if temp, err = temp.Run(&noopDriver{}); err != nil {
		t.Fatal(err)
	}

	cmdNodes := temp.CommandNodesIterator()

	cmdNodes[0].CmdResult = "vpc-123"
	cmdNodes[1].CmdResult = "sub-123"
	cmdNodes[2].CmdResult = struct{}{}
	cmdNodes[2].CmdErr = errors.New("cannot delete instance")

	executed := NewTemplateExecution(temp)

	if _, err := ulid.Parse(executed.ID); err != nil {
		t.Fatalf("parsing '%s': %s", executed.ID, err)
	}
	if got, want := len(executed.Executed), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := executed.Executed[0].Err, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[0].Line, "create vpc name=any"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[0].Result, "vpc-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[1].Err, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[1].Line, "create subnet ip=10.0.0.0"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[1].Result, "sub-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[2].Err, "cannot delete instance"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[2].Line, "delete instance id=i-5d678"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := executed.Executed[2].Result, ""; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestTemplateExecutionHasErrors(t *testing.T) {
	exec := &TemplateExecution{
		Executed: []*ExecutedStatement{
			{Line: "create vpc", Result: "vpc-56g4h", Err: ""},
			{Line: "create instance", Result: "", Err: "cannot create instance"},
		},
	}

	if got, want := exec.HasErrors(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	exec = &TemplateExecution{
		Executed: []*ExecutedStatement{
			{Line: "create vpc", Result: "vpc-56g4h", Err: ""},
			{Line: "create instance", Result: ""},
		},
	}

	if got, want := exec.HasErrors(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
}

func TestRevertTemplateExecution(t *testing.T) {
	exec := &TemplateExecution{
		Executed: []*ExecutedStatement{
			{Line: "attach policy arn=stuff user=mrT", Result: "", Err: ""},
			{Line: "create vpc", Result: "vpc-56g4h", Err: ""},
			{Line: "create subnet", Result: "sub-65bh4nj", Err: ""},
			{Line: "start instance id=i-54g3hj", Result: "i-54g3hj", Err: ""},
			{Line: "create tags", Result: "", Err: ""},
			{Line: "create instance", Result: "", Err: "cannot create instance"},
		},
	}

	tpl, err := exec.Revert()
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(tpl.Statements), 4; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	expr := tpl.Statements[0].Node.(*ast.CommandNode)
	if got, want := "stop", expr.Action; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := "instance", expr.Entity; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	expected := map[string]interface{}{"id": "i-54g3hj"}
	if got, want := expected, expr.Params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	expr = tpl.Statements[1].Node.(*ast.CommandNode)
	if got, want := "delete", expr.Action; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := "subnet", expr.Entity; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	expected = map[string]interface{}{"id": "sub-65bh4nj"}
	if got, want := expected, expr.Params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	expr = tpl.Statements[2].Node.(*ast.CommandNode)
	if got, want := "delete", expr.Action; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := "vpc", expr.Entity; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	expected = map[string]interface{}{"id": "vpc-56g4h"}
	if got, want := expected, expr.Params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	expr = tpl.Statements[3].Node.(*ast.CommandNode)
	if got, want := "detach", expr.Action; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := "policy", expr.Entity; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	expected = map[string]interface{}{"arn": "stuff", "user": "mrT"}
	if got, want := expected, expr.Params; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestExecutedStatementIsRevertible(t *testing.T) {
	tcases := []struct {
		line, result, err string
		revertible        bool
	}{
		{line: "update vpc", result: "any", revertible: false},
		{line: "delete vpc", result: "any", revertible: false},
		{line: "create vpc", result: "any", err: "any", revertible: false},
		{line: "create vpc", revertible: false},
		{line: "start instance", revertible: false},
		{line: "create vpc", result: "any", revertible: true},
		{line: "stop instance", result: "any", revertible: true},
		{line: "attach policy", result: "", revertible: true},
		{line: "detach policy", result: "", revertible: true},
	}

	for _, tc := range tcases {
		ex := &ExecutedStatement{Line: tc.line, Result: tc.result, Err: tc.err}
		if tc.revertible != ex.IsRevertible() {
			t.Fatalf("expected %#v to have revertible=%t", ex, tc.revertible)
		}
	}
}

func TestRunDriverOnTemplate(t *testing.T) {
	t.Run("Driver run TWICE multiline statement", func(t *testing.T) {
		s := &Template{AST: &ast.AST{}}

		s.Statements = append(s.Statements, &ast.Statement{Node: &ast.DeclarationNode{
			Ident: "createdvpc",
			Expr: &ast.CommandNode{
				Action: "create", Entity: "vpc",
				Params: map[string]interface{}{"count": 1},
			}}}, &ast.Statement{Node: &ast.DeclarationNode{
			Ident: "createdsubnet",
			Expr: &ast.CommandNode{
				Action: "create", Entity: "subnet",
				Refs: map[string]string{"vpc": "createdvpc"},
			}}}, &ast.Statement{Node: &ast.CommandNode{
			Action: "create", Entity: "instance",
			Refs: map[string]string{"subnet": "createdsubnet"},
		}},
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

		if _, err := s.Run(mDriver); err != nil {
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

		if _, err := s.Run(mDriver); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver visit expression nodes", func(t *testing.T) {
		s := &Template{AST: &ast.AST{}}

		n := &ast.Statement{Node: &ast.CommandNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 1},
		}}
		s.Statements = append(s.Statements, n)

		mDriver := &mockDriver{prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		if _, err := s.Run(mDriver); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver visit declaration nodes", func(t *testing.T) {
		s := &Template{AST: &ast.AST{}}

		decl := &ast.Statement{Node: &ast.DeclarationNode{
			Ident: "myvar",
			Expr: &ast.CommandNode{
				Action: "create", Entity: "vpc",
				Params: map[string]interface{}{"count": 1},
			},
		}}
		s.Statements = append(s.Statements, decl)

		mDriver := &mockDriver{prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		executedTemplate, err := s.Run(mDriver)
		if err != nil {
			t.Fatal(err)
		}

		modifiedDecl := executedTemplate.Statements[0].Node.(*ast.DeclarationNode)
		if got, want := modifiedDecl.Expr.Result(), "mynewvpc"; got != want {
			t.Fatalf("identifier: got %#v, want %#v", got, want)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetNormalisedAliases(t *testing.T) {
	tree := &ast.AST{}

	tree.Statements = append(tree.Statements, &ast.Statement{Node: &ast.DeclarationNode{
		Expr: &ast.CommandNode{
			Aliases: map[string]string{"1.1": "one"},
		}}}, &ast.Statement{Node: &ast.DeclarationNode{
		Expr: &ast.CommandNode{
			Entity:  "bear",
			Aliases: map[string]string{"2": "two", "3.3": "three"},
		}}}, &ast.Statement{Node: &ast.CommandNode{
		Entity: "shark", Aliases: map[string]string{"4": "four"},
	}},
	)
	s := &Template{AST: tree}
	expect := map[string]string{
		"1.1":     "one",
		"bear.2":  "two",
		"3.3":     "three",
		"shark.4": "four",
	}
	if got, want := s.GetNormalizedAliases(), expect; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestMergeParams(t *testing.T) {
	templ := &Template{AST: &ast.AST{}}

	templ.Statements = append(templ.Statements, &ast.Statement{Node: &ast.DeclarationNode{
		Expr: &ast.CommandNode{
			Action: "create", Entity: "vpc",
			Params: map[string]interface{}{"count": 1},
		}}}, &ast.Statement{Node: &ast.DeclarationNode{
		Expr: &ast.CommandNode{
			Action: "create", Entity: "subnet",
		}}}, &ast.Statement{Node: &ast.CommandNode{
		Action: "create", Entity: "instance",
		Params: map[string]interface{}{"type": "t1", "image": "image1"},
	}})
	templ.MergeParams(map[string]interface{}{
		"vpc.count":       10,
		"subnet.cidr":     "10.0.0.0/24",
		"instance.image":  "image2",
		"instance.subnet": "mysubnet",
	})

	expect := []*ast.Statement{
		{Node: &ast.DeclarationNode{
			Expr: &ast.CommandNode{
				Action: "create", Entity: "vpc",
				Params: map[string]interface{}{"count": 10},
			},
		}},
		{Node: &ast.DeclarationNode{
			Expr: &ast.CommandNode{
				Action: "create", Entity: "subnet",
				Params: map[string]interface{}{"cidr": "10.0.0.0/24"},
			},
		}},
		{Node: &ast.CommandNode{
			Action: "create", Entity: "instance",
			Params: map[string]interface{}{"type": "t1", "image": "image2", "subnet": "mysubnet"},
		}},
	}

	if got, want := templ.Statements, expect; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestResolveHoles(t *testing.T) {
	s := &Template{AST: &ast.AST{}}

	expr := &ast.CommandNode{
		Entity: "president",
		Holes:  map[string]string{"name": "presidentName", "rank": "presidentRank"},
	}
	s.Statements = append(s.Statements, &ast.Statement{Node: expr})

	decl := &ast.DeclarationNode{
		Expr: &ast.CommandNode{
			Entity: "wife",
			Holes:  map[string]string{"age": "wifeAge", "name": "wifeName"},
		},
	}
	s.Statements = append(s.Statements, &ast.Statement{Node: decl})

	fills := map[string]interface{}{
		"presidentName": "trump",
		"presidentRank": 45,
		"wifeAge":       40,
		"wifeName":      "melania",
	}

	filled, err := s.ResolveHoles(fills)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{"name": "trump", "rank": 45}
	if got, want := expr.Params, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := len(expr.Holes), 0; got != want {
		t.Fatalf("length of holes: got %d, want %d", got, want)
	}

	cmd := decl.Expr.(*ast.CommandNode)
	expected = map[string]interface{}{"age": 40, "name": "melania"}
	if got, want := cmd.Params, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := len(cmd.Holes), 0; got != want {
		t.Fatalf("length of holes: got %d, want %d", got, want)
	}

	expectedFilled := map[string]interface{}{
		"president.name": "trump",
		"president.rank": 45,
		"wife.age":       40,
		"wife.name":      "melania",
	}
	if got, want := filled, expectedFilled; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
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

func (r *mockDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	for _, expect := range r.expects {
		if lookups[0] == expect.action && lookups[1] == expect.entity {
			expect.lookupDone = true

			return func(params map[string]interface{}) (interface{}, error) {
				if got, want := expect.expectedParams, params; !reflect.DeepEqual(got, want) {
					return nil, fmt.Errorf("[%s %s] params mismatch: expected %v, got %v", expect.action, expect.entity, got, want)
				}
				return r.prefix + expect.entity, nil
			}, nil
		}
	}

	return func(params map[string]interface{}) (interface{}, error) {
		return nil, errors.New("Unexpected lookup fallthrough")
	}, nil
}

func (r *mockDriver) SetLogger(*logger.Logger) {}
func (r *mockDriver) SetDryRun(bool)           {}

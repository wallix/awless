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

func TestIsSameAsAst(t *testing.T) {
	tcases := []struct{ tpl string }{
		{tpl: "create instance subnet=@my-subnet count=4"},
		{tpl: "create vpc name=any\ncreate subnet ip=10.0.0.0\ndelete instance id=i-5d678\nstop instance id=i-5d678"},
		{tpl: "myvar = create vpc name={my.hole}\ndelete vpc id=$myvar\nid = create instance name=inst"},
		{tpl: "create vpc array=1,2,3"},
	}
	for i, tcase := range tcases {
		tpl := MustParse(tcase.tpl)
		if got, want := MustParse(tpl.String()), tpl; !want.IsSameAs(got) {
			t.Fatalf("%d: got \n%s\n, want \n%s\n", i+1, got, want)
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

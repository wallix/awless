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
	"github.com/wallix/awless/template/driver"
	"github.com/wallix/awless/template/internal/ast"
)

type noopDriver struct{}

func (d *noopDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	return func(driver.Context, map[string]interface{}) (interface{}, error) { return nil, nil }, nil
}
func (d *noopDriver) SetLogger(*logger.Logger) {}
func (d *noopDriver) SetDryRun(bool)           {}

type errorDriver struct {
	err error
}

func (d *errorDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	return func(driver.Context, map[string]interface{}) (interface{}, error) { return nil, d.err }, nil
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

		env := &Env{Driver: tcase.driver}
		ran, _ := templ.Run(env)

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

func TestRunDriverOnTemplate(t *testing.T) {
	t.Run("Driver run TWICE multiline statement", func(t *testing.T) {
		s, err := Parse(`createdvpc = create vpc count=1
createdsubnet = create subnet vpc=$createdvpc
create instance subnet=$createdsubnet`)
		if err != nil {
			t.Fatal(err)
		}

		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{{
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

		env := &Env{Driver: mDriver}
		if _, err := s.Run(env); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}

		mDriver = &mockDriver{t: t, prefix: "myother", expects: []*expectation{{
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

		env = &Env{Driver: mDriver}
		if _, err := s.Run(env); err != nil {
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
			Params: map[string]ast.CompositeValue{"count": ast.NewInterfaceValue(1)},
		}}
		s.Statements = append(s.Statements, n)

		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		env := &Env{Driver: mDriver}
		if _, err := s.Run(env); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver visit declaration nodes", func(t *testing.T) {
		s, err := Parse(`myvar = create vpc count=1`)
		if err != nil {
			t.Fatal(err)
		}

		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{{
			action: "create", entity: "vpc",
			expectedParams: map[string]interface{}{"count": 1},
		}},
		}

		env := &Env{Driver: mDriver}
		executedTemplate, err := s.Run(env)
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

	t.Run("Driver visit action nodes with list", func(t *testing.T) {
		s, err := Parse(`subnet1 = create subnet name=test
create loadbalancer subnets=[$subnet1,sub-1234]`)
		if err != nil {
			t.Fatal(err)
		}

		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{
			{action: "create", entity: "subnet", expectedParams: map[string]interface{}{"name": "test"}},
			{action: "create", entity: "loadbalancer", expectedParams: map[string]interface{}{"subnets": []interface{}{"mynewsubnet", "sub-1234"}}},
		},
		}

		env := &Env{Driver: mDriver}
		if _, err := s.Run(env); err != nil {
			t.Fatal(err)
		}
		if err := mDriver.lookupsCalled(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver stops when there is an error on a statement", func(t *testing.T) {
		s, err := Parse(`subnet1 = create subnet name=mysubnet
create instance name=instance-with-error
create vpc name=never-achieved`)
		if err != nil {
			t.Fatal(err)
		}
		errorMsg := "can not create instance"
		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{
			{action: "create", entity: "subnet", expectedParams: map[string]interface{}{"name": "mysubnet"}},
			{action: "create", entity: "instance", expectedParams: map[string]interface{}{"name": "instance-with-error"}, errorToReturn: errors.New(errorMsg)},
			{action: "create", entity: "vpc", expectedParams: map[string]interface{}{"name": "never-achieved"}},
		},
		}
		env := &Env{Driver: mDriver}
		if _, err := s.Run(env); err != nil {
			t.Fatal(err)
		}
		if expect := mDriver.expects[0]; !expect.lookupDone {
			t.Fatalf("expect %s %s done, got %t", expect.action, expect.entity, expect.lookupDone)
		}
		if expect := mDriver.expects[1]; !expect.lookupDone {
			t.Fatalf("expect %s %s done, got %t", expect.action, expect.entity, expect.lookupDone)
		}
		if expect := mDriver.expects[2]; expect.lookupDone {
			t.Fatalf("expect %s %s not done, got %t", expect.action, expect.entity, expect.lookupDone)
		}
	})

	t.Run("Dryrun and run are performed on cloned template", func(t *testing.T) {
		tplText := `subnet1 = create subnet name=mysubnet
create instance list=[test,test2] name=myinstance subnet=$subnet1
create vpc name=myvpc subnet=$subnet1`
		s, err := Parse(tplText)
		if err != nil {
			t.Fatal(err)
		}
		mDriver := &mockDriver{t: t, prefix: "mynew", expects: []*expectation{
			{action: "create", entity: "subnet", expectedParams: map[string]interface{}{"name": "mysubnet"}},
			{action: "create", entity: "instance", expectedParams: map[string]interface{}{"name": "myinstance", "subnet": "mynewsubnet", "list": []interface{}{"test", "test2"}}},
			{action: "create", entity: "vpc", expectedParams: map[string]interface{}{"name": "myvpc", "subnet": "mynewsubnet"}},
		},
		}
		env := &Env{Driver: mDriver}
		if err := s.DryRun(env); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Run(env); err != nil {
			t.Fatal(err)
		}
		if got, want := s.String(), tplText; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}
func TestGetTemplateUniqueDefinitions(t *testing.T) {
	text := "create instance name=nemo\ncreate keypair name=mykey\ncreate tag key=mine\ncreate instance\ncreate keypair"
	tpl := MustParse(text)

	lookup := func(key string) (t Definition, ok bool) {
		t, ok = DefsExample[key]
		return
	}
	defs := tpl.UniqueDefinitions(lookup)

	if got, want := len(defs), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := defs[0].Name(), "createinstance"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := defs[1].Name(), "createkeypair"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := defs[2].Name(), "createtag"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

type expectation struct {
	lookupDone     bool
	action, entity string
	expectedParams map[string]interface{}
	errorToReturn  error
}

type mockDriver struct {
	expects []*expectation
	prefix  string
	t       *testing.T
}

func (r *mockDriver) lookupsCalled() error {
	for _, expect := range r.expects {
		if expect.lookupDone == false {
			return fmt.Errorf("lookup for expectation '%s %s' with params '%+v' not called", expect.action, expect.entity, expect.expectedParams)
		}
	}

	return nil
}

func (r *mockDriver) Lookup(lookups ...string) (driver.DriverFn, error) {
	for _, expect := range r.expects {
		if lookups[0] == expect.action && lookups[1] == expect.entity {
			expect.lookupDone = true

			return func(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
				if got, want := expect.expectedParams, params; !reflect.DeepEqual(got, want) {
					r.t.Fatalf("[%s %s] params mismatch: expected %v, got %v", expect.action, expect.entity, got, want)
				}
				return r.prefix + expect.entity, expect.errorToReturn
			}, nil
		}
	}

	return func(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
		return nil, errors.New("Unexpected lookup fallthrough")
	}, nil
}

func (r *mockDriver) SetLogger(*logger.Logger) {}
func (r *mockDriver) SetDryRun(bool)           {}

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

package ast

import (
	"reflect"
	"testing"
)

func TestCloneAST(t *testing.T) {
	tree := &AST{}

	tree.Statements = append(tree.Statements, &Statement{Node: &DeclarationNode{
		Ident: "myvar",
		Expr: &CommandNode{
			Action: "create", Entity: "vpc",
			Refs:   map[string]string{"myname": "name"},
			Params: map[string]interface{}{"count": 1},
			Holes:  make(map[string]string),
		}}}, &Statement{Node: &DeclarationNode{
		Ident: "myothervar",
		Expr: &CommandNode{
			Action: "create", Entity: "subnet",
			Refs:   make(map[string]string),
			Params: make(map[string]interface{}),
			Holes:  map[string]string{"vpc": "myvar"},
		}}}, &Statement{Node: &CommandNode{
		Action: "create", Entity: "instance",
		Refs:   make(map[string]string),
		Params: make(map[string]interface{}),
		Holes:  map[string]string{"subnet": "myothervar"},
	}},
	)

	clone := tree.Clone()

	if got, want := clone, tree; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v", got, want)
	}

	clone.Statements[0].Node.(*DeclarationNode).Expr.(*CommandNode).Params["new"] = "trump"

	if got, want := clone.Statements, tree.Statements; reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %s\n\nwant %s", got, want)
	}
}

func TestNodesEqual(t *testing.T) {
	tcases := []struct {
		n1, n2 Node
		expect bool
	}{
		{
			n1: &CommandNode{Action: "create", Entity: "vpc",
				Refs:   map[string]string{"myname": "name"},
				Params: map[string]interface{}{"count": 1},
				Holes:  make(map[string]string)},
			n2: &CommandNode{Action: "create", Entity: "vpc",
				Refs:   map[string]string{"myname": "name"},
				Params: map[string]interface{}{"count": 1},
				Holes:  make(map[string]string)},
			expect: true,
		},
		{
			n1:     &CommandNode{Action: "create", Entity: "subnet"},
			n2:     &CommandNode{Action: "create", Entity: "vpc"},
			expect: false,
		},
		{
			n1:     &CommandNode{Refs: map[string]string{"myname": "name"}},
			n2:     &CommandNode{Refs: map[string]string{"myname": "name", "other": "test"}},
			expect: false,
		},
		{
			n1:     &CommandNode{Params: map[string]interface{}{"count": 1}},
			n2:     &CommandNode{Params: map[string]interface{}{"count": "1"}},
			expect: false,
		},
		{
			n1:     &DeclarationNode{Ident: "ident", Expr: &CommandNode{Action: "create", Entity: "subnet"}},
			n2:     &DeclarationNode{Ident: "ident", Expr: &CommandNode{Action: "create", Entity: "subnet"}},
			expect: true,
		},
		{
			n1:     &CommandNode{Action: "create", Entity: "subnet"},
			n2:     &DeclarationNode{Expr: &CommandNode{Action: "create", Entity: "subnet"}},
			expect: false,
		},
		{
			n1:     &DeclarationNode{Ident: "ident1"},
			n2:     &DeclarationNode{Ident: "ident2"},
			expect: false,
		},
	}
	for i, tcase := range tcases {
		if got, want := tcase.n1.Equal(tcase.n2), tcase.expect; got != want {
			t.Fatalf("%d: got %t, want %t", i, got, want)
		}
	}
}

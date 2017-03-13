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
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Node interface {
	clone() Node
	String() string
	Equal(Node) bool
}

type AST struct {
	Statements []*Statement

	// state to build the AST
	currentStatement *Statement
	currentKey       string
}

type Statement struct {
	Node
}

type DeclarationNode struct {
	Ident string
	Expr  ExpressionNode
}

func (n *DeclarationNode) Equal(n2 Node) bool {
	return reflect.DeepEqual(n, n2)
}

type ExpressionNode interface {
	Node
	Result() interface{}
	Err() error
}

type CommandNode struct {
	CmdResult interface{}
	CmdErr    error

	Action, Entity string
	Refs           map[string]string
	Params         map[string]interface{}
	Holes          map[string]string
}

func (n *CommandNode) Result() interface{} { return n.CmdResult }
func (n *CommandNode) Err() error          { return n.CmdErr }

func (n *CommandNode) Equal(n2 Node) bool {
	return reflect.DeepEqual(n, n2)
}

func (s *Statement) clone() *Statement {
	newStat := &Statement{}
	newStat.Node = s.Node.clone()

	return newStat
}

func (a *AST) String() string {
	var all []string
	for _, stat := range a.Statements {
		all = append(all, stat.String())
	}
	return strings.Join(all, "\n")
}

func (n *DeclarationNode) clone() Node {
	return &DeclarationNode{
		Ident: n.Ident,
		Expr:  n.Expr.clone().(ExpressionNode),
	}
}

func (n *DeclarationNode) String() string {
	return fmt.Sprintf("%s = %s", n.Ident, n.Expr)
}

func (n *CommandNode) clone() Node {
	cmd := &CommandNode{
		Action: n.Action, Entity: n.Entity,
		Refs:   make(map[string]string),
		Params: make(map[string]interface{}),
		Holes:  make(map[string]string),
	}

	for k, v := range n.Refs {
		cmd.Refs[k] = v
	}
	for k, v := range n.Params {
		cmd.Params[k] = v
	}
	for k, v := range n.Holes {
		cmd.Holes[k] = v
	}

	return cmd
}

func (n *CommandNode) String() string {
	var all []string
	for k, v := range n.Refs {
		all = append(all, fmt.Sprintf("%s=$%s", k, v))
	}
	for k, v := range n.Params {
		switch vv := v.(type) {
		case nil:
			continue
		case []string:
			all = append(all, fmt.Sprintf("%s=%s", k, strings.Join(vv, ",")))
		default:
			all = append(all, fmt.Sprintf("%s=%v", k, v))
		}

	}
	for k, v := range n.Holes {
		all = append(all, fmt.Sprintf("%s={%s}", k, v))
	}

	sort.Strings(all)

	var buff bytes.Buffer

	fmt.Fprintf(&buff, "%s %s", n.Action, n.Entity)

	if len(all) > 0 {
		fmt.Fprintf(&buff, " %s", strings.Join(all, " "))
	}

	return buff.String()

}

func (n *CommandNode) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	if n.Params == nil {
		n.Params = make(map[string]interface{})
	}
	for key, hole := range n.Holes {
		if val, ok := fills[hole]; ok {
			n.Params[key] = val
			processed[key] = val
			delete(n.Holes, key)
		}
	}
	return processed
}

func (n *CommandNode) ProcessRefs(fills map[string]interface{}) {
	if n.Params == nil {
		n.Params = make(map[string]interface{})
	}
	for key, ref := range n.Refs {
		if val, ok := fills[ref]; ok {
			n.Params[key] = val
			delete(n.Refs, key)
		}
	}
}

func (a *AST) Clone() *AST {
	clone := &AST{}
	for _, stat := range a.Statements {
		clone.Statements = append(clone.Statements, stat.clone())
	}
	return clone
}

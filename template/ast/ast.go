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
	"fmt"
	"strings"
)

type Node interface {
	clone() Node
	String() string
}

type Statement struct {
	Node
	Result interface{}
	Line   string
	Err    error
}

func (s *Statement) clone() *Statement {
	newStat := &Statement{}
	newStat.Node = s.Node.clone()
	newStat.Result = s.Result
	newStat.Err = s.Err

	return newStat
}

func (s *Statement) Action() string {
	switch n := s.Node.(type) {
	case *ExpressionNode:
		return n.Action
	case *DeclarationNode:
		return n.Right.Action
	default:
		panic(fmt.Sprintf("unknown type of node %T", s.Node))
	}
}

func (s *Statement) Entity() string {
	switch n := s.Node.(type) {
	case *ExpressionNode:
		return n.Entity
	case *DeclarationNode:
		return n.Right.Entity
	default:
		panic(fmt.Sprintf("unknown type of node %T", s.Node))
	}
}

func (s *Statement) Params() map[string]interface{} {
	switch n := s.Node.(type) {
	case *ExpressionNode:
		return n.Params
	case *DeclarationNode:
		return n.Right.Params
	default:
		panic(fmt.Sprintf("unknown type of node %T", s.Node))
	}
}

type AST struct {
	Statements []*Statement

	currentStatement *Statement
	currentKey       string
}

func (a *AST) String() string {
	var all []string
	for _, stat := range a.Statements {
		all = append(all, stat.String())
	}
	return strings.Join(all, "\n")
}

type IdentifierNode struct {
	Ident string
	Val   interface{}
}

func (n *IdentifierNode) String() string {
	return fmt.Sprintf("%s", n.Ident)
}

func (n *IdentifierNode) clone() Node {
	return &IdentifierNode{
		Ident: n.Ident,
		Val:   n.Val,
	}
}

type DeclarationNode struct {
	Left  *IdentifierNode
	Right *ExpressionNode
}

func (n *DeclarationNode) clone() Node {
	return &DeclarationNode{
		Left:  n.Left.clone().(*IdentifierNode),
		Right: n.Right.clone().(*ExpressionNode),
	}
}

func (n *DeclarationNode) String() string {
	return fmt.Sprintf("%s = %s", n.Left, n.Right)
}

type ExpressionNode struct {
	Action, Entity string
	Refs           map[string]string
	Params         map[string]interface{}
	Aliases        map[string]string
	Holes          map[string]string
}

func (n *ExpressionNode) clone() Node {
	expr := &ExpressionNode{
		Action: n.Action, Entity: n.Entity,
		Refs:    make(map[string]string),
		Params:  make(map[string]interface{}),
		Aliases: make(map[string]string),
		Holes:   make(map[string]string),
	}

	for k, v := range n.Refs {
		expr.Refs[k] = v
	}
	for k, v := range n.Params {
		expr.Params[k] = v
	}
	for k, v := range n.Aliases {
		expr.Aliases[k] = v
	}
	for k, v := range n.Holes {
		expr.Holes[k] = v
	}

	return expr
}

func (n *ExpressionNode) String() string {
	var all []string
	for k, v := range n.Refs {
		all = append(all, fmt.Sprintf("%s=$%v", k, v))
	}
	for k, v := range n.Params {
		all = append(all, fmt.Sprintf("%s=%v", k, v))
	}
	for k, v := range n.Aliases {
		all = append(all, fmt.Sprintf("%s=@%s", k, v))
	}
	for k, v := range n.Holes {
		all = append(all, fmt.Sprintf("%s={%s}", k, v))
	}
	return fmt.Sprintf("%s %s %s", n.Action, n.Entity, strings.Join(all, " "))
}

func (n *ExpressionNode) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
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

func (n *ExpressionNode) ProcessRefs(fills map[string]interface{}) {
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

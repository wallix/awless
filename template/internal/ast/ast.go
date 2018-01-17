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
	"sort"
	"strings"

	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"
)

type Node interface {
	clone() Node
	String() string
}

type AST struct {
	Statements []*Statement

	// state to build the AST
	stmtBuilder *statementBuilder
}

type Statement struct {
	Node
}

type DeclarationNode struct {
	Ident string
	Expr  ExpressionNode
}

type ExpressionNode interface {
	Node
	Result() interface{}
	Err() error
}

type Command interface {
	ParamsSpec() params.Spec
	Run(env.Running, map[string]interface{}) (interface{}, error)
}

func (c *CommandNode) Result() interface{} { return c.CmdResult }
func (c *CommandNode) Err() error          { return c.CmdErr }

func (c *CommandNode) Keys() (keys []string) {
	for k := range c.ParamNodes {
		keys = append(keys, k)
	}
	for k := range c.Refs {
		keys = append(keys, k)
	}
	return
}

func (c *CommandNode) String() string {
	var all []string

	for k, v := range c.ParamNodes {
		switch vv := v.(type) {
		case string:
			all = append(all, fmt.Sprintf("%s=%v", k, quoteStringIfNeeded(vv)))
		case []interface{}:
			var a []string
			for _, e := range vv {
				switch ee := e.(type) {
				case string:
					a = append(a, quoteStringIfNeeded(ee))
				default:
					a = append(a, fmt.Sprint(ee))
				}
			}
			all = append(all, fmt.Sprintf("%s=[%s]", k, strings.Join(a, ",")))
		default:
			all = append(all, fmt.Sprintf("%s=%v", k, v))
		}
	}
	for k, v := range c.Refs {
		all = append(all, fmt.Sprintf("%s=%v", k, v))
	}

	sort.Strings(all)

	var buff bytes.Buffer

	fmt.Fprintf(&buff, "%s %s", c.Action, c.Entity)

	if len(all) > 0 {
		fmt.Fprintf(&buff, " %s", strings.Join(all, " "))
	}

	return buff.String()
}

func (c *CommandNode) clone() Node {
	cmd := &CommandNode{
		Command: c.Command,
		Action:  c.Action, Entity: c.Entity,
		ParamNodes: make(map[string]interface{}),
		Refs:       make(map[string]interface{}),
	}

	for k, v := range c.ParamNodes {
		cmd.ParamNodes[k] = v
	}
	for k, v := range c.Refs {
		cmd.Refs[k] = v
	}
	return cmd
}

func (c *CommandNode) ProcessRefs(refs map[string]interface{}) {
	for paramKey, param := range c.Refs {
		if ref, ok := param.(RefNode); ok {
			for k, v := range refs {
				if k == ref.key {
					c.ParamNodes[paramKey] = v
				}
			}
		}

		if list, ok := param.(ListNode); ok {
			var new []interface{}
			for _, e := range list.arr {
				newElem := e
				if ref, isRef := e.(RefNode); isRef {
					for k, v := range refs {
						if k == ref.key {
							newElem = v
						}
					}
				}
				new = append(new, newElem)
			}
			c.ParamNodes[paramKey] = new
		}
	}
}

func (c *CommandNode) ToDriverParams() map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range c.ParamNodes {
		switch node := v.(type) {
		case InterfaceNode:
			params[k] = node.i
		case RefNode, HoleNode, AliasNode:
		default:
			params[k] = node
		}
	}
	return params
}

func (c *CommandNode) ToFillerParams() map[string]interface{} {
	params := make(map[string]interface{})
	fn := func(k string, v interface{}) interface{} {
		switch vv := v.(type) {
		case InterfaceNode:
			return vv.i
		case AliasNode:
			return v
		}
		return nil
	}

	for k, v := range c.ParamNodes {
		i := fn(k, v)
		if i != nil {
			params[k] = i
			continue
		}
		switch vv := v.(type) {
		case ListNode:
			var arr []interface{}
			for _, a := range vv.arr {
				arr = append(arr, fn(k, a))
			}
			params[k] = NewListNode(arr)
		}
	}
	return params
}

func (s *Statement) Clone() *Statement {
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
	decl := &DeclarationNode{
		Ident: n.Ident,
	}
	if n.Expr != nil {
		decl.Expr = n.Expr.clone().(ExpressionNode)
	}
	return decl
}

func (n *DeclarationNode) String() string {
	return fmt.Sprintf("%s = %s", n.Ident, n.Expr)
}

func (a *AST) Clone() *AST {
	clone := &AST{}
	for _, stat := range a.Statements {
		clone.Statements = append(clone.Statements, stat.Clone())
	}
	return clone
}

func (a *AST) clone() Node {
	return a.Clone()
}

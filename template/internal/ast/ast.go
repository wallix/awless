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
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Node interface {
	clone() Node
	String() string
}

type AST struct {
	Statements []*Statement

	// state to build the AST
	currentStatement   *Statement
	currentKey         string
	currentListBuilder *listValueBuilder
	stmtBuilder        *statementBuilder
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

type WithHoles interface {
	ProcessHoles(fills map[string]interface{}) (processed map[string]interface{})
	GetHoles() map[string][]string
}

type Command interface {
	Run(ctx map[string]interface{}, params map[string]interface{}) (interface{}, error)
	DryRun(ctx, params map[string]interface{}) (interface{}, error)
}

type CommandNode struct {
	Command
	CmdResult interface{}
	CmdErr    error

	Action, Entity string
	Params         map[string]CompositeValue
}

func (c *CommandNode) Result() interface{} { return c.CmdResult }
func (c *CommandNode) Err() error          { return c.CmdErr }

func (c *CommandNode) Keys() (keys []string) {
	for k := range c.Params {
		keys = append(keys, k)
	}
	return
}

func (c *CommandNode) String() string {
	var all []string

	for k, v := range c.Params {
		all = append(all, fmt.Sprintf("%s=%s", k, v.String()))
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
		Params: make(map[string]CompositeValue),
	}

	for k, v := range c.Params {
		cmd.Params[k] = v.Clone()
	}
	return cmd
}

func (c *CommandNode) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})

	for _, param := range c.Params {
		if withHoles, ok := param.(WithHoles); ok {
			paramProcessed := withHoles.ProcessHoles(fills)
			for k, v := range paramProcessed {
				processed[k] = v
			}
		}
	}
	return processed
}

func (c *CommandNode) GetHoles() map[string][]string {
	holes := make(map[string][]string)
	for paramKey, param := range c.Params {
		if withHoles, ok := param.(WithHoles); ok {
			for k := range withHoles.GetHoles() {
				holes[k] = append(holes[k], strings.Join([]string{c.Action, c.Entity, paramKey}, "."))
			}

		}
	}
	return holes
}

func (c *CommandNode) ProcessRefs(refs map[string]interface{}) {
	for _, param := range c.Params {
		if withRef, ok := param.(WithRefs); ok {
			withRef.ProcessRefs(refs)
		}
	}
}

func (c *CommandNode) GetRefs() (refs []string) {
	for _, param := range c.Params {
		if withRef, ok := param.(WithRefs); ok {
			refs = append(refs, withRef.GetRefs()...)
		}
	}
	return
}

func (c *CommandNode) ReplaceRef(key string, value CompositeValue) {
	for k, param := range c.Params {
		if withRef, ok := param.(WithRefs); ok {
			if withRef.IsRef(key) {
				c.Params[k] = value
			} else {
				withRef.ReplaceRef(key, value)
			}
		}
	}
}

func (c *CommandNode) IsRef(key string) bool {
	return false
}

func (c *CommandNode) ToDriverParams() map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range c.Params {
		if v.Value() != nil {
			params[k] = v.Value()
		}
	}
	return params
}

func (c *CommandNode) ToFillerParams() map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range c.Params {
		if v.Value() != nil {
			params[k] = v.Value()
		} else if _, ok := v.(WithAlias); ok {
			params[k] = v
		}
	}
	return params
}

type ValueNode struct {
	Value CompositeValue
}

func (n *ValueNode) clone() Node {
	return &ValueNode{
		Value: n.Value.Clone(),
	}
}

func (n *ValueNode) String() string {
	return n.Value.String()
}

func (n *ValueNode) Result() interface{} { return n.Value }
func (n *ValueNode) Err() error          { return nil }

func (n *ValueNode) IsResolved() bool {
	if withHoles, ok := n.Value.(WithHoles); ok {
		return len(withHoles.GetHoles()) == 0
	}
	return true
}

func (n *ValueNode) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	if withHoles, ok := n.Value.(WithHoles); ok {
		return withHoles.ProcessHoles(fills)
	}
	return make(map[string]interface{})
}

func (n *ValueNode) ProcessRefs(refs map[string]interface{}) {
	if withRef, ok := n.Value.(WithRefs); ok {
		withRef.ProcessRefs(refs)
	}
}

func (n *ValueNode) GetRefs() (refs []string) {
	if withRef, ok := n.Value.(WithRefs); ok {
		refs = append(refs, withRef.GetRefs()...)
	}
	return
}

func (n *ValueNode) ReplaceRef(key string, value CompositeValue) {
	if withRef, ok := n.Value.(WithRefs); ok {
		if withRef.IsRef(key) {
			n.Value = value
		} else {
			withRef.ReplaceRef(key, value)
		}
	}
}

func (n *ValueNode) IsRef(key string) bool {
	return false
}

func (n *ValueNode) GetHoles() map[string][]string {
	if withHoles, ok := n.Value.(WithHoles); ok {
		return withHoles.GetHoles()
	}
	return make(map[string][]string)
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

func printParamValue(i interface{}) string {
	switch ii := i.(type) {
	case nil:
		return ""
	case []string:
		return "[" + strings.Join(ii, ",") + "]"
	case []interface{}:
		var strs []string
		for _, val := range ii {
			strs = append(strs, fmt.Sprint(val))
		}
		return "[" + strings.Join(strs, ",") + "]"
	case string:
		return quoteStringIfNeeded(ii)
	default:
		return fmt.Sprintf("%v", i)
	}
}

func (a *AST) Clone() *AST {
	clone := &AST{}
	for _, stat := range a.Statements {
		clone.Statements = append(clone.Statements, stat.Clone())
	}
	return clone
}

var SimpleStringValue = regexp.MustCompile("^[a-zA-Z0-9-._:/+;~@<>*]+$") // in sync with [a-zA-Z0-9-._:/+;~@<>]+ in PEG (with ^ and $ around)

func quoteStringIfNeeded(input string) string {
	if _, err := strconv.Atoi(input); err == nil {
		return "'" + input + "'"
	}
	if _, err := strconv.ParseFloat(input, 64); err == nil {
		return "'" + input + "'"
	}
	if SimpleStringValue.MatchString(input) {
		return input
	} else {
		return quoteString(input)
	}
}

func quoteString(str string) string {
	if strings.ContainsRune(str, '\'') {
		return "\"" + str + "\""
	} else {
		return "'" + str + "'"
	}
}

func isQuoted(str string) bool {
	if len(str) < 2 {
		return false
	}
	if str[0] == '\'' && str[len(str)-1] == '\'' {
		return true
	}
	if str[0] == '"' && str[len(str)-1] == '"' {
		return true
	}
	return false
}

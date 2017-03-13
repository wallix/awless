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
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver"
)

type Template struct {
	ID string
	*ast.AST
}

func (s *Template) Run(d driver.Driver) (*Template, error) {
	vars := map[string]interface{}{}

	current := &Template{AST: s.Clone()}
	current.ID = ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()

	for _, sts := range current.Statements {
		switch sts.Node.(type) {
		case *ast.CommandNode:
			cmd := sts.Node.(*ast.CommandNode)
			fn, err := d.Lookup(cmd.Action, cmd.Entity)
			if err != nil {
				return current, err
			}
			cmd.ProcessRefs(vars)

			if cmd.CmdResult, cmd.CmdErr = fn(cmd.Params); cmd.CmdErr != nil {
				return current, cmd.CmdErr
			}
		case *ast.DeclarationNode:
			ident := sts.Node.(*ast.DeclarationNode).Ident
			expr := sts.Node.(*ast.DeclarationNode).Expr
			switch expr.(type) {
			case *ast.CommandNode:
				cmd := expr.(*ast.CommandNode)
				fn, err := d.Lookup(cmd.Action, cmd.Entity)
				if err != nil {
					return current, err
				}
				cmd.ProcessRefs(vars)

				if cmd.CmdResult, cmd.CmdErr = fn(cmd.Params); cmd.CmdErr != nil {
					return current, cmd.CmdErr
				}
				vars[ident] = cmd.CmdResult
			}
		}
	}

	return current, nil
}

func (s *Template) Compile(d driver.Driver) (*Template, error) {
	defer d.SetDryRun(false)
	d.SetDryRun(true)

	return s.Run(d)
}

func (s *Template) Validate(rules ...Validator) (all []error) {
	for _, rule := range rules {
		errs := rule.Execute(s)
		all = append(all, errs...)
	}

	return
}

func (t *Template) HasErrors() bool {
	for _, cmd := range t.CommandNodesIterator() {
		if cmd.CmdErr != nil {
			return true
		}
	}
	return false
}

func (s *Template) Visit(v Visitor) error {
	return v.Visit(s.CommandNodesIterator())
}

func (s *Template) visitCommandNodes(fn func(n *ast.CommandNode)) {
	for _, cmd := range s.CommandNodesIterator() {
		fn(cmd)
	}
}

func (s *Template) visitCommandNodesE(fn func(n *ast.CommandNode) error) error {
	for _, cmd := range s.CommandNodesIterator() {
		err := fn(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Template) CommandNodesIterator() (nodes []*ast.CommandNode) {
	for _, sts := range s.Statements {
		switch sts.Node.(type) {
		case *ast.CommandNode:
			nodes = append(nodes, sts.Node.(*ast.CommandNode))
		case *ast.DeclarationNode:
			expr := sts.Node.(*ast.DeclarationNode).Expr
			switch expr.(type) {
			case *ast.CommandNode:
				nodes = append(nodes, expr.(*ast.CommandNode))
			}
		}
	}
	return
}

func (s *Template) CmdNodesReverseIterator() (nodes []*ast.CommandNode) {
	for i := len(s.Statements) - 1; i >= 0; i-- {
		sts := s.Statements[i]
		switch sts.Node.(type) {
		case *ast.CommandNode:
			nodes = append(nodes, sts.Node.(*ast.CommandNode))
		case *ast.DeclarationNode:
			expr := sts.Node.(*ast.DeclarationNode).Expr
			switch expr.(type) {
			case *ast.CommandNode:
				nodes = append(nodes, expr.(*ast.CommandNode))
			}
		}
	}
	return
}

func (s *Template) IsSameAs(t2 *Template) bool {
	if s == t2 {
		return true
	}
	if s == nil || t2 == nil {
		return false
	}
	if len(s.Statements) != len(t2.Statements) {
		return false
	}
	for i := 0; i < len(s.Statements); i++ {
		s1 := s.Statements[i]
		s2 := t2.Statements[i]
		if !s1.Node.Equal(s2.Node) {
			return false
		}
	}

	return true
}

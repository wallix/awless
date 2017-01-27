package template

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver"
)

type Template struct {
	*ast.AST
}

type Operation struct {
	ID     string
	Line   string
	Output interface{}
	Err    error
}

func (op *Operation) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("operation[uid: %s, executed: ", op.ID))
	if op.Output != nil {
		out.WriteString(fmt.Sprintf("%v <- ", op.Output))
	}
	out.WriteString(fmt.Sprintf("%s", op.Line))
	if op.Err != nil {
		out.WriteString(fmt.Sprintf(": error: %s", op.Err))
	}
	out.WriteByte(']')

	return out.String()
}

func (s *Template) Run(d driver.Driver) (*Template, []*Operation, error) {
	vars := map[string]interface{}{}
	var operations []*Operation

	executedTemplate := &Template{s.Clone()}

	for _, sts := range executedTemplate.Statements {
		uid := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
		op := &Operation{ID: uid.String()}

		operations = append(operations, op)

		switch sts.(type) {
		case *ast.ExpressionNode:
			expr := sts.(*ast.ExpressionNode)
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)

			op.Line = expr.String()
			if op.Output, op.Err = fn(expr.Params); op.Err != nil {
				return executedTemplate, operations, op.Err
			}
		case *ast.DeclarationNode:
			ident := sts.(*ast.DeclarationNode).Left
			expr := sts.(*ast.DeclarationNode).Right
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)

			op.Output, op.Err = fn(expr.Params)
			ident.Val = op.Output
			op.Line = expr.String()
			if op.Err != nil {
				return executedTemplate, operations, op.Err
			}
			vars[ident.Ident] = ident.Val
		}
	}

	return executedTemplate, operations, nil
}

func (s *Template) Compile(d driver.Driver) (*Template, []*Operation, error) {
	defer d.SetDryRun(false)
	d.SetDryRun(true)

	return s.Run(d)
}

func (s *Template) GetAliases() map[string]string {
	aliases := make(map[string]string)
	each := func(expr *ast.ExpressionNode) {
		for k, v := range expr.Aliases {
			aliases[k] = v
		}
	}
	s.visitExpressionNodes(each)
	return aliases
}

func (s *Template) MergeParams(newParams map[string]interface{}) {
	each := func(expr *ast.ExpressionNode) {
		for k, v := range newParams {
			if strings.SplitN(k, ".", 2)[0] == expr.Entity {
				if expr.Params == nil {
					expr.Params = make(map[string]interface{})
				}
				expr.Params[strings.SplitN(k, ".", 2)[1]] = v
			}
		}
	}
	s.visitExpressionNodes(each)
}

func (s *Template) ResolveTemplate(refs map[string]interface{}) error {
	each := func(expr *ast.ExpressionNode) {
		expr.ProcessHoles(refs)
	}

	s.visitExpressionNodes(each)

	return nil
}

func (s *Template) InteractiveResolveTemplate(each func(question string) interface{}) error {
	fn := func(expr *ast.ExpressionNode) {
		for key, hole := range expr.Holes {
			if expr.Params == nil {
				expr.Params = make(map[string]interface{})
			}
			res := each(hole)
			expr.Params[key] = res
			delete(expr.Holes, key)
		}
	}

	s.visitExpressionNodes(fn)

	return nil
}

func (s *Template) visitExpressionNodes(fn func(n *ast.ExpressionNode)) {
	for _, sts := range s.Statements {
		var expr *ast.ExpressionNode

		switch sts.(type) {
		case *ast.ExpressionNode:
			expr = sts.(*ast.ExpressionNode)
		case *ast.DeclarationNode:
			expr = sts.(*ast.DeclarationNode).Right
		}

		if expr != nil {
			fn(expr)
		}
	}
}

package template

import (
	"strings"

	"github.com/wallix/awless/template/ast"
	"github.com/wallix/awless/template/driver"
)

type Template struct {
	*ast.AST
}

func (s *Template) Run(d driver.Driver) (*Template, error) {
	vars := map[string]interface{}{}

	executedTemplate := &Template{s.Clone()}

	for _, sts := range executedTemplate.Statements {
		switch sts.(type) {
		case *ast.ExpressionNode:
			expr := sts.(*ast.ExpressionNode)
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)
			if _, err := fn(expr.Params); err != nil {
				return executedTemplate, err
			}
		case *ast.DeclarationNode:
			ident := sts.(*ast.DeclarationNode).Left
			expr := sts.(*ast.DeclarationNode).Right
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)
			identVal, err := fn(expr.Params)
			ident.Val = identVal
			if err != nil {
				return executedTemplate, err
			}
			vars[ident.Ident] = ident.Val
		}
	}

	return executedTemplate, nil
}

func (s *Template) Compile(d driver.Driver) (*Template, error) {
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

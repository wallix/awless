package script

import (
	"strings"

	"github.com/wallix/awless/script/ast"
	"github.com/wallix/awless/script/driver"
)

func Visit(s *ast.Script, d driver.Driver) error {
	vars := map[string]interface{}{}

	for _, sts := range s.Statements {
		switch sts.(type) {
		case *ast.ExpressionNode:
			expr := sts.(*ast.ExpressionNode)
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessHoles(vars)
			_, err := fn(expr.Params)
			if err != nil {
				return err
			}
		case *ast.DeclarationNode:
			ident := sts.(*ast.DeclarationNode).Left
			expr := sts.(*ast.DeclarationNode).Right
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessHoles(vars)
			identVal, err := fn(expr.Params)
			ident.Val = identVal
			if err != nil {
				return err
			}
			vars[ident.Ident] = ident.Val
		}
	}

	return nil
}

func VisitExpressionNodes(s *ast.Script, fn func(n *ast.ExpressionNode)) {
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

func ResolveHolesWith(fills map[string]interface{}) func(expr *ast.ExpressionNode) {
	return func(expr *ast.ExpressionNode) {
		expr.ProcessHoles(fills)
	}
}

func InteractiveResolveHoles(fn func(question string) interface{}) func(expr *ast.ExpressionNode) {
	return func(expr *ast.ExpressionNode) {
		for key, hole := range expr.Holes {
			if expr.Params == nil {
				expr.Params = make(map[string]interface{})
			}
			res := fn(humanize(hole))
			expr.Params[key] = res
			delete(expr.Holes, key)
		}
	}
}

func humanize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}

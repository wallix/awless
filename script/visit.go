package script

import (
	"github.com/wallix/awless/script/ast"
	"github.com/wallix/awless/script/driver"
)

func Visit(s *ast.Script, d driver.Driver) (err error) {
	for _, sts := range s.Statements {
		switch sts.(type) {
		case *ast.ExpressionNode:
			expr := sts.(*ast.ExpressionNode)
			fn := d.Lookup(expr.Action, expr.Entity)
			_, err = fn(expr.Params)
		case *ast.DeclarationNode:
			ident := sts.(*ast.DeclarationNode).Left
			expr := sts.(*ast.DeclarationNode).Right
			fn := d.Lookup(expr.Action, expr.Entity)
			ident.Val, err = fn(expr.Params)
		}
	}

	return err
}

func VisitHoles(s *ast.Script, fills map[string]interface{}) {
	for _, sts := range s.Statements {
		var expr *ast.ExpressionNode

		switch sts.(type) {
		case *ast.ExpressionNode:
			expr = sts.(*ast.ExpressionNode)
		case *ast.DeclarationNode:
			expr = sts.(*ast.DeclarationNode).Right
		}

		if expr != nil {
			for key, hole := range expr.Holes {
				if val, ok := fills[hole]; ok {
					if expr.Params == nil {
						expr.Params = make(map[string]interface{})
					}
					expr.Params[key] = val
					delete(expr.Holes, key)
				}
			}
		}
	}
}

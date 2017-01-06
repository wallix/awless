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
			err = fn(expr.Params)
		}
	}

	return err
}

func VisitHoles(s *ast.Script, fills map[string]interface{}) {
	for _, sts := range s.Statements {
		switch sts.(type) {
		case *ast.ExpressionNode:
			expr := sts.(*ast.ExpressionNode)
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

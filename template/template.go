package template

import (
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

func (s *Template) Run(d driver.Driver) (*Template, error) {
	vars := map[string]interface{}{}

	current := &Template{AST: s.Clone()}

	for _, sts := range current.Statements {
		switch sts.Node.(type) {
		case *ast.ExpressionNode:
			expr := sts.Node.(*ast.ExpressionNode)
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)

			sts.Line = expr.String()
			if sts.Result, sts.Err = fn(expr.Params); sts.Err != nil {
				return current, sts.Err
			}
		case *ast.DeclarationNode:
			ident := sts.Node.(*ast.DeclarationNode).Left
			expr := sts.Node.(*ast.DeclarationNode).Right
			fn := d.Lookup(expr.Action, expr.Entity)
			expr.ProcessRefs(vars)

			sts.Result, sts.Err = fn(expr.Params)
			ident.Val = sts.Result
			sts.Line = expr.String()
			if sts.Err != nil {
				return current, sts.Err
			}
			vars[ident.Ident] = ident.Val
		}
	}

	return current, nil
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

		switch sts.Node.(type) {
		case *ast.ExpressionNode:
			expr = sts.Node.(*ast.ExpressionNode)
		case *ast.DeclarationNode:
			expr = sts.Node.(*ast.DeclarationNode).Right
		}

		if expr != nil {
			fn(expr)
		}
	}
}

type TemplateExecution struct {
	ID       string
	Executed []*ExecutedStatement
}

type ExecutedStatement struct {
	Line, Err, Result string
}

func NewTemplateExecution(tpl *Template) *TemplateExecution {
	out := &TemplateExecution{
		ID: ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String(),
	}

	for _, sts := range tpl.Statements {
		var errMsg string
		if sts.Err != nil {
			errMsg = sts.Err.Error()
		}
		var result string
		switch sts.Result.(type) {
		case string:
			result = sts.Result.(string)
		}
		out.Executed = append(out.Executed,
			&ExecutedStatement{Line: sts.Line, Result: result, Err: errMsg},
		)
	}

	return out
}

func (te *TemplateExecution) lines() (lines []string) {
	for _, ex := range te.Executed {
		lines = append(lines, ex.Line)
	}

	return
}

func (te *TemplateExecution) Revert() (*Template, error) {
	var lines []string

	for i := len(te.Executed) - 1; i >= 0; i-- {
		if exec := te.Executed[i]; exec.Err == "" && exec.Result != "" {
			n, err := ParseStatement(exec.Line)
			if err != nil {
				return nil, err
			}

			switch n.(type) {
			case *ast.ExpressionNode:
				node := n.(*ast.ExpressionNode)
				var revertAction string
				switch node.Action {
				case "create":
					revertAction = "delete"
				case "start":
					revertAction = "stop"
				case "stop":
					revertAction = "start"
				default:
					continue
				}

				lines = append(lines, fmt.Sprintf("%s %s id=%s\n", revertAction, node.Entity, exec.Result))
			default:
				return nil, fmt.Errorf("cannot parse [%s] as expression node", exec.Line)
			}
		}
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("revert: found nothing to revert from:\n%s\n(note: no revert provided for statement in error)", strings.Join(te.lines(), "\n"))
	}

	tpl, err := Parse(strings.Join(lines, "\n"))
	if err != nil {
		return nil, fmt.Errorf("revert: \n%s\n%s", strings.Join(lines, "\n"), err)
	}

	return tpl, nil
}

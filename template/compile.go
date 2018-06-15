package template

import (
	"errors"
	"fmt"
	"sort"

	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/internal/ast"
	"github.com/wallix/awless/template/params"
)

type Mode []compileFunc

var (
	TestCompileMode = []compileFunc{
		injectCommandsInNodesPass,
		failOnDeclarationWithNoResultPass,
		processAndValidateParamsPass,
		checkInvalidReferenceDeclarationsPass,
		resolveHolesPass,
		resolveMissingHolesPass,
		removeOptionalHolesPass,
		resolveAliasPass,
		inlineVariableValuePass,
		resolveParamsAndExtractRefsPass,
	}

	PreRevertCompileMode = []compileFunc{
		resolveParamsAndExtractRefsPass,
	}

	NewRunnerCompileMode = []compileFunc{
		injectCommandsInNodesPass,
		failOnDeclarationWithNoResultPass,
		processAndValidateParamsPass,
		checkInvalidReferenceDeclarationsPass,
		resolveHolesPass,
		resolveMissingHolesPass,
		removeOptionalHolesPass,
		resolveAliasPass,
		inlineVariableValuePass,
		failOnUnresolvedHolesPass,
		failOnUnresolvedAliasPass,
		resolveParamsAndExtractRefsPass,
		convertParamsPass,
		validateCommandsPass,
	}
)

func Compile(tpl *Template, cenv env.Compiling, mode ...Mode) (*Template, env.Compiling, error) {
	var pass *multiPass

	if len(mode) > 0 {
		pass = newMultiPass(mode[0]...)
	} else {
		pass = newMultiPass(NewRunnerCompileMode...)
	}

	return pass.compile(tpl, cenv)
}

type compileFunc func(*Template, env.Compiling) (*Template, env.Compiling, error)

// Leeloo Dallas
type multiPass struct {
	passes []compileFunc
}

func newMultiPass(passes ...compileFunc) *multiPass {
	return &multiPass{passes: passes}
}

func (p *multiPass) compile(tpl *Template, cenv env.Compiling) (newTpl *Template, newEnv env.Compiling, err error) {
	newTpl, newEnv = tpl, cenv
	for _, pass := range p.passes {
		newTpl, newEnv, err = pass(newTpl, newEnv)
		if err != nil {
			return
		}
	}

	return
}

func injectCommandsInNodesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	if cenv.LookupCommandFunc() == nil {
		return tpl, cenv, fmt.Errorf("command lookuper is undefined")
	}

	for _, node := range tpl.CommandNodesIterator() {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd, ok := cenv.LookupCommandFunc()(key).(ast.Command)
		if !ok {
			return tpl, cenv, fmt.Errorf("%s: casting: %v is not a command", key, cmd)
		}
		if cmd == nil {
			return tpl, cenv, fmt.Errorf("command for '%s' is nil", key)
		}
		node.Command = cmd
	}
	return tpl, cenv, nil
}

func failOnDeclarationWithNoResultPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	failOnDeclarationWithNoResult := func(node *ast.DeclarationNode) error {
		cmdNode, ok := node.Expr.(*ast.CommandNode)
		if !ok {
			return nil
		}
		type ER interface {
			ExtractResult(interface{}) string
		}
		if _, ok := cmdNode.Command.(ER); !ok {
			return cmdErr(cmdNode, "command does not return a result, cannot assign to a variable")
		}
		return nil
	}

	for _, dcl := range tpl.declarationNodesIterator() {
		if err := failOnDeclarationWithNoResult(dcl); err != nil {
			return tpl, cenv, err
		}
	}
	return tpl, cenv, nil
}

func processAndValidateParamsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	normalizeMissingRequiredParamsAsHoleAndValidate := func(node *ast.CommandNode) error {
		rule := node.ParamsSpec().Rule()

		missingRequired := rule.Missing(node.Keys())
		for _, e := range missingRequired {
			normalized := fmt.Sprintf("%s.%s", node.Entity, e)
			node.ParamNodes[e] = ast.NewHoleNode(normalized)
		}
		if err := params.Run(rule, node.Keys()); err != nil {
			return cmdErr(node, err)
		}

		_, optionals, suggested := params.List(rule)

		switch cenv.ParamsMode() {
		case env.REQUIRED_PARAMS_ONLY:
			return nil
		case env.REQUIRED_AND_SUGGESTED_PARAMS:
			suggested = excludeFromSlice(suggested, node.Keys())
		case env.ALL_PARAMS:
			suggested = excludeFromSlice(optionals, node.Keys())
		}

		for _, e := range suggested {
			key := fmt.Sprintf("%s.%s", node.Entity, e)
			node.ParamNodes[e] = ast.NewOptionalHoleNode(key)
		}
		return nil
	}

	err := tpl.visitCommandNodesE(normalizeMissingRequiredParamsAsHoleAndValidate)
	return tpl, cenv, err
}

func resolveParamsAndExtractRefsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	for _, node := range tpl.CommandNodesIterator() {
		for k, param := range node.ParamNodes {
			switch paramNode := param.(type) {
			case ast.InterfaceNode:
				node.ParamNodes[k] = paramNode.Value()
			case ast.RefNode:
				node.Refs[k] = paramNode
				delete(node.ParamNodes, k)
			case ast.ListNode:
				var hasRef bool
				var arr []interface{}
				for _, elem := range paramNode.Elems() {
					switch e := elem.(type) {
					case ast.InterfaceNode:
						arr = append(arr, e.Value())
					case ast.RefNode:
						hasRef = true
						arr = append(arr, e)
					case ast.ConcatenationNode:
						arr = append(arr, e.Concat())
					case ast.HoleNode, ast.AliasNode, ast.ListNode:
						return tpl, cenv, fmt.Errorf("%s: unresolved value in list of type %T", k, e)
					default:
						arr = append(arr, e)
					}
				}
				if hasRef {
					node.Refs[k] = ast.NewListNode(arr)
					delete(node.ParamNodes, k)
				} else {
					node.ParamNodes[k] = arr
				}
			case ast.ConcatenationNode:
				node.ParamNodes[k] = paramNode.Concat()
			case ast.HoleNode, ast.AliasNode:
				return tpl, cenv, fmt.Errorf("%s: unresolved value of type %T", k, paramNode)
			}
		}
	}
	return tpl, cenv, nil
}

func convertParamsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	convert := func(node *ast.CommandNode) error {
		for _, reducer := range node.ParamsSpec().Reducers() {
			params := make(map[string]interface{})
			for k, v := range node.ParamNodes {
				params[k] = v
			}
			for k, v := range node.Refs {
				params[k] = v
			}

			out, err := reducer.Reduce(params)
			if err != nil {
				return cmdErr(node, err)
			}
			for _, k := range reducer.Keys() {
				delete(node.ParamNodes, k)
				delete(node.Refs, k)
			}
			for k, v := range out {
				switch v.(type) {
				case ast.ListNode, ast.RefNode, ast.ConcatenationNode:
					node.Refs[k] = v
				default:
					node.ParamNodes[k] = v
				}
			}
		}
		return nil
	}
	err := tpl.visitCommandNodesE(convert)
	return tpl, cenv, err
}

func validateCommandsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	collectValidationErrs := func(node *ast.CommandNode) error {
		if err := params.Validate(node.ParamsSpec().Validators(), node.ParamNodes); err != nil {
			return cmdErr(node, err)
		}
		return nil
	}
	err := tpl.visitCommandNodesE(collectValidationErrs)
	return tpl, cenv, err
}

func checkInvalidReferenceDeclarationsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	return tpl, cenv, ast.VerifyRefs(tpl.AST)
}

func inlineVariableValuePass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	newTpl := &Template{ID: tpl.ID, AST: tpl.AST.Clone()}
	newTpl.Statements = []*ast.Statement{}

	for i, st := range tpl.Statements {
		decl, isDecl := st.Node.(*ast.DeclarationNode)
		if isDecl {
			if right, isRightExpr := decl.Expr.(*ast.RightExpressionNode); isRightExpr {
				if res := right.Result(); res != nil {
					cenv.Push(env.RESOLVED_VARS, map[string]interface{}{decl.Ident: res})
				}
				ast.ProcessRefs(
					&ast.AST{Statements: tpl.Statements[i+1:]},
					map[string]interface{}{decl.Ident: right.Node()},
				)
				continue
			}
		}

		newTpl.Statements = append(newTpl.Statements, st)
	}
	return newTpl, cenv, nil
}

func resolveHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	processed := ast.ProcessHoles(tpl.AST, cenv.Get(env.FILLERS))
	cenv.Push(env.PROCESSED_FILLERS, processed)

	return tpl, cenv, nil
}

func resolveMissingHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	uniqueHoles := ast.CollectUniqueHoles(tpl.AST)

	var sortedHoles []ast.HoleNode
	for hole := range uniqueHoles {
		sortedHoles = append(sortedHoles, hole)
	}
	sort.Slice(sortedHoles, func(i, j int) bool {
		a := sortedHoles[i]
		b := sortedHoles[j]

		if a.IsOptional() == b.IsOptional() {
			return a.Hole() < b.Hole()
		} else {
			if a.IsOptional() {
				return false
			}
			return true
		}
	})

	for _, hole := range sortedHoles {
		k := hole.Hole()
		if cenv.MissingHolesFunc() != nil {
			actual := cenv.MissingHolesFunc()(k, uniqueHoles[hole], hole.IsOptional())
			if actual == "" && hole.IsOptional() {
				continue
			}
			params, err := ParseParams(fmt.Sprintf("%s=%s", k, actual))
			if err != nil {
				if params, err = ParseParams(fmt.Sprintf("%s=%s", k, ast.Quote(actual))); err != nil {
					return tpl, cenv, err
				}
			}
			cenv.Push(env.FILLERS, map[string]interface{}{k: params[k]})
		}
	}

	processed := ast.ProcessHoles(tpl.AST, cenv.Get(env.FILLERS))
	cenv.Push(env.PROCESSED_FILLERS, processed)

	return tpl, cenv, nil
}

func removeOptionalHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	ast.RemoveOptionalHoles(tpl.AST)
	return tpl, cenv, nil
}

func resolveAliasPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	var emptyResolv []string
	resolvAliasFunc := func(action, entity string, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			if cenv.AliasFunc() == nil {
				return "", false
			}
			normalized := fmt.Sprintf("%s.%s.%s", action, entity, key)
			actual := cenv.AliasFunc()(normalized, alias)
			if actual == "" {
				emptyResolv = append(emptyResolv, alias)
				return "", false
			} else {
				cenv.Log().ExtraVerbosef("alias: resolved '%s' to '%s' for key %s", alias, actual, key)
				return actual, true
			}
		}
	}

	ast.ProcessAliases(tpl.AST, resolvAliasFunc)

	switch len(emptyResolv) {
	case 0:
		break
	case 1:
		return tpl, cenv, fmt.Errorf("cannot resolve alias \"%s\". Not found in locally synced data.", emptyResolv[0])
	default:
		return tpl, cenv, fmt.Errorf("cannot resolve aliases: %q. Not found in locally synced data.", emptyResolv)

	}

	return tpl, cenv, nil
}

func failOnUnresolvedHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	uniqueUnresolved := make(map[string]struct{})
	for _, hole := range ast.CollectHoles(tpl.AST) {
		uniqueUnresolved[hole.String()] = struct{}{}
	}

	var unresolved []string
	for k := range uniqueUnresolved {
		unresolved = append(unresolved, k)
	}

	if len(unresolved) > 0 {
		sort.Strings(unresolved)
		return tpl, cenv, fmt.Errorf("template contains unresolved holes: %v", unresolved)
	}

	return tpl, cenv, nil
}

func failOnUnresolvedAliasPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	uniqueUnresolved := make(map[string]struct{})

	for _, alias := range ast.CollectAliases(tpl.AST) {
		uniqueUnresolved[alias.String()] = struct{}{}
	}

	var unresolved []string
	for k := range uniqueUnresolved {
		unresolved = append(unresolved, k)
	}

	if len(unresolved) > 0 {
		sort.Strings(unresolved)
		return tpl, cenv, fmt.Errorf("template contains unresolved alias: %v", unresolved)
	}

	return tpl, cenv, nil
}

func cmdErr(cmd *ast.CommandNode, i interface{}, a ...interface{}) error {
	var prefix string
	if cmd != nil {
		prefix = fmt.Sprintf("%s %s: ", cmd.Action, cmd.Entity)
	}
	var msg string
	switch ii := i.(type) {
	case nil:
		return nil
	case string:
		msg = ii
	case error:
		msg = ii.Error()
	}
	if len(a) == 0 {
		return errors.New(prefix + msg)
	}
	return fmt.Errorf("%s"+msg, append([]interface{}{prefix}, a...)...)
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func excludeFromSlice(in []string, exclude []string) (out []string) {
	for _, v := range in {
		if !contains(exclude, v) {
			out = append(out, v)
		}
	}
	return out
}

package template

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/internal/ast"
	"github.com/wallix/awless/template/params"
)

type Mode []compileFunc

var (
	TestCompileMode = []compileFunc{
		resolveMetaPass,
		injectCommandsInNodesPass,
		failOnDeclarationWithNoResultPass,
		processAndValidateParamsPass,
		checkInvalidReferenceDeclarationsPass,
		resolveHolesPass,
		resolveMissingHolesPass,
		removeOptionalHolesPass,
		resolveAliasPass,
		inlineVariableValuePass,
	}

	NewRunnerCompileMode = append(TestCompileMode,
		failOnUnresolvedHolesPass,
		failOnUnresolvedAliasPass,
		convertParamsPass,
		validateCommandsPass,
	)
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

func resolveMetaPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	if cenv.LookupMetaCommandFunc() == nil {
		return tpl, cenv, nil
	}

	for _, node := range tpl.CommandNodesIterator() {
		meta := cenv.LookupMetaCommandFunc()(node.Action, node.Entity, node.Keys())
		if meta != nil {
			type R interface {
				Resolve(map[string]string) (*Template, error)
			}
			resolv, ok := meta.(R)
			if !ok {
				return tpl, cenv, errors.New("meta command can not be resolved")
			}
			paramsStr := make(map[string]string)
			for k, v := range node.Params {
				paramsStr[k] = v.String()
			}
			resolved, err := resolv.Resolve(paramsStr)
			if err != nil {
				return tpl, cenv, fmt.Errorf("%s %s: resolve meta command: %s", node.Action, node.Entity, err)
			}
			tpl.ReplaceNodeByTemplate(node, resolved)
		}
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
			fmt.Printf("%T, %#v", cmdNode.Command, cmdNode.Command)
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
			node.Params[e] = ast.NewHoleValue(normalized)
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
			node.Params[e] = ast.NewOptionalHoleValue(fmt.Sprintf("%s.%s", node.Entity, e))
		}
		return nil
	}

	err := tpl.visitCommandNodesE(normalizeMissingRequiredParamsAsHoleAndValidate)
	return tpl, cenv, err
}

func convertParamsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	convert := func(node *ast.CommandNode) error {
		refsParams := make(map[string]struct{})
		for _, reducer := range node.ParamsSpec().Reducers() {
			params := node.ToDriverParams()
			for k, v := range node.Params {
				if ref, isRef := v.(ast.WithRefs); isRef && len(ref.GetRefs()) > 0 {
					params[k] = v
					refsParams[k] = struct{}{}
				}
			}

			out, err := reducer.Reduce(params)
			if err != nil {
				return cmdErr(node, err)
			}
			for _, k := range reducer.Keys() {
				delete(node.Params, k)
			}
			for k, v := range out {
				if ref, isRef := v.(ast.CompositeValue); isRef {
					node.Params[k] = ref
				} else {
					node.Params[k] = ast.NewInterfaceValue(v)
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
		if err := params.Validate(node.ParamsSpec().Validators(), node.ToDriverParamsExcludingRefs()); err != nil {
			return cmdErr(node, err)
		}
		return nil
	}
	err := tpl.visitCommandNodesE(collectValidationErrs)
	return tpl, cenv, err
}

func checkInvalidReferenceDeclarationsPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	usedRefs := make(map[string]struct{})

	for _, withRef := range tpl.WithRefsIterator() {
		for _, ref := range withRef.GetRefs() {
			usedRefs[ref] = struct{}{}
		}
	}

	knownRefs := make(map[string]bool)

	var each = func(withRef ast.WithRefs) error {
		for _, ref := range withRef.GetRefs() {
			if _, ok := knownRefs[ref]; !ok {
				return fmt.Errorf("using reference '$%s' but '%s' is undefined in template\n", ref, ref)
			}
		}
		return nil
	}

	for _, st := range tpl.Statements {
		switch n := st.Node.(type) {
		case ast.WithRefs:
			if err := each(n); err != nil {
				return tpl, cenv, err
			}
		case *ast.DeclarationNode:
			expr := st.Node.(*ast.DeclarationNode).Expr
			switch nn := expr.(type) {
			case ast.WithRefs:
				if err := each(nn); err != nil {
					return tpl, cenv, err
				}
			}
		}
		if decl, isDecl := st.Node.(*ast.DeclarationNode); isDecl {
			ref := decl.Ident
			if _, ok := knownRefs[ref]; ok {
				return tpl, cenv, fmt.Errorf("using reference '$%s' but '%s' has already been assigned in template\n", ref, ref)
			}
			knownRefs[ref] = true
		}
	}

	return tpl, cenv, nil
}

func inlineVariableValuePass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	newTpl := &Template{ID: tpl.ID, AST: tpl.AST.Clone()}
	newTpl.Statements = []*ast.Statement{}

	for i, st := range tpl.Statements {
		decl, isDecl := st.Node.(*ast.DeclarationNode)
		if isDecl {
			value, isValue := decl.Expr.(*ast.ValueNode)
			if isValue {
				if val := value.Value.Value(); val != nil {
					cenv.Push(env.RESOLVED_VARS, map[string]interface{}{decl.Ident: val})
				}
				for j := i + 1; j < len(tpl.Statements); j++ {
					expr := extractExpressionNode(tpl.Statements[j])
					if expr != nil {
						if withRef, ok := expr.(ast.WithRefs); ok {
							withRef.ReplaceRef(decl.Ident, value.Value)
						}
					}
				}
				if value.IsResolved() {
					continue
				}
			}
		}
		newTpl.Statements = append(newTpl.Statements, st)
	}
	return newTpl, cenv, nil
}

func resolveHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	tpl.visitHoles(func(h ast.WithHoles) {
		processed := h.ProcessHoles(cenv.Get(env.FILLERS))
		cenv.Push(env.PROCESSED_FILLERS, processed)
	})

	return tpl, cenv, nil
}

func resolveMissingHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	uniqueHoles := make(map[string]*ast.Hole)
	tpl.visitHoles(func(h ast.WithHoles) {
		for k, v := range h.GetHoles() {
			if _, ok := uniqueHoles[k]; !ok {
				uniqueHoles[k] = v
			}
			for _, vv := range v.ParamPaths {
				if !contains(uniqueHoles[k].ParamPaths, vv) {
					uniqueHoles[k].ParamPaths = append(uniqueHoles[k].ParamPaths, vv)
				}
			}
		}
	})
	var sortedHoles []*ast.Hole
	for _, v := range uniqueHoles {
		sortedHoles = append(sortedHoles, v)
	}
	sort.Slice(sortedHoles, func(i, j int) bool {
		a := sortedHoles[i]
		b := sortedHoles[j]

		if a.IsOptional == b.IsOptional {
			return a.Name < b.Name
		} else {
			if a.IsOptional {
				return false
			}
			return true
		}
	})

	for _, hole := range sortedHoles {
		k := hole.Name
		if cenv.MissingHolesFunc() != nil {
			actual := cenv.MissingHolesFunc()(k, uniqueHoles[k].ParamPaths, uniqueHoles[k].IsOptional)
			if actual == "" && uniqueHoles[k].IsOptional {
				continue
			}
			params, err := ParseParams(fmt.Sprintf("%s=%s", k, actual))
			if err != nil {
				if params, err = ParseParams(fmt.Sprintf("%s=%s", k, quoteString(actual))); err != nil {
					return tpl, cenv, err
				}
			}
			cenv.Push(env.FILLERS, map[string]interface{}{k: params[k]})
		}
	}

	tpl.visitHoles(func(h ast.WithHoles) {
		processed := h.ProcessHoles(cenv.Get(env.FILLERS))
		cenv.Push(env.PROCESSED_FILLERS, processed)
	})

	return tpl, cenv, nil
}

func removeOptionalHolesPass(tpl *Template, cenv env.Compiling) (*Template, env.Compiling, error) {
	removeOptionalHoles := func(node *ast.CommandNode) error {
		for key, param := range node.Params {
			if param.Value() != nil {
				continue
			}
			if withHole, ok := param.(ast.WithHoles); ok {
				if len(withHole.GetHoles()) == 0 {
					continue
				}
				isOptional := true
				for _, h := range withHole.GetHoles() {
					isOptional = isOptional && h.IsOptional
				}
				if isOptional {
					delete(node.Params, key)
				}
			}
		}
		return nil
	}
	err := tpl.visitCommandNodesE(removeOptionalHoles)
	return tpl, cenv, err
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

	for _, expr := range tpl.expressionNodesIterator() {
		switch ee := expr.(type) {
		case *ast.CommandNode:
			for k, v := range ee.Params {
				if vv, ok := v.(ast.WithAlias); ok {
					vv.ResolveAlias(resolvAliasFunc(ee.Action, ee.Entity, k))
				}
			}
		case *ast.ValueNode:
			if vv, ok := ee.Value.(ast.WithAlias); ok {
				vv.ResolveAlias(resolvAliasFunc("", "", ""))
			}
		}
	}

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
	tpl.visitHoles(func(withHole ast.WithHoles) {
		for hole := range withHole.GetHoles() {
			uniqueUnresolved[hole] = struct{}{}
		}
	})

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
	var unresolved []string

	visitAliases := func(withAlias ast.WithAlias) {
		for _, alias := range withAlias.GetAliases() {
			unresolved = append(unresolved, alias)
		}
	}

	for _, n := range tpl.expressionNodesIterator() {
		switch nn := n.(type) {
		case *ast.ValueNode:
			if withAlias, ok := nn.Value.(ast.WithAlias); ok {
				visitAliases(withAlias)
			}
		case *ast.CommandNode:
			for _, param := range nn.Params {
				if withAlias, ok := param.(ast.WithAlias); ok {
					visitAliases(withAlias)
				}
			}
		}
	}

	if len(unresolved) > 0 {
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

func quoteString(str string) string {
	if strings.ContainsRune(str, '\'') {
		return "\"" + str + "\""
	} else {
		return "'" + str + "'"
	}
}

func excludeFromSlice(in []string, exclude []string) (out []string) {
	for _, v := range in {
		if !contains(exclude, v) {
			out = append(out, v)
		}
	}
	return out
}
func (t *Template) ReplaceNodeByTemplate(n ast.Node, tplToReplace *Template) error {
	nodeIndex := -1
	for i, st := range t.Statements {
		if st.Node == n {
			nodeIndex = i
		}
	}
	if nodeIndex == -1 {
		return fmt.Errorf("node '%v' not found", n)
	}
	after := make([]*ast.Statement, len(t.Statements[nodeIndex+1:]))
	copy(after, t.Statements[nodeIndex+1:])
	t.Statements = append(t.Statements[:nodeIndex], tplToReplace.Statements...)
	t.Statements = append(t.Statements, after...)
	return nil
}

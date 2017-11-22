package template

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/internal/ast"
)

type LookupFunc func(...string) interface{}

type Env struct {
	Lookuper LookupFunc
	IsDryRun bool

	ResolvedVariables map[string]interface{}

	Fillers          map[string]interface{}
	AliasFunc        func(entity, key, alias string) string
	MissingHolesFunc func(string, []string) interface{}
	Log              *logger.Logger

	processedFillers map[string]interface{}
}

func NewEnv() *Env {
	return &Env{
		AliasFunc:         nil,
		MissingHolesFunc:  nil,
		Lookuper:          func(...string) interface{} { return nil },
		Log:               logger.DiscardLogger,
		ResolvedVariables: make(map[string]interface{}),
		processedFillers:  make(map[string]interface{}),
	}
}

func (e *Env) AddFillers(fills ...map[string]interface{}) {
	if e.Fillers == nil {
		e.Fillers = make(map[string]interface{})
	}

	for _, f := range fills {
		for k, v := range f {
			e.Fillers[k] = v
		}
	}
}

func (e *Env) addToProcessedFillers(fills ...map[string]interface{}) {
	if e.processedFillers == nil {
		e.processedFillers = make(map[string]interface{})
	}

	for _, f := range fills {
		for k, v := range f {
			e.processedFillers[k] = v
		}
	}
}

func (e *Env) GetProcessedFillers() (copy map[string]interface{}) {
	copy = make(map[string]interface{}, 0)
	for k, v := range e.processedFillers {
		copy[k] = v
	}
	return
}

type Mode []compileFunc

var (
	TestCompileMode = []compileFunc{
		verifyCommandsDefinedPass,
		failOnDeclarationWithNoResultPass,
		validateCommandsParamsPass,
		normalizeMissingRequiredParamsAsHolePass,
		checkInvalidReferenceDeclarationsPass,
		resolveHolesPass,
		resolveMissingHolesPass,
		resolveAliasPass,
		inlineVariableValuePass,
	}

	NewRunnerCompileMode = []compileFunc{
		verifyCommandsDefinedPass,
		failOnDeclarationWithNoResultPass,
		validateCommandsParamsPass,
		normalizeMissingRequiredParamsAsHolePass,
		checkInvalidReferenceDeclarationsPass,
		resolveHolesPass,
		resolveMissingHolesPass,
		resolveAliasPass,
		inlineVariableValuePass,
		failOnUnresolvedHolesPass,
		failOnUnresolvedAliasPass,
		convertParamsPass,
		validateCommandsPass,
		injectCommandsPass,
	}
)

func Compile(tpl *Template, env *Env, mode ...Mode) (*Template, *Env, error) {
	var pass *multiPass

	if len(mode) > 0 {
		pass = newMultiPass(mode[0]...)
	} else {
		pass = newMultiPass(NewRunnerCompileMode...)
	}

	return pass.compile(tpl, env)
}

type compileFunc func(*Template, *Env) (*Template, *Env, error)

// Leeloo Dallas
type multiPass struct {
	passes []compileFunc
}

func newMultiPass(passes ...compileFunc) *multiPass {
	return &multiPass{passes: passes}
}

func (p *multiPass) compile(tpl *Template, env *Env) (newTpl *Template, newEnv *Env, err error) {
	newTpl, newEnv = tpl, env
	for _, pass := range p.passes {
		newTpl, newEnv, err = pass(newTpl, newEnv)
		if err != nil {
			return
		}
	}

	return
}

func verifyCommandsDefinedPass(tpl *Template, env *Env) (*Template, *Env, error) {
	if env.Lookuper == nil {
		return tpl, env, fmt.Errorf("command lookuper is undefined")
	}

	for _, node := range tpl.CommandNodesIterator() {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return tpl, env, fmt.Errorf("cannot find command for '%s'", key)
		}
	}
	return tpl, env, nil
}

func failOnDeclarationWithNoResultPass(tpl *Template, env *Env) (*Template, *Env, error) {
	failOnDeclarationWithNoResult := func(node *ast.DeclarationNode) error {
		cmdNode, ok := node.Expr.(*ast.CommandNode)
		if !ok {
			return nil
		}
		key := fmt.Sprintf("%s%s", cmdNode.Action, cmdNode.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return fmt.Errorf("validate: cannot find command for '%s'", key)
		}
		type ER interface {
			ExtractResult(interface{}) string
		}
		if _, ok := cmd.(ER); !ok {
			return cmdErr(cmdNode, "command does not return a result, cannot assign to a variable")
		}
		return nil
	}

	for _, dcl := range tpl.declarationNodesIterator() {
		if err := failOnDeclarationWithNoResult(dcl); err != nil {
			return tpl, env, err
		}
	}
	return tpl, env, nil
}

func validateCommandsParamsPass(tpl *Template, env *Env) (*Template, *Env, error) {
	verifyValidParamsOnly := func(node *ast.CommandNode) error {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return fmt.Errorf("validate: cannot find command for '%s'", key)
		}
		type VP interface {
			ValidateParams([]string) ([]string, error)
		}
		if v, ok := cmd.(VP); ok {
			if _, err := v.ValidateParams(node.Keys()); err != nil {
				return cmdErr(node, err)
			}
		} else {
			return cmdErr(node, "command does not implement param validation")
		}
		return nil
	}

	err := tpl.visitCommandNodesE(verifyValidParamsOnly)
	return tpl, env, err
}

func normalizeMissingRequiredParamsAsHolePass(tpl *Template, env *Env) (*Template, *Env, error) {
	normalize := func(node *ast.CommandNode) error {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return fmt.Errorf("normalize: cannot find command for '%s'", key)
		}
		type VP interface {
			ValidateParams([]string) ([]string, error)
		}
		if v, ok := cmd.(VP); ok {
			missing, err := v.ValidateParams(node.Keys())
			if err != nil {
				return cmdErr(node, err)
			}
			for _, e := range missing {
				normalized := fmt.Sprintf("%s.%s", node.Entity, e)
				node.Params[e] = ast.NewHoleValue(normalized)
			}
		} else {
			return cmdErr(node, "command does not implement param normalization")
		}
		return nil
	}

	err := tpl.visitCommandNodesE(normalize)
	return tpl, env, err
}

func convertParamsPass(tpl *Template, env *Env) (*Template, *Env, error) {
	convert := func(node *ast.CommandNode) error {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return fmt.Errorf("convert: cannot find command for '%s'", key)
		}

		type C interface {
			ConvertParams() ([]string, func(values map[string]interface{}) (map[string]interface{}, error))
		}
		if v, ok := cmd.(C); ok {
			keys, convFunc := v.ConvertParams()
			values := make(map[string]interface{})
			params := node.ToDriverParams()
			for _, k := range keys {
				if vv, ok := params[k]; ok {
					values[k] = vv
				}
			}
			converted, err := convFunc(values)
			if err != nil {
				return cmdErr(node, err)
			}
			for _, k := range keys {
				delete(node.Params, k)
			}
			for k, v := range converted {
				node.Params[k] = ast.NewInterfaceValue(v)
			}
		}
		return nil
	}
	err := tpl.visitCommandNodesE(convert)
	return tpl, env, err
}

func validateCommandsPass(tpl *Template, env *Env) (*Template, *Env, error) {
	var errs []error

	collectValidationErrs := func(node *ast.CommandNode) error {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		cmd := env.Lookuper(key)
		if cmd == nil {
			return fmt.Errorf("validate: cannot find command for '%s'", key)
		}
		type V interface {
			ValidateCommand(map[string]interface{}, []string) []error
		}
		if v, ok := cmd.(V); ok {
			var refsKey []string
			for k, p := range node.Params {
				if ref, isRef := p.(ast.WithRefs); isRef && len(ref.GetRefs()) > 0 {
					refsKey = append(refsKey, k)
				}
			}
			for _, validErr := range v.ValidateCommand(node.ToDriverParams(), refsKey) {
				errs = append(errs, fmt.Errorf("%s %s: %s", node.Action, node.Entity, validErr.Error()))
			}
		}
		return nil
	}
	if err := tpl.visitCommandNodesE(collectValidationErrs); err != nil {
		return tpl, env, err
	}
	switch len(errs) {
	case 0:
		return tpl, env, nil
	case 1:
		return tpl, env, fmt.Errorf("validation error: %s", errs[0])
	default:
		var errsSrings []string
		for _, err := range errs {
			if err != nil {
				errsSrings = append(errsSrings, err.Error())
			}
		}
		return tpl, env, fmt.Errorf("validation errors:\n\t- %s", strings.Join(errsSrings, "\n\t- "))
	}
}

func injectCommandsPass(tpl *Template, env *Env) (*Template, *Env, error) {
	for _, node := range tpl.CommandNodesIterator() {
		key := fmt.Sprintf("%s%s", node.Action, node.Entity)
		node.Command = env.Lookuper(key).(ast.Command)
		if node.Command == nil {
			return tpl, env, fmt.Errorf("inject: cannot find command for '%s'", key)
		}
	}
	return tpl, env, nil
}

func checkInvalidReferenceDeclarationsPass(tpl *Template, env *Env) (*Template, *Env, error) {
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
				return tpl, env, err
			}
		case *ast.DeclarationNode:
			expr := st.Node.(*ast.DeclarationNode).Expr
			switch nn := expr.(type) {
			case ast.WithRefs:
				if err := each(nn); err != nil {
					return tpl, env, err
				}
			}
		}
		if decl, isDecl := st.Node.(*ast.DeclarationNode); isDecl {
			ref := decl.Ident
			if _, ok := knownRefs[ref]; ok {
				return tpl, env, fmt.Errorf("using reference '$%s' but '%s' has already been assigned in template\n", ref, ref)
			}
			knownRefs[ref] = true
		}
	}

	return tpl, env, nil
}

func inlineVariableValuePass(tpl *Template, env *Env) (*Template, *Env, error) {
	newTpl := &Template{ID: tpl.ID, AST: tpl.AST.Clone()}
	newTpl.Statements = []*ast.Statement{}

	for i, st := range tpl.Statements {
		decl, isDecl := st.Node.(*ast.DeclarationNode)
		if isDecl {
			value, isValue := decl.Expr.(*ast.ValueNode)
			if isValue {
				if val := value.Value.Value(); val != nil {
					env.ResolvedVariables[decl.Ident] = val
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
	return newTpl, env, nil
}

func resolveHolesPass(tpl *Template, env *Env) (*Template, *Env, error) {
	tpl.visitHoles(func(h ast.WithHoles) {
		processed := h.ProcessHoles(env.Fillers)
		env.addToProcessedFillers(processed)
	})

	return tpl, env, nil
}

func resolveMissingHolesPass(tpl *Template, env *Env) (*Template, *Env, error) {
	uniqueHoles := make(map[string][]string)
	tpl.visitHoles(func(h ast.WithHoles) {
		for k, v := range h.GetHoles() {
			uniqueHoles[k] = nil
			for _, vv := range v {
				if !contains(uniqueHoles[k], vv) {
					uniqueHoles[k] = append(uniqueHoles[k], vv)
				}
			}
		}
	})
	var sortedHoles []string
	for k := range uniqueHoles {
		sortedHoles = append(sortedHoles, k)
	}
	sort.Strings(sortedHoles)
	fillers := make(map[string]interface{})

	for _, k := range sortedHoles {
		if env.MissingHolesFunc != nil {
			actual := env.MissingHolesFunc(k, uniqueHoles[k])
			fillers[k] = actual
		}
	}

	tpl.visitHoles(func(h ast.WithHoles) {
		processed := h.ProcessHoles(fillers)
		env.addToProcessedFillers(processed)
	})

	return tpl, env, nil
}

func resolveAliasPass(tpl *Template, env *Env) (*Template, *Env, error) {
	var emptyResolv []string
	resolvAliasFunc := func(entity string, key string) func(string) (string, bool) {
		return func(alias string) (string, bool) {
			if env.AliasFunc == nil {
				return "", false
			}
			actual := env.AliasFunc(entity, key, alias)
			if actual == "" {
				emptyResolv = append(emptyResolv, alias)
				return "", false
			} else {
				env.Log.ExtraVerbosef("alias: resolved '%s' to '%s' for key %s", alias, actual, key)
				return actual, true
			}
		}
	}

	for _, expr := range tpl.expressionNodesIterator() {
		switch ee := expr.(type) {
		case *ast.CommandNode:
			for k, v := range ee.Params {
				if vv, ok := v.(ast.WithAlias); ok {
					vv.ResolveAlias(resolvAliasFunc(ee.Entity, k))
				}
			}
		case *ast.ValueNode:
			if vv, ok := ee.Value.(ast.WithAlias); ok {
				vv.ResolveAlias(resolvAliasFunc("", ""))
			}
		}
	}

	if len(emptyResolv) > 0 {
		return tpl, env, fmt.Errorf("cannot resolve aliases: %q. Maybe you need to update your local model with `awless sync` ?", emptyResolv)
	}

	return tpl, env, nil
}

func failOnUnresolvedHolesPass(tpl *Template, env *Env) (*Template, *Env, error) {
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
		return tpl, env, fmt.Errorf("template contains unresolved holes: %v", unresolved)
	}

	return tpl, env, nil
}

func failOnUnresolvedAliasPass(tpl *Template, env *Env) (*Template, *Env, error) {
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
		return tpl, env, fmt.Errorf("template contains unresolved alias: %v", unresolved)
	}

	return tpl, env, nil
}

func foundIn(key string, slice []string) (found bool) {
	for _, k := range slice {
		if k == key {
			found = true
			break
		}
	}
	return
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

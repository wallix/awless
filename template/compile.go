package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/internal/ast"
)

type Env struct {
	Fillers map[string]interface{}

	Resolved         map[string]interface{}
	DefLookupFunc    DefinitionLookupFunc
	AliasFunc        func(entity, key, alias string) string
	MissingHolesFunc func(string) interface{}

	Log *logger.Logger
}

func NewEnv() *Env {
	return &Env{
		Resolved:         make(map[string]interface{}),
		AliasFunc:        func(e, k, v string) string { return v },
		MissingHolesFunc: func(s string) interface{} { return "" },
		Log:              logger.DiscardLogger,
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

func Compile(tpl *Template, env *Env) (*Template, *Env, error) {
	pass := newMultiPass(
		resolveAgainstDefinitions,
		checkReferencesDeclaration,
		resolveHolesPass,
		resolveMissingHolesPass,
		resolveAliasPass,
	)

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

func resolveAgainstDefinitions(tpl *Template, env *Env) (*Template, *Env, error) {
	each := func(cmd *ast.CommandNode) error {
		key := fmt.Sprintf("%s%s", cmd.Action, cmd.Entity)
		def, ok := env.DefLookupFunc(key)
		if !ok {
			return fmt.Errorf("cannot find template definition for '%s'", key)
		}

		for _, key := range cmd.Keys() {
			var found bool

			for _, k := range def.Required() {
				if k == key {
					found = true
					break
				}
			}

			for _, k := range def.Extra() {
				if k == key {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%s %s: unexpected param key '%s'\n\t- required params: %s\n\t- extra params: %s\n", cmd.Action, cmd.Entity, key, strings.Join(def.Required(), ", "), strings.Join(def.Extra(), ", "))
			}
		}

		return nil
	}

	if err := tpl.visitCommandNodesE(each); err != nil {
		return tpl, env, err
	}

	tpl.visitCommandNodes(func(cmd *ast.CommandNode) {
		if cmd.Holes == nil {
			cmd.Holes = make(map[string]string)
		}
		key := fmt.Sprintf("%s%s", cmd.Action, cmd.Entity)
		def, _ := env.DefLookupFunc(key)
		for _, required := range def.Required() {
			var isInParams bool
			var isInRefs bool

			for k := range cmd.Params {
				if k == required {
					isInParams = true
				}
			}
			for k := range cmd.Refs {
				if k == required {
					isInRefs = true
				}
			}
			normalized := fmt.Sprintf("%s.%s", cmd.Entity, required)

			if isInParams || isInRefs {
				delete(cmd.Holes, normalized)
				continue
			} else {
				if _, ok := cmd.Holes[required]; !ok {
					cmd.Holes[required] = normalized
				}
			}
		}
	})

	return tpl, env, nil
}

func checkReferencesDeclaration(tpl *Template, env *Env) (*Template, *Env, error) {
	usedRefs := make(map[string]struct{})
	tpl.visitCommandNodes(func(cmd *ast.CommandNode) {
		for _, v := range cmd.Refs {
			usedRefs[v] = struct{}{}
		}
	})

	declRefs := make(map[string]struct{})
	tpl.visitCommandDeclarationNodes(func(decl *ast.DeclarationNode) {
		declRefs[decl.Ident] = struct{}{}
	})

	for r := range usedRefs {
		if _, ok := declRefs[r]; !ok {
			return tpl, env, fmt.Errorf("using reference '$%s' but '%s' is undefined in template\n", r, r)
		}
	}

	for r := range declRefs {
		if _, ok := usedRefs[r]; !ok {
			return tpl, env, fmt.Errorf("unused reference '%s' in template\n", r)
		}
	}

	return tpl, env, nil
}

func resolveHolesPass(tpl *Template, env *Env) (*Template, *Env, error) {
	if env.Resolved == nil {
		env.Resolved = make(map[string]interface{})
	}

	each := func(cmd *ast.CommandNode) {
		processed := cmd.ProcessHoles(env.Fillers)
		for key, v := range processed {
			env.Resolved[cmd.Entity+"."+key] = v
		}
	}

	tpl.visitCommandNodes(each)

	env.Log.ExtraVerbosef("holes resolved: %v", env.Resolved)

	return tpl, env, nil
}

func resolveMissingHolesPass(tpl *Template, env *Env) (*Template, *Env, error) {
	uniqueHoles := make(map[string]struct{})
	tpl.visitCommandNodes(func(cmd *ast.CommandNode) {
		for _, v := range cmd.Holes {
			uniqueHoles[v] = struct{}{}
		}
	})
	var sortedHoles []string
	for k := range uniqueHoles {
		sortedHoles = append(sortedHoles, k)
	}
	sort.Strings(sortedHoles)

	fillers := make(map[string]interface{})
	for _, k := range sortedHoles {
		actual := env.MissingHolesFunc(k)
		fillers[k] = actual
	}

	tpl.visitCommandNodes(func(expr *ast.CommandNode) {
		expr.ProcessHoles(fillers)
	})

	return tpl, env, nil
}

func resolveAliasPass(tpl *Template, env *Env) (*Template, *Env, error) {
	var unresolved []string
	each := func(cmd *ast.CommandNode) {
		for k, v := range cmd.Params {
			if s, ok := v.(string); ok {
				if strings.HasPrefix(s, "@") {
					env.Log.ExtraVerbosef("alias resolving: %s for key %s", s, k)
					alias := strings.TrimPrefix(s, "@")
					actual := env.AliasFunc(cmd.Entity, k, alias)
					if actual == "" {
						unresolved = append(unresolved, alias)
					} else {
						env.Log.ExtraVerbosef("alias '%s' resolved to '%s' for key %s", alias, actual, k)
						cmd.Params[k] = actual
						delete(cmd.Holes, k)
					}
				}
			}
		}
	}

	tpl.visitCommandNodes(each)

	if len(unresolved) > 0 {
		return tpl, env, fmt.Errorf("cannot resolve aliases: %q", unresolved)
	}

	return tpl, env, nil
}

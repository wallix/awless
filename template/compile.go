package template

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/ast"
)

type Env struct {
	Fillers map[string]interface{}

	Resolved         map[string]interface{}
	DefLookupFunc    LookupTemplateDefFunc
	AliasFunc        func(key, alias string) string
	MissingHolesFunc func(string) interface{}

	Log *logger.Logger
}

func NewEnv() *Env {
	return &Env{
		Resolved:         make(map[string]interface{}),
		AliasFunc:        func(k, v string) string { return v },
		MissingHolesFunc: func(s string) interface{} { return s },
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
		resolveHolesPass,
		resolveAliasPass,
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

		for p, _ := range cmd.Params {
			var found bool

			for _, k := range def.Required() {
				if k == p {
					found = true
					break
				}
			}

			for _, k := range def.Extra() {
				if k == p {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%s %s: unexpected param '%s'", cmd.Action, cmd.Entity, p)
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

			for k, _ := range cmd.Params {
				if k == required {
					isInParams = true
				}
			}
			for k, _ := range cmd.Refs {
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

	fillers := make(map[string]interface{})
	for k := range uniqueHoles {
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
					actual := env.AliasFunc(k, alias)
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

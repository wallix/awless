package template

import (
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/env"
)

var (
	_ env.Running   = (*runEnv)(nil)
	_ env.Compiling = (*compileEnv)(nil)
)

type runEnv struct {
	log    *logger.Logger
	dryRun bool
	ctx    map[string]interface{}
}

func NewRunEnv(cenv env.Compiling, context ...map[string]interface{}) env.Running {
	renv := new(runEnv)
	renv.log = cenv.Log()
	renv.ctx = make(map[string]interface{})
	for _, m := range context {
		for k, v := range m {
			renv.ctx[k] = v
		}
	}
	renv.ctx["Variables"] = cenv.ResolvedVariables()
	renv.ctx["References"] = cenv.ResolvedVariables() // retro-compatibility with v0.1.2

	return renv
}

func (e *runEnv) IsDryRun() bool {
	return e.dryRun
}

func (e *runEnv) SetDryRun(b bool) {
	e.dryRun = b
}

func (e *runEnv) Context() (out map[string]interface{}) {
	out = make(map[string]interface{})
	for k, v := range e.ctx {
		out[k] = v
	}
	return
}

func (e *runEnv) Log() *logger.Logger {
	return e.log
}

type compileEnv struct {
	fillers, processedFillers, resolvedVariables map[string]interface{}
	lookupCommandFunc                            func(...string) interface{}
	aliasFunc                                    func(entity, key, alias string) string
	missingHolesFunc                             func(string, []string) interface{}
	log                                          *logger.Logger
}

func (e *compileEnv) Log() *logger.Logger {
	return e.log
}

func (e *compileEnv) ResolvedVariables() (out map[string]interface{}) {
	out = make(map[string]interface{})
	for k, v := range e.resolvedVariables {
		out[k] = v
	}
	return
}

func (e *compileEnv) addResolvedVariables(k string, i interface{}) {
	if e.resolvedVariables == nil {
		e.resolvedVariables = make(map[string]interface{})
	}
	e.resolvedVariables[k] = i
}

func (e *compileEnv) Fillers() (out map[string]interface{}) {
	out = make(map[string]interface{})
	for k, v := range e.fillers {
		out[k] = v
	}
	return
}

func (e *compileEnv) ProcessedFillers() (copy map[string]interface{}) {
	copy = make(map[string]interface{}, 0)
	for k, v := range e.processedFillers {
		copy[k] = v
	}
	return
}

func (e *compileEnv) addToProcessedFillers(fills ...map[string]interface{}) {
	if e.processedFillers == nil {
		e.processedFillers = make(map[string]interface{})
	}

	for _, f := range fills {
		for k, v := range f {
			e.processedFillers[k] = v
		}
	}
}

func (e *compileEnv) LookupCommandFunc() func(...string) interface{} {
	return e.lookupCommandFunc
}

func (e *compileEnv) AliasFunc() func(entity, key, alias string) string {
	return e.aliasFunc
}

func (e *compileEnv) MissingHolesFunc() func(string, []string) interface{} {
	return e.missingHolesFunc
}

func NewEnv() *envBuilder {
	b := &envBuilder{new(compileEnv)}
	b.E.lookupCommandFunc = func(...string) interface{} { return nil }
	b.E.log = logger.DiscardLogger
	b.E.fillers = make(map[string]interface{})
	b.E.processedFillers = make(map[string]interface{})
	b.E.resolvedVariables = make(map[string]interface{})
	return b
}

type envBuilder struct {
	E *compileEnv
}

func (b *envBuilder) WithAliasFunc(fn func(entity, key, alias string) string) *envBuilder {
	b.E.aliasFunc = fn
	return b
}

func (b *envBuilder) WithMissingHolesFunc(fn func(string, []string) interface{}) *envBuilder {
	b.E.missingHolesFunc = fn
	return b
}

func (b *envBuilder) WithLookupCommandFunc(fn func(...string) interface{}) *envBuilder {
	b.E.lookupCommandFunc = fn
	return b
}

func (b *envBuilder) WithFillers(maps ...map[string]interface{}) *envBuilder {
	b.E.fillers = make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			b.E.fillers[k] = v
		}
	}
	return b
}

func (b *envBuilder) WithLog(l *logger.Logger) *envBuilder {
	b.E.log = l
	return b
}

func (b *envBuilder) Build() env.Compiling {
	return b.E
}

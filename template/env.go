package template

import (
	"sync"

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

	renv.ctx["AWLESS"] = cenv.Get(env.RESOLVED_VARS)
	renv.ctx["Variables"] = cenv.Get(env.RESOLVED_VARS)  // retro-compatibility with > v0.1.9
	renv.ctx["References"] = cenv.Get(env.RESOLVED_VARS) // retro-compatibility with v0.1.2

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
	*dataMap
	lookupCommandFunc func(...string) interface{}
	aliasFunc         func(paramPath, alias string) string
	missingHolesFunc  func(string, []string, bool) string
	log               *logger.Logger
	paramsSuggested   int
}

func (e *compileEnv) LookupCommandFunc() func(...string) interface{} {
	return e.lookupCommandFunc
}

func (e *compileEnv) AliasFunc() func(paramPath, alias string) string {
	return e.aliasFunc
}

func (e *compileEnv) MissingHolesFunc() func(string, []string, bool) string {
	return e.missingHolesFunc
}

func (e *compileEnv) ParamsMode() int {
	return e.paramsSuggested
}

func (e *compileEnv) Log() *logger.Logger {
	return e.log
}

type noopCompileEnv struct{}

func (*noopCompileEnv) LookupCommandFunc() func(...string) interface{}        { return nil }
func (*noopCompileEnv) AliasFunc() func(paramPath, alias string) string       { return nil }
func (*noopCompileEnv) MissingHolesFunc() func(string, []string, bool) string { return nil }
func (*noopCompileEnv) ParamsMode() int                                       { return -1 }
func (*noopCompileEnv) Log() *logger.Logger                                   { return logger.DiscardLogger }
func (*noopCompileEnv) Push(int, ...map[string]interface{})                   {}
func (*noopCompileEnv) Get(int) map[string]interface{}                        { return make(map[string]interface{}) }

func NewEnv() *envBuilder {
	b := &envBuilder{new(compileEnv)}
	b.E.lookupCommandFunc = func(...string) interface{} { return nil }
	b.E.log = logger.DiscardLogger
	b.E.dataMap = new(dataMap)
	return b
}

type dataMap struct {
	mu sync.Mutex
	M  map[int]map[string]interface{}
}

func (d *dataMap) Push(typ int, data ...map[string]interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.M == nil {
		d.M = make(map[int]map[string]interface{})
	}
	if d.M[typ] == nil {
		d.M[typ] = make(map[string]interface{})
	}
	for _, m := range data {
		for k, v := range m {
			d.M[typ][k] = v
		}
	}
}

func (d *dataMap) Get(typ int) (out map[string]interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	out = make(map[string]interface{})
	if d.M[typ] == nil {
		return
	}
	for k, v := range d.M[typ] {
		out[k] = v
	}
	return
}

type envBuilder struct {
	E *compileEnv
}

func (b *envBuilder) WithAliasFunc(fn func(paramPath, alias string) string) *envBuilder {
	b.E.aliasFunc = fn
	return b
}

func (b *envBuilder) WithMissingHolesFunc(fn func(string, []string, bool) string) *envBuilder {
	b.E.missingHolesFunc = fn
	return b
}

func (b *envBuilder) WithLookupCommandFunc(fn func(...string) interface{}) *envBuilder {
	b.E.lookupCommandFunc = fn
	return b
}

func (b *envBuilder) WithLog(l *logger.Logger) *envBuilder {
	b.E.log = l
	return b
}

func (b *envBuilder) WithParamsMode(paramsSuggested int) *envBuilder {
	b.E.paramsSuggested = paramsSuggested
	return b
}

func (b *envBuilder) Build() env.Compiling {
	return b.E
}

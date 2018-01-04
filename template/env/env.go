package env

import (
	"github.com/wallix/awless/logger"
)

const (
	FILLERS = iota
	PROCESSED_FILLERS
	RESOLVED_VARS
)

const (
	DEFAULTS_SUGGESTED = iota
	NO_SUGGESTED
	ALL_SUGGESTED
)

type log interface {
	Log() *logger.Logger
}

type Running interface {
	log
	Context() map[string]interface{}
	IsDryRun() bool
	SetDryRun(b bool)
}

type Compiling interface {
	log
	LookupCommandFunc() func(...string) interface{}
	AliasFunc() func(paramPath, alias string) string
	MissingHolesFunc() func(string, []string, bool) string
	ParamsSuggested() int
	Push(int, ...map[string]interface{})
	Get(int) map[string]interface{}
}

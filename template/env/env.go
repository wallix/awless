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
	REQUIRED_AND_SUGGESTED_PARAMS = iota
	REQUIRED_PARAMS_ONLY
	ALL_PARAMS
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
	ParamsMode() int
	Push(int, ...map[string]interface{})
	Get(int) map[string]interface{}
}

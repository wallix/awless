package env

import (
	"github.com/wallix/awless/logger"
)

const (
	FILLERS = iota
	PROCESSED_FILLERS
	RESOLVED_VARS
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
	AliasFunc() func(entity, key, alias string) string
	MissingHolesFunc() func(string, []string) interface{}
	Push(int, ...map[string]interface{})
	Get(int) map[string]interface{}
}

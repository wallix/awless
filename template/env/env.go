package env

import "github.com/wallix/awless/logger"

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
	Fillers() map[string]interface{}
	ProcessedFillers() map[string]interface{}
	ResolvedVariables() map[string]interface{}
	MissingHolesFunc() func(string, []string) interface{}
}

package env

type Running interface {
	Context() map[string]interface{}
	IsDryRun() bool
}

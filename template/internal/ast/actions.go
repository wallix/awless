package ast

type Action string

const (
	UnknownAction Action = "unknown"
	NoneAction    Action = "none"

	Create Action = "create"
	Delete Action = "delete"
	Update Action = "update"

	Check Action = "check"

	Start   Action = "start"
	Restart Action = "restart"
	Stop    Action = "stop"

	Attach Action = "attach"
	Detach Action = "detach"

	Copy Action = "copy"

	Import       Action = "import"
	Authenticate Action = "authenticate"
)

var actions = map[Action]struct{}{
	NoneAction:   {},
	Create:       {},
	Delete:       {},
	Update:       {},
	Check:        {},
	Start:        {},
	Restart:      {},
	Stop:         {},
	Attach:       {},
	Detach:       {},
	Copy:         {},
	Import:       {},
	Authenticate: {},
}

func IsInvalidAction(s string) bool {
	_, ok := actions[Action(s)]
	return !ok
}

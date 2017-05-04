package ast

type Action string

const (
	UnknownAction Action = "unknown"
	NoneAction    Action = "none"

	Create Action = "create"
	Delete Action = "delete"
	Update Action = "update"

	Check Action = "check"

	Start Action = "start"
	Stop  Action = "stop"

	Attach Action = "attach"
	Detach Action = "detach"

	Copy Action = "copy"

	Import Action = "import"
)

var actions = map[Action]struct{}{
	NoneAction: {},
	Create:     {},
	Delete:     {},
	Update:     {},
	Check:      {},
	Start:      {},
	Stop:       {},
	Attach:     {},
	Detach:     {},
	Copy:       {},
	Import:     {},
}

func IsInvalidAction(s string) bool {
	_, ok := actions[Action(s)]
	return !ok
}

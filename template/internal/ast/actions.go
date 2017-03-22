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
)

var actions = map[Action]struct{}{
	NoneAction: struct{}{},
	Create:     struct{}{},
	Delete:     struct{}{},
	Update:     struct{}{},
	Check:      struct{}{},
	Start:      struct{}{},
	Stop:       struct{}{},
	Attach:     struct{}{},
	Detach:     struct{}{},
}

func IsInvalidAction(s string) bool {
	_, ok := actions[Action(s)]
	return !ok
}

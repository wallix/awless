package template

import (
	"fmt"

	"github.com/wallix/awless/graph"
)

type Validator interface {
	Execute(t *Template) []error
}

type LookupGraphFunc func(key string) (*graph.Graph, bool)

type UniqueNameValidator struct {
	LookupGraph LookupGraphFunc
}

func (v *UniqueNameValidator) Execute(t *Template) (errs []error) {
	for _, cmd := range t.CommandNodesIterator() {
		if cmd.Action == "create" {
			name := cmd.Params["name"]
			g, ok := v.LookupGraph(cmd.Entity)
			if !ok {
				continue
			}
			resources, err := g.FindResourcesByProperty("Name", name)
			if err != nil {
				errs = append(errs, err)
			}
			if len(resources) > 0 {
				errs = append(errs, fmt.Errorf("%s %s: name '%s' already exists", cmd.Action, cmd.Entity, name))
			}
		}
	}
	return
}

type ParamIsSetValidator struct {
	Entity, Action, Param, WarningMessage string
}

func (v *ParamIsSetValidator) Execute(t *Template) (errs []error) {
	for _, cmd := range t.CommandNodesIterator() {
		if cmd.Action == v.Action && cmd.Entity == v.Entity {
			_, hasParam := cmd.Params[v.Param]
			_, hasRef := cmd.Refs[v.Param]
			if !hasParam && !hasRef {
				errs = append(errs, fmt.Errorf(v.WarningMessage))
			}
		}
	}
	return
}

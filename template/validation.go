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
				errs = append(errs, fmt.Errorf("%s %s: name '%s' already exists\n", cmd.Action, cmd.Entity, name))
			}
		}
	}
	return
}

func sliceContains(s string, arrs ...[]string) bool {
	for _, arr := range arrs {
		for _, el := range arr {
			if el == s {
				return true
			}
		}
	}

	return false
}

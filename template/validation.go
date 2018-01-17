package template

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/wallix/awless/cloud"
)

type Validator interface {
	Execute(t *Template) []error
}

type LookupGraphFunc func(key string) (cloud.GraphAPI, bool)

type UniqueNameValidator struct {
	LookupGraph LookupGraphFunc
}

func (v *UniqueNameValidator) Execute(t *Template) (errs []error) {
	for _, cmd := range t.CommandNodesIterator() {
		if cmd.Action == "create" {
			name := cmd.ParamNodes["name"]
			g, ok := v.LookupGraph(cmd.Entity)
			if !ok {
				continue
			}
			resources, err := g.FindWithProperties(map[string]interface{}{"Name": name})
			if err != nil {
				errs = append(errs, err)
			}
			if len(resources) > 0 {
				for _, r := range resources {
					var buf bytes.Buffer
					buf.WriteString(fmt.Sprintf("'%s' name already used for %s %s", name, r.Type(), r.Id()))
					if state, ok := r.Properties()["State"].(string); ok {
						buf.WriteString(fmt.Sprintf(" (state: '%s')", state))
					}
					errs = append(errs, errors.New(buf.String()))
				}
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
			_, hasParam := cmd.ParamNodes[v.Param]
			if !hasParam {
				errs = append(errs, fmt.Errorf(v.WarningMessage))
			}
		}
	}
	return
}

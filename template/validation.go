package template

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/wallix/awless/graph"
)

type Validator interface {
	Execute(t *Template) []error
}

type DefinitionValidator struct {
	LookupDef LookupTemplateDefFunc
}

func (v *DefinitionValidator) Execute(t *Template) (errs []error) {
	for _, cmd := range t.CommandNodesIterator() {
		key := fmt.Sprintf("%s%s", cmd.Action, cmd.Entity)
		def, ok := v.LookupDef(key)
		if !ok {
			continue
		}

		var unexpected []string
		for p := range cmd.Params {
			if !sliceContains(p, def.Required(), def.Extra()) {
				unexpected = append(unexpected, fmt.Sprintf("'%s'", p))
			}
		}

		if len(unexpected) > 0 {
			var w bytes.Buffer
			w.WriteString(fmt.Sprintf("%s %s: unexpected params %s", cmd.Action, cmd.Entity, strings.Join(unexpected, ", ")))

			if len(def.Required()) > 0 {
				w.WriteString(fmt.Sprintf("\n\trequired: %s", strings.Join(def.Required(), ", ")))
			}
			if len(def.Extra()) > 0 {
				w.WriteString(fmt.Sprintf("\n\textra: %s", strings.Join(def.Extra(), ", ")))
			}

			w.WriteByte('\n')

			errs = append(errs, errors.New(w.String()))
		}
	}

	return
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

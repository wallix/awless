/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"fmt"

	"github.com/wallix/awless/template/ast"
)

type Visitor interface {
	Visit([]*ast.CommandNode) error
}

type CollectDefinitions struct {
	C []TemplateDefinition
	L LookupTemplateDefFunc
}

func (t *CollectDefinitions) Visit(cmds []*ast.CommandNode) error {
	for _, cmd := range cmds {
		key := fmt.Sprintf("%s%s", cmd.Action, cmd.Entity)
		if def, ok := t.L(key); ok {
			t.C = append(t.C, def)
		}
	}
	return nil
}

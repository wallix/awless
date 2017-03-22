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
	"strings"
)

type DefinitionLookupFunc func(key string) (Definition, bool)

type Definitions []Definition

func (defs Definitions) Map(fn func(Definition) string) (reduced []string) {
	for _, def := range defs {
		reduced = append(reduced, fn(def))
	}
	return
}

type Definition struct {
	Action, Entity, Api         string
	RequiredParams, ExtraParams []string
}

func (def Definition) Name() string {
	return fmt.Sprintf("%s%s", def.Action, def.Entity)
}

func (def Definition) String() string {
	var required []string
	for _, v := range def.Required() {
		required = append(required, fmt.Sprintf("%s = { %s.%s }", v, def.Entity, v))
	}
	return fmt.Sprintf("%s %s %s", def.Action, def.Entity, strings.Join(required, " "))
}

func (def Definition) GetTemplate() (*Template, error) {
	return Parse(def.String())
}

func (def Definition) Required() []string {
	return def.RequiredParams
}

func (def Definition) Extra() []string {
	return def.ExtraParams
}

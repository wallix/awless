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

type LookupTemplateDefFunc func(key string) (TemplateDefinition, bool)

type TemplateDefinition struct {
	Action, Entity, Api                      string
	RequiredParams, ExtraParams, TagsMapping []string
}

func (def TemplateDefinition) Name() string {
	return fmt.Sprintf("%s%s", def.Action, def.Entity)
}

func (def TemplateDefinition) String() string {
	var required []string
	for _, v := range def.Required() {
		required = append(required, fmt.Sprintf("%s = { %s.%s }", v, def.Entity, v))
	}
	var tags []string
	for _, v := range def.TagsMapping {
		tags = append(tags, fmt.Sprintf("%s = { %s.%s }", v, def.Entity, v))
	}
	return fmt.Sprintf("%s %s %s %s", def.Action, def.Entity, strings.Join(required, " "), strings.Join(tags, " "))
}

func (def TemplateDefinition) GetTemplate() (*Template, error) {
	return Parse(def.String())
}

func (def TemplateDefinition) Required() []string {
	return def.RequiredParams
}

func (def TemplateDefinition) Extra() []string {
	return append(def.ExtraParams, def.TagsMapping...)
}

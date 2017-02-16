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

package template_test

import (
	"testing"

	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver/aws"
)

func TestCollectTemplateDefinitions(t *testing.T) {
	text := "create instance name=nemo\ndelete subnet id=5678\nstop instance id=mine"
	tpl := template.MustParse(text)

	lookup := func(key string) (t template.TemplateDefinition, ok bool) {
		t, ok = aws.AWSTemplatesDefinitions[key]
		return
	}
	collector := &template.CollectDefinitions{L: lookup}
	tpl.Visit(collector)

	collected := collector.C

	if got, want := len(collected), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := collected[0].Name(), "createinstance"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := collected[1].Name(), "deletesubnet"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := collected[2].Name(), "stopinstance"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

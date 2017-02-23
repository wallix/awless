package template_test

import (
	"testing"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/driver/aws"
)

func TestValidation(t *testing.T) {
	t.Run("Validate definitions", func(t *testing.T) {
		text := `create instance name=nemo cidr=10.0.0.0
    delete subnet id=5678
    stop instance ip=10.0.0.0`

		tpl := template.MustParse(text)

		lookup := func(key string) (t template.TemplateDefinition, ok bool) {
			t, ok = aws.AWSTemplatesDefinitions[key]
			return
		}
		rule := &template.DefinitionValidator{lookup}

		errs := tpl.Validate(rule)
		if got, want := len(errs), 2; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		exp := "create instance: unexpected params 'cidr'\n\trequired: image, count, count, type, subnet\n\textra: key, ip, userdata, group, lock, name\n"
		if got, want := errs[0].Error(), exp; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
		exp = "stop instance: unexpected params 'ip'\n\trequired: id\n"
		if got, want := errs[1].Error(), exp; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("Validate name unique", func(t *testing.T) {
		text := "create instance name=instance1_name"

		g := graph.NewGraph()
		g.Unmarshal([]byte(`
      /instance<inst_1> "has_type"@[] "/instance"^^type:text
      /instance<inst_1> "property"@[] "{"Key":"Name","Value":"instance1_name"}"^^type:text
      /instance<inst_2> "has_type"@[] "/instance"^^type:text
      /instance<inst_2> "property"@[] "{"Key":"Id","Value":"inst_2"}"^^type:text
      /instance<inst_2> "property"@[] "{"Key":"Name","Value":"instance2_name"}"^^type:text
    `))

		tpl := template.MustParse(text)

		lookup := func(key string) (*graph.Graph, bool) { return g, true }
		rule := &template.UniqueNameValidator{lookup}

		errs := tpl.Validate(rule)
		if got, want := len(errs), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		exp := "create instance: name 'instance1_name' already exists\n"
		if got, want := errs[0].Error(), exp; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

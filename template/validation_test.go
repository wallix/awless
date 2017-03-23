package template_test

import (
	"testing"

	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
	"github.com/wallix/awless/template"
)

func TestValidation(t *testing.T) {
	t.Run("Validate name unique", func(t *testing.T) {
		text := "create instance name=instance1_name"

		g := graph.NewGraph()
		g.AddResource(
			resourcetest.Instance("inst_1").Prop("Name", "instance1_name").Build(),
			resourcetest.Instance("inst_2").Prop("Name", "instance2_name").Build(),
		)

		tpl := template.MustParse(text)

		lookup := func(key string) (*graph.Graph, bool) { return g, true }
		rule := &template.UniqueNameValidator{lookup}

		errs := tpl.Validate(rule)
		if got, want := len(errs), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		exp := "create instance: name 'instance1_name' already exists"
		if got, want := errs[0].Error(), exp; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

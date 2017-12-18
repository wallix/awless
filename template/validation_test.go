package template_test

import (
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/graph/resourcetest"
	"github.com/wallix/awless/template"
)

func TestValidation(t *testing.T) {
	t.Run("Run name unique", func(t *testing.T) {
		text := "create instance name=instance1_name"

		g := graph.NewGraph()
		g.AddResource(
			resourcetest.Instance("inst_1").Prop("Name", "instance1_name").Prop("State", "terminated").Build(),
			resourcetest.Instance("inst_2").Prop("Name", "instance2_name").Build(),
		)

		tpl := template.MustParse(text)

		lookup := func(key string) (cloud.GraphAPI, bool) { return g, true }
		rule := &template.UniqueNameValidator{lookup}

		errs := tpl.Validate(rule)
		if got, want := len(errs), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		exp := "'instance1_name' name already used for instance inst_1 (state: 'terminated')"
		if got, want := errs[0].Error(), exp; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("Run param is set", func(t *testing.T) {
		text := `create subnet name=subnet_name
		create instance name=instance1_name`
		tpl := template.MustParse(text)

		rule := &template.ParamIsSetValidator{Entity: "instance", Action: "create", Param: "keypair", WarningMessage: "no keypair set"}

		errs := tpl.Validate(rule)
		if got, want := len(errs), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		exp := "no keypair set"
		if got, want := errs[0].Error(), exp; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

		text = `create subnet name=subnet_name
		create instance keypair=$mykey`
		tpl = template.MustParse(text)

		rule = &template.ParamIsSetValidator{Entity: "instance", Action: "create", Param: "keypair", WarningMessage: "no keypair set"}

		errs = tpl.Validate(rule)
		if got, want := len(errs), 0; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	})
}

package template

import (
	"testing"
)

func TestGetTemplateFromDef(t *testing.T) {
	def := &TemplateDefinition{
		Action: "create",
		Entity: "instance",
	}

	tpl, _ := def.GetTemplate()

	exp := "create instance"
	if got, want := tpl.String(), exp; got != want {
		t.Fatalf("\ngot\n%q\n\nwant\n%q\n", got, want)
	}
}

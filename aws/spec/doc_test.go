package awsspec

import (
	"testing"

	"github.com/wallix/awless/aws/doc"
	"github.com/wallix/awless/template/params"
)

func TestDocForEachCommand(t *testing.T) {
	t.Skip()
	for name, def := range AWSTemplatesDefinitions {
		if doc := awsdoc.AwlessExamplesDoc(def.Action, def.Entity); len(doc) == 0 {
			t.Errorf("missing awless CLI examples for template '%s'", name)
		}
	}
}
func TestDocForEachParam(t *testing.T) {
	for name, def := range AWSTemplatesDefinitions {
		for param, _ := range params.Iter(def.Params) {
			if doc, ok := awsdoc.TemplateParamsDoc(name, param); !ok || doc == "" {
				t.Fatalf("missing documentation for param '%s' for '%s'", param, name)
			}
		}
	}
}

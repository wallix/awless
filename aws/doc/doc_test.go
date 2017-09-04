package awsdoc

import (
	"testing"

	"github.com/wallix/awless/aws/driver"
)

func TestDocForEachCommand(t *testing.T) {
	t.Skip()
	for name := range awsdriver.AWSTemplatesDefinitions {
		if doc := exampleDoc(name); len(doc) == 0 {
			t.Errorf("missing awless CLI examples for template '%s'", name)
		}
	}
}
func TestDocForEachParam(t *testing.T) {
	for name, def := range awsdriver.AWSTemplatesDefinitions {
		for _, param := range def.Required() {
			if doc, ok := TemplateParamsDoc(name, param); !ok || doc == "" {
				t.Fatalf("missing documentation for param '%s' for '%s'", param, name)
			}
		}

		for _, param := range def.Extra() {
			if doc, ok := TemplateParamsDoc(name, param); !ok || doc == "" {
				t.Fatalf("missing documentation for param '%s' for '%s'", param, name)
			}
		}

	}
}

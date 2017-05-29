package awsdoc

import (
	"testing"

	"github.com/wallix/awless/aws/driver"
)

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

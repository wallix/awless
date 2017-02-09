package validation

import (
	"fmt"

	"github.com/wallix/awless/graph"
)

var ValidatorsPerActions = map[string][]Validator{
	"create": {new(uniqueName)},
}

type Validator interface {
	Validate(*graph.Graph, map[string]interface{}) error
}

type uniqueName struct {
}

func (v *uniqueName) Validate(g *graph.Graph, params map[string]interface{}) error {
	resources, err := g.FindResourcesByProperty("Name", params["name"])
	if err != nil {
		return err
	}
	switch len(resources) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("name='%s' is alread used by resource %s[%s]", params["name"], resources[0].Id(), resources[0].Type())
	default:
		return fmt.Errorf("name='%s' is alread used by %d resource", params["name"], len(resources))
	}

	return nil
}

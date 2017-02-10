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

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

package validation_test

import (
	"testing"

	"github.com/wallix/awless/cloud/aws/validation"
	"github.com/wallix/awless/graph"
)

func TestValidation(t *testing.T) {
	g := graph.NewGraph()
	g.Unmarshal([]byte(`
    /instance<inst_1>	"has_type"@[]	"/instance"^^type:text
    /instance<inst_1>	"property"@[]	"{"Key":"Name","Value":"instance1_name"}"^^type:text
    /instance<inst_2>	"has_type"@[]	"/instance"^^type:text
    /instance<inst_2>	"property"@[]	"{"Key":"Id","Value":"inst_2"}"^^type:text
    /instance<inst_2>	"property"@[]	"{"Key":"Name","Value":"instance2_name"}"^^type:text
    /subnet<sub_1>	"has_type"@[]	"/subnet"^^type:text
    /subnet<sub_1>	"property"@[]	"{"Key":"Id","Value":"sub_1"}"^^type:text
    `))

	if err := validation.ValidatorsPerActions["create"][0].Validate(g, map[string]interface{}{"name": "instance3_name", "useless": "sub_1"}); err != nil {
		t.Fatal("expected unique")
	}

	if err := validation.ValidatorsPerActions["create"][0].Validate(g, map[string]interface{}{"name": "instance2_name"}); err == nil {
		t.Fatal("expected not unique")
	}

}

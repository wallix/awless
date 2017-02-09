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

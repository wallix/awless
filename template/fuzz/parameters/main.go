package parameters

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/template"
)

func Fuzz(data []byte) int {
	var ok bool
	for _, def := range awsspec.AWSTemplatesDefinitions {
		fillers := make(map[string]interface{})
		for _, param := range def.Params.Required() {
			fillers[param] = "default"
		}

		env := template.NewEnv().WithFillers(fillers).WithAliasFunc(func(e, k, v string) string { return "" }).
			WithLookupCommandFunc(func(tokens ...string) interface{} {
				return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
			}).Build()
		for _, param := range def.Params.Required() {
			inTpl, err := template.Parse(fmt.Sprintf("%s %s %s=%s", def.Action, def.Entity, param, string(data)))
			if err != nil {
				continue
			}

			_, _, err = template.Compile(inTpl, env, template.TestCompileMode)
			if err != nil {
				continue
			}
			ok = true
		}
	}
	if ok {
		return 1
	}
	return 0
}

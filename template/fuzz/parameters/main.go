package parameters

import (
	"fmt"
	"strings"

	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/template"
	"github.com/wallix/awless/template/env"
)

func Fuzz(data []byte) int {
	var ok bool
	for _, def := range awsspec.AWSTemplatesDefinitions {
		fillers := make(map[string]interface{})
		for _, param := range def.Params.Required() {
			fillers[param] = "default"
		}

		cenv := template.NewEnv().WithAliasFunc(func(p, v string) string { return "" }).
			WithLookupCommandFunc(func(tokens ...string) interface{} {
				return awsspec.MockAWSSessionFactory.Build(strings.Join(tokens, ""))()
			}).Build()
		cenv.Push(env.FILLERS, fillers)
		for _, param := range def.Params.Required() {
			inTpl, err := template.Parse(fmt.Sprintf("%s %s %s=%s", def.Action, def.Entity, param, string(data)))
			if err != nil {
				continue
			}

			_, _, err = template.Compile(inTpl, cenv, template.TestCompileMode)
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

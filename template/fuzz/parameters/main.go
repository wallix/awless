package parameters

import (
	"fmt"

	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/template"
)

func Fuzz(data []byte) int {
	var ok bool
	for _, def := range awsdriver.AWSTemplatesDefinitions {
		env := template.NewEnv()
		fillers := make(map[string]interface{})
		for _, param := range def.Required() {
			fillers[param] = "default"
		}
		env.AddFillers(fillers)

		env.AliasFunc = func(e, k, v string) string { return "" }
		env.DefLookupFunc = func(in string) (template.Definition, bool) {
			return def, true
		}
		for _, param := range def.Required() {
			inTpl, err := template.Parse(fmt.Sprintf("%s %s %s=%s", def.Action, def.Entity, param, string(data)))
			if err != nil {
				continue
			}

			_, _, err = template.Compile(inTpl, env, template.LenientCompileMode)
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

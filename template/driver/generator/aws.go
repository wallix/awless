//go:generate go run $GOFILE
package main

import (
	"html/template"
	"os"
	"strings"

	"github.com/wallix/awless/template/driver/aws"
)

func main() {
	generateDriverFuncs()
	generateTemplateTemplates()
}

func generateTemplateTemplates() {
	templ, err := template.New("templates_definitions").Parse(templateDefinitions)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("../aws/template_defs.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	err = templ.Execute(f, aws.DriverDefinitions)
	if err != nil {
		panic(err)
	}
}

func generateDriverFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"capitalize": capitalize,
	}).Parse(funcsTempl)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("../aws/driver_gen_funcs.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	err = templ.Execute(f, aws.DriverDefinitions)
	if err != nil {
		panic(err)
	}
}

const templateDefinitions = `// DO NOT EDIT
// This file was automatically generated with go generate
package aws

var AWSTemplatesDefinitions = map[string]string{
{{- range $index, $def := . }}
	"{{ $def.Action }}{{ $def.Entity }}": "{{ $def.Action }} {{ $def.Entity }}{{ range $awsField, $field := $def.RequiredParams }} {{ $field }}={ {{ $def.Entity }}.{{ $field }} }{{ end }} {{ range $awsField, $field := $def.TagsMapping }} {{ $field }}={ {{ $def.Entity }}.{{ $field }} }{{ end }}",
{{- end }}
}
`

const funcsTempl = `// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)
{{ range $index, $def := . }}
{{- if not $def.ManualFuncDefinition }}

{{- if $def.DryRunUnsupported }}
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}_DryRun(params map[string]interface{}) (interface{}, error) {
	d.logger.Println("!! {{ $def.Action }} {{ $def.Entity }}: dry run not supported")
	return nil, nil
}
{{ else }}
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &{{ $def.Api }}.{{ $def.Input }}{}

	input.DryRun = aws.Bool(true)
	{{- if gt (len $def.RequiredParams) 0 }}
	// Required params
	{{- range $awsField, $field := $def.RequiredParams }}
	setField(params["{{ $field }}"], input, "{{ $awsField }}")
	{{- end }}
	{{- end }}
	{{- if gt (len $def.ExtraParams) 0 }}
	// Extra params
	{{- range $awsField, $field := $def.ExtraParams }}
	if _, ok := params["{{ $field }}"]; ok {
		setField(params["{{ $field }}"], input, "{{ $awsField }}")
	}
	{{- end }}
	{{- end }}

	_, err := d.{{ $def.Api }}.{{ $def.ApiMethod }}(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			id := fmt.Sprintf("{{ $def.Entity }}_%d", rand.Intn(1e3))
			{{- if gt (len $def.TagsMapping) 0 }}
			tagsParams := map[string]interface{}{"resource": id}
			{{- range $tagName, $field := $def.TagsMapping }}
			tagsParams["{{ $tagName }}"] = params["{{ $field }}"]
			{{- end }}
			d.Create_Tags_DryRun(tagsParams)
			{{- end }}
			d.logger.Println("dry run: {{ $def.Action }} {{ $def.Entity }} ok")
			return id, nil
		}
	}

	d.logger.Printf("dry run: {{ $def.Action }} {{ $def.Entity }} error: %s", err)
	return nil, err
}
{{ end }}
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}(params map[string]interface{}) (interface{}, error) {
	input := &{{ $def.Api }}.{{ $def.Input }}{}
	{{- if gt (len $def.RequiredParams) 0 }}
	// Required params
	{{- range $awsField, $field := $def.RequiredParams }}
	setField(params["{{ $field }}"], input, "{{ $awsField }}")
	{{- end }}
	{{- end }}
	{{- if gt (len $def.ExtraParams) 0 }}
	// Extra params
	{{- range $awsField, $field := $def.ExtraParams }}
	if _, ok := params["{{ $field }}"]; ok {
		setField(params["{{ $field }}"], input, "{{ $awsField }}")
	}
	{{- end }}
	{{- end }}

	output, err := d.{{ $def.Api }}.{{ $def.ApiMethod }}(input)
	if err != nil {
		d.logger.Printf("{{ $def.Action }} {{ $def.Entity }} error: %s", err)
		return nil, err
	}

	{{- if ne $def.OutputExtractor "" }}
	id := aws.StringValue(output.{{ $def.OutputExtractor }})
	{{- if gt (len $def.TagsMapping) 0 }}
	tagsParams := map[string]interface{}{"resource": id}

	{{- range $tagName, $field := $def.TagsMapping }}
	tagsParams["{{ $tagName }}"] = params["{{ $field }}"]
	{{- end }}
	d.Create_Tags(tagsParams)
	{{- end }}
	d.logger.Printf("{{ $def.Action }} {{ $def.Entity }} '%s' done", id)
	return aws.StringValue(output.{{ $def.OutputExtractor }}), nil
	{{- else }}
	d.logger.Println("{{ $def.Action }} {{ $def.Entity }} done")
	return output, nil
	{{- end }}
}
{{ end }}
{{- end }}`

func capitalize(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return strings.ToUpper(s)
}

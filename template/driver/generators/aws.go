//go:generate go run $GOFILE
//go:generate gofmt -s -w ../aws/gen_template_defs.go
//go:generate gofmt -s -w ../aws/gen_driver_funcs.go

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

package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/wallix/awless/template/driver/aws/definitions"
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

	var buff bytes.Buffer
	err = templ.Execute(&buff, definitions.Driver)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("../aws/gen_template_defs.go", buff.Bytes(), 0666); err != nil {
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

	var buff bytes.Buffer
	err = templ.Execute(&buff, definitions.Driver)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("../aws/gen_driver_funcs.go", buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const templateDefinitions = `/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"github.com/wallix/awless/template"
)

var AWSTemplatesDefinitions = map[string]template.TemplateDefinition{
{{- range $index, $def := . }}
	"{{ $def.Action }}{{ $def.Entity }}": template.TemplateDefinition{
			Action: "{{ $def.Action }}",
			Entity: "{{ $def.Entity }}",
			Api: "{{ $def.Api }}",
			RequiredParams: []string{ {{- range $awsField, $field := $def.RequiredParams }}"{{ $field.TemplateName }}", {{- end}} },
			ExtraParams: []string{ {{- range $awsField, $field := $def.ExtraParams }}"{{ $field.TemplateName }}", {{- end}} },
			TagsMapping: []string{ {{- range $awsField, $field := $def.TagsMapping }}"{{ $field }}", {{- end}} },
		},
{{- end }}
}
`

const funcsTempl = `/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	dryRunOperation = "DryRunOperation"
	notFound = "NotFound"
)

{{ range $index, $def := . }}
{{- if not $def.ManualFuncDefinition }}

{{- if $def.DryRunUnsupported }}
// This function was auto generated
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}_DryRun(params map[string]interface{}) (interface{}, error) {
	{{- range $awsField, $field := $def.RequiredParams }}
	if _, ok := params["{{ $field.TemplateName }}"]; !ok {
		return nil, errors.New("{{ $def.Action }} {{ $def.Entity }}: missing required params '{{ $field.TemplateName }}'")
	}
	{{ end }}
	d.logger.Verbose("params dry run: {{ $def.Action }} {{ $def.Entity }} ok")
	return nil, nil
}
{{ else }}
// This function was auto generated
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &{{ $def.Api }}.{{ $def.Input }}{}
	input.DryRun = aws.Bool(true)
	var err error
	{{if gt (len $def.RequiredParams) 0 }}
	// Required params
	{{- range $i, $field := $def.RequiredParams }}
	err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }})
	if err != nil {
		return nil, err
	}
	{{- end }}
	{{- end }}
	{{if gt (len $def.ExtraParams) 0 }}
	// Extra params
	{{- range $awsField, $field := $def.ExtraParams }}
	if _, ok := params["{{ $field.TemplateName }}"]; ok {
		err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }})
		if err != nil {
			return nil, err
		}
	}
	{{- end }}
	{{- end }}

	_, err = d.{{ $def.Api }}.{{ $def.ApiMethod }}(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			id := fakeDryRunId("{{ $def.Entity }}")
			{{- if gt (len $def.TagsMapping) 0 }}
			tagsParams := map[string]interface{}{"resource": id}
			{{- range $tagName, $field := $def.TagsMapping }}
			if v, ok := params["{{ $field }}"]; ok {
				tagsParams["{{ $tagName }}"] = v
			}
			{{- end }}
			if len(tagsParams) > 1 {
				_, err = d.Create_Tags_DryRun(tagsParams)
				if err != nil {
					d.logger.Errorf("{{ $def.Action }} {{ $def.Entity }}: adding tags: error: %s", err)
					return nil, err
				}
			}
			{{- end }}
			d.logger.Verbose("full dry run: {{ $def.Action }} {{ $def.Entity }} ok")
			return id, nil
		}
	}

	d.logger.Errorf("dry run: {{ $def.Action }} {{ $def.Entity }} error: %s", err)
	return nil, err
}
{{ end }}
// This function was auto generated
func (d *AwsDriver) {{ capitalize $def.Action }}_{{ capitalize $def.Entity }}(params map[string]interface{}) (interface{}, error) {
	input := &{{ $def.Api }}.{{ $def.Input }}{}
	var err error
	{{if gt (len $def.RequiredParams) 0 }}
	// Required params
	{{- range $i, $field := $def.RequiredParams }}
	err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }})
	if err != nil {
		return nil, err
	}
	{{- end }}
	{{- end }}
	{{if gt (len $def.ExtraParams) 0 }}
	// Extra params
	{{- range $awsField, $field := $def.ExtraParams }}
	if _, ok := params["{{ $field.TemplateName }}"]; ok {
		err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }})
		if err != nil {
			return nil, err
		}
	}
	{{- end }}
	{{- end }}

	start := time.Now()
	var output *{{ $def.Api }}.{{ $def.Output }}
	output, err = d.{{ $def.Api }}.{{ $def.ApiMethod }}(input)
	output = output
	if err != nil {
		d.logger.Errorf("{{ $def.Action }} {{ $def.Entity }} error: %s", err)
		return nil, err
	}
	d.logger.ExtraVerbosef("{{ $def.Api }}.{{ $def.ApiMethod }} call took %s", time.Since(start))


	{{- if ne $def.OutputExtractor "" }}
	id := {{ $def.OutputExtractor }}
	{{- if gt (len $def.TagsMapping) 0 }}
	tagsParams := map[string]interface{}{"resource": id}
	{{- range $tagName, $field := $def.TagsMapping }}
	if v, ok := params["{{ $field }}"]; ok {
		tagsParams["{{ $tagName }}"] = v
	}
	{{- end }}
	if len(tagsParams) > 1 {
		_, err := d.Create_Tags(tagsParams)
		if err != nil {
			d.logger.Errorf("{{ $def.Action }} {{ $def.Entity }}: adding tags: error: %s", err)
			return nil, err
		}
	}
	{{- end }}
	d.logger.Verbosef("{{ $def.Action }} {{ $def.Entity }} '%s' done", id)
	return {{ $def.OutputExtractor }}, nil
	{{- else }}
	d.logger.Verbose("{{ $def.Action }} {{ $def.Entity }} done")
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

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
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

func generateTemplateTemplates() {
	templ, err := template.New("templates_definitions").Parse(templateDefinitions)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.DriversDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(DRIVERS_DIR, "gen_template_defs.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

func generateDriverFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title": strings.Title,
	}).Parse(driversTempl)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.DriversDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(DRIVERS_DIR, "gen_driver_funcs.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

func generateDriverTypes() {
	templ, err := template.New("types").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(typesTempl)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.DriversDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(DRIVERS_DIR, "gen_drivers.go"), buff.Bytes(), 0666); err != nil {
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
package awsdriver

import (
	"github.com/wallix/awless/template"
)


var APIPerTemplateDefName = map[string]string {
{{- range $, $service := . }}
  {{- range $, $def := $service.Drivers }}
  "{{ $def.Action }}{{ $def.Entity }}": "{{ $service.Api }}",
  {{- end }}
{{- end }}
}

var AWSTemplatesDefinitions = map[string]template.Definition{
{{- range $, $service := . }}
{{- range $index, $def := $service.Drivers }}
	"{{ $def.Action }}{{ $def.Entity }}": template.Definition{
			Action: "{{ $def.Action }}",
			Entity: "{{ $def.Entity }}",
			Api: "{{ $service.Api }}",
			RequiredParams: []string{ {{- range $key := $def.RequiredKeys }}"{{ $key }}", {{- end}} },
			ExtraParams: []string{ {{- range $key := $def.ExtraKeys }}"{{ $key }}", {{- end}} },
		},
{{- end }}
{{- end }}
}

func DriverSupportedActions() map[string][]string { 
	supported := make(map[string][]string)
{{- range $, $service := . }}
{{- range $index, $def := $service.Drivers }}
	supported["{{ $def.Action }}"] = append(supported["{{ $def.Action }}"], "{{ $def.Entity }}")
{{- end }}
{{- end }}
	return supported
}
`

const driversTempl = `/* Copyright 2017 WALLIX

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
package awsdriver

import (
	"errors"
	"strings"
	"time"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	{{- range $index, $service := . }}
	"github.com/aws/aws-sdk-go/service/{{ $service.Api }}"
	{{- end }}
)

const (
	dryRunOperation = "DryRunOperation"
	notFound = "NotFound"
)

{{- range $, $service := . }}
{{ range $index, $def := $service.Drivers }}
{{- if not $def.ManualFuncDefinition }}

{{- if $def.DryRunUnsupported }}
// This function was auto generated
func (d *{{ Title $service.Api }}Driver) {{ Title $def.Action }}_{{ Title $def.Entity }}_DryRun(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
	{{- range $awsField, $field := $def.RequiredParams }}
	if _, ok := params["{{ $field.TemplateName }}"]; !ok {
		return nil, errors.New("{{ $def.Action }} {{ $def.Entity }}: missing required params '{{ $field.TemplateName }}'")
	}
	{{ end }}
	d.logger.Verbose("params dry run: {{ $def.Action }} {{ $def.Entity }} ok")
	return fakeDryRunId("{{ $def.Entity }}"), nil
}
{{ else }}
// This function was auto generated
func (d *{{ Title $service.Api }}Driver) {{ Title $def.Action }}_{{ Title $def.Entity }}_DryRun(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
	input := &{{ $service.Api }}.{{ $def.Input }}{}
	input.DryRun = aws.Bool(true)
	var err error
	{{if gt (len $def.RequiredParams) 0 }}
	// Required params
		{{- range $i, $field := $def.RequiredParams }}
			{{- if not $field.AsAwsTag }}
	err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }}, ctx)
	if err != nil {
		return nil, err
	}
			{{- end }}
		{{- end }}
	{{- end }}
	{{if gt (len $def.ExtraParams) 0 }}
	// Extra params
		{{- range $awsField, $field := $def.ExtraParams }}
			{{- if not $field.AsAwsTag }}
	if _, ok := params["{{ $field.TemplateName }}"]; ok {
		err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }}, ctx)
		if err != nil {
			return nil, err
		}
	}
			{{- end }}
		{{- end }}
	{{- end }}

	_, err = d.{{ $def.ApiMethod }}(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound), strings.Contains(awsErr.Message(), "Invalid IAM Instance Profile name"):
			id := fakeDryRunId("{{ $def.Entity }}")
			{{- range $i, $field := $def.RequiredParams }}
				{{- if $field.AsAwsTag }}
				// Required param as tag
			_, err = d.Create_Tag_DryRun(ctx, map[string]interface{}{"key":"{{ $field.AwsField }}", "value":params["{{ $field.TemplateName }}"], "resource":id})
			if err != nil {
				return nil, fmt.Errorf("dry run: {{ $def.Action }} {{ $def.Entity }}: adding tags: %s",err)
			}
				{{- end }}
			{{- end }}
			{{- range $i, $field := $def.ExtraParams }}
				{{- if $field.AsAwsTag }}
				// Extra param as tag
			if v, ok := params["{{ $field.TemplateName }}"]; ok {
				_, err = d.Create_Tag_DryRun(ctx, map[string]interface{}{"key":"{{ $field.AwsField }}", "value":v, "resource":id})
				if err != nil {
					return nil, fmt.Errorf("dry run: {{ $def.Action }} {{ $def.Entity }}: adding tags: %s",err)
				}
			}
				{{- end }}
			{{- end }}
			d.logger.Verbose("dry run: {{ $def.Action }} {{ $def.Entity }} ok")
			return id, nil
		}
	}

	return nil, fmt.Errorf("dry run: {{ $def.Action }} {{ $def.Entity }}: %s", err)
}
{{ end }}
// This function was auto generated
func (d *{{ Title $service.Api }}Driver) {{ Title $def.Action }}_{{ Title $def.Entity }}(ctx driver.Context, params map[string]interface{}) (interface{}, error) {
	input := &{{ $service.Api }}.{{ $def.Input }}{}
	var err error
	{{if gt (len $def.RequiredParams) 0 }}
	// Required params
		{{- range $i, $field := $def.RequiredParams }}
			{{- if not $field.AsAwsTag }}
	err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }}, ctx)
	if err != nil {
		return nil, err
	}
			{{- end }}
		{{- end }}
	{{- end }}
	{{if gt (len $def.ExtraParams) 0 }}
	// Extra params
		{{- range $awsField, $field := $def.ExtraParams }}
			{{- if not $field.AsAwsTag }}
	if _, ok := params["{{ $field.TemplateName }}"]; ok {
		err = setFieldWithType(params["{{ $field.TemplateName }}"], input, "{{ $field.AwsField }}", {{ $field.AwsType }}, ctx)
		if err != nil {
			return nil, err
		}
	}
			{{- end }}
		{{- end }}
	{{- end }}

	start := time.Now()
	var output *{{ $service.Api }}.{{ $def.Output }}
	output, err = d.{{ $def.ApiMethod }}(input)
	output = output
	if err != nil {
		return nil, fmt.Errorf("{{ $def.Action }} {{ $def.Entity }}: %s", err)
	}
	d.logger.ExtraVerbosef("{{ $service.Api }}.{{ $def.ApiMethod }} call took %s", time.Since(start))

	{{- if ne $def.OutputExtractor "" }}
	id := {{ $def.OutputExtractor }}
	
	{{- range $i, $field := $def.RequiredParams }}
		{{- if $field.AsAwsTag }}
		// Required param as tag
	_, err = d.Create_Tag(ctx, map[string]interface{}{"key":"{{ $field.AwsField }}", "value":params["{{ $field.TemplateName }}"], "resource":id})
	if err != nil {
		return nil, fmt.Errorf("{{ $def.Action }} {{ $def.Entity }}: adding tags: %s",err)
	}
		{{- end }}
	{{- end }}
	{{- range $i, $field := $def.ExtraParams }}
		{{- if $field.AsAwsTag }}
		// Extra param as tag
	if v, ok := params["{{ $field.TemplateName }}"]; ok {
		_, err = d.Create_Tag(ctx, map[string]interface{}{"key":"{{ $field.AwsField }}", "value":v, "resource":id})
		if err != nil {
			return nil, fmt.Errorf("{{ $def.Action }} {{ $def.Entity }}: adding tags: %s",err)
		}
	}
		{{- end }}
	{{- end }}
	
	d.logger.Infof("{{ $def.Action }} {{ $def.Entity }} '%s' done", id)
	return id, nil
	{{- else }}
	d.logger.Info("{{ $def.Action }} {{ $def.Entity }} done")
	return output, nil
	{{- end }}
}
{{ end }}
{{- end }}
{{- end }}`

const typesTempl = `/* Copyright 2017 WALLIX

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
package awsdriver

import (
	"strings"
	"github.com/wallix/awless/template/driver"
	"github.com/wallix/awless/logger"
	{{- range $index, $service := . }}
  "github.com/aws/aws-sdk-go/service/{{ $service.Api }}/{{ $service.Api }}iface"
	{{- end }}
)

{{ range $, $service := . }}
type {{ Title $service.Api }}Driver struct {
	dryRun bool
	logger *logger.Logger
	{{ $service.Api }}iface.{{ ApiToInterface $service.Api }}
}

func (d *{{ Title $service.Api }}Driver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *{{ Title $service.Api }}Driver) SetLogger(l *logger.Logger) { d.logger = l }
func New{{ Title $service.Api }}Driver(api {{ $service.Api }}iface.{{ ApiToInterface $service.Api }}) driver.Driver{
	return &{{ Title $service.Api }}Driver{false, logger.DiscardLogger, api}
}

func (d *{{ Title $service.Api }}Driver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {
	{{ range $, $def := $service.Drivers }}
case "{{ $def.Action}}{{ $def.Entity}}":
		if d.dryRun {
			return d.{{ Title $def.Action }}_{{ Title $def.Entity }}_DryRun, nil
		}
		return d.{{ Title $def.Action }}_{{ Title $def.Entity }}, nil
	{{ end }}
	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

{{ end }}`

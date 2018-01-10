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

func generateTestMocks() {
	templ, err := template.New("mocks").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"Join":           strings.Join,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(mocksTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.Mocks())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(SERVICES_DIR, "gen_mocks_test.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const mocksTempl = `// Auto generated implementation for the AWS cloud service

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

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

{{ range $, $mock := . }}

type {{ $mock.Name }} struct {
	{{ $mock.Api }}iface.{{ ApiToInterface $mock.Api }}
	{{- range $, $func := $mock.Funcs }}
	{{- if eq $func.MockFieldType "mapslice" }}
		{{ $func.MockField}} map[string][]*{{ $func.AWSType }}
	{{- else if eq $func.MockFieldType "map" }}
			{{ $func.MockField}} map[string]{{ $func.AWSType }}
	{{- else }}
		{{ $func.MockField}} []*{{ $func.AWSType }}
	{{- end }}
	{{- end }}
}

func (m * {{ $mock.Name }}) Name() string {
	return ""
}

func (m * {{ $mock.Name }}) Region() string {
	return ""
}

func (m * {{ $mock.Name }}) Profile() string {
	return ""
}

func (m * {{ $mock.Name }}) Provider() string {
	return ""
}

func (m * {{ $mock.Name }}) ProviderAPI() string {
	return ""
}

func (m * {{ $mock.Name }}) ResourceTypes() []string {
	return []string{}
}

func (m * {{ $mock.Name }}) Fetch(context.Context) (cloud.GraphAPI, error) {
	return nil, nil
}

func (m * {{ $mock.Name }}) IsSyncDisabled() bool {
	return false
}

func (m * {{ $mock.Name }}) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return nil, nil
}

{{ range $, $func := $mock.Funcs }}
	{{- if not $func.Manual }}
		{{- if eq $func.FuncType "list" }}
			{{- if $func.Multipage }}
				func (m * {{ $mock.Name }}) {{ $func.ApiMethod }}(input *{{ $func.Input }}, fn func(p *{{ $func.Output}}, lastPage bool) (shouldContinue bool)) error {
					var pages [][]*{{ $func.AWSType }}
					for i := 0; i < len(m.{{ $func.MockField}}); i += 2 {
						page := []*{{ $func.AWSType }}{m.{{ $func.MockField}}[i]}
						if i+1 < len(m.{{ $func.MockField}}) {
							page = append(page, m.{{ $func.MockField}}[i+1])
						}
						pages = append(pages, page)
					}
					for i, page := range pages {
						fn(&{{ $func.Output }} { {{ $func.OutputsExtractor }}: page, {{ $func.NextPageMarker }}: aws.String(strconv.Itoa(i + 1))},
							i < len(pages),
						)
					}
					return nil
				}
			{{ else }}
				func (m * {{ $mock.Name }}) {{ $func.ApiMethod }}(input *{{ $func.Input }}) (*{{ $func.Output}}, error ){
					return &{{ $func.Output}}{ {{ $func.OutputsExtractor }}: m.{{ $func.MockField}} }, nil
				}
			{{ end }}
		{{- end }}
	{{- end }}

{{ end }}

{{- end }}
`

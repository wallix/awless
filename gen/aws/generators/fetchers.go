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
	"strings"
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

func generateFetcherFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"Join":           strings.Join,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(fetchersTempl)

	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, aws.FetchersDefs, FETCHERS_DIR, "gen_fetchers.go")
}

const fetchersTempl = `// Auto generated implementation for the AWS cloud service

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

package awsfetch

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "context"
 
  awssdk "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  {{- range $, $api := $service.Api }}
  "github.com/aws/aws-sdk-go/service/{{ $api }}"
  "github.com/aws/aws-sdk-go/service/{{ $api }}/{{ $api }}iface"
  {{- end }}
  {{- end }}
  "github.com/wallix/awless/fetch"
  "github.com/wallix/awless/graph"
  "github.com/wallix/awless/aws/conv"
)

{{- range $index, $service := . }}
func Build{{ Title $service.Name }}FetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManual{{ Title $service.Name }}FetchFuncs(conf, funcs)
	
{{- range $index, $fetcher := $service.Fetchers }}
	{{- if not $fetcher.ManualFetcher }}

	funcs["{{ $fetcher.ResourceType }}"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*{{ $fetcher.AWSType }}

		if !conf.getBoolDefaultTrue("aws.{{ $service.Name }}.{{ $fetcher.ResourceType }}.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource {{ $service.Name }}[{{ $fetcher.ResourceType }}]")
			return resources, objects, nil
		}
		
		{{- if $fetcher.Multipage }}
		var badResErr error
		err := conf.APIs.{{ Title $fetcher.Api}}.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }},
			func(out *{{ $fetcher.Output }}, lastPage bool) (shouldContinue bool) {
				{{- if ne $fetcher.OutputsContainers "" }}
				for _, all := range out.{{ $fetcher.OutputsContainers }} {
				{{- end }}
					for _, output := range {{ if ne $fetcher.OutputsContainers "" }}all{{ else }}out{{ end }}.{{ $fetcher.OutputsExtractor }} {
						if badResErr != nil {
							return false
						}
						objects = append(objects, output)
						var res *graph.Resource
						if res, badResErr = awsconv.NewResource(output); badResErr != nil {
							return false
						}
						resources = append(resources, res)
					}
				{{- if ne $fetcher.OutputsContainers "" }}
				}
				{{- end }}
				return out.{{ $fetcher.NextPageMarker }} != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
		{{- else }}
		
		out, err := conf.APIs.{{ Title $fetcher.Api}}.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.{{ $fetcher.OutputsExtractor }} {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}
			
		return resources, objects, nil{{ end }}
	}
{{- end }}
{{- end }}
	return funcs
}
{{- end }}`

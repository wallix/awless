//go:generate go run $GOFILE
//go:generate gofmt -s -w ../aws/gen_api.go
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

	"github.com/wallix/awless/cloud/aws/definitions"
)

func main() {
	generateFetcherFuncs()
}

func generateFetcherFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":   strings.Title,
		"ToUpper": strings.ToUpper,
		"Join":    strings.Join,
	}).Parse(funcsTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, definitions.Services)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("../aws/gen_api.go", buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const funcsTempl = `// Auto generated implementation for the AWS cloud service

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

package aws

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "fmt"
	"sync"

  awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  "github.com/aws/aws-sdk-go/service/{{ $service.Api }}"
  "github.com/aws/aws-sdk-go/service/{{ $service.Api }}/{{ $service.Api }}iface"
  {{- end }}
	"github.com/wallix/awless/cloud"
  "github.com/wallix/awless/graph"
)

func init() {
  {{- range $index, $service := . }}
  ServiceNames = append(ServiceNames, "{{ $service.Name }}")
  {{- end }}
}

var ServiceNames = []string{}

var ResourceTypesPerAPI = map[string][]string {
{{- range $index, $service := . }}
  "{{ $service.Api }}": []string{
    {{- range $idx, $fetcher := $service.Fetchers }}
      "{{ $fetcher.ResourceType }}",
    {{- end }}
  },
{{- end }}
}

var ServicePerAPI = map[string]string {
{{- range $index, $service := . }}
  "{{ $service.Api }}": "{{ $service.Name }}",
{{- end }}
}

var ServicePerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $service.Name }}",
  {{- end }}
{{- end }}
}

{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
	once oncer
  region string
  {{ $service.Api }}iface.{{ ToUpper $service.Api }}API
}

func New{{ Title $service.Name }}(sess *session.Session) *{{ Title $service.Name }} {
  region := awssdk.StringValue(sess.Config.Region)
	return &{{ Title $service.Name }}{ {{ ToUpper $service.Api }}API: {{ $service.Api }}.New(sess), region: region }
}

func (s *{{ Title $service.Name }}) Name() string {
  return "{{ $service.Name }}"
}

func (s *{{ Title $service.Name }}) Provider() string {
  return "aws"
}

func (s *{{ Title $service.Name }}) ProviderAPI() string {
  return "{{ $service.Api }}"
}

func (s *{{ Title $service.Name }}) ProviderRunnableAPI() interface{} {
  return s.{{ ToUpper $service.Api }}API
}

func (s *{{ Title $service.Name }}) ResourceTypes() (all []string) {
  {{- range $index, $fetcher := $service.Fetchers }}
  all = append(all, "{{ $fetcher.ResourceType }}")
  {{- end }}
  return
}

func (s *{{ Title $service.Name }}) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	regionN := graph.InitResource(s.region, graph.Region)
	g.AddResource(regionN)

	{{- range $index, $fetcher := $service.Fetchers }}
  var {{ $fetcher.ResourceType }}List []*{{ $fetcher.AWSType }}
  {{- end }}

	errc := make(chan error)
	var wg sync.WaitGroup

	{{- range $index, $fetcher := $service.Fetchers }}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, {{ $fetcher.ResourceType }}List, err = s.fetch_all_{{ $fetcher.ResourceType }}_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
  {{- end }}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case "Access Denied":
				return g, cloud.ErrFetchAccessDenied
			default:
				return g, ee
			}
		case nil:
			continue
		default:
			return g, ee
		}
	}

	errc = make(chan error)
	{{- range $index, $fetcher := $service.Fetchers }}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range {{ $fetcher.ResourceType }}List {
			for _, fn := range addParentsFns["{{ $fetcher.ResourceType }}"] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()
  {{- end }}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return g, err
		}
	}

	return g, nil
}

func (s *{{ Title $service.Name }}) FetchByType(t string) (*graph.Graph, error) {
  switch t {
  {{- range $index, $fetcher := $service.Fetchers }}
  case "{{ $fetcher.ResourceType }}":
		graph, _, err := s.fetch_all_{{ $fetcher.ResourceType }}_graph()
    return graph, err
  {{- end }}
  default:
    return nil, fmt.Errorf("aws {{ $service.Name }}: unsupported fetch for type %s", t)
  }
}

{{ range $index, $fetcher := $service.Fetchers }}
{{- if not $fetcher.ManualFetcher }}
func (s *{{ Title $service.Name }}) fetch_all_{{ $fetcher.ResourceType }}_graph() (*graph.Graph, []*{{ $fetcher.AWSType }}, error) {
  g := graph.NewGraph()
	var cloudResources []*{{ $fetcher.AWSType }}
	{{- if $fetcher.Multipage }}
	var badResErr error
	err := s.{{ $fetcher.ApiMethod }}(&{{ $service.Api }}.{{ $fetcher.Input }},
		func(out *{{ $service.Api }}.{{ $fetcher.Output }}, lastPage bool) (shouldContinue bool) {
			{{- if ne $fetcher.OutputsContainers "" }}
			for _, all := range out.{{ $fetcher.OutputsContainers }} {
	      for _, output := range all.{{ $fetcher.OutputsExtractor }} {
					cloudResources = append(cloudResources, output)
					var res *graph.Resource
					res, badResErr = newResource(output)
					if badResErr != nil {
						return false
					}
	        g.AddResource(res)
	      }
	    }
			{{- else }}
			for _, output := range out.{{ $fetcher.OutputsExtractor }} {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			{{- end }}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
	{{- else }}
  out, err := s.{{ $fetcher.ApiMethod }}(&{{ $service.Api }}.{{ $fetcher.Input }})
  if err != nil {
    return nil, cloudResources, err
  }
  	{{ if ne $fetcher.OutputsContainers "" }}
    for _, all := range out.{{ $fetcher.OutputsContainers }} {
      for _, output := range all.{{ $fetcher.OutputsExtractor }} {
				cloudResources = append(cloudResources, output)
        res, err := newResource(output)
        if err != nil {
          return g, cloudResources, err
        }
        g.AddResource(res)
      }
    }
  	{{ else }}
    for _, output := range out.{{ $fetcher.OutputsExtractor }} {
			cloudResources = append(cloudResources, output)
      res, err := newResource(output)
      if err != nil {
        return g, cloudResources, err
      }
      g.AddResource(res)
    }
  	{{ end }}
  return g, cloudResources, nil
	{{ end }}
}
{{- end }}
{{ end }}

{{ end }}`

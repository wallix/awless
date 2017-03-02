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

func generateFetcherFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":   strings.Title,
		"ToUpper": strings.ToUpper,
		"Join":    strings.Join,
	}).Parse(fetchersTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.FetchersDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(FETCHERS_DIR, "gen_api.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
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

package aws

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "fmt"
	"sync"

  awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  {{- range $, $api := $service.Api }}
  "github.com/aws/aws-sdk-go/service/{{ $api }}"
  "github.com/aws/aws-sdk-go/service/{{ $api }}/{{ $api }}iface"
  {{- end }}
  {{- end }}
	"github.com/wallix/awless/cloud"
  "github.com/wallix/awless/graph"
	"github.com/wallix/awless/template/driver"
	awsdriver "github.com/wallix/awless/aws/driver"
)

func init() {
  {{- range $index, $service := . }}
  ServiceNames = append(ServiceNames, "{{ $service.Name }}")
  {{- end }}
}

var ServiceNames = []string{}

var ResourceTypes = []string {
{{- range $index, $service := . }}
    {{- range $idx, $fetcher := $service.Fetchers }}
      "{{ $fetcher.ResourceType }}",
    {{- end }}
{{- end }}
}

var ServicePerAPI = map[string]string {
{{- range $index, $service := . }}
{{- range $, $api := $service.Api }}
  "{{ $api }}": "{{ $service.Name }}",
{{- end }}
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
	{{- range $, $api := $service.Api }}
  {{ $api }}iface.{{ ToUpper $api }}API
	{{- end }}
}

func New{{ Title $service.Name }}(sess *session.Session) *{{ Title $service.Name }} {
  region := awssdk.StringValue(sess.Config.Region)
	return &{{ Title $service.Name }}{ 
	{{- range $, $api := $service.Api }}
	{{ ToUpper $api }}API: {{ $api }}.New(sess),
	{{- end }}
	 region: region,
  }
}

func (s *{{ Title $service.Name }}) Name() string {
  return "{{ $service.Name }}"
}

func (s *{{ Title $service.Name }}) Drivers() []driver.Driver {
  return []driver.Driver{ 
		{{- range $, $api := $service.Api }}
		awsdriver.New{{ Title $api }}Driver(s.{{ ToUpper $api }}API),
		{{- end }}
	}
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
	err := s.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }},
		func(out *{{ $fetcher.Output }}, lastPage bool) (shouldContinue bool) {
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
			return out.{{ $fetcher.NextPageMarker }} != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
	{{- else }}
  out, err := s.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }})
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

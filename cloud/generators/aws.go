//go:generate go run $GOFILE
package main

import (
	"bytes"
	"go/format"
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

	formatted, err := format.Source(buff.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("../aws/generated_api.go", formatted, 0666); err != nil {
		panic(err)
	}
}

const funcsTempl = `// Auto generated implementation for the AWS cloud service
package aws

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "fmt"

  awssdk "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  "github.com/aws/aws-sdk-go/service/{{ $service.Api }}"
  "github.com/aws/aws-sdk-go/service/{{ $service.Api }}/{{ $service.Api }}iface"
  {{- end }}
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

var ServicePerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $service.Name }}",
  {{- end }}
{{- end }}
}

{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
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

func (s *{{ Title $service.Name }}) FetchByType(t string) (*graph.Graph, error) {
  switch t {
  {{- range $index, $fetcher := $service.Fetchers }}
  case "{{ $fetcher.ResourceType }}":
    return s.fetch_all_{{ $fetcher.ResourceType }}_graph()
  {{- end }}
  default:
    return nil, fmt.Errorf("aws {{ $service.Name }}: unsupported fetch for type %s", t)
  }
}

{{ range $index, $fetcher := $service.Fetchers }}
func (s *{{ Title $service.Name }}) fetch_all_{{ $fetcher.ResourceType }}() (interface{}, error) {
  return s.{{ $fetcher.ApiMethod }}(&{{ $service.Api }}.{{ $fetcher.Input }})
}
{{- end }}

{{ range $index, $fetcher := $service.Fetchers }}
func (s *{{ Title $service.Name }}) fetch_all_{{ $fetcher.ResourceType }}_graph() (*graph.Graph, error) {
  g := graph.NewGraph()
  out, err := s.fetch_all_{{ $fetcher.ResourceType }}()
  if err != nil {
    return nil, err
  }
  {{ if ne $fetcher.OutputsContainers "" }}
    for _, all := range out.(*{{ $service.Api }}.{{ $fetcher.Output }}).{{ $fetcher.OutputsContainers }} {
      for _, output := range all.{{ $fetcher.OutputsExtractor }} {
        res, err := newResource(output)
        if err != nil {
          return g, err
        }
        g.AddResource(res)
      }
    }
  {{ else }}
    for _, output := range out.(*{{ $service.Api }}.{{ $fetcher.Output }}).{{ $fetcher.OutputsExtractor }} {
      res, err := newResource(output)
      if err != nil {
        return g, err
      }
      g.AddResource(res)
    }
  {{ end }}
  return g, nil
}
{{ end }}

{{- if not $service.ManualGlobalFetch }}
type Aws{{ Title $service.Name }} struct {
{{- range $index, $fetcher := $service.Fetchers }}
  {{ $fetcher.ResourceType }}List   []*{{ $service.Api }}.{{ $fetcher.AWSType }}
{{- end }}
}

func (s *{{ Title $service.Name }}) global_fetch() (*Aws{{ Title $service.Name }}, error) {
	resultc, errc := multiFetch(
    {{- range $index, $fetcher := $service.Fetchers }}
      s.fetch_all_{{ $fetcher.ResourceType }},
    {{- end }}
  )

	awsService := &Aws{{ Title $service.Name }}{}

	for r := range resultc {
		switch rr := r.(type) {
    {{- range $index, $fetcher := $service.Fetchers }}
  case *{{ $service.Api }}.{{ $fetcher.Output }}:
      {{- if eq $fetcher.OutputsContainers "" }}
        awsService.{{ $fetcher.ResourceType }}List = append(awsService.{{ $fetcher.ResourceType }}List, rr.{{ $fetcher.OutputsExtractor }}...)
      {{- else }}
      for _, c := range rr.{{ $fetcher.OutputsContainers }} {
        awsService.{{ $fetcher.ResourceType }}List = append(awsService.{{ $fetcher.ResourceType }}List, c.{{ $fetcher.OutputsExtractor }}...)
      }
      {{- end }}
      {{- end }}
  	}
  }

	return awsService, <-errc
}
{{- end }}
{{ end }}`

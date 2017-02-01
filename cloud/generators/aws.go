//go:generate go run $GOFILE
package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/wallix/awless/cloud/aws"
)

func main() {
	generateFetcherFuncs()
}

func generateFetcherFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":   strings.Title,
		"ToUpper": strings.ToUpper,
	}).Parse(funcsTempl)

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("../aws/fetcher_gen_funcs.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	err = templ.Execute(f, aws.ServicesDefinitions)
	if err != nil {
		panic(err)
	}
}

const funcsTempl = `// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
  "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)
{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
  {{ $service.Api }}iface.{{ ToUpper $service.Api }}API
}

func New{{ Title $service.Name }}(sess *session.Session) *{{ Title $service.Name }} {
	return &{{ Title $service.Name }}{ {{ $service.Api }}.New(sess) }
}
{{ range $index, $fetcher := $service.Fetchers }}
func (s *{{ Title $service.Name }}) {{ Title $fetcher.ResourceType.PluralString }}() (interface{}, error) {
	return s.{{ $fetcher.ApiMethod }}(&{{ $service.Api }}.{{ $fetcher.Input }}{})
}
{{- end }}

type Aws{{ Title $service.Name }} struct {
{{ range $index, $fetcher := $service.Fetchers }}
  {{ Title $fetcher.ResourceType.PluralString }}   []*{{ $service.Api }}.{{ $fetcher.AWSType }}
{{- end }}
}

func (s *{{ Title $service.Name }}) FetchAws{{ Title $service.Name }}() (*Aws{{ Title $service.Name }}, error) {
	resultc, errc := multiFetch(
    {{- range $index, $fetcher := $service.Fetchers }}
      s.{{ Title $fetcher.ResourceType.PluralString }},
    {{- end }}
  )

	awsService := &Aws{{ Title $service.Name }}{}

	for r := range resultc {
		switch rr := r.(type) {
    {{- range $index, $fetcher := $service.Fetchers }}
  case *{{ $service.Api }}.{{ $fetcher.Output }}:
      {{- if eq $fetcher.OutputsContainers "" }}
        awsService.{{ Title $fetcher.ResourceType.PluralString }} = append(awsService.{{ Title $fetcher.ResourceType.PluralString }}, rr.{{ $fetcher.OutputsExtractor }}...)
      {{- else }}
      for _, c := range rr.{{ $fetcher.OutputsContainers }} {
        awsService.{{ Title $fetcher.ResourceType.PluralString }} = append(awsService.{{ Title $fetcher.ResourceType.PluralString }}, c.{{ $fetcher.OutputsExtractor }}...)
      }
      {{- end }}
      {{- end }}
  	}
  }

	return awsService, <-errc
}
{{- end }}`

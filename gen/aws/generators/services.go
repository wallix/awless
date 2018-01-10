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

func generateServicesFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"Join":           strings.Join,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(servicesTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.FetchersDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(SERVICES_DIR, "gen_services.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const servicesTempl = `// Auto generated implementation for the AWS cloud service

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
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/aws/fetch"
	tstore "github.com/wallix/triplestore"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	{{- range $index, $service := . }}
  "{{ $service.Name }}",
  {{- end }}
}

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

var APIPerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $fetcher.Api }}",
  {{- end }}
{{- end }}
}

{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
	fetcher fetch.Fetcher
  region, profile string
	config map[string]interface{}
	log *logger.Logger
	{{- range $, $api := $service.Api }}
		{{ $api }}iface.{{ ApiToInterface $api }}
	{{- end }}
}

func New{{ Title $service.Name }}(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
  {{- if $service.Global }}
	region := "global"
	{{- else}}
	region := awssdk.StringValue(sess.Config.Region)
	{{- end}}	

	{{- range $, $api := $service.Api }}
		{{$api }}API := {{ $api }}.New(sess)
	{{- end }}

	fetchConfig := awsfetch.NewConfig(
		{{- range $, $api := $service.Api }}
			{{$api }}API,
		{{- end }}
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &{{ Title $service.Name }}{ 
	{{- range $, $api := $service.Api }}
		{{ApiToInterface $api }}: {{ $api }}API,
	{{- end }}
		fetcher: fetch.NewFetcher(awsfetch.Build{{ Title $service.Name }}FetchFuncs(fetchConfig)),
		config: extraConf,
		region: region,
		profile: profile,
		log: log,
  }
}

func (s *{{ Title $service.Name }}) Name() string {
  return "{{ $service.Name }}"
}

func (s *{{ Title $service.Name }}) Region() string {
  return s.region
}

func (s *{{ Title $service.Name }}) Profile() string {
  return s.profile
}

func (s *{{ Title $service.Name }}) ResourceTypes() []string {
	return []string{
	{{- range $index, $fetcher := $service.Fetchers }}
		"{{ $fetcher.ResourceType }}",
	{{- end }}
	}
}

func (s *{{ Title $service.Name }}) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

  gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()
	
	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup

	{{- range $index, $fetcher := $service.Fetchers }}
	if getBool(s.config, "aws.{{ $service.Name }}.{{ $fetcher.ResourceType }}.sync", true) {
		list, err := s.fetcher.Get("{{ $fetcher.ResourceType }}_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*{{ $fetcher.AWSType }}); !ok {
			return gph, errors.New("cannot cast to '[]*{{ $fetcher.AWSType }}' type from fetch context")
		}
		for _, r := range list.([]*{{ $fetcher.AWSType }}) {
			for _, fn := range addParentsFns["{{ $fetcher.ResourceType }}"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *{{ $fetcher.AWSType }}) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
  {{- end }}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
				allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *{{ Title $service.Name }}) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
  return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *{{ Title $service.Name }}) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.{{ $service.Name }}.sync", true)
}

{{ end }}`

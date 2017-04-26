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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

func generateParamsDocLookup() {
	templ, err := template.New("paramsdocs").Parse(paramsdocsTempl)
	if err != nil {
		panic(err)
	}

	paramsDoc := loadAllRefs()

	doc := make(map[string]map[string]string)
	for _, def := range aws.DriversDefs {
		for _, driv := range def.Drivers {
			key := fmt.Sprintf("%s%s", driv.Action, driv.Entity)
			if doc[key] == nil {
				doc[key] = make(map[string]string)
			}

			params := doc[key]
			for _, p := range driv.RequiredParams {
				if s, ok := paramsDoc[fmt.Sprintf("%s$%s", driv.Input, p.AwsField)]; ok {
					params[p.TemplateName] = fmt.Sprint(s)
				}
			}

			for _, p := range driv.ExtraParams {
				if s, ok := paramsDoc[fmt.Sprintf("%s$%s", driv.Input, p.AwsField)]; ok {
					params[p.TemplateName] = fmt.Sprint(s)
				}
			}
		}
	}

	file, err := os.Create(filepath.Join(DOC_DIR, "gen_paramsdoc.go"))
	if err != nil {
		panic(err)
	}
	if err = templ.Execute(file, doc); err != nil {
		panic(err)
	}
}

var simpleTagRegex = regexp.MustCompile(`</?\w+>`)

func trimVal(v interface{}) (out string) {
	out = strings.SplitN(fmt.Sprint(v), ".", 2)[0]
	out = simpleTagRegex.ReplaceAllString(out, "")
	return
}

func cleanKey(s string) string {
	return strings.Replace(s, "Request$", "Input$", 1)
}

type entries struct {
	Shapes map[string]interface{} `json:"shapes"`
}

func loadAllRefs() map[string]string {
	fileRefs := []string{
		filepath.Join("ec2", "2016-11-15", "docs-2.json"),
		filepath.Join("iam", "2010-05-08", "docs-2.json"),
		filepath.Join("autoscaling", "2011-01-01", "docs-2.json"),
		filepath.Join("s3", "2006-03-01", "docs-2.json"),
		filepath.Join("sns", "2010-03-31", "docs-2.json"),
		filepath.Join("sqs", "2012-11-05", "docs-2.json"),
		filepath.Join("lambda", "2015-03-31", "docs-2.json"),
		filepath.Join("rds", "2014-10-31", "docs-2.json"),
		filepath.Join("route53", "2013-04-01", "docs-2.json"),
		filepath.Join("monitoring", "2010-08-01", "docs-2.json"),
		filepath.Join("elasticloadbalancingv2", "2015-12-01", "docs-2.json"),
		filepath.Join("sts", "2011-06-15", "docs-2.json"),
	}

	entriesC := make(chan *entries)

	var wg sync.WaitGroup

	for _, fileRef := range fileRefs {
		wg.Add(1)
		go func(ref string) {
			defer wg.Done()
			file, err := os.Open(filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "aws", "aws-sdk-go", "models", "apis", ref))
			if err != nil {
				panic(err)
			}

			all := new(entries)
			if err := json.NewDecoder(file).Decode(all); err != nil {
				panic(err)
			}

			entriesC <- all
		}(fileRef)
	}

	go func() {
		wg.Wait()
		close(entriesC)
	}()

	refs := make(map[string]string)

	for e := range entriesC {
		for _, val := range e.Shapes {
			if all, ok := val.(map[string]interface{}); ok {
				if allRefs, ok := all["refs"].(map[string]interface{}); ok {
					for k, v := range allRefs {
						refs[cleanKey(k)] = trimVal(v)
					}
				}
			}
		}
	}

	return refs
}

const paramsdocsTempl = `/* Copyright 2017 WALLIX

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
package awsdoc

var TemplateParamsDoc = map[string]map[string]string{
  {{- range $tplKey, $paramsDoc := . }}
  "{{ $tplKey }}": map[string]string {
    {{- range $param, $doc := $paramsDoc }}
    "{{$param}}": "{{$doc}}",
    {{- end }}
  },
  {{- end }}
}`

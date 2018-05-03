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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"unicode"
)

func generateParamsDocLookup() {
	templ, err := template.New("paramsdocs").Parse(paramsdocsTempl)
	if err != nil {
		panic(err)
	}

	paramsDoc := loadAllRefs()

	cmdsData := loadCommandStructs()

	doc := make(map[string]map[string]string)
	for _, cmd := range cmdsData {
		key := fmt.Sprintf("%s.%s", cmd.Action, cmd.Entity)
		if doc[key] == nil {
			doc[key] = make(map[string]string)
		}

		params := doc[key]
		for _, p := range cmd.Params {
			if s, ok := searchParamInDoc(paramsDoc, trimBeforeFirstDot(cmd.Input), p.AwsField); ok {
				params[p.Name] = fmt.Sprint(s)
			}
		}
	}

	writeTemplateToFile(templ, doc, DOC_DIR, "gen_paramsdoc.go")
}

func searchParamInDoc(paramsDoc map[string]string, input, field string) (string, bool) {
	var lowerField string
	if len(field) > 0 {
		fieldRules := []rune(field)
		fieldRules[0] = unicode.ToLower(fieldRules[0])
		lowerField = string(fieldRules)
	}
	if s, ok := paramsDoc[fmt.Sprintf("%s$%s", input, field)]; ok {
		return s, ok
	}
	if s, ok := paramsDoc[fmt.Sprintf("%s$%s", inputToRequestKey(input), field)]; ok {
		return s, ok
	}
	if s, ok := paramsDoc[fmt.Sprintf("%s$%s", inputToRequestKey(input), lowerField)]; ok {
		return s, ok
	}
	if s, ok := paramsDoc[fmt.Sprintf("%s$%s", inputToTypeKey(input), field)]; ok {
		return s, ok
	}
	if s, ok := paramsDoc[fmt.Sprintf("%s$%s", dbInstanceKey(input), field)]; ok {
		return s, ok
	}
	return "", false
}

var simpleTagRegex = regexp.MustCompile(`</?\w+>`)
var bracketTextRegex = regexp.MustCompile(`\[[\w-]+\]`)

func trimVal(v interface{}) (out string) {
	out = strings.SplitN(fmt.Sprint(v), ".", 2)[0]
	out = simpleTagRegex.ReplaceAllString(out, "")
	out = bracketTextRegex.ReplaceAllString(out, "")
	out = strings.TrimSpace(out)
	return
}

func trimBeforeFirstDot(s string) string {
	splits := strings.Split(s, ".")
	switch len(splits) {
	case 0, 1:
		return s
	default:
		return strings.Join(splits[1:], ".")

	}
}

func inputToRequestKey(s string) string {
	return strings.Replace(s, "Input", "Request", 1)
}

func inputToTypeKey(s string) string {
	return strings.Replace(s, "Input", "Type", 1)
}

func dbInstanceKey(s string) string {
	if strings.Contains(s, "DBInstanceInput") {
		return "DBInstance"
	}
	return s
}

type entries struct {
	fileRef string
	Shapes  map[string]interface{} `json:"shapes"`
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
		filepath.Join("elasticloadbalancing", "2012-06-01", "docs-2.json"),
		filepath.Join("sts", "2011-06-15", "docs-2.json"),
		filepath.Join("cloudformation", "2010-05-15", "docs-2.json"),
		filepath.Join("ecr", "2015-09-21", "docs-2.json"),
		filepath.Join("ecs", "2014-11-13", "docs-2.json"),
		filepath.Join("application-autoscaling", "2016-02-06", "docs-2.json"),
		filepath.Join("acm", "2015-12-08", "docs-2.json"),
	}

	apisPath := filepath.Join(ROOT_DIR, "vendor", "github.com", "aws", "aws-sdk-go", "models", "apis")
	log.Printf("looking up local AWS services doc at %s", relativePathToRoot(apisPath))

	entriesC := make(chan *entries)

	var wg sync.WaitGroup
	for _, fileRef := range fileRefs {
		wg.Add(1)
		go func(ref string) {
			defer wg.Done()
			file, err := os.Open(filepath.Join(apisPath, ref))
			if err != nil {
				panic(err)
			}

			all := &entries{fileRef: ref}
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
						refs[k] = trimVal(v)
					}
				}
			}
		}
		log.Printf("\t %s", e.fileRef)
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

var generatedParamsDoc = map[string]map[string]string{
  {{- range $tplKey, $paramsDoc := . }}
  "{{ $tplKey }}": map[string]string {
    {{- range $param, $doc := $paramsDoc }}
    "{{$param}}": "{{$doc}}",
    {{- end }}
  },
  {{- end }}
}`

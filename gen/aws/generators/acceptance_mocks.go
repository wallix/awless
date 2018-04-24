package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

var AWS_SDK_PATH = filepath.Join(ROOT_DIR, "vendor", "github.com", "aws", "aws-sdk-go")

func generateAcceptanceMocks() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, SPEC_DIR, func(os.FileInfo) bool { return true }, 0)
	if err != nil {
		panic(err)
	}

	finder := &findStructs{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Walk(finder, f)
		}
	}

	usedApis := make(map[string]string)
	for _, cmd := range finder.result {
		if cmd.API == "" {
			continue
		}
		usedApis[cmd.API] = aws.ApiToInterface(cmd.API)
	}

	apis := make(map[string]apiInfo)

	for api, ifaceName := range usedApis {
		var functions []functionInfo
		apiPath := filepath.Join(AWS_SDK_PATH, "service", api, api+"iface")
		pkgs, err := parser.ParseDir(fset, apiPath, func(os.FileInfo) bool { return true }, 0)
		if err != nil {
			panic(err)
		}
		conf := types.Config{Importer: importer.Default()}

		//ifaceFinder := &findIfaces{}
		for _, pkg := range pkgs {
			files := make([]*ast.File, len(pkg.Files))
			i := 0
			for _, file := range pkg.Files {
				files[i] = file
				i++
			}
			tpkg, err := conf.Check(api, fset, files, nil)
			if err != nil {
				panic(err)
			}

			iface := tpkg.Scope().Lookup(ifaceName)
			if iface == nil {
				panic(fmt.Sprintf("cannot find interface %s", ifaceName))
			}
			if !types.IsInterface(iface.Type()) {
				panic(fmt.Sprintf("%s (%s) not an interface", ifaceName, iface.Type().String()))
			}
			isPointer := func(s string) bool { return len(s) > 0 && s[0] == '*' }
			iiface := iface.Type().Underlying().(*types.Interface).Complete()
			for i := 0; i < iiface.NumMethods(); i++ {
				meth := iiface.Method(i)
				sig := meth.Type().(*types.Signature)
				if strings.HasSuffix(meth.Name(), "Pages") || strings.HasSuffix(meth.Name(), "PagesWithContext") {
					continue
				}
				var paramBuff bytes.Buffer
				var paramNames []string
				for j := 0; j < sig.Params().Len(); j++ {
					p := sig.Params().At(j)
					paramName := fmt.Sprintf("param%d", j)
					paramBuff.WriteString(paramName)
					paramBuff.WriteRune(' ')
					t := p.Type().String()
					if found := strings.LastIndexByte(t, '/'); found != -1 {
						t = t[found+1:]
					}
					if isPointer(p.Type().String()) {
						t = "*" + t
					}
					if sig.Variadic() && j == sig.Params().Len()-1 {
						t = "..." + t
						paramName = paramName + "..."
					}
					paramNames = append(paramNames, paramName)
					paramBuff.WriteString(t)
					if j < sig.Params().Len()-1 {
						paramBuff.WriteString(", ")
					}
				}
				var returnsBuff bytes.Buffer
				if sig.Results().Len() > 1 {
					returnsBuff.WriteByte('(')
				}
				for j := 0; j < sig.Results().Len(); j++ {
					p := sig.Results().At(j)
					t := p.Type().String()
					if found := strings.LastIndexByte(t, '/'); found != -1 {
						t = t[found+1:]
					}
					if isPointer(p.Type().String()) {
						t = "*" + t
					}
					if sig.Variadic() && j == sig.Params().Len()-1 {
						t = "..." + t
					}
					returnsBuff.WriteString(t)
					if j < sig.Results().Len()-1 {
						returnsBuff.WriteString(", ")
					}
				}
				if sig.Results().Len() > 1 {
					returnsBuff.WriteByte(')')
				}
				functions = append(functions, functionInfo{
					Name:         meth.Name(),
					Sig:          fmt.Sprintf("func (m *%sMock) %s(%s) %s", api, meth.Name(), paramBuff.String(), returnsBuff.String()),
					AnonymousSig: fmt.Sprintf("func (%s) %s", paramBuff.String(), returnsBuff.String()),
					ParamNames:   paramNames,
				})
			}

		}
		apis[api] = apiInfo{Name: api, IfaceName: ifaceName, Methods: functions}
	}

	templ, err := template.New("mocks").Funcs(
		template.FuncMap{"Join": strings.Join},
	).Parse(atMocksTemplate)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, apis, AWSAT_DIR, "gen_mocks.go")
}

func generateAcceptanceFactory() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, SPEC_DIR, func(os.FileInfo) bool { return true }, 0)
	if err != nil {
		panic(err)
	}

	finder := &findStructs{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Walk(finder, f)
		}
	}

	templ, err := template.New("acceptanceFactory").Funcs(
		template.FuncMap{
			"ApiToInterface": aws.ApiToInterface,
		},
	).Parse(atMocksCmdBuilders)
	if err != nil {
		panic(err)
	}

	writeTemplateToFile(templ, finder.result, AWSAT_DIR, "gen_factory.go")
}

type apiInfo struct {
	Name      string
	IfaceName string
	Methods   []functionInfo
}

type functionInfo struct {
	Name         string
	Sig          string
	AnonymousSig string
	ParamNames   []string
}

const atMocksCmdBuilders = `/* Copyright 2017 WALLIX

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
package awsat

import (
  "github.com/wallix/awless/aws/spec"
)

type AcceptanceFactory struct {
	Mock   interface{}
	Logger *logger.Logger
	Graph cloud.GraphAPI
}

func NewAcceptanceFactory(mock interface{}, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	logger := logger.DiscardLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph:g, Logger: logger}
}

func (f *AcceptanceFactory) Build(key string) func() interface{} {
	switch key {
		{{- range $cmdName, $cmd := . }}
		case "{{ $cmd.Action }}{{ $cmd.Entity }}":
			return func() interface{} {
				cmd := awsspec.New{{ $cmdName }}(nil, f.Graph, f.Logger)
				cmd.SetApi(f.Mock.({{$cmd.API}}iface.{{ ApiToInterface $cmd.API }}))
				return cmd
			}
		{{- end}}
	}
	return nil
}
`

const atMocksTemplate = `/* Copyright 2017 WALLIX

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
package awsat


{{ range $api, $apiInfo := . }}

type {{ $api }}Mock struct {
  basicMock
  {{ $api }}iface.{{ $apiInfo.IfaceName }}
  {{- range $, $method := $apiInfo.Methods }}
  {{ $method.Name }}Func {{ $method.AnonymousSig }}
  {{- end }}
}

{{- range $, $method := $apiInfo.Methods }}
{{ $method.Sig }} {
	m.addCall("{{$method.Name}}")
	m.verifyInput("{{$method.Name}}", param0)
	return m.{{$method.Name}}Func({{Join $method.ParamNames ","}})
}
{{ end }}
{{- end }}

`

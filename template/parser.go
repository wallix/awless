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

package template

import "github.com/wallix/awless/template/ast"

func Parse(text string) (*Template, error) {
	p := &ast.Peg{AST: &ast.AST{}, Buffer: string(text), Pretty: true}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	return &Template{AST: p.AST}, nil
}

func ParseStatement(text string) (ast.Node, error) {
	templ, err := Parse(text)
	if err != nil {
		return nil, err
	}

	return templ.Statements[0].Node, nil
}

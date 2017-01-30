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

package script

import "github.com/wallix/awless/script/ast"

func Parse(text string) (*Script, error) {
	p := &ast.Peg{AST: &ast.AST{}, Buffer: string(text), Pretty: true}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	return &Script{p.AST}, nil
}

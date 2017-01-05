package script

import "github.com/wallix/awless/script/ast"

func Parse(text string) (*ast.Script, error) {
	p := &ast.Peg{Script: &ast.Script{}, Buffer: string(text), Pretty: true}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil, err
	}
	p.Execute()

	return p.Script, nil
}

package template

import (
	"encoding/json"
	"errors"

	"github.com/wallix/awless/template/ast"
)

type toJSON struct {
	ID       string    `json:"id"`
	Commands []command `json:"commands"`
}

type command struct {
	Line    string   `json:"line"`
	Errors  []string `json:"errors,omitempty"`
	Results []string `json:"results,omitempty"`
}

func (t *Template) MarshalJSON() ([]byte, error) {
	out := &toJSON{}
	out.ID = t.ID
	out.Commands = []command{}

	for _, cmd := range t.CommandNodesIterator() {
		newCmd := command{}
		newCmd.Line = cmd.String()
		if cmd.CmdErr != nil {
			newCmd.Errors = append(newCmd.Errors, cmd.CmdErr.Error())
		}
		if cmd.CmdResult != nil {
			if s, ok := cmd.CmdResult.(string); ok {
				newCmd.Results = append(newCmd.Results, s)
			}
		}
		out.Commands = append(out.Commands, newCmd)
	}

	return json.Marshal(out)
}

func (t *Template) UnmarshalJSON(b []byte) error {
	var v toJSON

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	tt := &Template{ID: v.ID, AST: &ast.AST{
		Statements: make([]*ast.Statement, 0),
	}}

	for _, c := range v.Commands {
		node, err := parseStatement(c.Line)
		if err != nil {
			return err
		}

		switch node.(type) {
		case *ast.CommandNode:
			n := node.(*ast.CommandNode)
			if len(c.Results) > 0 {
				n.CmdResult = c.Results[0]
			}
			if len(c.Errors) > 0 {
				n.CmdErr = errors.New(c.Errors[0])
			}
			tt.Statements = append(tt.Statements, &ast.Statement{n})
		}
	}

	*t = *tt

	return nil
}

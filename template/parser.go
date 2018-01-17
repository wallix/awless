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

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wallix/awless/template/internal/ast"
)

func Parse(text string) (tmpl *Template, err error) {
	defer func() { // as peg lib does not allow errors in Execute, we use panic to build the AST
		if rerr := recover(); rerr != nil {
			switch rerr.(type) {
			case error:
				err = fmt.Errorf("template parsing: %s", rerr.(error))
			default:
				panic(rerr)
			}
		}
	}()

	if clean := strings.TrimSpace(text); clean == "" {
		return nil, errors.New("empty template")
	}

	tmpl = &Template{}

	p := &ast.Peg{AST: &ast.AST{}, Buffer: string(text)}
	p.Init()

	if err = p.Parse(); err != nil {
		err = newParseError(text, err.Error())
		return
	}
	p.Execute()

	tmpl.AST = p.AST

	return
}

func MustParse(text string) *Template {
	t, err := Parse(text)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseParams(text string) (map[string]interface{}, error) {
	node, err := parseParamsAsCommandNode(text)
	if err != nil {
		return nil, err
	}
	return node.ToFillerParams(), nil
}

func parseParamsAsCommandNode(text string) (*ast.CommandNode, error) {
	full := fmt.Sprintf("none none %s", text)
	n, err := parseStatement(full)
	if err != nil {
		return nil, fmt.Errorf("parse params: %s", err)
	}

	switch n.(type) {
	case *ast.CommandNode:
		return (n.(*ast.CommandNode)), nil
	default:
		return nil, fmt.Errorf("parse params: expected a command node")
	}
}

func parseStatement(text string) (ast.Node, error) {
	templ, err := Parse(text)
	if err != nil {
		return nil, err
	}

	return templ.Statements[0].Node, nil
}

type parseError struct {
	origMsg          string
	lines            []string
	line, start, end int
}

func newParseError(templText, pegErrMsg string) (perr *parseError) {
	perr = buildParseError(pegErrMsg)
	perr.origMsg = pegErrMsg
	perr.lines = strings.Split(templText, "\n")
	return
}

func (pe *parseError) Error() string {
	if pe.invalidIndexes() {
		return pe.origMsg
	}

	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("error parsing template at line %d (char %d):\n", pe.line, pe.start))

	for i, l := range pe.lines {
		buff.WriteByte('\t')
		if pe.line == i+1 {
			buff.WriteString("-> ")
			buff.WriteString(l[0:pe.start])
		} else {
			buff.WriteString("   ")
			buff.WriteString(l)
		}
		if i < len(pe.lines)-1 {
			buff.WriteByte('\n')
		}
	}

	return buff.String()
}

func (pe *parseError) invalidIndexes() bool {
	if pe.line == 0 {
		return true
	}
	if pe.line > len(pe.lines) {
		return true
	}
	if pe.start > len(pe.lines[pe.line-1]) {
		return true
	}

	return false
}

// ex: parse error near Equal (line 1 symbol 21 - line 1 symbol 22)
var indexesRegex = regexp.MustCompile(`line\s+(\d{1,})\s+symbol\s+(\d{1,})\s+-\s+line \d{1,}\s+symbol\s+(\d{1,3})`)

func buildParseError(s string) (perr *parseError) {
	perr = &parseError{}

	matches := indexesRegex.FindStringSubmatch(s)
	if len(matches) != 4 {
		return
	}

	toInt := func(s string) (i int) {
		i, _ = strconv.Atoi(s)
		return
	}

	perr.line, perr.start, perr.end = toInt(matches[1]), toInt(matches[2]), toInt(matches[3])

	return
}

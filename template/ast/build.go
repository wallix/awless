package ast

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func (a *AST) addAction(text string) {
	cmd := &CommandNode{Action: text}
	decl := a.currentDeclaration()
	if decl != nil {
		decl.Expr = cmd
	} else {
		node := a.currentCommand()
		if node == nil {
			a.addStatement(cmd)
		} else {
			node.Action = text
		}
	}
}

func (a *AST) addEntity(text string) {
	node := a.currentCommand()
	node.Entity = text
}

func (a *AST) addDeclarationIdentifier(text string) {
	a.addStatement(&DeclarationNode{Ident: text})
}

func (a *AST) LineDone() {
	a.currentStatement = nil
	a.currentKey = ""
}

func (a *AST) addParamKey(text string) {
	node := a.currentCommand()
	if node.Params == nil {
		node.Refs = make(map[string]string)
		node.Params = make(map[string]interface{})
		node.Holes = make(map[string]string)
	}
	a.currentKey = text
}

func (a *AST) addParamValue(text string) {
	node := a.currentCommand()
	node.Params[a.currentKey] = text
}

func (a *AST) addCsvValue(text string) {
	var csv []string
	for _, val := range strings.Split(text, ",") {
		csv = append(csv, strings.TrimSpace(val))
	}
	node := a.currentCommand()
	node.Params[a.currentKey] = csv
}

func (a *AST) addParamIntValue(text string) {
	node := a.currentCommand()
	num, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to int", text))
	}
	node.Params[a.currentKey] = num
}

func (a *AST) addParamCidrValue(text string) {
	node := a.currentCommand()
	_, ipnet, err := net.ParseCIDR(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to net cidr", text))
	}
	node.Params[a.currentKey] = ipnet.String()
}

func (a *AST) addParamIpValue(text string) {
	node := a.currentCommand()
	ip := net.ParseIP(text)
	if ip == nil {
		panic(fmt.Sprintf("cannot convert '%s' to net ip", text))
	}
	node.Params[a.currentKey] = ip.String()
}

func (a *AST) addParamRefValue(text string) {
	node := a.currentCommand()
	node.Refs[a.currentKey] = text
}

func (a *AST) addParamHoleValue(text string) {
	node := a.currentCommand()
	node.Holes[a.currentKey] = text
}

func (a *AST) currentDeclaration() *DeclarationNode {
	st := a.currentStatement
	if st == nil {
		return nil
	}

	switch st.Node.(type) {
	case *DeclarationNode:
		return st.Node.(*DeclarationNode)
	}

	return nil
}

func (a *AST) currentCommand() *CommandNode {
	st := a.currentStatement
	if st == nil {
		return nil
	}

	switch st.Node.(type) {
	case *CommandNode:
		return st.Node.(*CommandNode)
	case *DeclarationNode:
		expr := st.Node.(*DeclarationNode).Expr
		switch expr.(type) {
		case *CommandNode:
			return expr.(*CommandNode)
		}
		return nil
	default:
		panic("last expression: unexpected node type")
	}
}

func (a *AST) addStatement(n Node) {
	stat := &Statement{Node: n}
	a.currentStatement = stat
	a.Statements = append(a.Statements, stat)
}

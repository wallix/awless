package ast

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func (a *AST) addAction(text string) {
	if IsInvalidAction(text) {
		panic(fmt.Errorf("unknown action '%s'", text))
	}

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
	if IsInvalidEntity(text) {
		panic(fmt.Errorf("unknown entity '%s'", text))
	}
	node := a.currentCommand()
	node.Entity = text
}

func (a *AST) addValue() {
	val := &ValueNode{}

	decl := a.currentDeclaration()
	if decl != nil {
		decl.Expr = val
	}
}

func (a *AST) addDeclarationIdentifier(text string) {
	a.addStatement(&DeclarationNode{Ident: text})
}

func (a *AST) LineDone() {
	a.currentStatement = nil
	a.currentKey = ""
}

func (a *AST) addParam(i interface{}) {
	if node := a.currentCommand(); node != nil {
		node.Params[a.currentKey] = i
	} else {
		varDecl := a.currentDeclarationValue()
		varDecl.Value = i
	}
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

func (a *AST) addAliasParam(text string) {
	a.addParam("@" + text)
}

func (a *AST) addParamValue(text string) {
	a.addParam(text)
}

func (a *AST) addCsvValue(text string) {
	var csv []string
	for _, val := range strings.Split(text, ",") {
		csv = append(csv, strings.TrimSpace(val))
	}
	a.addParam(csv)
}

func (a *AST) addParamIntValue(text string) {
	num, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to int", text))
	}
	a.addParam(num)
}

func (a *AST) addParamCidrValue(text string) {
	_, ipnet, err := net.ParseCIDR(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to net cidr", text))
	}
	a.addParam(ipnet.String())
}

func (a *AST) addParamIpValue(text string) {
	ip := net.ParseIP(text)
	if ip == nil {
		panic(fmt.Sprintf("cannot convert '%s' to net ip", text))
	}
	a.addParam(ip.String())
}

func (a *AST) addParamRefValue(text string) {
	if node := a.currentCommand(); node != nil {
		node.Refs[a.currentKey] = text
	}
}

func (a *AST) addParamHoleValue(text string) {
	if node := a.currentCommand(); node != nil {
		node.Holes[a.currentKey] = text
	} else {
		varDecl := a.currentDeclarationValue()
		varDecl.Hole = text
	}
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
		return nil
	}
}

func (a *AST) currentDeclarationValue() *ValueNode {
	st := a.currentStatement
	if st == nil {
		return nil
	}

	switch st.Node.(type) {
	case *DeclarationNode:
		expr := st.Node.(*DeclarationNode).Expr
		switch expr.(type) {
		case *ValueNode:
			return expr.(*ValueNode)
		}
		return nil
	default:
		return nil
	}
}

func (a *AST) addStatement(n Node) {
	stat := &Statement{Node: n}
	a.currentStatement = stat
	a.Statements = append(a.Statements, stat)
}

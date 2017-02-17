package ast

import (
	"fmt"
	"net"
	"strconv"
)

func (a *AST) addAction(text string) {
	expr := a.currentExpression()
	if expr == nil {
		a.addStatement(&ExpressionNode{Action: text})
	} else {
		expr.Action = text
	}
}

func (a *AST) addEntity(text string) {
	expr := a.currentExpression()
	expr.Entity = text
}

func (a *AST) addDeclarationIdentifier(text string) {
	decl := &DeclarationNode{
		Left:  &IdentifierNode{Ident: text},
		Right: &ExpressionNode{},
	}
	a.addStatement(decl)
}

func (a *AST) LineDone() {
	a.currentStatement = nil
	a.currentKey = ""
}

func (a *AST) addParamKey(text string) {
	expr := a.currentExpression()
	if expr.Params == nil {
		expr.Refs = make(map[string]string)
		expr.Params = make(map[string]interface{})
		expr.Aliases = make(map[string]string)
		expr.Holes = make(map[string]string)
	}
	a.currentKey = text
}

func (a *AST) addParamValue(text string) {
	expr := a.currentExpression()
	expr.Params[a.currentKey] = text
}

func (a *AST) addParamIntValue(text string) {
	expr := a.currentExpression()
	num, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to int", text))
	}
	expr.Params[a.currentKey] = num
}

func (a *AST) addParamCidrValue(text string) {
	expr := a.currentExpression()
	_, ipnet, err := net.ParseCIDR(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to net cidr", text))
	}
	expr.Params[a.currentKey] = ipnet.String()
}

func (a *AST) addParamIpValue(text string) {
	expr := a.currentExpression()
	ip := net.ParseIP(text)
	if ip == nil {
		panic(fmt.Sprintf("cannot convert '%s' to net ip", text))
	}
	expr.Params[a.currentKey] = ip.String()
}

func (a *AST) addParamRefValue(text string) {
	expr := a.currentExpression()
	expr.Refs[a.currentKey] = text
}

func (a *AST) addParamAliasValue(text string) {
	expr := a.currentExpression()
	expr.Aliases[a.currentKey] = text
}

func (a *AST) addParamHoleValue(text string) {
	expr := a.currentExpression()
	expr.Holes[a.currentKey] = text
}

func (a *AST) currentExpression() *ExpressionNode {
	st := a.currentStatement
	if st == nil {
		return nil
	}

	switch st.Node.(type) {
	case *ExpressionNode:
		return st.Node.(*ExpressionNode)
	case *DeclarationNode:
		return st.Node.(*DeclarationNode).Right
	default:
		panic("last expression: unexpected node type")
	}
}

func (a *AST) addStatement(n Node) {
	stat := &Statement{Node: n}
	a.currentStatement = stat
	a.Statements = append(a.Statements, stat)
}

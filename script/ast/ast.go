package ast

import (
	"fmt"
	"net"
	"strconv"
)

type Node interface {
	clone() Node
}

type AST struct {
	Statements []Node

	currentStatement Node
	currentKey       string
}

type IdentifierNode struct {
	Ident string
	Val   interface{}
}

func (n *IdentifierNode) clone() Node {
	return &IdentifierNode{
		Ident: n.Ident,
		Val:   n.Val,
	}
}

type DeclarationNode struct {
	Left  *IdentifierNode
	Right *ExpressionNode
}

func (n *DeclarationNode) clone() Node {
	return &DeclarationNode{
		Left:  n.Left.clone().(*IdentifierNode),
		Right: n.Right.clone().(*ExpressionNode),
	}
}

type ExpressionNode struct {
	Action, Entity string
	Params         map[string]interface{}
	Holes          map[string]string
}

func (n *ExpressionNode) clone() Node {
	expr := &ExpressionNode{
		Action: n.Action, Entity: n.Entity,
		Params: make(map[string]interface{}),
		Holes:  make(map[string]string),
	}

	for k, v := range n.Params {
		expr.Params[k] = v
	}
	for k, v := range n.Holes {
		expr.Holes[k] = v
	}

	return expr
}

func (n *ExpressionNode) ProcessHoles(fills map[string]interface{}) {
	for key, hole := range n.Holes {
		if val, ok := fills[hole]; ok {
			if n.Params == nil {
				n.Params = make(map[string]interface{})
			}
			n.Params[key] = val
			delete(n.Holes, key)
		}
	}
}

func (s *AST) AddAction(text string) {
	expr := s.currentExpression()
	if expr == nil {
		s.addStatement(&ExpressionNode{Action: text})
	} else {
		expr.Action = text
	}
}

func (s *AST) AddEntity(text string) {
	expr := s.currentExpression()
	expr.Entity = text
}

func (s *AST) AddDeclarationIdentifier(text string) {
	decl := &DeclarationNode{
		Left:  &IdentifierNode{Ident: text},
		Right: &ExpressionNode{},
	}
	s.addStatement(decl)
}

func (s *AST) EndOfParams() {
	s.currentStatement = nil
	s.currentKey = ""
}

func (s *AST) AddParamKey(text string) {
	expr := s.currentExpression()
	if expr.Params == nil {
		expr.Params = make(map[string]interface{})
		expr.Holes = make(map[string]string)
	}
	s.currentKey = text
}

func (s *AST) AddParamValue(text string) {
	expr := s.currentExpression()
	expr.Params[s.currentKey] = text
}

func (s *AST) AddParamIntValue(text string) {
	expr := s.currentExpression()
	num, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to int", text))
	}
	expr.Params[s.currentKey] = num
}

func (s *AST) AddParamCidrValue(text string) {
	expr := s.currentExpression()
	_, ipnet, err := net.ParseCIDR(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert '%s' to net cidr", text))
	}
	expr.Params[s.currentKey] = ipnet.String()
}

func (s *AST) AddParamIpValue(text string) {
	expr := s.currentExpression()
	ip := net.ParseIP(text)
	if ip == nil {
		panic(fmt.Sprintf("cannot convert '%s' to net ip", text))
	}
	expr.Params[s.currentKey] = ip.String()
}

func (s *AST) AddParamHoleValue(text string) {
	expr := s.currentExpression()
	expr.Holes[s.currentKey] = text
}

func (s *AST) currentExpression() *ExpressionNode {
	st := s.currentStatement
	if st == nil {
		return nil
	}

	switch st.(type) {
	case *ExpressionNode:
		return st.(*ExpressionNode)
	case *DeclarationNode:
		return st.(*DeclarationNode).Right
	default:
		panic("last expression: unexpected node type")
	}
}

func (a *AST) Clone() *AST {
	clone := &AST{}
	for _, node := range a.Statements {
		clone.Statements = append(clone.Statements, node.clone())
	}
	return clone
}

func (s *AST) addStatement(n Node) {
	s.currentStatement = n
	s.Statements = append(s.Statements, n)
}

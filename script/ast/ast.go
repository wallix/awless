package ast

import (
	"fmt"
	"strconv"
)

type Node interface{}

type Script struct {
	Statements []Node

	currentStatement Node
	currentKey       string
}

type IdentifierNode struct {
	Ident string
	Val   interface{}
}

type DeclarationNode struct {
	Left  *IdentifierNode
	Right *ExpressionNode
}

type ExpressionNode struct {
	Action, Entity string
	Params         map[string]interface{}
	Holes          map[string]string
}

func (s *Script) AddAction(text string) {
	expr := s.currentExpression()
	if expr == nil {
		s.addStatement(&ExpressionNode{Action: text})
	} else {
		expr.Action = text
	}
}

func (s *Script) AddEntity(text string) {
	expr := s.currentExpression()
	expr.Entity = text
}

func (s *Script) AddDeclarationIdentifier(text string) {
	decl := &DeclarationNode{
		Left:  &IdentifierNode{Ident: text},
		Right: &ExpressionNode{},
	}
	s.addStatement(decl)
}

func (s *Script) EndOfParams() {
	s.currentStatement = nil
	s.currentKey = ""
}

func (s *Script) AddParamKey(text string) {
	expr := s.currentExpression()
	if expr.Params == nil {
		expr.Params = make(map[string]interface{})
		expr.Holes = make(map[string]string)
	}
	s.currentKey = text
}

func (s *Script) AddParamValue(text string) {
	expr := s.currentExpression()
	expr.Params[s.currentKey] = text
}

func (s *Script) AddParamIntValue(text string) {
	expr := s.currentExpression()
	num, err := strconv.Atoi(text)
	if err != nil {
		panic(fmt.Sprintf("cannot convert %s to int", text))
	}
	expr.Params[s.currentKey] = num
}

func (s *Script) AddParamHoleValue(text string) {
	expr := s.currentExpression()
	expr.Holes[s.currentKey] = text
}

func (s *Script) currentExpression() *ExpressionNode {
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

func (s *Script) addStatement(n Node) {
	s.currentStatement = n
	s.Statements = append(s.Statements, n)
}

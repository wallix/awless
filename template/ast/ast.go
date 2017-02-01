package ast

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Node interface {
	clone() Node
	String() string
}

type Statement struct {
	Node
	Result interface{}
	Line   string
	Err    error
}

func (s *Statement) clone() *Statement {
	newStat := &Statement{}
	newStat.Node = s.Node.clone()
	newStat.Result = s.Result
	newStat.Err = s.Err

	return newStat
}

type AST struct {
	Statements []*Statement

	currentStatement *Statement
	currentKey       string
}

func (a *AST) String() string {
	var all []string
	for _, stat := range a.Statements {
		all = append(all, stat.String())
	}
	return strings.Join(all, "\n")
}

type IdentifierNode struct {
	Ident string
	Val   interface{}
}

func (n *IdentifierNode) String() string {
	return fmt.Sprintf("%s", n.Ident)
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

func (n *DeclarationNode) String() string {
	return fmt.Sprintf("%s = %s", n.Left, n.Right)
}

type ExpressionNode struct {
	Action, Entity string
	Refs           map[string]string
	Params         map[string]interface{}
	Aliases        map[string]string
	Holes          map[string]string
}

func (n *ExpressionNode) clone() Node {
	expr := &ExpressionNode{
		Action: n.Action, Entity: n.Entity,
		Refs:    make(map[string]string),
		Params:  make(map[string]interface{}),
		Aliases: make(map[string]string),
		Holes:   make(map[string]string),
	}

	for k, v := range n.Refs {
		expr.Refs[k] = v
	}
	for k, v := range n.Params {
		expr.Params[k] = v
	}
	for k, v := range n.Aliases {
		expr.Aliases[k] = v
	}
	for k, v := range n.Holes {
		expr.Holes[k] = v
	}

	return expr
}

func (n *ExpressionNode) String() string {
	var all []string
	for k, v := range n.Refs {
		all = append(all, fmt.Sprintf("%s=$%v", k, v))
	}
	for k, v := range n.Params {
		all = append(all, fmt.Sprintf("%s=%v", k, v))
	}
	for k, v := range n.Aliases {
		all = append(all, fmt.Sprintf("%s=@%s", k, v))
	}
	for k, v := range n.Holes {
		all = append(all, fmt.Sprintf("%s={%s}", k, v))
	}
	return fmt.Sprintf("%s %s %s", n.Action, n.Entity, strings.Join(all, " "))
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

func (n *ExpressionNode) ProcessRefs(fills map[string]interface{}) {
	for key, ref := range n.Refs {
		if val, ok := fills[ref]; ok {
			if n.Params == nil {
				n.Params = make(map[string]interface{})
			}
			n.Params[key] = val
			delete(n.Refs, key)
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

func (s *AST) LineDone() {
	s.currentStatement = nil
	s.currentKey = ""
}

func (s *AST) AddParamKey(text string) {
	expr := s.currentExpression()
	if expr.Params == nil {
		expr.Refs = make(map[string]string)
		expr.Params = make(map[string]interface{})
		expr.Aliases = make(map[string]string)
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

func (s *AST) AddParamRefValue(text string) {
	expr := s.currentExpression()
	expr.Refs[s.currentKey] = text
}

func (s *AST) AddParamAliasValue(text string) {
	expr := s.currentExpression()
	expr.Aliases[s.currentKey] = text
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

	switch st.Node.(type) {
	case *ExpressionNode:
		return st.Node.(*ExpressionNode)
	case *DeclarationNode:
		return st.Node.(*DeclarationNode).Right
	default:
		panic("last expression: unexpected node type")
	}
}

func (a *AST) Clone() *AST {
	clone := &AST{}
	for _, stat := range a.Statements {
		clone.Statements = append(clone.Statements, stat.clone())
	}
	return clone
}

func (s *AST) addStatement(n Node) {
	stat := &Statement{Node: n}
	s.currentStatement = stat
	s.Statements = append(s.Statements, stat)
}

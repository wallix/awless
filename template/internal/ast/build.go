package ast

import (
	"fmt"
	"strconv"
)

type statementBuilder struct {
	action                string
	entity                string
	declarationIdentifier string
	isValue               bool
	newparams             map[string]interface{}
	currentKey            string
	currentNode           interface{}
	listBuilder           *listValueBuilder
	concatenationBuilder  *concatenationValueBuilder
}

func (b *statementBuilder) build() *Statement {
	if b.action == "" && b.entity == "" && b.declarationIdentifier == "" && !b.isValue {
		return nil
	}
	var expr ExpressionNode
	if b.isValue {
		expr = &RightExpressionNode{i: b.currentNode}
	} else {
		if b.newparams == nil {
			b.newparams = make(map[string]interface{})
		}
		expr = &CommandNode{
			Action:     b.action,
			Entity:     b.entity,
			ParamNodes: b.newparams,
			Refs:       make(map[string]interface{}),
		}
	}
	if b.declarationIdentifier != "" {
		decl := &DeclarationNode{Ident: b.declarationIdentifier, Expr: expr}
		return &Statement{Node: decl}
	}
	return &Statement{Node: expr}
}

func (b *statementBuilder) addParamKey(key string) *statementBuilder {
	b.currentKey = key
	return b
}

func (b *statementBuilder) addParamValue(node interface{}) *statementBuilder {
	if b.newparams == nil {
		b.newparams = make(map[string]interface{})
	}
	b.currentNode = node
	if b.concatenationBuilder != nil {
		b.concatenationBuilder.add(node)
		b.currentNode = nil
	} else if b.listBuilder != nil {
		b.listBuilder.add(node)
		b.currentNode = nil
	} else {
		if b.currentKey != "" {
			b.newparams[b.currentKey] = node
			b.currentKey = ""
			b.currentNode = nil
		}
	}

	return b
}

func (b *statementBuilder) newList() *statementBuilder {
	b.listBuilder = &listValueBuilder{}
	return b
}

func (b *statementBuilder) buildList() *statementBuilder {
	if b.listBuilder != nil {
		node := b.listBuilder.build()
		b.listBuilder = nil
		b.addParamValue(node)
	}
	return b
}

func (a *AST) addAction(text string) {
	if IsInvalidAction(text) {
		panic(fmt.Errorf("unknown action '%s'", text))
	}
	a.stmtBuilder.action = text
}

func (a *AST) addEntity(text string) {
	if IsInvalidEntity(text) {
		panic(fmt.Errorf("unknown entity '%s'", text))
	}
	a.stmtBuilder.entity = text
}

func (a *AST) addValue() {
	a.stmtBuilder.isValue = true
}

func (a *AST) addDeclarationIdentifier(text string) {
	a.stmtBuilder.declarationIdentifier = text
}

func (a *AST) NewStatement() {
	a.stmtBuilder = &statementBuilder{}
}

func (a *AST) StatementDone() {

	if stmt := a.stmtBuilder.build(); stmt != nil {
		a.Statements = append(a.Statements, stmt)
	}
	a.stmtBuilder = nil
}

func (a *AST) addParamKey(text string) {
	a.stmtBuilder.addParamKey(text)
}

func (a *AST) addParamValue(text string) {
	var val interface{}
	i, err := strconv.Atoi(text)
	if err == nil {
		if len(text) > 1 && text[0] == '0' {
			// We want an integer beginning with '0' to keep its initial '0' (so as string)
			val = text
		} else {
			val = i
		}
	} else {
		f, err := strconv.ParseFloat(text, 64)
		if err == nil {
			val = f
		} else {
			val = text
		}
	}
	a.stmtBuilder.addParamValue(InterfaceNode{i: val})
}

func (a *AST) addFirstValueInList() {
	a.stmtBuilder.newList()
}
func (a *AST) lastValueInList() {
	a.stmtBuilder.buildList()
}

func (a *AST) addFirstValueInConcatenation() {
	a.stmtBuilder.concatenationBuilder = &concatenationValueBuilder{}
}

func (a *AST) lastValueInConcatenation() {
	if a.stmtBuilder.concatenationBuilder != nil {
		node := a.stmtBuilder.concatenationBuilder.build()
		a.stmtBuilder.concatenationBuilder = nil
		a.stmtBuilder.addParamValue(node)
	}
}

func (a *AST) addStringValue(text string) {
	a.stmtBuilder.addParamValue(InterfaceNode{i: text})
}

func (a *AST) addParamRefValue(text string) {
	a.stmtBuilder.addParamValue(RefNode{key: text})
}

func (a *AST) addParamHoleValue(text string) {
	a.stmtBuilder.addParamValue(HoleNode{key: text})
}

func (a *AST) addAliasParam(text string) {
	a.stmtBuilder.addParamValue(AliasNode{key: text})
}

type listValueBuilder struct {
	elements []interface{}
}

func (c *listValueBuilder) add(node interface{}) *listValueBuilder {
	c.elements = append(c.elements, node)
	return c
}

func (c *listValueBuilder) build() ListNode {
	node := ListNode{arr: c.elements}
	return node
}

type concatenationValueBuilder struct {
	elements []interface{}
}

func (c *concatenationValueBuilder) add(node interface{}) *concatenationValueBuilder {
	c.elements = append(c.elements, node)
	return c
}

func (c *concatenationValueBuilder) build() ConcatenationNode {
	node := ConcatenationNode{arr: c.elements}
	return node
}

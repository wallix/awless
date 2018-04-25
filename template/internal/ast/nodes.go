package ast

import (
	"errors"
	"fmt"
	"strings"
)

var (
	_ ExpressionNode = (*RightExpressionNode)(nil)

	_ Node = (*HoleNode)(nil)
	_ Node = (*AliasNode)(nil)
	_ Node = (*RefNode)(nil)
	_ Node = (*ConcatenationNode)(nil)
	_ Node = (*ListNode)(nil)
	_ Node = (*InterfaceNode)(nil)
)

type RightExpressionNode struct {
	i interface{}
}

func (n *RightExpressionNode) Node() interface{} {
	return n.i
}

func (n *RightExpressionNode) Result() interface{} {
	switch v := n.i.(type) {
	case InterfaceNode:
		return v.i
	case RefNode, AliasNode, HoleNode:
		return nil
	case ListNode:
		var arr []interface{}
		for _, e := range v.arr {
			switch ev := e.(type) {
			case InterfaceNode:
				arr = append(arr, ev.i)
			case RefNode, AliasNode, HoleNode:
				return nil
			default:
				arr = append(arr, ev)
			}
		}
		return arr
	case ConcatenationNode:
		return v.Concat()
	default:
		return n.i
	}
}

func (n *RightExpressionNode) Err() error {
	switch n.i.(type) {
	case InterfaceNode:
		return nil
	default:
		return errors.New("right expr node is not an interface node")
	}
}

func (n *RightExpressionNode) String() string {
	return fmt.Sprint(n.i)
}

func (n *RightExpressionNode) clone() Node {
	return &RightExpressionNode{
		i: n.i,
	}
}

type CommandNode struct {
	Command
	CmdResult interface{}
	CmdErr    error

	Action, Entity string
	ParamNodes     map[string]interface{}
	Refs           map[string]interface{}
}

type RefNode struct {
	key string
}

func NewRefNode(s string) RefNode {
	return RefNode{key: s}
}

func (n RefNode) Ref() string {
	return n.key
}

func (n RefNode) clone() Node {
	return n
}

func (n RefNode) String() string {
	return "$" + n.key
}

type AliasNode struct {
	key string
}

func NewAliasNode(s string) AliasNode {
	return AliasNode{key: s}
}

func (n AliasNode) clone() Node {
	return n
}

func (n AliasNode) Alias() string {
	return n.key
}

func (n AliasNode) String() string {
	return "@" + n.key
}

type HoleNode struct {
	key      string
	optional bool
}

func NewHoleNode(s string) HoleNode {
	return HoleNode{key: s}
}

func NewOptionalHoleNode(s string) HoleNode {
	return HoleNode{key: s, optional: true}
}

func (n HoleNode) IsOptional() bool {
	return n.optional
}

func (n HoleNode) Hole() string {
	return n.key
}

func (n HoleNode) String() string {
	return "{" + n.key + "}"
}

func (n HoleNode) clone() Node {
	return n
}

type ListNode struct {
	arr []interface{}
}

func NewListNode(arr []interface{}) ListNode {
	return ListNode{arr: arr}
}

func (n ListNode) String() string {
	var a []string
	for _, e := range n.arr {
		a = append(a, fmt.Sprint(e))
	}
	return "[" + strings.Join(a, ",") + "]"
}

func (n ListNode) Elems() []interface{} {
	return n.arr
}

func (n ListNode) clone() Node {
	return n
}

type ConcatenationNode struct {
	arr []interface{}
}

func NewConcatenationNode(arr []interface{}) ConcatenationNode {
	return ConcatenationNode{arr: arr}
}

func (n ConcatenationNode) Concat() string {
	var arr []string
	for _, e := range n.arr {
		switch ee := e.(type) {
		case InterfaceNode:
			arr = append(arr, fmt.Sprint(ee.i))
		default:
			arr = append(arr, fmt.Sprint(ee))
		}
	}
	return strings.Join(arr, "")
}

func (n ConcatenationNode) String() string {
	var hasUnresolvedHole bool
	var elems []string
	for _, val := range n.arr {
		if _, has := val.(HoleNode); has {
			hasUnresolvedHole = true
			break
		}
	}
	for _, val := range n.arr {
		switch node := val.(type) {
		case InterfaceNode:
			if str, isStr := node.i.(string); isStr {
				if hasUnresolvedHole {
					elems = append(elems, Quote(str))
				} else {
					elems = append(elems, str)
				}
			} else {
				panic(fmt.Sprintf("concatenation node expects only strings and holes: got %T", node.i))
			}
		default:
			elems = append(elems, fmt.Sprint(val))
		}
	}
	if hasUnresolvedHole {
		return strings.Join(elems, "+")
	} else {
		return quoteStringIfNeeded(strings.Join(elems, ""))
	}
}

func (n ConcatenationNode) clone() Node {
	return n
}

type InterfaceNode struct {
	i interface{}
}

func (n InterfaceNode) Value() interface{} {
	return n.i
}

func (n InterfaceNode) String() string {
	switch v := n.i.(type) {
	case []string:
		return "[" + strings.Join(v, ",") + "]"
	case string:
		return quoteStringIfNeeded(v)
	default:
		return fmt.Sprint(v)
	}
}

func (n InterfaceNode) clone() Node {
	return n
}

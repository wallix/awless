package ast

import (
	"errors"
	"fmt"
	"strings"
)

func VerifyRefs(tree Node) error {
	var errs []string
	addErr := func(err string) {
		errs = append(errs, err)
	}

	v := newVisitor()
	v.onRefs = func(parent interface{}, node RefNode) {
		if !contains(v.declaredVariables, node.key) {
			addErr(fmt.Sprintf("using reference '$%s' but '%[1]s' is undefined in template", node.key))
		}
	}
	v.visit(tree)

	for i, declared := range v.declaredVariables {
		if contains(v.declaredVariables[:i], declared) {
			addErr(fmt.Sprintf("using reference '$%s' but '%[1]s' has already been assigned in template", declared))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func ProcessRefs(tree Node, fillers map[string]interface{}) {
	v := newVisitor()
	v.onRefs = func(parent interface{}, node RefNode) {
		var done bool
		var val interface{}
		for k, v := range fillers {
			if k == node.key {
				done = true
				val = v
			}
		}

		if done {
			switch p := parent.(type) {
			case ListNode:
				p.arr[v.listIndex] = val
			case *CommandNode:
				p.ParamNodes[v.key] = val
			case *RightExpressionNode:
				p.i = val
			}
		}
	}
	v.visit(tree)
}

func RemoveOptionalHoles(tree Node) {
	v := newVisitor()
	v.onHoles = func(parent interface{}, node HoleNode) {
		if node.IsOptional() {
			switch p := parent.(type) {
			case ListNode:
				p.arr = append(p.arr[:v.listIndex], p.arr[v.listIndex+1:]...)
			case *CommandNode:
				delete(p.ParamNodes, v.key)
			case *RightExpressionNode:
				p.i = nil
			}
		}
	}
	v.visit(tree)
}

func CollectUniqueHoles(tree Node) map[HoleNode][]string {
	uniqueHoles := make(map[HoleNode][]string)

	v := newVisitor()
	v.onHoles = func(parent interface{}, node HoleNode) {
		if _, ok := uniqueHoles[node]; !ok {
			uniqueHoles[node] = []string{}
		}

		if v.action != "" && v.entity != "" && v.key != "" {
			paramPath := fmt.Sprintf("%s.%s.%s", v.action, v.entity, v.key)
			if !contains(uniqueHoles[node], paramPath) {
				uniqueHoles[node] = append(uniqueHoles[node], paramPath)
			}
		}
	}
	v.visit(tree)
	return uniqueHoles
}

func CollectHoles(tree Node) (holes []HoleNode) {
	v := newVisitor()
	v.onHoles = func(parent interface{}, node HoleNode) {
		holes = append(holes, node)
	}
	v.visit(tree)
	return
}

func ProcessHoles(tree Node, fillers map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})

	v := newVisitor()
	v.onHoles = func(parent interface{}, node HoleNode) {
		var done bool
		var val interface{}
		for k, v := range fillers {
			if k == node.key {
				done = true
				val = v
				switch vv := v.(type) {
				case AliasNode, RefNode, HoleNode, ConcatenationNode:
					processed[k] = fmt.Sprint(v)
				case ListNode:
					var arr []interface{}
					for _, a := range vv.arr {
						switch e := a.(type) {
						case AliasNode, RefNode, HoleNode:
							arr = append(arr, fmt.Sprint(e))
						default:
							arr = append(arr, e)
						}
					}
					processed[k] = arr
				default:
					processed[k] = v
				}
			}
		}

		if done {
			switch p := parent.(type) {
			case ConcatenationNode:
				p.arr[v.concatItemIndex] = val
			case ListNode:
				p.arr[v.listIndex] = val
			case *CommandNode:
				p.ParamNodes[v.key] = val
			case *RightExpressionNode:
				p.i = val
			}
		}
	}

	v.visit(tree)

	return processed
}

func CollectAliases(tree Node) (aliases []AliasNode) {
	v := newVisitor()
	v.onAliases = func(parent interface{}, node AliasNode) {
		aliases = append(aliases, node)
	}
	v.visit(tree)
	return
}

func ProcessAliases(tree Node, aliasFunc func(action, entity string, key string) func(string) (string, bool)) {
	v := newVisitor()
	v.onAliases = func(parent interface{}, node AliasNode) {
		if resolv, hasResolv := aliasFunc(v.action, v.entity, v.key)(node.key); hasResolv {
			switch p := parent.(type) {
			case ListNode:
				p.arr[v.listIndex] = resolv
			case ConcatenationNode:
				p.arr[v.concatItemIndex] = resolv
			case *CommandNode:
				p.ParamNodes[v.key] = resolv
			case *RightExpressionNode:
				p.i = resolv
			}
		}
	}

	v.visit(tree)
	return
}

type visitor struct {
	onRefs    func(parent interface{}, n RefNode)
	onAliases func(parent interface{}, n AliasNode)
	onHoles   func(parent interface{}, n HoleNode)

	parent                     Node
	declaredVariables          []string
	action, entity, key        string
	listIndex, concatItemIndex int
}

func newVisitor() *visitor {
	return &visitor{
		onRefs:    func(interface{}, RefNode) {},
		onAliases: func(interface{}, AliasNode) {},
		onHoles:   func(interface{}, HoleNode) {},
	}
}

func (v *visitor) visit(tree Node) {
	switch t := tree.(type) {
	case InterfaceNode:
		return
	case HoleNode:
		v.onHoles(v.parent, t)
		return
	case AliasNode:
		v.onAliases(v.parent, t)
		return
	case RefNode:
		v.onRefs(v.parent, t)
		return
	}

	switch t := tree.(type) {
	case *AST:
		v.parent = tree
		for _, st := range t.Statements {
			v.visit(st)
		}
	case *Statement:
		v.parent = tree
		v.visit(t.Node)
	case *CommandNode:
		v.action, v.entity = t.Action, t.Entity
		for key, param := range t.ParamNodes {
			if n, ok := param.(Node); ok {
				v.parent = tree
				v.key = key
				v.visit(n)
			}
		}
	case *DeclarationNode:
		v.key = t.Ident
		v.parent = tree
		v.visit(t.Expr)
		v.declaredVariables = append(v.declaredVariables, t.Ident)
	case *RightExpressionNode:
		if n, ok := t.i.(Node); ok {
			v.parent = tree
			v.visit(n)
		}

	case ListNode:
		for i, el := range t.arr {
			if n, ok := el.(Node); ok {
				v.listIndex = i
				v.parent = tree
				v.visit(n)
			}
		}
	case ConcatenationNode:
		for i, el := range t.arr {
			if n, ok := el.(Node); ok {
				v.concatItemIndex = i
				v.parent = tree
				v.visit(n)
			}
		}
	default:
		panic(fmt.Sprintf("unsupported AST type %T", t))
	}
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

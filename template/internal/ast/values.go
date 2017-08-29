package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type CompositeValue interface {
	String() string
	Value() interface{}
	Clone() CompositeValue
}

type WithRefs interface {
	GetRefs() []string
	ProcessRefs(map[string]interface{})
	ReplaceRef(string, CompositeValue)
	IsRef(string) bool
}

type WithAlias interface {
	GetAliases() []string
	ResolveAlias(func(string) (string, bool))
}

type listValue struct {
	vals []CompositeValue
}

func (l *listValue) GetHoles() (res []string) {
	for _, val := range l.vals {
		if withHoles, ok := val.(WithHoles); ok {
			res = append(res, withHoles.GetHoles()...)
		}
	}
	return
}

func (l *listValue) GetRefs() (res []string) {
	for _, val := range l.vals {
		if withRefs, ok := val.(WithRefs); ok {
			res = append(res, withRefs.GetRefs()...)
		}
	}
	return
}

func (l *listValue) Value() interface{} {
	var res []interface{}
	for _, val := range l.vals {
		if v := val.Value(); v != nil {
			res = append(res, v)
		}
	}
	return res
}

func (l *listValue) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	for _, val := range l.vals {
		if withHoles, ok := val.(WithHoles); ok {
			valProc := withHoles.ProcessHoles(fills)
			for k, v := range valProc {
				processed[k] = v
			}
		}
	}
	return processed
}

func (l *listValue) ProcessRefs(fills map[string]interface{}) {
	for _, val := range l.vals {
		if withRefs, ok := val.(WithRefs); ok {
			withRefs.ProcessRefs(fills)
		}
	}
}

func (l *listValue) ReplaceRef(key string, value CompositeValue) {
	for k, param := range l.vals {
		if withRef, ok := param.(WithRefs); ok {
			if withRef.IsRef(key) {
				l.vals[k] = value
			} else {
				withRef.ReplaceRef(key, value)
			}
		}
	}
}

func (l *listValue) IsRef(key string) bool {
	return false
}

func (l *listValue) String() string {
	var buff bytes.Buffer
	buff.WriteRune('[')
	for i, val := range l.vals {
		buff.WriteString(val.String())
		if i < len(l.vals)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteRune(']')
	return buff.String()
}

func (l *listValue) GetAliases() (res []string) {
	for _, val := range l.vals {
		if alias, ok := val.(WithAlias); ok {
			res = append(res, alias.GetAliases()...)
		}
	}
	return
}

func (l *listValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	for _, val := range l.vals {
		if alias, ok := val.(WithAlias); ok {
			alias.ResolveAlias(resolvFunc)
		}
	}
}

func (l *listValue) Clone() CompositeValue {
	clone := &listValue{}
	for _, val := range l.vals {
		clone.vals = append(clone.vals, val.Clone())
	}
	return clone
}

type interfaceValue struct {
	val interface{}
}

func NewInterfaceValue(i interface{}) CompositeValue {
	return &interfaceValue{val: i}
}

func (i *interfaceValue) Value() interface{} {
	return i.val
}

func (i *interfaceValue) String() string {
	return printParamValue(i.val)
}

func (i *interfaceValue) Clone() CompositeValue {
	return &interfaceValue{val: i.val}
}

type holeValue struct {
	hole  string
	val   interface{}
	alias string
}

func NewHoleValue(hole string) CompositeValue {
	return &holeValue{hole: hole}
}

func (h *holeValue) GetHoles() (res []string) {
	if h.val == nil && h.alias == "" {
		res = append(res, h.hole)
	}
	return
}

func (h *holeValue) Value() interface{} {
	return h.val
}

func (h *holeValue) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	if fill, ok := fills[h.hole]; ok {
		if strFill, ok := fill.(string); ok && strings.HasPrefix(strFill, "@") {
			h.alias = strFill[1:]
			processed[h.hole] = strFill
		} else if aliasValue, ok := fill.(*aliasValue); ok {
			h.alias = aliasValue.alias
			h.val = aliasValue.val
			processed[h.hole] = aliasValue.String()
		} else {
			h.val = fill
			processed[h.hole] = fill
		}
	}
	return processed
}

func (h *holeValue) String() string {
	if h.val != nil {
		return printParamValue(h.val)
	} else {
		return fmt.Sprintf("{%s}", h.hole)
	}
}

func (h *holeValue) GetAliases() (aliases []string) {
	if h.val == nil && h.alias != "" {
		aliases = append(aliases, h.alias)
	}
	return
}

func (h *holeValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	if h.val != nil || h.alias == "" {
		return
	}
	if val, ok := resolvFunc(h.alias); ok {
		h.val = val
	}
}

func (h *holeValue) Clone() CompositeValue {
	return &holeValue{val: h.val, hole: h.hole, alias: h.alias}
}

type holesStringValue struct {
	holes []*holeValue
	input string
}

func (h *holesStringValue) GetHoles() (res []string) {
	for _, hole := range h.holes {
		res = append(res, hole.GetHoles()...)
	}
	return
}

func (h *holesStringValue) Value() interface{} {
	out := h.input
	for _, hole := range h.holes {
		if hole.Value() != nil {
			out = strings.Replace(out, "{"+hole.hole+"}", fmt.Sprint(hole.Value()), 1)
		}
	}
	return out
}

func (h *holesStringValue) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	for _, hole := range h.holes {
		valProc := hole.ProcessHoles(fills)
		for k, v := range valProc {
			processed[k] = v
		}
	}
	return processed
}

func (h *holesStringValue) String() string {
	out := h.input
	for _, hole := range h.holes {
		if hole.Value() != nil {
			out = strings.Replace(out, "{"+hole.hole+"}", hole.String(), 1)
		}
	}
	return out
}

func (h *holesStringValue) GetAliases() (aliases []string) {
	for _, hole := range h.holes {
		aliases = append(aliases, hole.GetAliases()...)
	}
	return
}

func (h *holesStringValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	for _, hole := range h.holes {
		hole.ResolveAlias(resolvFunc)
	}
}

func (h *holesStringValue) Clone() CompositeValue {
	clone := &holesStringValue{input: h.input}
	for _, hole := range h.holes {
		clone.holes = append(clone.holes, hole.Clone().(*holeValue))
	}
	return clone
}

type aliasValue struct {
	alias string
	val   interface{}
}

func NewAliasValue(alias string) CompositeValue {
	return &aliasValue{alias: alias}
}

func (a *aliasValue) Value() interface{} {
	return a.val
}

func (a *aliasValue) String() string {
	if a.val != nil {
		return printParamValue(a.val)
	} else {
		return fmt.Sprintf("@%s", a.alias)
	}
}

func (a *aliasValue) GetAliases() (aliases []string) {
	if a.val == nil {
		aliases = append(aliases, a.alias)
	}
	return
}

func (a *aliasValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	if val, ok := resolvFunc(a.alias); ok {
		a.val = val
	}
}

func (a *aliasValue) Clone() CompositeValue {
	return &aliasValue{val: a.val, alias: a.alias}
}

type referenceValue struct {
	ref   string
	val   interface{}
	alias string
}

func (r *referenceValue) GetRefs() (refs []string) {
	if r.val == nil {
		refs = append(refs, r.ref)
	}
	return refs
}

func (r *referenceValue) ProcessRefs(fills map[string]interface{}) {
	if fill, ok := fills[r.ref]; ok {
		if strFill, ok := fill.(string); ok && strings.HasPrefix(strFill, "@") {
			r.alias = strFill[1:]
		} else {
			r.val = fill
		}
	}
}

func (r *referenceValue) ReplaceRef(key string, value CompositeValue) {

}

func (r *referenceValue) IsRef(key string) bool {
	return key == r.ref
}

func (r *referenceValue) Value() interface{} {
	return r.val
}

func (r *referenceValue) String() string {
	if r.val != nil {
		return printParamValue(r.val)
	} else {
		return fmt.Sprintf("$%s", r.ref)
	}
}

func (r *referenceValue) Clone() CompositeValue {
	return &referenceValue{val: r.val, ref: r.ref, alias: r.alias}
}

func (r *referenceValue) GetAliases() (aliases []string) {
	if r.val == nil && r.alias != "" {
		aliases = append(aliases, r.alias)
	}
	return
}

func (r *referenceValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	if r.val != nil || r.alias == "" {
		return
	}
	if val, ok := resolvFunc(r.alias); ok {
		r.val = val
	}
}

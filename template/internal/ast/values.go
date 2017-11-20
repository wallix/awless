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

func (l *listValue) GetHoles() map[string][]string {
	res := make(map[string][]string)
	for _, val := range l.vals {
		if withHoles, ok := val.(WithHoles); ok {
			for k, v := range withHoles.GetHoles() {
				res[k] = v
			}
		}
	}
	return res
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
	if len(res) == 0 {
		return nil
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
	alias WithAlias
}

func NewHoleValue(hole string) CompositeValue {
	return &holeValue{hole: hole}
}

func (h *holeValue) GetHoles() map[string][]string {
	if h.val == nil && h.alias == nil {
		return map[string][]string{h.hole: nil}
	}
	return make(map[string][]string)
}

func (h *holeValue) Value() interface{} {
	if h.alias != nil {
		return h.alias.(CompositeValue).Value()
	}
	return h.val
}

func (h *holeValue) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	if fill, ok := fills[h.hole]; ok {
		if withAlias, ok := fill.(WithAlias); ok {
			h.alias = withAlias
			switch vv := withAlias.(type) {
			case *listValue:
				var processedAliases []string
				for _, alias := range vv.vals {
					processedAliases = append(processedAliases, alias.String())
				}
				processed[h.hole] = processedAliases
			case *aliasValue:
				processed[h.hole] = vv.String()
			}

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
	} else if h.alias != nil {
		return fmt.Sprint(h.alias)
	} else {
		return fmt.Sprintf("{%s}", h.hole)
	}
}

func (h *holeValue) GetAliases() (aliases []string) {
	if h.val == nil && h.alias != nil {
		return h.alias.GetAliases()
	}
	return
}

func (h *holeValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	if h.val == nil && h.alias != nil {
		h.alias.ResolveAlias(resolvFunc)
	}
	return
}

func (h *holeValue) Clone() CompositeValue {
	return &holeValue{val: h.val, hole: h.hole, alias: h.alias}
}

type concatenationValue struct {
	vals []CompositeValue
}

func (c *concatenationValue) GetHoles() map[string][]string {
	res := make(map[string][]string)
	for _, val := range c.vals {
		if withHoles, ok := val.(WithHoles); ok {
			for k, v := range withHoles.GetHoles() {
				res[k] = v
			}
		}
	}
	return res
}

func (c *concatenationValue) GetRefs() (res []string) {
	for _, val := range c.vals {
		if withRefs, ok := val.(WithRefs); ok {
			res = append(res, withRefs.GetRefs()...)
		}
	}
	return
}

func (c *concatenationValue) Value() interface{} {
	var buff bytes.Buffer
	for _, val := range c.vals {
		if val.Value() != nil {
			buff.WriteString(fmt.Sprint(val.Value()))
		}
	}
	return buff.String()
}

func (c *concatenationValue) ProcessHoles(fills map[string]interface{}) map[string]interface{} {
	processed := make(map[string]interface{})
	for _, val := range c.vals {
		if withHoles, ok := val.(WithHoles); ok {
			valProc := withHoles.ProcessHoles(fills)
			for k, v := range valProc {
				processed[k] = v
			}
		}
	}
	return processed
}

func (c *concatenationValue) ProcessRefs(fills map[string]interface{}) {
	for _, val := range c.vals {
		if withRefs, ok := val.(WithRefs); ok {
			withRefs.ProcessRefs(fills)
		}
	}
}

func (c *concatenationValue) ReplaceRef(key string, value CompositeValue) {
	for k, param := range c.vals {
		if withRef, ok := param.(WithRefs); ok {
			if withRef.IsRef(key) {
				c.vals[k] = value
			} else {
				withRef.ReplaceRef(key, value)
			}
		}
	}
}

func (c *concatenationValue) IsRef(key string) bool {
	return false
}

func (c *concatenationValue) String() string {
	if len(c.GetHoles())+len(c.GetRefs())+len(c.GetAliases()) == 0 {
		return quoteStringIfNeeded(c.Value().(string))
	}
	var elems []string
	for _, val := range c.vals {
		str := val.String()
		if val.Value() != nil && !isQuoted(str) {
			elems = append(elems, quoteString(str))
			continue
		}
		elems = append(elems, str)
	}
	return strings.Join(elems, "+")
}

func (c *concatenationValue) GetAliases() (res []string) {
	for _, val := range c.vals {
		if alias, ok := val.(WithAlias); ok {
			res = append(res, alias.GetAliases()...)
		}
	}
	return
}

func (c *concatenationValue) ResolveAlias(resolvFunc func(string) (string, bool)) {
	for _, val := range c.vals {
		if alias, ok := val.(WithAlias); ok {
			alias.ResolveAlias(resolvFunc)
		}
	}
}

func (c *concatenationValue) Clone() CompositeValue {
	clone := &concatenationValue{}
	for _, val := range c.vals {
		clone.vals = append(clone.vals, val.Clone())
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

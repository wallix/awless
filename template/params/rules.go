package params

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	DEFAULTS_SUGGESTED = iota
	NO_SUGGESTED
	ALL_SUGGESTED
)

type Rule interface {
	Visit(func(Rule))
	Run(input []string) error
	Required() []string
	Missing(input []string) []string
	Suggested(input []string, paramsSuggested int) []string
	String() string
}

func List(r Rule) ([]string, []string) {
	return collect(r)
}

func Run(r Rule, input []string) error {
	if err := unexpectedParam(r, input); err != nil {
		return err
	}
	return r.Run(input)
}

func unexpectedParam(r Rule, input []string) (err error) {
	var unex []string
	params, opts := collect(r)
	all := append(params, opts...)
	for _, s := range input {
		if !contains(all, s) {
			unex = append(unex, s)
		}
	}
	if len(unex) > 0 {
		err = fmt.Errorf("unexpected param(s): %s", strings.Join(unex, ", "))
	}
	return
}

func collect(r Rule) (out []string, opts []string) {
	r.Visit(func(r Rule) {
		switch v := r.(type) {
		case key:
			out = append(out, v.String())
		case opt:
			opts = append(opts, v.optionals...)
		}
	})
	sort.Strings(out)
	sort.Strings(opts)
	return
}

type allOf struct {
	defaultRule
}

func AllOf(rules ...Rule) Rule {
	return allOf{build(rules)}
}

func (n allOf) Run(input []string) (err error) {
	for _, r := range n.rules {
		err = r.Run(input)
		if err != optErr && err != nil {
			return errors.New(n.String())
		}
	}
	return nil
}

func (n allOf) Missing(input []string) (miss []string) {
	for _, r := range n.rules {
		miss = append(miss, r.Missing(input)...)
	}
	return
}

func (n allOf) Required() (all []string) {
	for _, r := range n.rules {
		all = append(all, r.Required()...)
	}
	return
}

func (n allOf) String() string {
	return n.rules.join(" + ")
}

type onlyOneOf struct {
	defaultRule
}

func OnlyOneOf(rules ...Rule) Rule {
	return onlyOneOf{build(rules)}
}

func (n onlyOneOf) Run(input []string) error {
	if len(n.rules) == 0 {
		return nil
	}
	var pass int
	for _, r := range n.rules {
		if err := r.Run(input); err == nil {
			pass++
		}
	}
	if pass != 1 {
		return fmt.Errorf("only %s", n.rules)
	}
	return nil
}

func (n onlyOneOf) Missing(input []string) (miss []string) {
	if err := n.Run(input); err != nil && len(n.rules) > 0 {
		miss = append(miss, n.rules[0].Missing(input)...)
	}
	return
}

func (n onlyOneOf) Required() (all []string) {
	if len(n.rules) > 0 {
		all = append(all, n.rules[0].Required()...)
	}
	return
}

func (n onlyOneOf) String() string {
	return "(" + n.rules.join(" | ") + ")"
}

type atLeastOneOf struct {
	defaultRule
}

func AtLeastOneOf(rules ...Rule) Rule {
	return atLeastOneOf{build(rules)}
}

func (n atLeastOneOf) Run(input []string) error {
	if len(n.rules) == 0 {
		return nil
	}
	var pass int
	for _, r := range n.rules {
		if err := r.Run(input); err == nil || err == optErr {
			pass++
		}
	}
	if pass < 1 {
		return fmt.Errorf("at least one of %v", n.rules)
	}
	return nil
}

func (n atLeastOneOf) Missing(input []string) (miss []string) {
	if err := n.Run(input); err != nil && len(n.rules) > 0 {
		miss = append(miss, n.rules[0].Missing(input)...)
	}
	return
}

func (n atLeastOneOf) Required() (all []string) {
	if len(n.rules) > 0 {
		all = append(all, n.rules[0].Required()...)
	}
	return
}

func (n atLeastOneOf) String() string {
	return "(" + n.rules.join(" / ") + ")"
}

type opt struct {
	optionals   []string
	isSuggested bool
}

func Opt(s ...string) Rule {
	o := opt{}
	o.optionals = append(o.optionals, s...)
	return o
}

func Suggested(s ...string) Rule {
	o := opt{isSuggested: true}
	o.optionals = append(o.optionals, s...)
	return o
}

var optErr = errors.New("opt err")

func (n opt) Visit(fn func(r Rule)) {
	fn(n)
}

func (n opt) Run(input []string) error {
	return optErr
}

func (n opt) Missing(input []string) (miss []string) {
	return
}

func (n opt) Required() []string {
	return []string{}
}

func (n opt) String() string {
	return "[" + strings.Join(n.optionals, " ") + "]"
}

func (n opt) Suggested(input []string, paramsSuggested int) (miss []string) {
	for _, p := range n.optionals {
		if !contains(input, p) {
			if paramsSuggested == ALL_SUGGESTED || (paramsSuggested == DEFAULTS_SUGGESTED && n.isSuggested) {
				miss = append(miss, p)
			}
		}
	}
	return
}

func Key(k string, isSuggested ...bool) Rule {
	var suggested bool
	if len(isSuggested) > 0 {
		suggested = isSuggested[0]
	}
	return key{key: k, isSuggested: suggested}
}

type key struct {
	key         string
	isSuggested bool
}

func (n key) Visit(fn func(r Rule)) {
	fn(n)
}

func (n key) Run(input []string) error {
	if !contains(input, n.key) {
		return errors.New(n.key)
	}
	return nil
}

func (n key) Missing(input []string) (miss []string) {
	if !contains(input, n.key) {
		miss = append(miss, n.key)
	}
	return
}

func (n key) Suggested(input []string, paramsSuggested int) (miss []string) {
	if n.isSuggested {
		if !contains(input, n.key) {
			miss = append(miss, n.key)
		}
	}
	return
}

func (n key) Required() []string {
	return []string{n.key}
}

func (n key) String() string {
	return n.key
}

type none struct{}

func None() Rule {
	return none{}
}

func (n none) Visit(func(Rule))                 {}
func (n none) Run(input []string) error         { return nil }
func (n none) Required() []string               { return []string{} }
func (n none) Missing([]string) []string        { return []string{} }
func (n none) String() string                   { return "none" }
func (n none) Suggested([]string, int) []string { return nil }

func build(rules []Rule) (d defaultRule) {
	for _, n := range rules {
		d.rules = append(d.rules, n)
	}
	return
}

type defaultRule struct {
	rules rules
}

func (d defaultRule) Suggested(input []string, paramsSuggested int) (miss []string) {
	for _, r := range d.rules {
		miss = append(miss, r.Suggested(input, paramsSuggested)...)
	}
	return
}

type rules []Rule

func (rs rules) join(sep string) string {
	var arr []string
	for _, r := range rs {
		arr = append(arr, fmt.Sprint(r))
	}
	return strings.Join(arr, sep)
}

func (r defaultRule) Visit(fn func(r Rule)) {
	for _, n := range r.rules {
		n.Visit(fn)
	}
}

func contains(arr []string, s string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}
	return false
}

package params

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Rule interface {
	Visit(func(Rule))
	Run(input []string) error
	Required() []string
	Missing(input []string) []string
	String() string
}

func List(r Rule) (required []string, optionals []string, suggested []string) {
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
	params, opts, _ := collect(r)
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

func collect(r Rule) (out []string, opts []string, suggs []string) {
	r.Visit(func(r Rule) {
		switch v := r.(type) {
		case key:
			out = append(out, v.String())
		case opt:
			opts = append(opts, v.keys()...)
			for _, p := range v.suggested {
				suggs = append(suggs, string(p))
			}
		}
	})
	sort.Strings(out)
	sort.Strings(opts)
	sort.Strings(suggs)
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
			return fmt.Errorf("%s: expecting %s", err, n.String())
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
	optionals []string
	suggested []suggested
}

func Opt(i ...interface{}) Rule {
	o := opt{}
	for _, v := range i {
		switch vv := v.(type) {
		case string:
			o.optionals = append(o.optionals, vv)
		case []suggested:
			o.suggested = append(o.suggested, vv...)
		default:
			panic("invalid type for optional param")
		}
	}
	return o
}

func Suggested(s ...string) (sugs []suggested) {
	for _, v := range s {
		sugs = append(sugs, suggested(v))
	}
	return
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
	return "[" + strings.Join(n.keys(), " ") + "]"
}

func (n opt) keys() (keys []string) {
	for _, k := range n.optionals {
		keys = append(keys, string(k))
	}
	for _, k := range n.suggested {
		keys = append(keys, string(k))
	}
	return
}

func Key(k string) Rule {
	return key(k)
}

type key string
type suggested key

func (n key) Visit(fn func(r Rule)) {
	fn(n)
}

func (n key) Run(input []string) error {
	if !contains(input, string(n)) {
		return errors.New(string(n))
	}
	return nil
}

func (n key) Missing(input []string) (miss []string) {
	if !contains(input, string(n)) {
		miss = append(miss, string(n))
	}
	return
}

func (n key) Required() []string {
	return []string{string(n)}
}

func (n key) String() string {
	return string(n)
}

type none struct{}

func None() Rule {
	return none{}
}

func (n none) Visit(func(Rule))          {}
func (n none) Run(input []string) error  { return nil }
func (n none) Required() []string        { return []string{} }
func (n none) Missing([]string) []string { return []string{} }
func (n none) String() string            { return "none" }

func build(rules []Rule) (d defaultRule) {
	for _, n := range rules {
		d.rules = append(d.rules, n)
	}
	return
}

type defaultRule struct {
	rules rules
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

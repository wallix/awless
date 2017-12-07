package graph

import (
	"fmt"
	"strings"
)

type FilterFn func(*Resource) bool

func (g *Graph) Filter(entity string, filters ...FilterFn) (*Graph, error) {
	return g.filter(applyAnd, entity, filters...)
}

func (g *Graph) OrFilter(entity string, filters ...FilterFn) (*Graph, error) {
	return g.filter(applyOr, entity, filters...)
}

func (g *Graph) filter(apply func(filters ...FilterFn) FilterFn, entity string, filters ...FilterFn) (*Graph, error) {
	filtered := NewGraph()

	all, err := g.GetAllResources(entity)
	if err != nil {
		return filtered, err
	}

	for _, r := range all {
		if apply(filters...)(r) {
			filtered.AddResource(r)
		}
	}

	return filtered, nil
}

func BuildPropertyFilterFunc(key, val string) FilterFn {
	return func(r *Resource) bool {
		return strings.Contains(strings.ToLower(fmt.Sprint(r.properties[key])), strings.ToLower(val))
	}
}

func BuildTagFilterFunc(key, val string) FilterFn {
	return func(r *Resource) bool {
		tags, ok := r.properties["Tags"].([]string)
		if !ok {
			return false
		}
		for _, t := range tags {
			if fmt.Sprintf("%s=%s", key, val) == t {
				return true
			}
		}
		return false
	}
}

func BuildTagKeyFilterFunc(key string) FilterFn {
	return func(r *Resource) bool {
		tags, ok := r.properties["Tags"].([]string)
		if !ok {
			return false
		}
		for _, t := range tags {
			splits := strings.Split(t, "=")
			if len(splits) > 0 {
				if splits[0] == key {
					return true
				}
			}
		}
		return false
	}
}

func BuildTagValueFilterFunc(value string) FilterFn {
	return func(r *Resource) bool {
		tags, ok := r.properties["Tags"].([]string)
		if !ok {
			return false
		}
		for _, t := range tags {
			splits := strings.Split(t, "=")
			if len(splits) > 1 {
				if splits[1] == value {
					return true
				}
			}
		}
		return false
	}
}

func applyAnd(filters ...FilterFn) FilterFn {
	return func(r *Resource) bool {
		include := true
		for _, f := range filters {
			include = include && f(r)
		}
		return include
	}
}

func applyOr(filters ...FilterFn) FilterFn {
	return func(r *Resource) bool {
		if len(filters) == 0 {
			return true
		}
		for _, f := range filters {
			if f(r) {
				return true
			}
		}
		return false
	}
}

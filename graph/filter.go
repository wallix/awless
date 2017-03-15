package graph

import (
	"fmt"
	"strings"
)

type FilterFn func(*Resource) bool

func (g *Graph) Filter(entity string, filters ...FilterFn) (*Graph, error) {
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
		return strings.Contains(strings.ToLower(fmt.Sprint(r.Properties[key])), strings.ToLower(val))
	}
}

func apply(filters ...FilterFn) FilterFn {
	return func(r *Resource) bool {
		include := true
		for _, f := range filters {
			include = include && f(r)
		}
		return include
	}
}

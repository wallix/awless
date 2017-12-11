/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Visitor interface {
	Visit(*Graph) error
}

type visitEachFunc func(res *Resource, depth int) error

func VisitorCollectFunc(collect *[]*Resource) visitEachFunc {
	return func(res *Resource, depth int) error {
		*collect = append(*collect, res)
		return nil
	}
}

type ParentsVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
	Relation    string
}

func (v *ParentsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	if v.Relation == "" {
		v.Relation = rdf.ParentOf
	}
	return tstore.NewTree(g.store.Snapshot(), v.Relation).TraverseAncestors(startNode, foreach)
}

type ChildrenVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
	Relation    string
}

func (v *ChildrenVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	if v.Relation == "" {
		v.Relation = rdf.ParentOf
	}
	return tstore.NewTree(g.store.Snapshot(), v.Relation).TraverseDFS(startNode, foreach)
}

type SiblingsVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *SiblingsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}

	return tstore.NewTree(g.store.Snapshot(), rdf.ParentOf).TraverseSiblings(startNode, resolveResourceType, foreach)
}

func prepareRDFVisit(g *Graph, root *Resource, each visitEachFunc, includeRoot bool) (string, func(g tstore.RDFGraph, n string, i int) error, error) {
	rootNode := root.Id()

	foreach := func(rdfG tstore.RDFGraph, n string, i int) error {
		rT, err := resolveResourceType(rdfG, n)
		if err != nil {
			return err
		}
		res, err := g.GetResource(rT, n)
		if err != nil {
			return err
		}
		if includeRoot || !root.Same(res) {
			if err := each(res, i); err != nil {
				return err
			}
		}
		return nil
	}
	return rootNode, foreach, nil
}

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
	"fmt"
	"sort"

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
}

func (v *ParentsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}

	return visitBottomUp(g.store.Snapshot(), startNode, foreach)
}

type ChildrenVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *ChildrenVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	return visitTopDown(g.store.Snapshot(), startNode, foreach)
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

	return visitSiblings(g.store.Snapshot(), startNode, foreach)
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

func visitTopDown(snap tstore.RDFGraph, root string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(snap, root, dist); err != nil {
		return err
	}

	triples := snap.WithSubjPred(root, rdf.ParentOf)

	var childs []string
	for _, tri := range triples {
		n, ok := tri.Object().Resource()
		if !ok {
			return fmt.Errorf("object is not a resource identifier")
		}
		childs = append(childs, n)
	}

	sort.Strings(childs)

	for _, child := range childs {
		visitTopDown(snap, child, each, dist+1)
	}

	return nil
}

func visitBottomUp(snap tstore.RDFGraph, startNode string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(snap, startNode, dist); err != nil {
		return err
	}
	triples := snap.WithPredObj(rdf.ParentOf, tstore.Resource(startNode))
	var parents []string
	for _, tri := range triples {
		parents = append(parents, tri.Subject())
	}

	sort.Strings(parents)

	for _, child := range parents {
		visitBottomUp(snap, child, each, dist+1)
	}

	return nil
}

func visitSiblings(snap tstore.RDFGraph, start string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	triples := snap.WithPredObj(rdf.ParentOf, tstore.Resource(start))

	var parents []string
	for _, tri := range triples {
		parents = append(parents, tri.Subject())
	}

	if len(parents) == 0 {
		return each(snap, start, 0)
	}

	sort.Strings(parents)

	for _, parent := range parents {
		parentTs := snap.WithSubjPred(parent, rdf.ParentOf)

		var childs []string
		for _, parentT := range parentTs {
			child, ok := parentT.Object().Resource()
			if !ok {
				return fmt.Errorf("object is not a resource identifier")
			}
			childs = append(childs, child)
		}

		sort.Strings(childs)

		startType, err := resolveResourceType(snap, start)
		if err != nil {
			return err
		}

		for _, child := range childs {
			rt, err := resolveResourceType(snap, child)
			if err != nil {
				return err
			}
			sameType := rt == startType
			if sameType {
				if err := each(snap, child, 0); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

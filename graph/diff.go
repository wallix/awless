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

import "github.com/wallix/awless/graph/internal/rdf"

type Diff struct {
	*rdf.Diff
}

func NewDiff(fromG, toG *Graph) *Diff {
	return &Diff{rdf.NewDiff(fromG.rdfG, toG.rdfG)}
}

func (d *Diff) FromGraph() *Graph {
	return &Graph{d.Diff.FromGraph()}
}

func (d *Diff) ToGraph() *Graph {
	return &Graph{d.Diff.ToGraph()}
}

func (d *Diff) MergedGraph() *Graph {
	return &Graph{d.Diff.MergedGraph()}
}

var Differ = &differ{rdf.DefaultDiffer}

type differ struct {
	rdf.Differ
}

func (d *differ) Run(root *Resource, from *Graph, to *Graph) (*Diff, error) {
	rootNode, err := root.toRDFNode()
	if err != nil {
		return nil, err
	}
	diff, err := d.Differ.Run(rootNode, from.rdfG, to.rdfG)
	return &Diff{diff}, err
}

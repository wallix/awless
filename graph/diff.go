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

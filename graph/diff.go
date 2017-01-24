package graph

import "github.com/wallix/awless/graph/internal/rdf"

type Diff struct {
	*rdf.Diff
}

func NewDiff(g *Graph) *Diff {
	return &Diff{rdf.NewDiff(g.Graph)}
}

func (d *Diff) FullGraph() *Graph {
	return &Graph{d.Diff.FullGraph()}
}

var HierarchicalDiffer = rdf.NewHierarchicalDiffer()
var DefaultDiffer = rdf.DefaultDiffer

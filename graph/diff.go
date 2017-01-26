package graph

import (
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Diff struct {
	*rdf.Diff
}

func NewDiff(g *Graph) *Diff {
	return &Diff{rdf.NewDiff(g.rdfG)}
}

func (d *Diff) FullGraph() *Graph {
	return &Graph{d.Diff.FullGraph()}
}

var HierarchicalDiffer = rdf.NewHierarchicalDiffer()
var Differ = &differ{rdf.DefaultDiffer}

type differ struct {
	rdf.Differ
}

func (d *differ) Run(root *node.Node, local *Graph, remote *Graph) (*Diff, error) {
	diff, err :=  d.Differ.Run(root, local.rdfG, remote.rdfG)
  return &Diff{diff}, err
}

package graph

import (
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Diff struct {
	*rdf.Diff
}

func NewDiff(localG, remoteG *Graph) *Diff {
	return &Diff{rdf.NewDiff(localG.rdfG, remoteG.rdfG)}
}

func (d *Diff) LocalGraph() *Graph {
	return &Graph{d.Diff.LocalGraph()}
}

func (d *Diff) RemoteGraph() *Graph {
	return &Graph{d.Diff.RemoteGraph()}
}

func (d *Diff) MergedGraph() *Graph {
	return &Graph{d.Diff.MergedGraph()}
}

var Differ = &differ{rdf.DefaultDiffer}

type differ struct {
	rdf.Differ
}

func (d *differ) Run(root *node.Node, local *Graph, remote *Graph) (*Diff, error) {
	diff, err := d.Differ.Run(root, local.rdfG, remote.rdfG)
	return &Diff{diff}, err
}

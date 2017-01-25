package sync

import (
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph"
)

// Diff represents the deleted/inserted RDF triples of a revision
type Diff struct {
	From       *Rev
	To         *Rev
	InfraDiff  *graph.Diff
	AccessDiff *graph.Diff
}

func BuildDiff(from, to *Rev, root *node.Node) (*Diff, error) {
	infraDiff, err := graph.DefaultDiffer.Run(root, to.Infra.Graph, from.Infra.Graph)
	if err != nil {
		return nil, err
	}

	accessDiff, err := graph.DefaultDiffer.Run(root, to.Access.Graph, from.Access.Graph)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		From:       from,
		To:         to,
		InfraDiff:  &graph.Diff{infraDiff},
		AccessDiff: &graph.Diff{accessDiff},
	}

	return res, nil
}

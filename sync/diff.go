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
	infraDiff, err := graph.Differ.Run(root, to.Infra, from.Infra)
	if err != nil {
		return nil, err
	}

	accessDiff, err := graph.Differ.Run(root, to.Access, from.Access)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		From:       from,
		To:         to,
		InfraDiff:  infraDiff,
		AccessDiff: accessDiff,
	}

	return res, nil
}

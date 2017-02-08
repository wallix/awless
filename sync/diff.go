package sync

import (
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/sync/repo"
)

// Diff represents the deleted/inserted RDF triples of a revision
type Diff struct {
	From       *repo.Rev
	To         *repo.Rev
	InfraDiff  *graph.Diff
	AccessDiff *graph.Diff
}

func BuildDiff(from, to *repo.Rev, root *graph.Resource) (*Diff, error) {
	infraDiff, err := graph.Differ.Run(root, from.Infra, to.Infra)
	if err != nil {
		return nil, err
	}

	accessDiff, err := graph.Differ.Run(root, from.Access, to.Access)
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

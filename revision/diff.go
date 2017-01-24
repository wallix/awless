package revision

import (
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph"
)

// Diff represents the deleted/inserted RDF triples of a revision
type Diff struct {
	From      *Revision
	To        *Revision
	GraphDiff *graph.Diff
}

func (rr *Repository) newDiff(from, to *Revision, root *node.Node, forFiles []string) (*Diff, error) {
	fromG, err := rr.revisionToRDFGraph(from, forFiles...)
	if err != nil {
		return nil, err
	}
	toG, err := rr.revisionToRDFGraph(to, forFiles...)
	if err != nil {
		return nil, err
	}
	diff, err := graph.DefaultDiffer.Run(root, toG.Graph, fromG.Graph)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		From:      from,
		To:        to,
		GraphDiff: &graph.Diff{diff},
	}

	return res, nil
}

func (rr *Repository) generateDiffs(pairs []*revisionPair, root *node.Node, forFiles []string) ([]*Diff, error) {
	var res []*Diff
	for _, revisionPair := range pairs {
		diff, err := rr.newDiff(revisionPair.from, revisionPair.to, root, forFiles)
		if err != nil {
			return res, err
		}
		res = append(res, diff)
	}
	return res, nil
}

package revision

import (
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

// Diff represents the deleted/inserted RDF triples of a revision
type Diff struct {
	From      *Revision
	To        *Revision
	GraphDiff *rdf.Diff
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
	diff, err := rdf.DefaultDiffer.Run(root, toG, fromG)
	if err != nil {
		return nil, err
	}

	res := &Diff{
		From:      from,
		To:        to,
		GraphDiff: diff,
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

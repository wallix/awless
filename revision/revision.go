package revision

import (
	"time"

	"github.com/wallix/awless/graph"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Revision references a commit of the local RDF resources
type Revision struct {
	Time time.Time
	ID   string
}

type revisionPair struct {
	from *Revision
	to   *Revision
}

func NewRevision(c *object.Commit) *Revision {
	return &Revision{Time: c.Committer.When, ID: c.ID().String()}
}

var initRevision = &Revision{}

func generateRevisionPairs(revisions []*Revision, param fetchParameter) []*revisionPair {
	var res []*revisionPair
	var groupF func(t1, t2 time.Time) bool

	switch param {
	case GroupAll:
		return []*revisionPair{{from: revisions[len(revisions)-1], to: revisions[0]}}
	case GroupByDay:
		groupF = func(t1, t2 time.Time) bool {
			return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
		}
	case GroupByWeek:
		groupF = func(t1, t2 time.Time) bool {
			y1, w1 := t1.ISOWeek()
			y2, w2 := t2.ISOWeek()
			return y1 == y2 && w1 == w2
		}
	default:
		for i := 0; i < len(revisions)-1; i++ {
			res = append(res, &revisionPair{from: revisions[i+1], to: revisions[i]})
		}
	}
	if groupF != nil && len(revisions) > 0 {
		r := revisions[0]
		time := r.Time
		lastAddedR := r
		for i := 1; i < len(revisions); i++ {
			r = revisions[i]
			if r != initRevision && !groupF(time, r.Time) {
				res = append(res, &revisionPair{from: revisions[i-1], to: lastAddedR})
				time = r.Time
				lastAddedR = revisions[i-1]
			}
		}
		res = append(res, &revisionPair{from: r, to: lastAddedR})
	}
	return res
}

func (rr *Repository) revisionToRDFGraph(revision *Revision, files ...string) (*graph.Graph, error) {
	g := graph.NewGraph()
	if revision == initRevision {
		return g, nil
	}

	commit, err := rr.gitRepository.Commit(plumbing.NewHash(revision.ID))
	if err != nil {
		return g, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return g, err
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		if contains(files, f.Name) {
			content, e := f.Contents()
			if e != nil {
				return e
			}
			return g.Unmarshal([]byte(content))
		}
		return nil
	})
	return g, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

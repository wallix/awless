package triplestore_test

import (
	"bytes"
	"fmt"
	"testing"

	tstore "github.com/wallix/triplestore"
)

func TestGraphTraverse(t *testing.T) {
	s := tstore.NewSource()
	s.Add(
		tstore.SubjPred("1", "->").Resource("2"),
		tstore.SubjPred("2", "->").Resource("3"),
		tstore.SubjPred("2", "->").Resource("4"),
		tstore.SubjPred("3", "->").Resource("5"),
		tstore.SubjPred("3", "->").Resource("6"),
		tstore.SubjPred("3", "->").Resource("7"),
		tstore.SubjPred("4", "->").Resource("8"),
	)
	g := s.Snapshot()

	var result bytes.Buffer
	each := func(gph tstore.RDFGraph, subj string, depth int) error {
		result.WriteString(fmt.Sprintf("(%d)%s ", depth, subj))
		return nil
	}

	tree := tstore.NewTree(g, "->")

	t.Run("depth forst search", func(t *testing.T) {
		tree.TraverseDFS("1", each)
		if got, want := result.String(), "(0)1 (1)2 (2)3 (3)5 (3)6 (3)7 (2)4 (3)8 "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		result.Reset()
		tree.TraverseDFS("8", each)
		if got, want := result.String(), "(0)8 "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		result.Reset()
		tree.TraverseDFS("4", each)
		if got, want := result.String(), "(0)4 (1)8 "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		result.Reset()
		tree.TraverseDFS("none", each)
		if got, want := result.String(), "(0)none "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})

	t.Run("ancestors", func(t *testing.T) {
		result.Reset()
		tree.TraverseAncestors("6", each)
		if got, want := result.String(), "(0)6 (1)3 (2)2 (3)1 "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		result.Reset()
		tree.TraverseAncestors("1", each)
		if got, want := result.String(), "(0)1 "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}

		result.Reset()
		tree.TraverseAncestors("none", each)
		if got, want := result.String(), "(0)none "; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	})
}

func TestTraverseSiblings(t *testing.T) {
	s := tstore.NewSource()
	s.Add(
		tstore.SubjPred("1", "->").Resource("2"),
		tstore.SubjPred("1", "->").Resource("3"),
		tstore.SubjPred("3", "->").Resource("4"),
		tstore.SubjPred("3", "->").Resource("5"),
		tstore.SubjPred("3", "->").Resource("6"),
		tstore.SubjPred("3", "->").Resource("7"),
		tstore.SubjPred("3", "->").Resource("8"),
		tstore.SubjPred("3", "->").Resource("9"),
		tstore.SubjPred("5", "type").StringLiteral("donkey"),
		tstore.SubjPred("7", "type").StringLiteral("donkey"),
		tstore.SubjPred("9", "type").StringLiteral("donkey"),
	)
	g := s.Snapshot()

	var result bytes.Buffer
	each := func(gph tstore.RDFGraph, subj string, depth int) error {
		result.WriteString(fmt.Sprintf("(%d)%s ", depth, subj))
		return nil
	}

	tree := tstore.NewTree(g, "->")

	siblingCriteria := func(g tstore.RDFGraph, node string) (string, error) {
		tris := g.WithSubjPred(node, "type")
		if len(tris) > 0 {
			return tstore.ParseString(tris[0].Object())
		}
		return "", nil
	}

	tree.TraverseSiblings("5", siblingCriteria, each)
	if got, want := result.String(), "(0)5 (0)7 (0)9 "; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

package rdf

import (
	"sort"
	"testing"

	"github.com/google/badwolf/triple/node"
)

func TestEmptyGraphDiff(t *testing.T) {
	any, _ := node.NewNodeFromStrings("/any", "any")
	diffGraph, _ := Diff(any, NewGraph(), NewGraph())

	if got, want := diffGraph.size(), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestGraphDiffGivenNilRootNode(t *testing.T) {
	diffGraph, _ := Diff(nil, NewGraph(), NewGraph())

	if got, want := diffGraph.size(), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestGraphDiff(t *testing.T) {
	local := NewGraph()
	//       1
	//   2       3       4
	// 5   6   7  8        9
	//                      10
	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")
	three, _ := node.NewNodeFromStrings("/three", "3")
	four, _ := node.NewNodeFromStrings("/four", "4")
	five, _ := node.NewNodeFromStrings("/five", "5")
	six, _ := node.NewNodeFromStrings("/six", "6")
	seven, _ := node.NewNodeFromStrings("/seven", "7")
	eight, _ := node.NewNodeFromStrings("/eight", "8")
	nine, _ := node.NewNodeFromStrings("/nine", "9")
	ten, _ := node.NewNodeFromStrings("/ten", "10")

	local.Add(noErrTriple(one, ParentOfPredicate, two))
	local.Add(noErrTriple(one, ParentOfPredicate, three))
	local.Add(noErrTriple(one, ParentOfPredicate, four))
	local.Add(noErrTriple(two, ParentOfPredicate, five))
	local.Add(noErrTriple(two, ParentOfPredicate, six))
	local.Add(noErrTriple(three, ParentOfPredicate, seven))
	local.Add(noErrTriple(three, ParentOfPredicate, eight))
	local.Add(noErrTriple(four, ParentOfPredicate, nine))
	local.Add(noErrTriple(nine, ParentOfPredicate, ten))

	remote := NewGraph()
	//       1
	//   2       3       4
	// 5   6   7           9
	//                  10   11
	rone, _ := node.NewNodeFromStrings("/one", "1")
	rtwo, _ := node.NewNodeFromStrings("/two", "2")
	rthree, _ := node.NewNodeFromStrings("/three", "3")
	rfour, _ := node.NewNodeFromStrings("/four", "4")
	rfive, _ := node.NewNodeFromStrings("/five", "5")
	rsix, _ := node.NewNodeFromStrings("/six", "6")
	rseven, _ := node.NewNodeFromStrings("/seven", "7")
	rnine, _ := node.NewNodeFromStrings("/nine", "9")
	rten, _ := node.NewNodeFromStrings("/ten", "10")
	releven, _ := node.NewNodeFromStrings("/eleven", "11")

	remote.Add(noErrTriple(rone, ParentOfPredicate, rtwo))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rthree))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rfour))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rfive))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rsix))
	remote.Add(noErrTriple(rthree, ParentOfPredicate, rseven))
	remote.Add(noErrTriple(rfour, ParentOfPredicate, rnine))
	remote.Add(noErrTriple(rnine, ParentOfPredicate, rten))
	remote.Add(noErrTriple(rnine, ParentOfPredicate, releven))

	diffGraph, err := Diff(one, local, remote)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := diffGraph.size(), 12; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	diffTriples, _ := diffGraph.TriplesForPredicateName("diff")
	sort.Sort(&tripleSorter{diffTriples})

	if got, want := len(diffTriples), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := diffTriples[0].Subject().ID(), eight.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ := diffTriples[0].Object().Literal()
	if got, want := lit.String(), `"extra"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := diffTriples[1].Subject().ID(), releven.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[1].Object().Literal()
	if got, want := lit.String(), `"missing"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	// Diff the other way around

	diffGraph, err = Diff(one, remote, local)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := diffGraph.size(), 12; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	diffTriples, _ = diffGraph.TriplesForPredicateName("diff")
	sort.Sort(&tripleSorter{diffTriples})

	if got, want := len(diffTriples), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := diffTriples[0].Subject().ID(), eight.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[0].Object().Literal()
	if got, want := lit.String(), `"missing"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := diffTriples[1].Subject().ID(), releven.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[1].Object().Literal()
	if got, want := lit.String(), `"extra"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestGraphDiffStoppingShortOnDifferentNode(t *testing.T) {
	local := NewGraph()
	//       1
	//   2       3       9
	// 5   6   7  8        10
	//                       11
	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")
	three, _ := node.NewNodeFromStrings("/three", "3")
	five, _ := node.NewNodeFromStrings("/five", "5")
	six, _ := node.NewNodeFromStrings("/six", "6")
	seven, _ := node.NewNodeFromStrings("/seven", "7")
	eight, _ := node.NewNodeFromStrings("/eight", "8")
	nine, _ := node.NewNodeFromStrings("/nine", "9")
	ten, _ := node.NewNodeFromStrings("/ten", "10")
	eleven, _ := node.NewNodeFromStrings("/eleven", "11")

	local.Add(noErrTriple(one, ParentOfPredicate, two))
	local.Add(noErrTriple(one, ParentOfPredicate, three))
	local.Add(noErrTriple(one, ParentOfPredicate, nine))
	local.Add(noErrTriple(two, ParentOfPredicate, five))
	local.Add(noErrTriple(two, ParentOfPredicate, six))
	local.Add(noErrTriple(three, ParentOfPredicate, seven))
	local.Add(noErrTriple(three, ParentOfPredicate, eight))
	local.Add(noErrTriple(nine, ParentOfPredicate, ten))
	local.Add(noErrTriple(ten, ParentOfPredicate, eleven))

	remote := NewGraph()
	//       1
	//   2       3       4
	// 5   6   7   8        9
	//                        10
	rone, _ := node.NewNodeFromStrings("/one", "1")
	rtwo, _ := node.NewNodeFromStrings("/two", "2")
	rthree, _ := node.NewNodeFromStrings("/three", "3")
	rfour, _ := node.NewNodeFromStrings("/four", "4")
	rfive, _ := node.NewNodeFromStrings("/five", "5")
	rsix, _ := node.NewNodeFromStrings("/six", "6")
	rseven, _ := node.NewNodeFromStrings("/seven", "7")
	reight, _ := node.NewNodeFromStrings("/eight", "8")
	rnine, _ := node.NewNodeFromStrings("/nine", "9")
	rten, _ := node.NewNodeFromStrings("/ten", "10")

	remote.Add(noErrTriple(rone, ParentOfPredicate, rtwo))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rthree))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rfour))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rfive))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rsix))
	remote.Add(noErrTriple(rthree, ParentOfPredicate, rseven))
	remote.Add(noErrTriple(rthree, ParentOfPredicate, reight))
	remote.Add(noErrTriple(rfour, ParentOfPredicate, rnine))
	remote.Add(noErrTriple(rnine, ParentOfPredicate, rten))

	diffGraph, err := Diff(one, local, remote)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := diffGraph.size(), 10; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	diffTriples, _ := diffGraph.TriplesForPredicateName("diff")
	sort.Sort(&tripleSorter{diffTriples})

	if got, want := len(diffTriples), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := diffTriples[0].Subject().ID(), rfour.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ := diffTriples[0].Object().Literal()
	if got, want := lit.String(), `"missing"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := diffTriples[1].Subject().ID(), nine.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[1].Object().Literal()
	if got, want := lit.String(), `"extra"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	// Diff the other way around

	diffGraph, err = Diff(one, remote, local)
	if err != nil {
		t.Fatal(err)
	}

	diffTriples, _ = diffGraph.TriplesForPredicateName("diff")
	sort.Sort(&tripleSorter{diffTriples})

	if got, want := len(diffTriples), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := diffTriples[0].Subject().ID(), rfour.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[0].Object().Literal()
	if got, want := lit.String(), `"extra"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := diffTriples[1].Subject().ID(), nine.ID(); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	lit, _ = diffTriples[1].Object().Literal()
	if got, want := lit.String(), `"missing"^^type:text`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

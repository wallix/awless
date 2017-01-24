package rdf

import (
	"reflect"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
)

func TestEmptyGraphDiffOrGivenNilRootNode(t *testing.T) {
	differ := &defaultDiffer{ParentOfPredicate}

	any, _ := node.NewNodeFromStrings("/any", "any")
	diff, _ := differ.Run(any, NewGraph(), NewGraph())

	if got, want := diff.HasResourceDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := diff.HasDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	diff, _ = differ.Run(nil, NewGraph(), NewGraph())
	if got, want := diff.HasResourceDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
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
	literalA, err := literal.DefaultBuilder().Build(literal.Text, "a")
	if err != nil {
		t.Fatal(err)
	}
	literalB, err := literal.DefaultBuilder().Build(literal.Text, "b")
	if err != nil {
		t.Fatal(err)
	}
	local.Add(noErrLiteralTriple(one, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(two, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(three, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(four, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(five, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(six, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(seven, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(eight, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(nine, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(ten, HasTypePredicate, literalB))

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

	remote.Add(noErrLiteralTriple(rone, HasTypePredicate, literalA))
	remote.Add(noErrLiteralTriple(rtwo, HasTypePredicate, literalA))
	remote.Add(noErrLiteralTriple(rthree, HasTypePredicate, literalB))
	remote.Add(noErrLiteralTriple(rfour, HasTypePredicate, literalA))
	remote.Add(noErrLiteralTriple(rfive, HasTypePredicate, literalA))
	remote.Add(noErrLiteralTriple(rsix, HasTypePredicate, literalB))
	remote.Add(noErrLiteralTriple(rseven, HasTypePredicate, literalB))
	remote.Add(noErrLiteralTriple(rnine, HasTypePredicate, literalB))
	remote.Add(noErrLiteralTriple(rten, HasTypePredicate, literalB))
	remote.Add(noErrLiteralTriple(releven, HasTypePredicate, literalB))

	remote.Add(noErrTriple(rone, ParentOfPredicate, rtwo))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rthree))
	remote.Add(noErrTriple(rone, ParentOfPredicate, rfour))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rfive))
	remote.Add(noErrTriple(rtwo, ParentOfPredicate, rsix))
	remote.Add(noErrTriple(rthree, ParentOfPredicate, rseven))
	remote.Add(noErrTriple(rfour, ParentOfPredicate, rnine))
	remote.Add(noErrTriple(rnine, ParentOfPredicate, rten))
	remote.Add(noErrTriple(rnine, ParentOfPredicate, releven))

	differ := &defaultDiffer{ParentOfPredicate}

	diff, err := differ.Run(one, local, remote)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := diff.FullGraph().size(), 24; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := diff.HasResourceDiff(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := diff.HasDiff(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	missing := noErrLiteralTriple(releven, DiffPredicate, MissingLiteral)
	extra := noErrLiteralTriple(eight, DiffPredicate, ExtraLiteral)
	if got, want := diff.TriplesInDiff(), []*triple.Triple{extra, missing}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	t1 := noErrTriple(three, ParentOfPredicate, eight)
	t2 := noErrTriple(rnine, ParentOfPredicate, releven)
	t3 := noErrLiteralTriple(eight, HasTypePredicate, literalB)
	t4 := noErrLiteralTriple(releven, HasTypePredicate, literalB)
	t5 := noErrLiteralTriple(three, HasTypePredicate, literalA)
	t6 := noErrLiteralTriple(rthree, HasTypePredicate, literalB)

	if got, want := diff.Inserted(), []*triple.Triple{t3, t5, t1}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := diff.Deleted(), []*triple.Triple{t4, t2, t6}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	commons := []*triple.Triple{
		noErrLiteralTriple(five, HasTypePredicate, literalA),
		noErrLiteralTriple(four, HasTypePredicate, literalA),
		noErrTriple(four, ParentOfPredicate, nine),
		noErrLiteralTriple(nine, HasTypePredicate, literalB),
		noErrTriple(nine, ParentOfPredicate, ten),
		noErrLiteralTriple(one, HasTypePredicate, literalA),
		noErrTriple(one, ParentOfPredicate, four),
		noErrTriple(one, ParentOfPredicate, three),
		noErrTriple(one, ParentOfPredicate, two),
		noErrLiteralTriple(seven, HasTypePredicate, literalB),
		noErrLiteralTriple(six, HasTypePredicate, literalB),
		noErrLiteralTriple(ten, HasTypePredicate, literalB),
		noErrTriple(three, ParentOfPredicate, seven),
		noErrLiteralTriple(two, HasTypePredicate, literalA),
		noErrTriple(two, ParentOfPredicate, five),
		noErrTriple(two, ParentOfPredicate, six),
	}
	if got, want := diff.Common(), commons; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	// Diff the other way around

	diff, err = differ.Run(one, remote, local)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := diff.FullGraph().size(), 24; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := diff.HasResourceDiff(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := diff.HasDiff(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	missing = noErrLiteralTriple(eight, DiffPredicate, MissingLiteral)
	extra = noErrLiteralTriple(releven, DiffPredicate, ExtraLiteral)
	if got, want := diff.TriplesInDiff(), []*triple.Triple{missing, extra}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	if got, want := diff.Common(), commons; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestDiffBuildFromGraph(t *testing.T) {
	local := NewGraph()
	//       1
	//   2   3   4
	// 5
	one, _ := node.NewNodeFromStrings("/a", "1")
	two, _ := node.NewNodeFromStrings("/a", "2")
	three, _ := node.NewNodeFromStrings("/b", "3")
	four, _ := node.NewNodeFromStrings("/b", "4")
	five, _ := node.NewNodeFromStrings("/b", "5")

	local.Add(noErrTriple(one, ParentOfPredicate, two))
	local.Add(noErrTriple(one, ParentOfPredicate, three))
	local.Add(noErrTriple(one, ParentOfPredicate, four))
	local.Add(noErrTriple(two, ParentOfPredicate, five))

	literalA, err := literal.DefaultBuilder().Build(literal.Text, "a")
	if err != nil {
		t.Fatal(err)
	}
	literalB, err := literal.DefaultBuilder().Build(literal.Text, "b")
	if err != nil {
		t.Fatal(err)
	}
	local.Add(noErrLiteralTriple(one, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(two, HasTypePredicate, literalA))
	local.Add(noErrLiteralTriple(three, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(four, HasTypePredicate, literalB))
	local.Add(noErrLiteralTriple(five, HasTypePredicate, literalB))

	d := NewDiff(local)

	if got, want := d.FullGraph().MustMarshal(), local.MustMarshal(); got != want {
		t.Fatalf("got \n%s\nwant\n%s\n", got, want)
	}
	if got, want := d.HasResourceDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	t1 := noErrTriple(one, ParentOfPredicate, three)
	t2 := noErrLiteralTriple(three, HasTypePredicate, literalB)
	t3 := noErrLiteralTriple(three, HasTypePredicate, literalA) //Triple not in graph
	d.AddDeleted(t2, ParentOfPredicate)
	if got, want := d.HasResourceDiff(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasDiff(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	d.AddDeleted(t1, ParentOfPredicate)
	if got, want := d.HasResourceDiff(true), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasDiff(true), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	d.AddDeleted(t3, ParentOfPredicate)
	if got, want := d.FullGraph().MustMarshal(), local.MustMarshal(); got != want {
		t.Fatalf("got \n%s\nwant\n%s\n", got, want)
	}
	if got, want := len(d.Inserted()), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := d.Deleted(), []*triple.Triple{t1, t2}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := d.HasDeletedTriple(t1), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasDeletedTriple(t2), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasDeletedTriple(t3), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	six, _ := node.NewNodeFromStrings("/b", "6")
	t4 := noErrTriple(three, ParentOfPredicate, six)
	t5 := noErrLiteralTriple(six, HasTypePredicate, literalB)
	t6 := noErrLiteralTriple(three, HasTypePredicate, literalB) //Triple already in graph
	d.AddInserted(t4, ParentOfPredicate)
	d.AddInserted(t5, ParentOfPredicate)
	d.AddInserted(t6, ParentOfPredicate)
	local.Add(t4)
	if got, want := d.FullGraph().MustMarshal(), local.MustMarshal(); got != want {
		t.Fatalf("got \n%s\nwant\n%s\n", got, want)
	}
	if got, want := d.Inserted(), []*triple.Triple{t4, t5}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := d.HasInsertedTriple(t4), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasInsertedTriple(t5), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := d.HasInsertedTriple(t6), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	deleted := noErrLiteralTriple(three, DiffPredicate, MissingLiteral)
	created := noErrLiteralTriple(six, DiffPredicate, ExtraLiteral)
	if got, want := d.TriplesInDiff(true), []*triple.Triple{deleted, created}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

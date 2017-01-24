package rdf

import (
	"testing"

	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
)

func TestGetTriplesAndNodesForType(t *testing.T) {
	graph := NewGraph()
	iLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/instance")
	if err != nil {
		t.Fatal(err)
	}
	sLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/subnet")
	if err != nil {
		t.Fatal(err)
	}
	azero, _ := node.NewNodeFromStrings("/instance", "0")
	graph.Add(noErrLiteralTriple(azero, HasTypePredicate, iLiteral))
	aone, _ := node.NewNodeFromStrings("/instance", "1")
	graph.Add(noErrLiteralTriple(aone, HasTypePredicate, iLiteral))

	bzero, _ := node.NewNodeFromStrings("/subnet", "0")
	graph.Add(noErrLiteralTriple(bzero, HasTypePredicate, sLiteral))
	bone, _ := node.NewNodeFromStrings("/subnet", "1")
	graph.Add(noErrLiteralTriple(bone, HasTypePredicate, sLiteral))
	btwo, _ := node.NewNodeFromStrings("/subnet", "2")
	graph.Add(noErrLiteralTriple(btwo, HasTypePredicate, sLiteral))

	atwo, _ := node.NewNodeFromStrings("/instance", "2")
	graph.Add(noErrLiteralTriple(atwo, HasTypePredicate, iLiteral))
	athree, _ := node.NewNodeFromStrings("/instance", "3")
	graph.Add(noErrLiteralTriple(athree, HasTypePredicate, iLiteral))

	graph.Add(noErrTriple(azero, ParentOfPredicate, aone))
	graph.Add(noErrTriple(aone, ParentOfPredicate, bone))
	graph.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	graph.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	graph.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	aTriples, err := graph.TriplesForType("/instance")
	if err != nil {
		t.Fatal(err)
	}
	result := NewGraphFromTriples(aTriples)

	expect := NewGraph()
	expect.Add(noErrLiteralTriple(azero, HasTypePredicate, iLiteral))
	expect.Add(noErrLiteralTriple(aone, HasTypePredicate, iLiteral))
	expect.Add(noErrLiteralTriple(atwo, HasTypePredicate, iLiteral))
	expect.Add(noErrLiteralTriple(athree, HasTypePredicate, iLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	bTriples, err := graph.TriplesForType("/subnet")
	if err != nil {
		t.Fatal(err)
	}
	result = NewGraphFromTriples(bTriples)

	expect = NewGraph()

	expect.Add(noErrLiteralTriple(bzero, HasTypePredicate, sLiteral))
	expect.Add(noErrLiteralTriple(bone, HasTypePredicate, sLiteral))
	expect.Add(noErrLiteralTriple(btwo, HasTypePredicate, sLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	aNodesExpect := []*node.Node{azero, aone, atwo, athree}
	aNodes, err := graph.NodesForType("/instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aNodes, aNodesExpect; !sameElementsInSlice(got, want) {
		t.Fatalf("got %#v\nwant%#v\n", got, want)
	}
	bNodesExpect := []*node.Node{bzero, bone, btwo}
	bNodes, err := graph.NodesForType("/subnet")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := bNodes, bNodesExpect; !sameElementsInSlice(got, want) {
		t.Fatalf("got %#v\nwant%#v\n", got, want)
	}
}

func TestGetTriplesForPredicateName(t *testing.T) {
	aLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/a")
	if err != nil {
		t.Fatal(err)
	}
	bLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/b")
	if err != nil {
		t.Fatal(err)
	}
	g := NewGraph()
	azero, _ := node.NewNodeFromStrings("/a", "0")
	g.Add(noErrLiteralTriple(azero, HasTypePredicate, aLiteral))
	aone, _ := node.NewNodeFromStrings("/a", "1")
	g.Add(noErrLiteralTriple(aone, HasTypePredicate, aLiteral))
	atwo, _ := node.NewNodeFromStrings("/a", "2")
	g.Add(noErrLiteralTriple(atwo, HasTypePredicate, aLiteral))

	bzero, _ := node.NewNodeFromStrings("/b", "0")
	g.Add(noErrLiteralTriple(bzero, HasTypePredicate, bLiteral))
	bone, _ := node.NewNodeFromStrings("/b", "1")
	g.Add(noErrLiteralTriple(bone, HasTypePredicate, bLiteral))
	btwo, _ := node.NewNodeFromStrings("/b", "2")
	g.Add(noErrLiteralTriple(btwo, HasTypePredicate, bLiteral))
	bthree, _ := node.NewNodeFromStrings("/b", "3")
	g.Add(noErrLiteralTriple(bthree, HasTypePredicate, bLiteral))

	g.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	g.Add(noErrTriple(aone, ParentOfPredicate, bone))
	g.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	triples, err := g.TriplesForPredicateName(string(ParentOfPredicate.ID()))
	if err != nil {
		t.Fatal(err)
	}
	result := NewGraphFromTriples(triples)

	expect := NewGraph()
	expect.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	expect.Add(noErrTriple(aone, ParentOfPredicate, bone))
	expect.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	triples, err = g.TriplesForPredicateName(string(HasTypePredicate.ID()))
	if err != nil {
		t.Fatal(err)
	}
	result = NewGraphFromTriples(triples)

	expect = NewGraph()
	expect.Add(noErrLiteralTriple(azero, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(aone, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(atwo, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(bzero, HasTypePredicate, bLiteral))
	expect.Add(noErrLiteralTriple(bone, HasTypePredicate, bLiteral))
	expect.Add(noErrLiteralTriple(btwo, HasTypePredicate, bLiteral))
	expect.Add(noErrLiteralTriple(bthree, HasTypePredicate, bLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestCountTriples(t *testing.T) {
	g := NewGraph()

	iLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/instance")
	if err != nil {
		t.Fatal(err)
	}
	sLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/subnet")
	if err != nil {
		t.Fatal(err)
	}
	one, _ := node.NewNodeFromStrings("/one", "1")
	g.Add(noErrLiteralTriple(one, HasTypePredicate, iLiteral))
	two, _ := node.NewNodeFromStrings("/two", "2")
	g.Add(noErrLiteralTriple(two, HasTypePredicate, iLiteral))
	three, _ := node.NewNodeFromStrings("/three", "3")
	g.Add(noErrLiteralTriple(three, HasTypePredicate, iLiteral))
	four, _ := node.NewNodeFromStrings("/four", "4")
	g.Add(noErrLiteralTriple(four, HasTypePredicate, iLiteral))
	five, _ := node.NewNodeFromStrings("/five", "5")
	g.Add(noErrLiteralTriple(five, HasTypePredicate, iLiteral))
	six, _ := node.NewNodeFromStrings("/six", "6")
	g.Add(noErrLiteralTriple(six, HasTypePredicate, sLiteral))
	seven, _ := node.NewNodeFromStrings("/seven", "7")
	g.Add(noErrLiteralTriple(seven, HasTypePredicate, sLiteral))
	eight, _ := node.NewNodeFromStrings("/eight", "8")
	g.Add(noErrLiteralTriple(eight, HasTypePredicate, sLiteral))
	nine, _ := node.NewNodeFromStrings("/nine", "9")
	g.Add(noErrLiteralTriple(nine, HasTypePredicate, sLiteral))
	ten, _ := node.NewNodeFromStrings("/ten", "10")
	g.Add(noErrLiteralTriple(ten, HasTypePredicate, sLiteral))

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	g.Add(noErrTriple(one, ParentOfPredicate, three))
	g.Add(noErrTriple(one, ParentOfPredicate, four))
	g.Add(noErrTriple(two, ParentOfPredicate, five))
	g.Add(noErrTriple(two, ParentOfPredicate, six))
	g.Add(noErrTriple(three, ParentOfPredicate, seven))
	g.Add(noErrTriple(three, ParentOfPredicate, eight))
	g.Add(noErrTriple(four, ParentOfPredicate, nine))
	g.Add(noErrTriple(nine, ParentOfPredicate, ten))

	count, err := g.CountTriplesForSubjectAndPredicate(one, ParentOfPredicate)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 3; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
	count, err = g.CountTriplesForSubjectAndPredicate(two, ParentOfPredicate)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 2; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicate(four, ParentOfPredicate)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicate(ten, ParentOfPredicate)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 0; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(four, ParentOfPredicate, "/instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 0; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(four, ParentOfPredicate, "/subnet")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(two, ParentOfPredicate, "/instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(two, ParentOfPredicate, "/subnet")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
}

func sameElementsInSlice(a, b []*node.Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for _, ae := range a {
		found := false
		for _, be := range a {
			if ae == be {
				found = true
			}
		}
		if found == false {
			return false
		}
	}
	return true
}

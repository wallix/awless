package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
)

func TestVisitDepthFirstGraph(t *testing.T) {
	g := NewGraph()

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

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	g.Add(noErrTriple(one, ParentOfPredicate, three))
	g.Add(noErrTriple(one, ParentOfPredicate, four))
	g.Add(noErrTriple(two, ParentOfPredicate, five))
	g.Add(noErrTriple(two, ParentOfPredicate, six))
	g.Add(noErrTriple(three, ParentOfPredicate, seven))
	g.Add(noErrTriple(three, ParentOfPredicate, eight))
	g.Add(noErrTriple(four, ParentOfPredicate, nine))
	g.Add(noErrTriple(nine, ParentOfPredicate, ten))

	var result bytes.Buffer
	each := func(g *Graph, n *node.Node, distance int) {
		for i := 0; i < distance; i++ {
			result.WriteByte('/')
		}
		result.WriteString(n.ID().String())
	}

	g.VisitDepthFirst(one, each)
	if got, want := result.String(), "1/2//5//6/3//7//8/4//9///10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirst(four, each)
	if got, want := result.String(), "4/9//10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirst(three, each)
	if got, want := result.String(), "3/7/8"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}
}

func TestIntersectGraph(t *testing.T) {
	azero, _ := node.NewNodeFromStrings("/a", "0")
	aone, _ := node.NewNodeFromStrings("/a", "1")
	atwo, _ := node.NewNodeFromStrings("/a", "2")
	athree, _ := node.NewNodeFromStrings("/a", "3")
	afour, _ := node.NewNodeFromStrings("/a", "4")
	afive, _ := node.NewNodeFromStrings("/a", "5")
	asix, _ := node.NewNodeFromStrings("/a", "6")

	bzero, _ := node.NewNodeFromStrings("/b", "0")
	bone, _ := node.NewNodeFromStrings("/b", "1")
	btwo, _ := node.NewNodeFromStrings("/b", "2")
	bthree, _ := node.NewNodeFromStrings("/b", "3")
	bfour, _ := node.NewNodeFromStrings("/b", "4")
	bfive, _ := node.NewNodeFromStrings("/b", "5")
	bsix, _ := node.NewNodeFromStrings("/b", "6")

	g1 := NewGraph()
	g1.Add(noErrTriple(aone, ParentOfPredicate, bone))
	g1.Add(noErrTriple(atwo, ParentOfPredicate, btwo))
	g1.Add(noErrTriple(athree, ParentOfPredicate, bthree))
	g1.Add(noErrTriple(afour, ParentOfPredicate, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	g2.Add(noErrTriple(atwo, ParentOfPredicate, btwo))
	g2.Add(noErrTriple(athree, ParentOfPredicate, bthree))
	g2.Add(noErrTriple(afive, ParentOfPredicate, bfive))
	g2.Add(noErrTriple(asix, ParentOfPredicate, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(atwo, ParentOfPredicate, btwo))
	expect.Add(noErrTriple(athree, ParentOfPredicate, bthree))

	result := g1.Intersect(g2)
	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = g2.Intersect(g1)
	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestSubstractGraph(t *testing.T) {
	azero, _ := node.NewNodeFromStrings("/a", "0")
	aone, _ := node.NewNodeFromStrings("/a", "1")
	atwo, _ := node.NewNodeFromStrings("/a", "2")
	athree, _ := node.NewNodeFromStrings("/a", "3")
	afour, _ := node.NewNodeFromStrings("/a", "4")
	afive, _ := node.NewNodeFromStrings("/a", "5")
	asix, _ := node.NewNodeFromStrings("/a", "6")

	bzero, _ := node.NewNodeFromStrings("/b", "0")
	bone, _ := node.NewNodeFromStrings("/b", "1")
	btwo, _ := node.NewNodeFromStrings("/b", "2")
	bthree, _ := node.NewNodeFromStrings("/b", "3")
	bfour, _ := node.NewNodeFromStrings("/b", "4")
	bfive, _ := node.NewNodeFromStrings("/b", "5")
	bsix, _ := node.NewNodeFromStrings("/b", "6")

	g1 := NewGraph()
	g1.Add(noErrTriple(aone, ParentOfPredicate, bone))
	g1.Add(noErrTriple(atwo, ParentOfPredicate, btwo))
	g1.Add(noErrTriple(athree, ParentOfPredicate, bthree))
	g1.Add(noErrTriple(afour, ParentOfPredicate, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	g2.Add(noErrTriple(atwo, ParentOfPredicate, btwo))
	g2.Add(noErrTriple(athree, ParentOfPredicate, bthree))
	g2.Add(noErrTriple(afive, ParentOfPredicate, bfive))
	g2.Add(noErrTriple(asix, ParentOfPredicate, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(aone, ParentOfPredicate, bone))
	expect.Add(noErrTriple(afour, ParentOfPredicate, bfour))

	result := g1.Substract(g2)
	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	expect = NewGraph()
	expect.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	expect.Add(noErrTriple(afive, ParentOfPredicate, bfive))
	expect.Add(noErrTriple(asix, ParentOfPredicate, bsix))

	result = g2.Substract(g1)
	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func TestGetTriplesAndNodesForType(t *testing.T) {
	graph := NewGraph()
	aLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/a")
	if err != nil {
		t.Fatal(err)
	}
	bLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/b")
	if err != nil {
		t.Fatal(err)
	}
	azero, _ := node.NewNodeFromStrings("/a", "0")
	graph.Add(noErrLiteralTriple(azero, HasTypePredicate, aLiteral))
	aone, _ := node.NewNodeFromStrings("/a", "1")
	graph.Add(noErrLiteralTriple(aone, HasTypePredicate, aLiteral))

	bzero, _ := node.NewNodeFromStrings("/b", "0")
	graph.Add(noErrLiteralTriple(bzero, HasTypePredicate, bLiteral))
	bone, _ := node.NewNodeFromStrings("/b", "1")
	graph.Add(noErrLiteralTriple(bone, HasTypePredicate, bLiteral))
	btwo, _ := node.NewNodeFromStrings("/b", "2")
	graph.Add(noErrLiteralTriple(btwo, HasTypePredicate, bLiteral))

	atwo, _ := node.NewNodeFromStrings("/a", "2")
	graph.Add(noErrLiteralTriple(atwo, HasTypePredicate, aLiteral))
	athree, _ := node.NewNodeFromStrings("/a", "3")
	graph.Add(noErrLiteralTriple(athree, HasTypePredicate, aLiteral))

	graph.Add(noErrTriple(azero, ParentOfPredicate, aone))
	graph.Add(noErrTriple(aone, ParentOfPredicate, bone))
	graph.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	graph.Add(noErrTriple(azero, ParentOfPredicate, bzero))
	graph.Add(noErrTriple(atwo, ParentOfPredicate, btwo))

	aTriples, err := graph.TriplesForType("/a")
	if err != nil {
		t.Fatal(err)
	}
	result := NewGraphFromTriples(aTriples)

	expect := NewGraph()
	expect.Add(noErrLiteralTriple(azero, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(aone, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(atwo, HasTypePredicate, aLiteral))
	expect.Add(noErrLiteralTriple(athree, HasTypePredicate, aLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	bTriples, err := graph.TriplesForType("/b")
	if err != nil {
		t.Fatal(err)
	}
	result = NewGraphFromTriples(bTriples)

	expect = NewGraph()

	expect.Add(noErrLiteralTriple(bzero, HasTypePredicate, bLiteral))
	expect.Add(noErrLiteralTriple(bone, HasTypePredicate, bLiteral))
	expect.Add(noErrLiteralTriple(btwo, HasTypePredicate, bLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	aNodesExpect := []*node.Node{azero, aone, atwo, athree}
	aNodes, err := graph.NodesForType("/a")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := aNodes, aNodesExpect; !sameElementsInSlice(got, want) {
		t.Fatalf("got %#v\nwant%#v\n", got, want)
	}
	bNodesExpect := []*node.Node{bzero, bone, btwo}
	bNodes, err := graph.NodesForType("/b")
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

	aLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/a")
	if err != nil {
		t.Fatal(err)
	}
	bLiteral, err := literal.DefaultBuilder().Build(literal.Text, "/b")
	if err != nil {
		t.Fatal(err)
	}
	one, _ := node.NewNodeFromStrings("/one", "1")
	g.Add(noErrLiteralTriple(one, HasTypePredicate, aLiteral))
	two, _ := node.NewNodeFromStrings("/two", "2")
	g.Add(noErrLiteralTriple(two, HasTypePredicate, aLiteral))
	three, _ := node.NewNodeFromStrings("/three", "3")
	g.Add(noErrLiteralTriple(three, HasTypePredicate, aLiteral))
	four, _ := node.NewNodeFromStrings("/four", "4")
	g.Add(noErrLiteralTriple(four, HasTypePredicate, aLiteral))
	five, _ := node.NewNodeFromStrings("/five", "5")
	g.Add(noErrLiteralTriple(five, HasTypePredicate, aLiteral))
	six, _ := node.NewNodeFromStrings("/six", "6")
	g.Add(noErrLiteralTriple(six, HasTypePredicate, bLiteral))
	seven, _ := node.NewNodeFromStrings("/seven", "7")
	g.Add(noErrLiteralTriple(seven, HasTypePredicate, bLiteral))
	eight, _ := node.NewNodeFromStrings("/eight", "8")
	g.Add(noErrLiteralTriple(eight, HasTypePredicate, bLiteral))
	nine, _ := node.NewNodeFromStrings("/nine", "9")
	g.Add(noErrLiteralTriple(nine, HasTypePredicate, bLiteral))
	ten, _ := node.NewNodeFromStrings("/ten", "10")
	g.Add(noErrLiteralTriple(ten, HasTypePredicate, bLiteral))

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

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(four, ParentOfPredicate, "/a")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 0; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(four, ParentOfPredicate, "/b")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(two, ParentOfPredicate, "/a")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicateObjectOfType(two, ParentOfPredicate, "/b")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
}

func TestGraphSize(t *testing.T) {
	g := NewGraph()
	if got, want := g.IsEmpty(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	one, _ := node.NewNodeFromStrings("/one", "1")
	two, _ := node.NewNodeFromStrings("/two", "2")

	g.Add(noErrTriple(one, ParentOfPredicate, two))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	g.Add(noErrTriple(two, ParentOfPredicate, one))
	if got, want := g.IsEmpty(), false; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := g.size(), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
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

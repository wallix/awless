package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
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

	g.Add(noErrTriple(one, ParentOf, two))
	g.Add(noErrTriple(one, ParentOf, three))
	g.Add(noErrTriple(one, ParentOf, four))
	g.Add(noErrTriple(two, ParentOf, five))
	g.Add(noErrTriple(two, ParentOf, six))
	g.Add(noErrTriple(three, ParentOf, seven))
	g.Add(noErrTriple(three, ParentOf, eight))
	g.Add(noErrTriple(four, ParentOf, nine))
	g.Add(noErrTriple(nine, ParentOf, ten))

	var result bytes.Buffer
	each := func(n *node.Node, distance int) {
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
	g1.Add(noErrTriple(aone, ParentOf, bone))
	g1.Add(noErrTriple(atwo, ParentOf, btwo))
	g1.Add(noErrTriple(athree, ParentOf, bthree))
	g1.Add(noErrTriple(afour, ParentOf, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, ParentOf, bzero))
	g2.Add(noErrTriple(atwo, ParentOf, btwo))
	g2.Add(noErrTriple(athree, ParentOf, bthree))
	g2.Add(noErrTriple(afive, ParentOf, bfive))
	g2.Add(noErrTriple(asix, ParentOf, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(atwo, ParentOf, btwo))
	expect.Add(noErrTriple(athree, ParentOf, bthree))

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
	g1.Add(noErrTriple(aone, ParentOf, bone))
	g1.Add(noErrTriple(atwo, ParentOf, btwo))
	g1.Add(noErrTriple(athree, ParentOf, bthree))
	g1.Add(noErrTriple(afour, ParentOf, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, ParentOf, bzero))
	g2.Add(noErrTriple(atwo, ParentOf, btwo))
	g2.Add(noErrTriple(athree, ParentOf, bthree))
	g2.Add(noErrTriple(afive, ParentOf, bfive))
	g2.Add(noErrTriple(asix, ParentOf, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(aone, ParentOf, bone))
	expect.Add(noErrTriple(afour, ParentOf, bfour))

	result := g1.Substract(g2)
	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	expect = NewGraph()
	expect.Add(noErrTriple(azero, ParentOf, bzero))
	expect.Add(noErrTriple(afive, ParentOf, bfive))
	expect.Add(noErrTriple(asix, ParentOf, bsix))

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
	graph.Add(noErrLiteralTriple(azero, HasType, aLiteral))
	aone, _ := node.NewNodeFromStrings("/a", "1")
	graph.Add(noErrLiteralTriple(aone, HasType, aLiteral))

	bzero, _ := node.NewNodeFromStrings("/b", "0")
	graph.Add(noErrLiteralTriple(bzero, HasType, bLiteral))
	bone, _ := node.NewNodeFromStrings("/b", "1")
	graph.Add(noErrLiteralTriple(bone, HasType, bLiteral))
	btwo, _ := node.NewNodeFromStrings("/b", "2")
	graph.Add(noErrLiteralTriple(btwo, HasType, bLiteral))

	atwo, _ := node.NewNodeFromStrings("/a", "2")
	graph.Add(noErrLiteralTriple(atwo, HasType, aLiteral))
	athree, _ := node.NewNodeFromStrings("/a", "3")
	graph.Add(noErrLiteralTriple(athree, HasType, aLiteral))

	graph.Add(noErrTriple(azero, ParentOf, aone))
	graph.Add(noErrTriple(aone, ParentOf, bone))
	graph.Add(noErrTriple(atwo, ParentOf, btwo))

	graph.Add(noErrTriple(azero, ParentOf, bzero))
	graph.Add(noErrTriple(atwo, ParentOf, btwo))

	aTriples, err := graph.TriplesForType("/a")
	if err != nil {
		t.Fatal(err)
	}
	result := NewGraphFromTriples(aTriples)

	expect := NewGraph()
	expect.Add(noErrLiteralTriple(azero, HasType, aLiteral))
	expect.Add(noErrLiteralTriple(aone, HasType, aLiteral))
	expect.Add(noErrLiteralTriple(atwo, HasType, aLiteral))
	expect.Add(noErrLiteralTriple(athree, HasType, aLiteral))

	if got, want := result.MustMarshal(), expect.MustMarshal(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	bTriples, err := graph.TriplesForType("/b")
	if err != nil {
		t.Fatal(err)
	}
	result = NewGraphFromTriples(bTriples)

	expect = NewGraph()

	expect.Add(noErrLiteralTriple(bzero, HasType, bLiteral))
	expect.Add(noErrLiteralTriple(bone, HasType, bLiteral))
	expect.Add(noErrLiteralTriple(btwo, HasType, bLiteral))

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

func TestCountTriples(t *testing.T) {
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

	g.Add(noErrTriple(one, ParentOf, two))
	g.Add(noErrTriple(one, ParentOf, three))
	g.Add(noErrTriple(one, ParentOf, four))
	g.Add(noErrTriple(two, ParentOf, five))
	g.Add(noErrTriple(two, ParentOf, six))
	g.Add(noErrTriple(three, ParentOf, seven))
	g.Add(noErrTriple(three, ParentOf, eight))
	g.Add(noErrTriple(four, ParentOf, nine))
	g.Add(noErrTriple(nine, ParentOf, ten))

	count, err := g.CountTriplesForSubjectAndPredicate(one, ParentOf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 3; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
	count, err = g.CountTriplesForSubjectAndPredicate(two, ParentOf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 2; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicate(four, ParentOf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 1; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}

	count, err = g.CountTriplesForSubjectAndPredicate(ten, ParentOf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := count, 0; got != want {
		t.Fatalf("got %d; want%d\n", got, want)
	}
}

func noErrTriple(s *node.Node, p *predicate.Predicate, o *node.Node) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewNodeObject(o))
	if err != nil {
		panic(err)
	}
	return tri
}

func noErrLiteralTriple(s *node.Node, p *predicate.Predicate, l *literal.Literal) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewLiteralObject(l))
	if err != nil {
		panic(err)
	}
	return tri
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

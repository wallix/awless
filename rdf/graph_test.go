package rdf

import (
	"bytes"
	"testing"

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

func TestVisitDepthFirstUnique(t *testing.T) {
	g := NewGraph()

	//       1
	//   2       3       4    8
	// 5   6   7  8        9    6
	//            6         10
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
	g.Add(noErrTriple(one, ParentOfPredicate, eight))
	g.Add(noErrTriple(two, ParentOfPredicate, five))
	g.Add(noErrTriple(two, ParentOfPredicate, six))
	g.Add(noErrTriple(three, ParentOfPredicate, seven))
	g.Add(noErrTriple(three, ParentOfPredicate, eight))
	g.Add(noErrTriple(eight, ParentOfPredicate, six))
	g.Add(noErrTriple(four, ParentOfPredicate, nine))
	g.Add(noErrTriple(nine, ParentOfPredicate, ten))

	var result bytes.Buffer
	each := func(g *Graph, n *node.Node, distance int) error {
		for i := 0; i < distance; i++ {
			result.WriteByte('/')
		}
		result.WriteString(n.ID().String())
		return nil
	}

	g.VisitDepthFirstUnique(one, each)
	if got, want := result.String(), "1/2//5//6/3//7//8/4//9///10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirstUnique(four, each)
	if got, want := result.String(), "4/9//10"; got != want {
		t.Fatalf("got '%s', want '%s'", got, want)
	}

	result.Reset()
	g.VisitDepthFirstUnique(three, each)
	if got, want := result.String(), "3/7/8//6"; got != want {
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

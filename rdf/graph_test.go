package rdf

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
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

	g.Add(noErrTriple(one, parentOf, two))
	g.Add(noErrTriple(one, parentOf, three))
	g.Add(noErrTriple(one, parentOf, four))
	g.Add(noErrTriple(two, parentOf, five))
	g.Add(noErrTriple(two, parentOf, six))
	g.Add(noErrTriple(three, parentOf, seven))
	g.Add(noErrTriple(three, parentOf, eight))
	g.Add(noErrTriple(four, parentOf, nine))
	g.Add(noErrTriple(nine, parentOf, ten))

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
	g1.Add(noErrTriple(aone, parentOf, bone))
	g1.Add(noErrTriple(atwo, parentOf, btwo))
	g1.Add(noErrTriple(athree, parentOf, bthree))
	g1.Add(noErrTriple(afour, parentOf, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, parentOf, bzero))
	g2.Add(noErrTriple(atwo, parentOf, btwo))
	g2.Add(noErrTriple(athree, parentOf, bthree))
	g2.Add(noErrTriple(afive, parentOf, bfive))
	g2.Add(noErrTriple(asix, parentOf, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(atwo, parentOf, btwo))
	expect.Add(noErrTriple(athree, parentOf, bthree))

	result := g1.Intersect(g2)
	if got, want := result.FlushString(), expect.FlushString(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	result = g2.Intersect(g1)
	if got, want := result.FlushString(), expect.FlushString(); got != want {
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
	g1.Add(noErrTriple(aone, parentOf, bone))
	g1.Add(noErrTriple(atwo, parentOf, btwo))
	g1.Add(noErrTriple(athree, parentOf, bthree))
	g1.Add(noErrTriple(afour, parentOf, bfour))

	g2 := NewGraph()
	g2.Add(noErrTriple(azero, parentOf, bzero))
	g2.Add(noErrTriple(atwo, parentOf, btwo))
	g2.Add(noErrTriple(athree, parentOf, bthree))
	g2.Add(noErrTriple(afive, parentOf, bfive))
	g2.Add(noErrTriple(asix, parentOf, bsix))

	expect := NewGraph()
	expect.Add(noErrTriple(aone, parentOf, bone))
	expect.Add(noErrTriple(afour, parentOf, bfour))

	result := g1.Substract(g2)
	if got, want := result.FlushString(), expect.FlushString(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}

	expect = NewGraph()
	expect.Add(noErrTriple(azero, parentOf, bzero))
	expect.Add(noErrTriple(afive, parentOf, bfive))
	expect.Add(noErrTriple(asix, parentOf, bsix))

	result = g2.Substract(g1)
	if got, want := result.FlushString(), expect.FlushString(); got != want {
		t.Fatalf("got %s\nwant%s\n", got, want)
	}
}

func noErrTriple(s *node.Node, p *predicate.Predicate, o *node.Node) *triple.Triple {
	tri, err := triple.New(s, p, triple.NewNodeObject(o))
	if err != nil {
		panic(err)
	}
	return tri
}

package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func NewTripleFromStrings(sub *node.Node, pred string, obj string) (*triple.Triple, error) {
	p, err := predicate.NewImmutable(pred)
	if err != nil {
		return nil, err
	}
	objL, err := literal.DefaultBuilder().Build(literal.Text, obj)
	if err != nil {
		return nil, err
	}
	return triple.New(sub, p, triple.NewLiteralObject(objL))
}

func attachLiteralToNode(g *Graph, n *node.Node, p *predicate.Predicate, lit *literal.Literal) error {
	tri, err := triple.New(n, p, triple.NewLiteralObject(lit))
	if err != nil {
		return err
	}

	g.Add(tri)
	return nil
}

func intersectTriples(a, b []*triple.Triple) []*triple.Triple {
	var inter []*triple.Triple

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				inter = append(inter, a[i])
			}
		}
	}

	return inter
}

func substractTriples(a, b []*triple.Triple) []*triple.Triple {
	var sub []*triple.Triple

	for i := 0; i < len(a); i++ {
		var found bool
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				found = true
			}
		}
		if !found {
			sub = append(sub, a[i])
		}
	}

	return sub
}

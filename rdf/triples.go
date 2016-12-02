package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/predicate"
)

func AttachLiteralToAllTriples(g *Graph, p *predicate.Predicate, lit *literal.Literal) error {
	all, _ := g.allTriples()
	for _, t := range all {
		err := attachLiteralToTriple(g, t, p, lit)
		if err != nil {
			return err
		}
	}
	return nil
}

func attachLiteralToTriple(g *Graph, t *triple.Triple, p *predicate.Predicate, lit *literal.Literal) error {
	node, err := t.Object().Node()
	if err != nil {
		return err
	}
	tri, err := triple.New(node, p, triple.NewLiteralObject(lit))
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

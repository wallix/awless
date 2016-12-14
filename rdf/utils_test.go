package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func noErrNode(typ, id string) *node.Node {
	n, _ := node.NewNodeFromStrings(typ, id)
	return n
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

func MarshalTriples(triples []*triple.Triple) string {
	g := NewGraphFromTriples(triples)
	return g.MustMarshal()
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}

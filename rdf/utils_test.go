package rdf

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
)

func marshalTriples(triples []*triple.Triple) string {
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

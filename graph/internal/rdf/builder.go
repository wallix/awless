package rdf

import (
	"strings"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

type tripleBuilder struct {
	sub  *node.Node
	pred *predicate.Predicate
}

func Subject(s string, ns ...string) *tripleBuilder {
	return &tripleBuilder{sub: MustBuildNode(addNs(s, ns...))}
}

func (b *tripleBuilder) Predicate(s string, ns ...string) *tripleBuilder {
	b.pred = MustBuildPredicate(addNs(s, ns...))
	return b
}

func (b *tripleBuilder) Literal(s string) *triple.Triple {
	t, err := triple.New(b.sub, b.pred, triple.NewLiteralObject(MustBuildLiteral(s)))
	if err != nil {
		panic(err)
	}
	return t
}

func (b *tripleBuilder) Object(s string, ns ...string) *triple.Triple {
	return b.ObjectNode(triple.NewNodeObject(MustBuildNode(addNs(s, ns...))))
}

func (b *tripleBuilder) ObjectNode(obj *triple.Object) *triple.Triple {
	t, err := triple.New(b.sub, b.pred, obj)
	if err != nil {
		panic(err)
	}
	return t
}

func MustBuildPredicate(name string) *predicate.Predicate {
	pred, err := predicate.NewImmutable(name)
	if err != nil {
		panic(err)
	}
	return pred
}

func MustBuildNode(name string) *node.Node {
	node, err := node.NewNodeFromStrings("/node", name)
	if err != nil {
		panic(err)
	}
	return node
}

func MustBuildLiteral(str string) *literal.Literal {
	lit, err := literal.DefaultBuilder().Build(literal.Text, str)
	if err != nil {
		panic(err)
	}
	return lit
}

func TrimNS(s string) string {
	spl := strings.Split(s, ":")
	if len(spl) == 0 {
		return s
	}
	return spl[len(spl)-1]
}

func addNs(s string, ns ...string) string {
	return strings.Join(append(ns, s), ":")
}

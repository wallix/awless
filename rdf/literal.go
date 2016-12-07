package rdf

import "github.com/google/badwolf/triple/literal"

var (
	ExtraLiteral   *literal.Literal
	MissingLiteral *literal.Literal
)

func init() {
	var err error
	if ExtraLiteral, err = literal.DefaultBuilder().Build(literal.Text, "extra"); err != nil {
		panic(err)
	}
	if MissingLiteral, err = literal.DefaultBuilder().Build(literal.Text, "missing"); err != nil {
		panic(err)
	}
}

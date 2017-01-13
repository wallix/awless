package rdf

import "github.com/google/badwolf/triple/literal"

var (
	ExtraLiteral    *literal.Literal
	MissingLiteral  *literal.Literal
	RegionLiteral   *literal.Literal
	InstanceLiteral *literal.Literal
)

func init() {
	var err error

	if RegionLiteral, err = literal.DefaultBuilder().Build(literal.Text, Region.ToRDFString()); err != nil {
		panic(err)
	}
	if ExtraLiteral, err = literal.DefaultBuilder().Build(literal.Text, "extra"); err != nil {
		panic(err)
	}
	if MissingLiteral, err = literal.DefaultBuilder().Build(literal.Text, "missing"); err != nil {
		panic(err)
	}
	if InstanceLiteral, err = literal.DefaultBuilder().Build(literal.Text, Instance.ToRDFString()); err != nil {
		panic(err)
	}
}

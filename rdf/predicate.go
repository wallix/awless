package rdf

import "github.com/google/badwolf/triple/predicate"

var (
	ParentOfPredicate *predicate.Predicate
	HasTypePredicate  *predicate.Predicate
	DiffPredicate     *predicate.Predicate
	PropertyPredicate *predicate.Predicate
)

func init() {
	var err error
	if ParentOfPredicate, err = predicate.NewImmutable("parent_of"); err != nil {
		panic(err)
	}
	if HasTypePredicate, err = predicate.NewImmutable("has_type"); err != nil {
		panic(err)
	}
	if DiffPredicate, err = predicate.NewImmutable("diff"); err != nil {
		panic(err)
	}
	if PropertyPredicate, err = predicate.NewImmutable("property"); err != nil {
		panic(err)
	}
	DefaultDiffer = &defaultDiffer{ParentOfPredicate}
}

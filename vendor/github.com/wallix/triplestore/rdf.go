// Package triplestore provides APIs to manage, store and query triples, sources and RDFGraphs
package triplestore

// A triple consists of a subject, a predicate and a object
type Triple interface {
	Subject() string
	Predicate() string
	Object() Object
	Equal(Triple) bool
}

// An object is a resource (i.e. IRI) or a literal. Note that blank node are not supported.
type Object interface {
	Literal() (Literal, bool)
	Resource() (string, bool)
	Equal(Object) bool
}

// A literal is a unicode string associated with a datatype (ex: string, integer, ...).
type Literal interface {
	Type() XsdType
	Value() string
}

type subject string
type predicate string

type triple struct {
	sub    subject
	pred   predicate
	obj    object
	triKey string
}

func (t *triple) Object() Object {
	return t.obj
}

func (t *triple) Subject() string {
	return string(t.sub)
}

func (t *triple) Predicate() string {
	return string(t.pred)
}

func (t *triple) key() string {
	if t.triKey == "" {
		t.triKey = "<" + string(t.sub) + "><" + string(t.pred) + ">" + t.obj.key()
	}
	return t.triKey
}

func (t *triple) clone() *triple {
	return &triple{
		sub:    t.sub,
		pred:   t.pred,
		obj:    t.obj,
		triKey: t.triKey,
	}
}

func (t *triple) Equal(other Triple) bool {
	switch {
	case t == nil:
		return other == nil
	case other == nil:
		return false
	default:
		otherT, ok := other.(*triple)
		if !ok {
			return false
		}
		return t.key() == otherT.key()
	}
}

type object struct {
	isLit    bool
	resource string
	lit      literal
}

func (o object) Literal() (Literal, bool) {
	return o.lit, o.isLit
}

func (o object) Resource() (string, bool) {
	return o.resource, !o.isLit
}

func (o object) key() string {
	if o.isLit {
		return "\"" + o.lit.val + "\"^^" + string(o.lit.typ)
	}
	return "<" + o.resource + ">"
}

func (o object) Equal(other Object) bool {
	lit, ok := o.Literal()
	otherLit, otherOk := other.Literal()
	if ok != otherOk {
		return false
	}
	if ok {
		return lit.Type() == otherLit.Type() && lit.Value() == otherLit.Value()
	}
	res, ok := o.Resource()
	otherRes, otherOk := other.Resource()
	if ok != otherOk {
		return false
	}
	if ok {
		return res == otherRes
	}
	return true
}

type literal struct {
	typ XsdType
	val string
}

func (l literal) Type() XsdType {
	return l.typ
}

func (l literal) Value() string {
	return l.val
}

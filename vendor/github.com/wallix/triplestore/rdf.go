// Package triplestore provides APIs to manage, store and query triples, sources and RDFGraphs
package triplestore

// Triple consists of a subject, a predicate and a object
type Triple interface {
	Subject() string
	Predicate() string
	Object() Object
	Equal(Triple) bool
}

// Object is a resource (i.e. IRI), a literal or a blank node.
type Object interface {
	Literal() (Literal, bool)
	Resource() (string, bool)
	Bnode() (string, bool)
	Equal(Object) bool
}

// Literal is a unicode string associated with a datatype (ex: string, integer, ...).
type Literal interface {
	Type() XsdType
	Value() string
	Lang() string
}

type triple struct {
	sub, pred  string
	isSubBnode bool
	obj        object
	triKey     string
}

func (t *triple) Object() Object {
	return t.obj
}

func (t *triple) Subject() string {
	return t.sub
}

func (t *triple) Predicate() string {
	return t.pred
}

func (t *triple) key() string {
	if t.triKey == "" {
		var sub string
		if t.isSubBnode {
			sub = "_:" + t.sub
		} else {
			sub = "<" + t.sub + ">"
		}
		t.triKey = sub + "<" + t.pred + ">" + t.obj.key()
		return t.triKey
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
	isLit, isBnode  bool
	resource, bnode string
	lit             literal
}

func (o object) Literal() (Literal, bool) {
	return o.lit, o.isLit
}

func (o object) Resource() (string, bool) {
	return o.resource, !o.isLit
}

func (o object) Bnode() (string, bool) {
	return o.bnode, o.isBnode
}

func (o object) key() string {
	if o.isLit {
		if o.lit.langtag != "" {
			return "\"" + o.lit.val + "\"@" + o.lit.langtag
		}
		return "\"" + o.lit.val + "\"^^<" + string(o.lit.typ) + ">"
	}
	if o.isBnode {
		return "_:" + o.bnode
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
	typ          XsdType
	val, langtag string
}

func (l literal) Type() XsdType {
	return l.typ
}

func (l literal) Value() string {
	return l.val
}

func (l literal) Lang() string {
	return l.langtag
}

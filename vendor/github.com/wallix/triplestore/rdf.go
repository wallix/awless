package triplestore

type Triple interface {
	Subject() string
	Predicate() string
	Object() Object
	Equal(Triple) bool
}

type Object interface {
	Literal() (Literal, bool)
	ResourceID() (string, bool)
	Equal(Object) bool
}

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
	isLit      bool
	resourceID string
	lit        literal
}

func (o object) Literal() (Literal, bool) {
	return o.lit, o.isLit
}

func (o object) ResourceID() (string, bool) {
	return o.resourceID, !o.isLit
}

func (o object) key() string {
	if o.isLit {
		return "\"" + o.lit.val + "\"^^" + string(o.lit.typ)
	}
	return "<" + o.resourceID + ">"
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
	resId, ok := o.ResourceID()
	otherResId, otherOk := other.ResourceID()
	if ok != otherOk {
		return false
	}
	if ok {
		return resId == otherResId
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

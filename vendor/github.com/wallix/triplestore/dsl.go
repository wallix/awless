package triplestore

import (
	"fmt"
	"time"
)

type tripleBuilder struct {
	sub, pred string
}

func SubjPred(s, p string) *tripleBuilder {
	return &tripleBuilder{sub: s, pred: p}
}

func Resource(s string) Object {
	return object{resourceID: s}

}

func (b *tripleBuilder) Resource(s string) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  Resource(s).(object),
	}
}

func (b *tripleBuilder) Object(o Object) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  o.(object),
	}
}

func ObjectLiteral(i interface{}) (Object, error) {
	switch ii := i.(type) {
	case string:
		return StringLiteral(ii), nil
	case bool:
		return BooleanLiteral(ii), nil
	case int:
		return IntegerLiteral(ii), nil
	case int64:
		return IntegerLiteral(int(ii)), nil
	case time.Time:
		return DateTimeLiteral(ii), nil
	case *time.Time:
		return DateTimeLiteral(*ii), nil
	default:
		return nil, fmt.Errorf("unsupported literal type %T", i)
	}
}

func BooleanLiteral(bl bool) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdBoolean, val: fmt.Sprint(bl)},
	}
}

func (b *tripleBuilder) BooleanLiteral(bl bool) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  BooleanLiteral(bl).(object),
	}
}

func IntegerLiteral(i int) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdInteger, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) IntegerLiteral(i int) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  IntegerLiteral(i).(object),
	}
}

func StringLiteral(s string) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdString, val: s},
	}
}

func (b *tripleBuilder) StringLiteral(s string) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  StringLiteral(s).(object),
	}
}

func DateTimeLiteral(tm time.Time) Object {
	text, err := tm.UTC().MarshalText()
	if err != nil {
		panic(fmt.Errorf("date time literal: %s", err))
	}

	return object{
		isLit: true,
		lit:   literal{typ: XsdDateTime, val: string(text)},
	}
}

func (b *tripleBuilder) DateTimeLiteral(tm time.Time) *triple {
	return &triple{
		sub:  subject(b.sub),
		pred: predicate(b.pred),
		obj:  DateTimeLiteral(tm).(object),
	}
}

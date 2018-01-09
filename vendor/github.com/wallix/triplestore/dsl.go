package triplestore

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func SubjPredRes(s, p, r string) *triple {
	return &triple{
		sub: s, pred: p,
		obj: Resource(r).(object),
	}
}

func BnodePredRes(s, p, r string) *triple {
	return &triple{
		isSubBnode: true,
		sub:        s,
		pred:       p,
		obj:        Resource(r).(object),
	}
}

func SubjPredBnode(s, p, r string) *triple {
	return &triple{
		sub:  s,
		pred: p,
		obj:  object{bnode: r, isBnode: true},
	}
}

func SubjPredLit(s, p string, l interface{}) (*triple, error) {
	o, err := ObjectLiteral(l)
	if err != nil {
		return nil, err
	}
	return &triple{
		sub:  s,
		pred: p,
		obj:  o.(object),
	}, nil
}

type tripleBuilder struct {
	sub, pred  string
	isSubBnode bool
	langtag    string
}

func SubjPred(s, p string) *tripleBuilder {
	return &tripleBuilder{sub: s, pred: p}
}

func BnodePred(s, p string) *tripleBuilder {
	return &tripleBuilder{sub: s, isSubBnode: true, pred: p}
}

func Resource(s string) Object {
	return object{resource: s}
}

func (b *tripleBuilder) Lang(l string) *tripleBuilder {
	b.langtag = l
	return b
}

func (b *tripleBuilder) Resource(s string) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Resource(s).(object),
	}
}

func (b *tripleBuilder) Object(o Object) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        o.(object),
	}
}

func (b *tripleBuilder) Bnode(s string) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        object{bnode: s, isBnode: true},
	}
}

type UnsupportedLiteralTypeError struct {
	i interface{}
}

func (e UnsupportedLiteralTypeError) Error() string {
	return fmt.Sprintf("unsupported literal type %T", e.i)
}

func ObjectLiteral(i interface{}) (Object, error) {
	switch ii := i.(type) {
	case string:
		return StringLiteral(ii), nil
	case bool:
		return BooleanLiteral(ii), nil
	case int:
		return IntegerLiteral(ii), nil
	case int64, int32:
		r := reflect.ValueOf(ii)
		return IntegerLiteral(int(r.Int())), nil
	case int16:
		return Int16Literal(ii), nil
	case int8:
		return Int8Literal(ii), nil
	case float32:
		return Float32Literal(ii), nil
	case float64:
		return Float64Literal(ii), nil
	case uint:
		return UintegerLiteral(uint(ii)), nil
	case uint64, uint32:
		r := reflect.ValueOf(ii)
		return UintegerLiteral(uint(r.Uint())), nil
	case uint16:
		return Uint16Literal(ii), nil
	case uint8:
		return Uint8Literal(ii), nil
	case time.Time:
		return DateTimeLiteral(ii), nil
	case *time.Time:
		return DateTimeLiteral(*ii), nil
	case fmt.Stringer:
		return StringLiteral(ii.String()), nil
	default:
		return nil, UnsupportedLiteralTypeError{i}
	}
}

func ParseLiteral(obj Object) (interface{}, error) {
	if lit, ok := obj.Literal(); ok {
		switch lit.Type() {
		case XsdBoolean:
			return ParseBoolean(obj)
		case XsdDateTime:
			return ParseDateTime(obj)
		case XsdInteger:
			return ParseInteger(obj)
		case XsdByte:
			return ParseInt8(obj)
		case XsdShort:
			return ParseInt16(obj)
		case XsdUinteger:
			return ParseUinteger(obj)
		case XsdUnsignedByte:
			return ParseUint8(obj)
		case XsdUnsignedShort:
			return ParseUint16(obj)
		case XsdDouble:
			return ParseFloat64(obj)
		case XsdFloat:
			return ParseFloat32(obj)
		case XsdString:
			return ParseString(obj)
		default:
			return nil, fmt.Errorf("unknown literal type: %s", lit.Type())
		}
	}
	return nil, errors.New("cannot parse literal: object is not literal")
}

func BooleanLiteral(bl bool) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdBoolean, val: fmt.Sprint(bl)},
	}
}

func (b *tripleBuilder) BooleanLiteral(bl bool) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        BooleanLiteral(bl).(object),
	}
}

func ParseBoolean(obj Object) (bool, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdBoolean {
			return false, fmt.Errorf("literal is not an %s but %s", XsdBoolean, lit.Type())
		}

		return strconv.ParseBool(lit.Value())
	}

	return false, fmt.Errorf("cannot parse %s: object is not literal", XsdBoolean)
}

func IntegerLiteral(i int) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdInteger, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) IntegerLiteral(i int) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        IntegerLiteral(i).(object),
	}
}

func ParseInteger(obj Object) (int, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdInteger {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdInteger, lit.Type())
		}

		return strconv.Atoi(lit.Value())
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdInteger)
}

func Int8Literal(i int8) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdByte, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Int8Literal(i int8) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Int8Literal(i).(object),
	}
}

func ParseInt8(obj Object) (int8, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdByte {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdByte, lit.Type())
		}

		num, err := strconv.ParseInt(lit.Value(), 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(num), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdByte)
}

func Int16Literal(i int16) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdShort, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Int16Literal(i int16) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Int16Literal(i).(object),
	}
}

func ParseInt16(obj Object) (int16, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdShort {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdShort, lit.Type())
		}

		num, err := strconv.ParseInt(lit.Value(), 10, 16)
		if err != nil {
			return 0, err
		}
		return int16(num), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdShort)
}

func UintegerLiteral(i uint) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdUinteger, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) UintegerLiteral(i uint) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        UintegerLiteral(i).(object),
	}
}

func ParseUinteger(obj Object) (uint, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdUinteger {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdUinteger, lit.Type())
		}

		num, err := strconv.ParseUint(lit.Value(), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(num), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdUinteger)
}

func Uint8Literal(i uint8) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdUnsignedByte, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Uint8(i uint8) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Uint8Literal(i).(object),
	}
}

func ParseUint8(obj Object) (uint8, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdUnsignedByte {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdUnsignedByte, lit.Type())
		}

		num, err := strconv.ParseUint(lit.Value(), 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(num), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdUnsignedByte)
}

func Uint16Literal(i uint16) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdUnsignedShort, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Uint16(i uint16) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Uint16Literal(i).(object),
	}
}

func ParseUint16(obj Object) (uint16, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdUnsignedShort {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdUnsignedShort, lit.Type())
		}

		num, err := strconv.ParseUint(lit.Value(), 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(num), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdUnsignedShort)
}

func Float64Literal(i float64) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdDouble, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Float64Literal(i float64) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Float64Literal(i).(object),
	}
}

func ParseFloat64(obj Object) (float64, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdDouble {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdDouble, lit.Type())
		}

		return strconv.ParseFloat(lit.Value(), 64)
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdDouble)
}

func Float32Literal(i float32) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdFloat, val: fmt.Sprint(i)},
	}
}

func (b *tripleBuilder) Float32Literal(i float32) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        Float32Literal(i).(object),
	}
}

func ParseFloat32(obj Object) (float32, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdFloat {
			return 0, fmt.Errorf("literal is not an %s but %s", XsdFloat, lit.Type())
		}

		conv, err := strconv.ParseFloat(lit.Value(), 32)
		if err != nil {
			return 0, err
		}
		return float32(conv), nil
	}

	return 0, fmt.Errorf("cannot parse %s: object is not literal", XsdFloat)
}

func StringLiteral(s string) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdString, val: s},
	}
}

func StringLiteralWithLang(s, l string) Object {
	return object{
		isLit: true,
		lit:   literal{typ: XsdString, val: s, langtag: l},
	}
}

func (b *tripleBuilder) StringLiteral(s string) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        StringLiteral(s).(object),
	}
}

func (b *tripleBuilder) StringLiteralWithLang(s, l string) *triple {
	return &triple{
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        StringLiteralWithLang(s, l).(object),
	}
}

func ParseString(obj Object) (string, error) {
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdString {
			return "", fmt.Errorf("literal is not a %s but %s", XsdString, lit.Type())
		}

		return lit.Value(), nil
	}

	return "", fmt.Errorf("cannot parse %s: object is not literal", XsdString)
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
		isSubBnode: b.isSubBnode,
		sub:        b.sub,
		pred:       b.pred,
		obj:        DateTimeLiteral(tm).(object),
	}
}

func ParseDateTime(obj Object) (time.Time, error) {
	var t time.Time
	if lit, ok := obj.Literal(); ok {
		if lit.Type() != XsdDateTime {
			return t, fmt.Errorf("literal is not an %s but %s", XsdDateTime, lit.Type())
		}

		err := t.UnmarshalText([]byte(lit.Value()))
		if err != nil {
			return t, err
		}
		return t, nil
	}

	return t, fmt.Errorf("cannot parse %s: object is not literal", XsdDateTime)
}

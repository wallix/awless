package triplestore

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const (
	predTag  = "predicate"
	bnodeTag = "bnode"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Convert a Struct or ptr to Struct into triples
// using field tags.
// For each struct's field a triple is created:
// - Subject: function first argument
// - Predicate: tag value
// - Literal: actual field value according to field's type
// Unsupported types are ignored
func TriplesFromStruct(sub string, i interface{}, bnodes ...bool) (out []Triple) {
	var isBnode bool
	if len(bnodes) > 0 {
		isBnode = bnodes[0]
	}
	val := reflect.ValueOf(i)

	var ok bool
	val, ok = getStructOrPtrToStruct(val)
	if !ok {
		return
	}

	st := val.Type()

	for i := 0; i < st.NumField(); i++ {
		field, fVal := st.Field(i), val.Field(i)
		if !fVal.CanInterface() {
			continue
		}

		intValue := reflect.ValueOf(fVal.Interface())
		if intValue.Kind() == reflect.Ptr && intValue.IsNil() {
			continue
		}

		pred := field.Tag.Get(predTag)
		if tri, ok := buildTripleFromVal(sub, pred, fVal, isBnode); ok {
			out = append(out, tri)
		}

		bnode, embedded := field.Tag.Lookup(bnodeTag)
		fVal, ok := getStructOrPtrToStruct(fVal)
		if embedded && ok {
			if bnode == "" {
				bnode = fmt.Sprintf("%x", rand.Uint32())
			}
			tris := TriplesFromStruct(bnode, fVal.Interface(), true)
			out = append(out, tris...)
			if embedPred, hasPred := field.Tag.Lookup(predTag); hasPred {
				out = append(out, SubjPred(sub, embedPred).Bnode(bnode))
			}
			continue
		}

		switch fVal.Kind() {
		case reflect.Slice:
			length := fVal.Len()
			for i := 0; i < length; i++ {
				sliceVal := fVal.Index(i)
				if tri, ok := buildTripleFromVal(sub, pred, sliceVal, isBnode); ok {
					out = append(out, tri)
				}
			}
		}

	}

	return
}

func buildTripleFromVal(sub, pred string, v reflect.Value, bnode bool) (Triple, bool) {
	if !v.CanInterface() {
		return nil, false
	}
	if pred == "" {
		return nil, false
	}
	objLit, err := ObjectLiteral(v.Interface())
	if err != nil {
		return nil, false
	}

	if bnode {
		return BnodePred(sub, pred).Object(objLit), true
	}
	return SubjPred(sub, pred).Object(objLit), true
}

func getStructOrPtrToStruct(v reflect.Value) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Struct:
		return v, true
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Struct {
			return v.Elem(), true
		}
	}

	return v, false
}

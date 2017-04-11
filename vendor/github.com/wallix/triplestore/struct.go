package triplestore

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const (
	predTag = "predicate"
	subTag  = "subject"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// Convert a Struct or ptr to Struct into triples
// using field tags.
// For each struct's field a triple is created:
// - Subject: function first argument
// - Predicate: tag value
// - Literal: actual field value according to field's type
// Unsupported types are ignored
func TriplesFromStruct(sub string, i interface{}) (out []Triple) {
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

		tag, embedded := field.Tag.Lookup(subTag)
		fVal, ok := getStructOrPtrToStruct(fVal)
		if ok && embedded {
			sub := tag
			if tag == "rand" {
				sub = fmt.Sprintf("%x", random.Uint32())
			}
			tris := TriplesFromStruct(sub, fVal.Interface())
			out = append(out, tris...)
			continue
		}

		pred := field.Tag.Get(predTag)

		switch fVal.Kind() {
		case reflect.Slice:
			length := fVal.Len()
			for i := 0; i < length; i++ {
				sliceVal := fVal.Index(i)
				if tri, ok := buildTripleFromVal(sub, pred, sliceVal); ok {
					out = append(out, tri)
				}
			}
		}

		if tri, ok := buildTripleFromVal(sub, pred, fVal); ok {
			out = append(out, tri)
		}
	}

	return
}

func buildTripleFromVal(sub, pred string, v reflect.Value) (Triple, bool) {
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

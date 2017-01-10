package aws

import (
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
)

type convertFn func(s interface{}) reflect.Value

func setField(s, i interface{}, fieldName string) {
	sval := reflect.ValueOf(s)
	ival := reflect.ValueOf(i)

	if !ival.IsValid() || !sval.IsValid() {
		return
	}

	if ival.Kind() != reflect.Ptr && ival.Kind() != reflect.Struct {
		panic("only support setting field on ptr to struct\n")
	}

	fieldVal := ival.Elem().FieldByName(fieldName)

	if fieldVal.Type() == sval.Type() {
		fieldVal.Set(sval)
		return
	}

	var stringptr *string
	var int64ptr *int64

	if fieldVal.Kind() == reflect.Ptr {
		switch fieldVal.Type() {
		case reflect.TypeOf(stringptr):
			fieldVal.Set(reflect.ValueOf(aws.String(s.(string))))
		case reflect.TypeOf(int64ptr):
			r, err := strconv.Atoi(s.(string))
			if err != nil {
				panic(err)
			}
			fieldVal.Set(reflect.ValueOf(aws.Int64(int64(r))))
		}
	}

	if fieldVal.Kind() == reflect.Slice {
		switch s.(type) {
		case string:
			slice := []*string{aws.String(s.(string))}
			fieldVal.Set(reflect.ValueOf(slice))
		case int64:
			slice := []*int64{aws.Int64(s.(int64))}
			fieldVal.Set(reflect.ValueOf(slice))
		case int:
			slice := []*int64{aws.Int64(int64(s.(int)))}
			fieldVal.Set(reflect.ValueOf(slice))
		}
	}
}

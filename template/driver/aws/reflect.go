package aws

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	var boolptr *bool
	var boolval *ec2.AttributeBooleanValue
	var stringval *ec2.AttributeValue

	if fieldVal.Kind() == reflect.Ptr {
		switch fieldVal.Type() {
		case reflect.TypeOf(stringptr):
			fieldVal.Set(reflect.ValueOf(aws.String(s.(string))))
		case reflect.TypeOf(boolval), reflect.TypeOf(boolptr):
			var b bool
			var err error
			switch ss := s.(type) {
			case string:
				b, err = strconv.ParseBool(ss)
				if err != nil {
					panic(err)
				}
			case bool:
				b = ss
			}
			if fieldVal.Type() == reflect.TypeOf(boolval) {
				boolval = &ec2.AttributeBooleanValue{Value: aws.Bool(b)}
				fieldVal.Set(reflect.ValueOf(boolval))
			}
			if fieldVal.Type() == reflect.TypeOf(boolptr) {
				fieldVal.Set(reflect.ValueOf(aws.Bool(b)))
			}
		case reflect.TypeOf(stringval):
			stringval = &ec2.AttributeValue{Value: aws.String(fmt.Sprint(s))}
			fieldVal.Set(reflect.ValueOf(stringval))
		case reflect.TypeOf(int64ptr):
			var r int64
			var err error
			switch s.(type) {
			case string:
				r, err = strconv.ParseInt(s.(string), 10, 64)
				if err != nil {
					panic(err)
				}
			case int:
				r = int64(s.(int))
			case int64:
				r = s.(int64)
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

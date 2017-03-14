/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	awsstr = iota
	awsint
	awsint64
	awsbool
	awsboolattribute
	awsstringattribute
	awsint64slice
	awsstringslice
	awsstringpointermap
	awsslicestruct
)

var (
	mapAttributeRegex = regexp.MustCompile(`(.+)\[(.+)\].*`)
	sliceStructRegex  = regexp.MustCompile(`(.+)\[0\](.*)`)
)

func setFieldWithType(v, i interface{}, fieldPath string, destType int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("set field %s for %T object: %s", fieldPath, i, e)
		}
	}()
	if v == nil || i == nil {
		return nil
	}
	switch destType {
	case awsstr:
		v = fmt.Sprint(v)
	case awsint64:
		v, err = castInt64(v)
		if err != nil {
			return
		}
	case awsint:
		v, err = castInt(v)
		if err != nil {
			return
		}
	case awsbool:
		v, err = castBool(v)
		if err != nil {
			return
		}
	case awsstringslice:
		switch vv := v.(type) {
		case string:
			v = []*string{&vv}
		case *string:
			v = []*string{vv}
		case []*string:
			v = vv
		case []string:
			v = aws.StringSlice(vv)
		default:
			str := fmt.Sprint(v)
			v = []*string{&str}
		}
	case awsint64slice:
		var awsint int64
		awsint, err = castInt64(v)
		if err != nil {
			return
		}
		v = []*int64{&awsint}
	case awsboolattribute:
		var b bool
		b, err = castBool(v)
		if err != nil {
			return
		}
		v = &ec2.AttributeBooleanValue{Value: &b}
	case awsstringattribute:
		str := fmt.Sprint(v)
		v = &ec2.AttributeValue{Value: &str}
	case awsstringpointermap:
		matches := mapAttributeRegex.FindStringSubmatch(fieldPath)
		if len(matches) < 2 {
			err = fmt.Errorf("set field awsstringmap: path %s does not start with mymap[key]", fieldPath)
			return
		}
		strcr := reflect.Indirect(reflect.ValueOf(i))
		if strcr.Kind() != reflect.Struct {
			err = fmt.Errorf("set field awsstringmap: %T is not a struct, but a %s", i, strcr.Kind())
			return
		}
		field := strcr.FieldByName(matches[1])
		if field.Kind() != reflect.Map {
			err = fmt.Errorf("set field awsstringmap: field %s is not a map, but a %s", matches[0], field.Kind())
			return
		}
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}
		str := fmt.Sprint(v)
		field.SetMapIndex(reflect.ValueOf(matches[2]), reflect.ValueOf(&str))
		return nil
	case awsslicestruct:
		matches := sliceStructRegex.FindStringSubmatch(fieldPath)
		if len(matches) < 2 {
			err = fmt.Errorf("set field awsslicestruct: path %s does not start with slice[0]", fieldPath)
			return
		}
		strcr := reflect.Indirect(reflect.ValueOf(i))
		if strcr.Kind() != reflect.Struct {
			err = fmt.Errorf("set field awsslicestruct: %T is not a struct, but a %s", i, strcr.Kind())
			return
		}
		sliceField := strcr.FieldByName(matches[1])
		if sliceField.Kind() != reflect.Slice {
			err = fmt.Errorf("set field awsslicestruct: field %s is not a slice, but a %s", matches[0], sliceField.Kind())
			return
		}
		var elemToSet reflect.Value
		if sliceField.Len() > 0 {
			elemToSet = sliceField.Index(0)
		} else {
			elemToSet = reflect.New(sliceField.Type().Elem().Elem())
			sliceField.Set(reflect.Append(sliceField, elemToSet))
		}
		if sliceField.Type().Elem().Kind() != reflect.Ptr {
			err = fmt.Errorf("set field awsslicestruct: field %s is not a slice of struct pointer, but a %s", matches[0], sliceField.Kind())
			return
		}
		awsutil.SetValueAtPath(elemToSet.Interface(), matches[2], v)

		return nil
	}
	awsutil.SetValueAtPath(i, fieldPath, v)
	return nil
}

func castInt(v interface{}) (int, error) {
	switch vv := v.(type) {
	case string:
		return strconv.Atoi(vv)
	case int:
		return vv, nil
	case int64:
		return int(vv), nil
	default:
		return 0, fmt.Errorf("cannot cast %T to int", v)
	}
}

func castBool(v interface{}) (bool, error) {
	switch vv := v.(type) {
	case string:
		return strconv.ParseBool(vv)
	case bool:
		return vv, nil
	default:
		return false, fmt.Errorf("cannot cast %T to bool", v)
	}
}

func castInt64(v interface{}) (int64, error) {
	switch vv := v.(type) {
	case string:
		in, err := strconv.Atoi(vv)
		return int64(in), err
	case int:
		return int64(vv), nil
	case int64:
		return vv, nil
	default:
		return int64(0), fmt.Errorf("cannot cast %T to int64", v)
	}
}

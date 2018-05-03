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

package awsspec

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	gotemplate "text/template"
	"time"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/wallix/awless/logger"
)

const (
	awsstr                   = "awsstr"
	awsint                   = "awsint"
	awsint64                 = "awsint64"
	awsfloat                 = "awsfloat"
	awsbool                  = "awsbool"
	awsboolattribute         = "awsboolattribute"
	awsstringattribute       = "awsstringattribute"
	awsint64slice            = "awsint64slice"
	awsstringslice           = "awsstringslice"
	awsstringpointermap      = "awsstringpointermap"
	awsslicestruct           = "awsslicestruct"
	awsslicestructint64      = "awsslicestructint64"
	awsuserdatatobase64      = "awsuserdatatobase64"
	awsfiletobyteslice       = "awsfiletobyteslice"
	awsfiletostring          = "awsfiletostring"
	awsdimensionslice        = "awsdimensionslice"
	awsparameterslice        = "awsparameterslice"
	awsecskeyvalue           = "awsecskeyvalue"
	awsportmappings          = "awsportmappings"
	awssubnetmappings        = "awssubnetmappings"
	awsclassicloadblisteners = "awsclassicloadblisteners"
	awsstepadjustments       = "awsstepadjustments"
	awscsvstr                = "awscsvstr"
	aws6digitsstring         = "aws6digitsstring"
	awsbyteslice             = "awsbyteslice"
	awstagslice              = "awstagslice"
	awsalarmrollbacktriggers = "awsalarmrollbacktriggers"
)

var (
	mapAttributeRegex = regexp.MustCompile(`(.+)\[(.+)\].*`)
	sliceStructRegex  = regexp.MustCompile(`(.+)\[0\](.*)`)
)

type setter struct {
	val       interface{}
	fieldPath string
	fieldType string
}

func (s setter) set(i interface{}) error {
	return setFieldWithType(s.val, i, s.fieldPath, s.fieldType)
}

func setFieldWithType(v, i interface{}, fieldPath string, destType string, interfs ...interface{}) (err error) {
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
		v = castString(v)
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
	case aws6digitsstring:
		v, err = castInt(v)
		if err != nil {
			return
		}
		v = fmt.Sprintf("%06d", v)
	case awsfloat:
		v, err = castFloat(v)
		if err != nil {
			return
		}
	case awsbool:
		v, err = castBool(v)
		if err != nil {
			return
		}
	case awsstringslice:
		v = castStringPointerSlice(v)
	case awsbyteslice:
	case awscsvstr:
		v = strings.Join(castStringSlice(v), ",")
	case awsdimensionslice:
		if dimensions, isDim := v.([]*cloudwatch.Dimension); isDim {
			v = dimensions
		} else {
			dimensions = []*cloudwatch.Dimension{}
			sl := castStringSlice(v)
			for _, s := range sl {
				splits := strings.SplitN(s, ":", 2)
				if len(splits) != 2 {
					return fmt.Errorf("invalid dimension '%s', expected 'key:value'", s)
				}
				dimensions = append(dimensions, &cloudwatch.Dimension{Name: aws.String(splits[0]), Value: aws.String(splits[1])})
				v = dimensions
			}
		}
	case awsecskeyvalue:
		sl := castStringSlice(v)
		var keyvalues []*ecs.KeyValuePair
		for _, s := range sl {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid keyvalue '%s', expected 'key:value'", s)
			}
			keyvalues = append(keyvalues, &ecs.KeyValuePair{Name: aws.String(splits[0]), Value: aws.String(splits[1])})
		}
		v = keyvalues
	case awsparameterslice:
		sl := castStringSlice(v)
		var parameters []*cloudformation.Parameter
		for _, s := range sl {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid parameter '%s', expected 'key:value'", s)
			}
			parameters = append(parameters, &cloudformation.Parameter{ParameterKey: aws.String(splits[0]), ParameterValue: aws.String(splits[1])})
		}
		v = parameters
	case awssubnetmappings:
		sl := castStringSlice(v)
		var subnetMappings []*elbv2.SubnetMapping
		for i, s := range sl {
			splits := strings.Split(s, ":")
			if len(splits) != 2 {
				return fmt.Errorf("invalid element %d in subnet mapping %v, expect format [subnet-123:eipalloc-321, subnet-234:eipalloc-678, ...]", i+1, splits)
			}
			subnetMappings = append(subnetMappings, &elbv2.SubnetMapping{SubnetId: aws.String(splits[0]), AllocationId: aws.String(splits[1])})
		}
		v = subnetMappings
	case awsclassicloadblisteners:
		var listeners []*elb.Listener
		for _, s := range castStringSlice(v) {
			splits := strings.Split(s, ":")
			if len(splits) != 4 {
				return fmt.Errorf("missing value in listeners param '%s', expect format like HTTP:80:HTTP:80", splits)
			}
			loadbPort, err := strconv.ParseInt(splits[1], 10, 64)
			if err != nil {
				return fmt.Errorf("expecting numerical port value for loadbalancer port in '%s', (expect format like HTTP:80:HTTP:80)", splits)
			}
			instancePort, err := strconv.ParseInt(splits[3], 10, 64)
			if err != nil {
				return fmt.Errorf("expecting numerical port value for instance port in '%s', (expect format like HTTP:80:HTTP:80)", splits)
			}
			listeners = append(listeners, &elb.Listener{
				Protocol:         aws.String(splits[0]),
				LoadBalancerPort: aws.Int64(loadbPort),
				InstanceProtocol: aws.String(splits[2]),
				InstancePort:     aws.Int64(instancePort),
			})
		}
		v = listeners
	case awsportmappings:
		sl := castStringSlice(v)
		var portMappings []*ecs.PortMapping
		for _, s := range sl {
			portMapping := &ecs.PortMapping{}
			if strings.Contains(s, "-") {
				return fmt.Errorf("invalid port mapping '%s', AWS do not support portrange (from-to)", s)
			}
			var protocol string
			if strings.Contains(s, "/") {
				splits := strings.Split(s, "/")
				protocol = splits[1]
				if protocol != "tcp" && protocol != "udp" {
					return fmt.Errorf("invalid port mapping '%s', invalid protocol, expect tcp or udp, got %s", s, protocol)
				}
				s = strings.TrimRight(s, "/"+protocol)
				portMapping.Protocol = aws.String(protocol)
			}
			splits := strings.Split(s, ":")
			switch len(splits) {
			case 1:
				containerPort, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				portMapping.ContainerPort = aws.Int64(containerPort)
			case 2:
				hostPort, err := strconv.ParseInt(splits[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				containerPort, err := strconv.ParseInt(splits[1], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
				}
				portMapping.HostPort = aws.Int64(hostPort)
				portMapping.ContainerPort = aws.Int64(containerPort)
			default:
				return fmt.Errorf("invalid port mapping '%s', expect from[:to][/protocol]", s)
			}

			portMappings = append(portMappings, portMapping)
		}
		v = portMappings
	case awsstepadjustments:
		sl := castStringSlice(v)
		var stepAdjustments []*applicationautoscaling.StepAdjustment
		for _, s := range sl {
			splits := strings.Split(s, ":")
			if len(splits) != 3 {
				return fmt.Errorf("invalid step adjustment '%s', expect from:to:scaling-adjustment", s)
			}
			stepAdjustment := &applicationautoscaling.StepAdjustment{}
			if splits[0] != "" {
				lower, err := strconv.ParseFloat(splits[0], 64)
				if err != nil {
					return fmt.Errorf("invalid from '%s' in step adjustment '%s', expect from:to:scaling-adjustment", splits[0], s)
				}
				stepAdjustment.MetricIntervalLowerBound = aws.Float64(lower)
			}
			if splits[1] != "" {
				upper, err := strconv.ParseFloat(splits[1], 64)
				if err != nil {
					return fmt.Errorf("invalid to '%s' in step adjustment '%s', expect from:to:scaling-adjustment", splits[1], s)
				}
				stepAdjustment.MetricIntervalUpperBound = aws.Float64(upper)
			}
			adjustment, err := strconv.ParseInt(splits[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid adjustment-adjustment '%s' in step adjustmentstep adjustment '%s', expect from:to:scaling-adjustment", splits[2], s)
			}
			stepAdjustment.ScalingAdjustment = aws.Int64(adjustment)
			stepAdjustments = append(stepAdjustments, stepAdjustment)
		}
		v = stepAdjustments
	case awsuserdatatobase64:
		var tplData interface{}
		if len(interfs) > 0 {
			tplData = interfs[0]
		}
		v, err = userDataContentAsBase64(v, tplData)
		if err != nil {
			return err
		}
	case awsfiletobyteslice:
		v, err = ioutil.ReadFile(castString(v))
		if err != nil {
			return err
		}
	case awsfiletostring:
		var b []byte
		b, err = ioutil.ReadFile(castString(v))
		if err != nil {
			return err
		}
		v = string(b)
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
		str := castString(v)
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
		str := castString(v)
		field.SetMapIndex(reflect.ValueOf(matches[2]), reflect.ValueOf(&str))
		return nil
	case awsslicestruct, awsslicestructint64:
		if destType == awsslicestructint64 {
			v, err = castInt64(v)
			if err != nil {
				return
			}
		}
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
	case awstagslice:
		var (
			elbTags    []*elb.Tag
			cfTags     []*cloudformation.Tag
			appendFunc func(s1, s2 string)
			assignFunc func()
		)
		switch i.(type) {
		case *elb.CreateLoadBalancerInput:
			appendFunc = func(s1, s2 string) {
				elbTags = append(elbTags, &elb.Tag{Key: aws.String(s1), Value: aws.String(s2)})
			}
			assignFunc = func() { v = elbTags }
		case *cloudformation.CreateStackInput, *cloudformation.UpdateStackInput:
			appendFunc = func(s1, s2 string) {
				cfTags = append(cfTags, &cloudformation.Tag{Key: aws.String(s1), Value: aws.String(s2)})
			}
			assignFunc = func() { v = cfTags }
		}
		for _, s := range castStringSlice(v) {
			splits := strings.SplitN(s, ":", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid tag '%s', expected 'key:value'", s)
			}
			appendFunc(splits[0], splits[1])
		}
		assignFunc()
	case awsalarmrollbacktriggers:
		var triggers []*cloudformation.RollbackTrigger
		if list := castStringSlice(v); len(list) > 0 {
			for _, t := range list {
				triggers = append(triggers, &cloudformation.RollbackTrigger{
					Arn:  aws.String(t),
					Type: aws.String("AWS::CloudWatch::Alarm"),
				})
			}
		}
		v = triggers
	}

	awsutil.SetValueAtPath(i, fieldPath, v)
	return nil
}

func castString(v interface{}) string {
	switch vv := v.(type) {
	case []string:
		return strings.Join(vv, ",")
	case *string:
		return *vv
	default:
		return fmt.Sprint(v)
	}
}

func castFloat(v interface{}) (float64, error) {
	switch vv := v.(type) {
	case string:
		f, err := strconv.ParseFloat(vv, 64)
		if err != nil {
			return f, fmt.Errorf("invalid float value '%s'", vv)
		}
		return f, nil
	case float32:
		return float64(vv), nil
	case float64:
		return vv, nil
	case *float64:
		return aws.Float64Value(vv), nil
	case int:
		return float64(vv), nil
	case int64:
		return float64(vv), nil
	default:
		return 0, fmt.Errorf("cannot cast %T to float64", v)
	}
}

func castInt(v interface{}) (int, error) {
	switch vv := v.(type) {
	case *string:
		i, err := strconv.Atoi(aws.StringValue(vv))
		if err != nil {
			return i, fmt.Errorf("invalid integer value '%s'", aws.StringValue(vv))
		}
		return i, nil
	case string:
		i, err := strconv.Atoi(vv)
		if err != nil {
			return i, fmt.Errorf("invalid integer value '%s'", vv)
		}
		return i, nil
	case *int:
		return aws.IntValue(vv), nil
	case int:
		return vv, nil
	case int64:
		return int(vv), nil
	case *int64:
		return int(aws.Int64Value(vv)), nil
	default:
		return 0, fmt.Errorf("cannot cast %T to int", v)
	}
}

func castBool(v interface{}) (bool, error) {
	switch vv := v.(type) {
	case string:
		b, err := strconv.ParseBool(vv)
		if err != nil {
			return b, fmt.Errorf("invalid integer value '%s'", vv)
		}
		return b, nil
	case bool:
		return vv, nil
	case *bool:
		return aws.BoolValue(vv), nil
	default:
		return false, fmt.Errorf("cannot cast %T to bool", v)
	}
}

func castInt64(v interface{}) (int64, error) {
	switch vv := v.(type) {
	case string:
		i, err := strconv.Atoi(vv)
		if err != nil {
			return int64(i), fmt.Errorf("invalid integer value '%s'", vv)
		}
		return int64(i), nil
	case int:
		return int64(vv), nil
	case *int:
		return int64(aws.IntValue(vv)), nil
	case int64:
		return vv, nil
	case *int64:
		return aws.Int64Value(vv), nil
	default:
		return int64(0), fmt.Errorf("cannot cast %T to int64", v)
	}
}

func castStringSlice(v interface{}) []string {
	switch vv := v.(type) {
	case string:
		return []string{vv}
	case *string:
		return []string{aws.StringValue(vv)}
	case []*string:
		return aws.StringValueSlice(vv)
	case []string:
		return vv
	case []interface{}:
		var slice []string
		for _, i := range vv {
			switch ii := i.(type) {
			case string:
				slice = append(slice, ii)
			case *string:
				slice = append(slice, *ii)
			default:
				slice = append(slice, fmt.Sprint(ii))
			}
		}
		return slice
	default:
		return []string{fmt.Sprint(v)}
	}
}

func castStringPointerSlice(v interface{}) []*string {
	switch vv := v.(type) {
	case string:
		return []*string{&vv}
	case *string:
		return []*string{vv}
	case []*string:
		return vv
	case []string:
		return aws.StringSlice(vv)
	case []interface{}:
		var slice []*string
		for _, i := range vv {
			switch ii := i.(type) {
			case string:
				slice = append(slice, &ii)
			case *string:
				slice = append(slice, ii)
			default:
				str := fmt.Sprint(ii)
				slice = append(slice, &str)
			}
		}
		return slice
	default:
		str := fmt.Sprint(v)
		return []*string{&str}
	}
}

func userDataContentAsBase64(v interface{}, tplData interface{}) (string, error) {
	userdata := castString(v)

	var readErr error
	var content []byte

	if strings.HasPrefix(strings.TrimSpace(userdata), "#") { // userdata are bash content or yml cloud script content (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts)
		r := strings.NewReplacer("\\a", "\a", "\\b", "\b", "\\f", "\f", "\\n", "\n", "\\t", "\t", "\\r", "\r", "\\v", "\v")
		content = []byte(r.Replace(userdata))
	} else if strings.HasPrefix(userdata, "http") {
		client := &http.Client{Timeout: 5 * time.Second}

		logger.ExtraVerbosef("fetching remote userdata at '%s'", userdata)
		resp, err := client.Get(userdata)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
			return "", fmt.Errorf("'%s' when fetching userdata at '%s'", resp.Status, userdata)
		}

		content, readErr = ioutil.ReadAll(resp.Body)
	} else {
		content, readErr = ioutil.ReadFile(userdata)
	}

	if readErr != nil {
		return "", fmt.Errorf("got userdata from '%s' but cannot read content: %s", userdata, readErr)
	}

	if tpl, err := gotemplate.New("userdata").Parse(string(content)); err != nil {
		logger.Warningf("cannot parse userdata as Go template: %s", err)
	} else {
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, tplData); err == nil {
			content = buf.Bytes()
		}
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func structSetter(s interface{}, params map[string]interface{}) error {
	if params == nil {
		return nil
	}
	val := reflect.ValueOf(s).Elem()
	stru := val.Type()

	for i := 0; i < stru.NumField(); i++ {
		field := stru.Field(i)
		tplName := field.Tag.Get("templateName")
		var fieldType string
		if v, ok := params[tplName]; ok {
			kind := field.Type.Kind()
			if kind == reflect.Ptr {
				switch field.Type.Elem().Kind() {
				case reflect.String:
					fieldType = awsstr
				case reflect.Int64:
					fieldType = awsint64
				case reflect.Bool:
					fieldType = awsbool
				case reflect.Float64:
					fieldType = awsfloat
				default:
					return fmt.Errorf("unknown type %s for parameter %s in struct setter", tplName, field.Type.String())
				}
			} else if kind == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr {
				switch field.Type.Elem().Elem().Kind() {
				case reflect.String:
					fieldType = awsstringslice
				case reflect.Int64:
					fieldType = awsint64slice
				default:
					return fmt.Errorf("unknown type in slice %s for parameter %s", field.Type.String(), tplName)
				}
			}
			if err := setFieldWithType(v, s, field.Name, fieldType); err != nil {
				return fmt.Errorf("%s: %s", tplName, err)
			}
		}
	}
	return nil
}

func structInjector(src, dest interface{}, ctx map[string]interface{}) error {
	val := reflect.ValueOf(src).Elem()
	stru := val.Type()

	for i := 0; i < stru.NumField(); i++ {
		field := stru.Field(i)
		if dstNames, ok := field.Tag.Lookup("awsName"); ok {
			splits := strings.Split(dstNames, ",")
			for _, destName := range splits {
				destName = strings.TrimSpace(destName)
				if dstType, tok := field.Tag.Lookup("awsType"); tok {
					fieldValue := val.Field(i)
					if fieldValue.IsValid() && fieldValue.Interface() != nil && !fieldValue.IsNil() {
						if err := setFieldWithType(fieldValue.Interface(), dest, destName, dstType, ctx); err != nil {
							fieldName := field.Name
							if tplName, ok := field.Tag.Lookup("templateName"); ok {
								fieldName = tplName
							}
							return fmt.Errorf("%s: %s", fieldName, err)
						}
					}
				}
			}
		}
	}
	return nil
}

func contains(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}

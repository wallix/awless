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
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func TestGoTemplatingInUserdata(t *testing.T) {
	text := []byte("file content {{ .name }}")
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	finfo, _ := f.Stat()
	err = ioutil.WriteFile(f.Name(), text, finfo.Mode().Perm())
	if err != nil {
		t.Fatal(f)
	}

	awsparams := &ec2.RunInstancesInput{}

	err = setFieldWithType(f.Name(), awsparams, "UserData", awsuserdatatobase64, map[string]string{"name": "johndoe"})
	if err != nil {
		t.Fatal(err)
	}
	expText := []byte("file content johndoe")
	if got, want := awssdk.StringValue(awsparams.UserData), base64.StdEncoding.EncodeToString(expText); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSetFieldWithTypeAWSFile(t *testing.T) {
	text := []byte("file content")
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	finfo, _ := f.Stat()
	err = ioutil.WriteFile(f.Name(), text, finfo.Mode().Perm())
	if err != nil {
		t.Fatal(f)
	}

	awsparams := &ec2.RunInstancesInput{}

	err = setFieldWithType(f.Name(), awsparams, "UserData", awsuserdatatobase64)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := awssdk.StringValue(awsparams.UserData), base64.StdEncoding.EncodeToString(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	functionInput := &lambda.CreateFunctionInput{}

	err = setFieldWithType(f.Name(), functionInput, "Code.ZipFile", awsfiletobyteslice)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(functionInput.Code.ZipFile), string(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	stackInput := &cloudformation.CreateStackInput{}

	err = setFieldWithType(f.Name(), stackInput, "TemplateBody", awsfiletostring)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := awssdk.StringValue(stackInput.TemplateBody), string(text); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSetFieldsOnAwsStruct(t *testing.T) {
	awsparams := &ec2.RunInstancesInput{}

	err := setFieldWithType("ami", awsparams, "ImageId", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType("t2.micro", awsparams, "InstanceType", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType("5", awsparams, "MaxCount", awsint64)
	if err != nil {
		t.Fatal(err)
	}
	err = setFieldWithType(3, awsparams, "MinCount", awsint64)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := awssdk.StringValue(awsparams.ImageId), "ami"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := awssdk.StringValue(awsparams.InstanceType), "t2.micro"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := awssdk.Int64Value(awsparams.MaxCount), int64(5); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := awssdk.Int64Value(awsparams.MinCount), int64(3); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestSetFieldWithMultiType(t *testing.T) {
	any := struct {
		Field               string
		IntField            int
		FloatField          *float64
		BoolPointerField    *bool
		BoolField           bool
		StringArrayField    []*string
		Int64ArrayField     []*int64
		BooleanValueField   *ec2.AttributeBooleanValue
		StringValueField    *ec2.AttributeValue
		DimensionSliceField []*cloudwatch.Dimension
		KeyValueSliceField  []*ecs.KeyValuePair
		StructAttribute     struct {
			Str  *string
			Bool *bool
		}
		SliceStructPointerAttribute []*struct {
			Str1, Str2 *string
			Integer    *int64
		}
		MapAttribute          map[string]*string
		EmptyMapAttribute     map[string]*string
		ParameterList         []*cloudformation.Parameter
		PortMappings          []*ecs.PortMapping
		SubnetMappings        []*elbv2.SubnetMapping
		LoadBalancerListeners []*elb.Listener
		StepAdjustments       []*applicationautoscaling.StepAdjustment
		CSVString             *string
		SixDigitsString       *string
		ByteSlice             []byte
	}{Field: "initial", MapAttribute: map[string]*string{"test": awssdk.String("1234")}}

	err := setFieldWithType("expected", &any, "Field", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.Field, "expected"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(5, &any, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(42.21, &any, "FloatField", awsfloat)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.FloatField, 42.21; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("5", &any, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(nil, &any, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("first", &any, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.StringArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[0]), "first"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType([]string{"one", "two", "three"}, &any, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.StringArrayField), 3; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[0]), "one"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[1]), "two"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[2]), "three"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType([]interface{}{"four", "five"}, &any, "StringArrayField", awsstringslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.StringArrayField), 2; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[0]), "four"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := awssdk.StringValue(any.StringArrayField[1]), "five"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(int64(321), &any, "Int64ArrayField", awsint64slice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := awssdk.Int64Value(any.Int64ArrayField[0]), int64(321); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType(567, &any, "Int64ArrayField", awsint64slice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := awssdk.Int64Value(any.Int64ArrayField[0]), int64(567); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("any", nil, "IntField", awsint)
	if err != nil {
		t.Fatal(err)
	}

	err = setFieldWithType("true", &any, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := awssdk.BoolValue(any.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	err = setFieldWithType(nil, &any, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := awssdk.BoolValue(any.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	err = setFieldWithType(false, &any, "BooleanValueField", awsboolattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := awssdk.BoolValue(any.BooleanValueField.Value), false; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}

	err = setFieldWithType("true", &any, "BooleanValueField", awsbool)
	if err == nil {
		t.Fatalf("expected error got nil")
	}
	if got, want := err.Error(), "reflect.Set: value of type bool is not assignable to type ec2.AttributeBooleanValue"; !strings.HasSuffix(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType("abcd", &any, "StringValueField", awsstringattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := awssdk.StringValue(any.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}
	err = setFieldWithType(nil, &any, "StringValueField", awsstringattribute)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := awssdk.StringValue(any.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}

	err = setFieldWithType(true, &any, "BoolField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.BoolField, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(false, &any, "BoolField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.BoolField, false; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("true", &any, "BoolPointerField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.BoolPointerField, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType(false, &any, "BoolPointerField", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.BoolPointerField, false; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("fieldValue", &any, "StructAttribute.Str", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.StructAttribute.Str, "fieldValue"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType([]string{"one", "two", "three"}, &any, "StructAttribute.Str", awsstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.StructAttribute.Str, "one,two,three"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType("true", &any, "StructAttribute.Bool", awsbool)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.StructAttribute.Bool, true; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	err = setFieldWithType("abc", &any, "MapAttribute[Field1]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.MapAttribute), 1+1; got != want { //First "test" key + Field1
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.MapAttribute["Field1"], "abc"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType("def", &any, "MapAttribute[Field2]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.MapAttribute), 1+2; got != want { //First "test" key + Field1 and Field2
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.MapAttribute["Field1"], "abc"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := *any.MapAttribute["Field2"], "def"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType("abcd", &any, "EmptyMapAttribute[Field1]", awsstringpointermap)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.EmptyMapAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.EmptyMapAttribute["Field1"], "abcd"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	err = setFieldWithType("tata", &any, "SliceStructPointerAttribute[0]Str1", awsslicestruct)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.SliceStructPointerAttribute[0].Str1, "tata"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType("toto", &any, "SliceStructPointerAttribute[0]Str2", awsslicestruct)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.SliceStructPointerAttribute[0].Str2, "toto"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType(10, &any, "SliceStructPointerAttribute[0]Integer", awsslicestructint64)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.SliceStructPointerAttribute), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.SliceStructPointerAttribute[0].Integer, int64(10); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	err = setFieldWithType("key:value", &any, "DimensionSliceField", awsdimensionslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.DimensionSliceField), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.DimensionSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.DimensionSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType([]string{"key:value", "key1:value1:with:"}, &any, "DimensionSliceField", awsdimensionslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.DimensionSliceField), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.DimensionSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.DimensionSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.DimensionSliceField[1].Name, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.DimensionSliceField[1].Value, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType([]string{"key:value", "key1:value1:with:"}, &any, "ParameterList", awsparameterslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.ParameterList), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.ParameterList[0].ParameterKey, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.ParameterList[0].ParameterValue, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.ParameterList[1].ParameterKey, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.ParameterList[1].ParameterValue, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType([]string{"key:value", "key1:value1:with:"}, &any, "KeyValueSliceField", awsecskeyvalue)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.KeyValueSliceField), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.KeyValueSliceField[0].Name, "key"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.KeyValueSliceField[0].Value, "value"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.KeyValueSliceField[1].Name, "key1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.KeyValueSliceField[1].Value, "value1:with:"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType([]string{"80:8080", "8082", "1234:8083/udp"}, &any, "PortMappings", awsportmappings)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.PortMappings), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[0].HostPort, int64(80); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[0].ContainerPort, int64(8080); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[1].ContainerPort, int64(8082); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[2].HostPort, int64(1234); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[2].ContainerPort, int64(8083); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.PortMappings[2].Protocol, "udp"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType([]string{"subnet-123:eipalloc-123", "subnet-456:eipalloc-456"}, &any, "SubnetMappings", awssubnetmappings)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.SubnetMappings), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.SubnetMappings[0].SubnetId, "subnet-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.SubnetMappings[0].AllocationId, "eipalloc-123"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.SubnetMappings[1].SubnetId, "subnet-456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.SubnetMappings[1].AllocationId, "eipalloc-456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	err = setFieldWithType([]string{"HTTP:80:UDP:8080", "HTTPS:443:TCP:12345"}, &any, "LoadBalancerListeners", awsclassicloadblisteners)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.LoadBalancerListeners), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.LoadBalancerListeners[0].Protocol, "HTTP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.LoadBalancerListeners[0].LoadBalancerPort, int64(80); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.LoadBalancerListeners[0].InstanceProtocol, "UDP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.LoadBalancerListeners[0].InstancePort, int64(8080); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.LoadBalancerListeners[1].Protocol, "HTTPS"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.LoadBalancerListeners[1].LoadBalancerPort, int64(443); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.LoadBalancerListeners[1].InstanceProtocol, "TCP"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := *any.LoadBalancerListeners[1].InstancePort, int64(12345); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	err = setFieldWithType([]string{"0:0.25:-1", "0.75:1:+1"}, &any, "StepAdjustments", awsstepadjustments)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(any.StepAdjustments), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.StepAdjustments[0].MetricIntervalLowerBound, float64(0); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *any.StepAdjustments[0].MetricIntervalUpperBound, float64(0.25); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *any.StepAdjustments[0].ScalingAdjustment, int64(-1); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := *any.StepAdjustments[1].MetricIntervalLowerBound, float64(0.75); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *any.StepAdjustments[1].MetricIntervalUpperBound, float64(1); got != want {
		t.Fatalf("got %f, want %f", got, want)
	}
	if got, want := *any.StepAdjustments[1].ScalingAdjustment, int64(+1); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	err = setFieldWithType([]interface{}{"abcdef", "ghijk"}, &any, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType([]string{"abcdef", "ghijk"}, &any, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType("abcdef", &any, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.CSVString, "abcdef"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType("abcdef,ghijk", &any, "CSVString", awscsvstr)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.CSVString, "abcdef,ghijk"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType("123456", &any, "SixDigitsString", aws6digitsstring)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.SixDigitsString, "123456"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType("2345", &any, "SixDigitsString", aws6digitsstring)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *any.SixDigitsString, "002345"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	err = setFieldWithType([]byte("hello"), &any, "ByteSlice", awsbyteslice)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := any.ByteSlice, []byte("hello"); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

type TestStruct struct {
	FieldStringRequired *string   `awsName:"CloudStringRequired" awsType:"awsstr" templateName:"fstringrequired"`
	FieldString         *string   `awsName:"CloudString" awsType:"awsstr" templateName:"fstring"`
	FieldInt64          *int64    `awsName:"CloudInt64" awsType:"awsint64" templateName:"fint"`
	FieldBool           *bool     `awsName:"Embedded.CloudBool" awsType:"awsbool" templateName:"fbool"`
	FieldStringSlice    []*string `awsName:"CloudStringSlice" awsType:"awsstringslice" templateName:"fstrslice"`
	NilField            *string   `awsName:"CloudNilString" awsType:"awsstr" templateName:"fnilstring"`
	MultiCloudField     *int64    `awsName:"CloudField1,CloudField2" awsType:"awsint64" templateName:"fmultistring"`
}

func (ts *TestStruct) Validate_FieldStringRequired() (err error) {
	if len(*ts.FieldStringRequired) == 0 {
		err = errors.New("fstringrequired should not be empty")
	}
	return
}

func (ts *TestStruct) Validate_FieldString() (err error) {
	if len(*ts.FieldString) != 10 {
		err = errors.New("fstring should be 10 chars")
	}
	return
}

func (ts *TestStruct) Validate_FieldInt64() (err error) {
	if *ts.FieldInt64 > 10 {
		err = errors.New("fint should not exceed 10")
	}
	return
}

func TestStructDynamicSetter(t *testing.T) {
	params := map[string]interface{}{
		"fstringrequired": "jdoe",
		"fint":            "345",
		"fbool":           "true",
		"fstrslice":       []interface{}{"one", "two", 3},
	}

	in := &TestStruct{}
	err := structSetter(in, params)
	if err != nil {
		t.Fatal(err)
	}

	exp := &TestStruct{
		FieldStringRequired: awssdk.String("jdoe"),
		FieldInt64:          awssdk.Int64(345),
		FieldBool:           awssdk.Bool(true),
		FieldStringSlice:    awssdk.StringSlice([]string{"one", "two", "3"}),
	}

	if got, want := in, exp; !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot %#v\n\nwant %#v\n", got, want)
	}
}

func TestStructInjector(t *testing.T) {
	in := &TestStruct{
		FieldStringRequired: awssdk.String("jdoe"),
		FieldInt64:          awssdk.Int64(345),
		FieldBool:           awssdk.Bool(true),
		FieldStringSlice:    awssdk.StringSlice([]string{"one", "two", "3"}),
		MultiCloudField:     awssdk.Int64(12345),
	}

	type embStruct struct {
		CloudBool *bool
	}
	type outStruct struct {
		CloudStringRequired      *string
		CloudInt64               *int64
		Embedded                 *embStruct
		CloudStringSlice         []*string
		CloudField1, CloudField2 *int64
	}

	out := new(outStruct)

	err := structInjector(in, out, nil)
	if err != nil {
		t.Fatal(err)
	}

	exp := &outStruct{
		CloudStringRequired: awssdk.String("jdoe"),
		CloudInt64:          awssdk.Int64(345),
		Embedded:            &embStruct{CloudBool: awssdk.Bool(true)},
		CloudStringSlice:    awssdk.StringSlice([]string{"one", "two", "3"}),
		CloudField1:         awssdk.Int64(12345),
		CloudField2:         awssdk.Int64(12345),
	}

	if got, want := out, exp; !reflect.DeepEqual(got, want) {
		// pretty.Print(got)
		// pretty.Print(want)
		t.Fatalf("\ngot %#v\n\nwant %#v\n", got, want)
	}
}

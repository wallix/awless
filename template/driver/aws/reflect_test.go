package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSetFieldsOnAwsStruct(t *testing.T) {
	awsparams := &ec2.RunInstancesInput{}

	setField("ami", awsparams, "ImageId")
	setField("t2.micro", awsparams, "InstanceType")
	setField("5", awsparams, "MaxCount")
	setField(3, awsparams, "MinCount")

	if got, want := aws.StringValue(awsparams.ImageId), "ami"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := aws.StringValue(awsparams.InstanceType), "t2.micro"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := aws.Int64Value(awsparams.MaxCount), int64(5); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := aws.Int64Value(awsparams.MinCount), int64(3); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestSetFieldWithMultiType(t *testing.T) {
	any := struct {
		Field             string
		IntField          int
		StringArrayField  []*string
		Int64ArrayField   []*int64
		BooleanValueField *ec2.AttributeBooleanValue
		StringValueField  *ec2.AttributeValue
	}{Field: "initial"}

	setField("expected", &any, "Field")
	if got, want := any.Field, "expected"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField(5, &any, "IntField")
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField("5", &any, "IntField")
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField(nil, &any, "IntField")
	if got, want := any.IntField, 5; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField("first", &any, "StringArrayField")
	if got, want := len(any.StringArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.StringValue(any.StringArrayField[0]), "first"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField(int64(321), &any, "Int64ArrayField")
	if got, want := len(any.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.Int64Value(any.Int64ArrayField[0]), int64(321); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField(567, &any, "Int64ArrayField")
	if got, want := len(any.Int64ArrayField), 1; got != want {
		t.Fatalf("len: got %d, want %d", got, want)
	}
	if got, want := aws.Int64Value(any.Int64ArrayField[0]), int64(567); got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	setField("any", nil, "IntField")

	setField("true", &any, "BooleanValueField")
	if got, want := aws.BoolValue(any.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	setField(nil, &any, "BooleanValueField")
	if got, want := aws.BoolValue(any.BooleanValueField.Value), true; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}
	setField(false, &any, "BooleanValueField")
	if got, want := aws.BoolValue(any.BooleanValueField.Value), false; got != want {
		t.Fatalf("len: got %t, want %t", got, want)
	}

	setField("abcd", &any, "StringValueField")
	if got, want := aws.StringValue(any.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}
	setField(nil, &any, "StringValueField")
	if got, want := aws.StringValue(any.StringValueField.Value), "abcd"; got != want {
		t.Fatalf("len: got %s, want %s", got, want)
	}
}

func TestCanOnlySetFieldOnStructPtr(t *testing.T) {
	defer func() {
		if panicked := recover(); panicked == nil {
			t.Fatal("expected panic to occur")
		}
	}()

	setField("", struct{}{}, "")
}

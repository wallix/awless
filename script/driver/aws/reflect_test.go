package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSetFieldsOnAwsStruct(t *testing.T) {
	awsparams := &ec2.RunInstancesInput{
	//  ImageId:      aws.String(ami),
	//  MaxCount:     aws.Int64(1),
	//  MinCount:     aws.Int64(1),
	//  InstanceType: aws.String("t2.micro"),
	}

	setField("ami", awsparams, "ImageId")
	setField("t2.micro", awsparams, "InstanceType")
	setField("5", awsparams, "MaxCount")
	setField("3", awsparams, "MinCount")

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
		Field    string
		IntField int
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

	setField("any", nil, "IntField")
}

func TestCanOnlySetFieldOnStructPtr(t *testing.T) {
	defer func() {
		if panicked := recover(); panicked == nil {
			t.Fatal("expected panic to occur")
		}
	}()

	setField("", struct{}{}, "")
}

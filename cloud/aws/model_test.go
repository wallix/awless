package aws

import (
	"fmt"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestTransformFunctions(t *testing.T) {
	tag := []*ec2.Tag{
		{Key: awssdk.String("Name"), Value: awssdk.String("instance-name")},
		{Key: awssdk.String("Created with"), Value: awssdk.String("awless")},
	}

	val, _ := extractTagFn("Name")(tag)
	if got, want := fmt.Sprint(val), "instance-name"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	val, _ = extractTagFn("Created with")(tag)
	if got, want := fmt.Sprint(val), "awless"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	val, _ = extractValueFn(awssdk.String("any"))
	if got, want := val.(string), "any"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	val, _ = extractValueFn(awssdk.Int(2))
	if got, want := val.(int), 2; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	val, _ = extractValueFn(awssdk.Int64(4))
	if got, want := val.(int64), int64(4); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	data := &ec2.InstanceState{Code: awssdk.Int64(12), Name: awssdk.String("running")}

	val, _ = extractFieldFn("Code")(data)
	if got, want := val.(int64), int64(12); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	val, _ = extractFieldFn("Name")(data)
	if got, want := val.(string), "running"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

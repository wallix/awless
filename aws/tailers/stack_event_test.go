package awstailers

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type stackEventTest struct {
	*cloudformation.StackEvent
}

func TestReflect(t *testing.T) {
	a := stackEventTest{
		&cloudformation.StackEvent{
			EventId:        aws.String("sdasda"),
			ResourceStatus: aws.String("dsada"),
		},
	}
	v, _ := reflect.TypeOf(a).Elem().FieldByName("EventId")

	v.Tag = reflect.StructTag(`fixed:"12,22"`)

	t.Log(v.Tag.Get("fixed"))
	// t.Log(v.FieldByName("type"))
}

package awsfetch

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

func TestAssignAPIs(t *testing.T) {
	type iamMock struct {
		iamiface.IAMAPI
	}
	type ec2Mock struct {
		ec2iface.EC2API
	}
	type any struct {
	}

	conf := NewConfig(&iamMock{}, ec2Mock{}, nil, nil, any{}, new(any))
	if conf.APIs.Iam == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.Ec2 == nil {
		t.Fatal("unexpected nil")
	}
	if conf.APIs.Rds != nil {
		t.Fatal("expected nil")
	}
	if conf.APIs.Ecr != nil {
		t.Fatal("expected nil")
	}
}

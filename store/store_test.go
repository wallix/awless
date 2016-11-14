package store

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestBuildInfraTreeFromAwsReservations(t *testing.T) {
	instances := []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("inst_1"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_2"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_3"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_4"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_5"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_6"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_3")},
		&ec2.Instance{InstanceId: aws.String("inst_7"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_3")},
	}

	vpcs := []*ec2.Vpc{
		&ec2.Vpc{VpcId: aws.String("vpc_1")},
		&ec2.Vpc{VpcId: aws.String("vpc_2")},
		&ec2.Vpc{VpcId: aws.String("vpc_3")},
	}

	subnets := []*ec2.Subnet{
		&ec2.Subnet{SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_3")},
	}

	region := BuildRegionTree("us-trump", vpcs, subnets, instances)

	if got, want := region.Id, "us-trump"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.Vpcs[0].Id, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[0].Subnets), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.Vpcs[1].Id, "vpc_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[1].Subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.Vpcs[2].Id, "vpc_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[2].Subnets), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.Vpcs[0].Subnets[0].Id, "sub_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[0].VpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[0].Subnets[0].Instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[1].Id, "sub_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[1].VpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[0].Subnets[1].Instances), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.Vpcs[2].Subnets[0].Id, "sub_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[2].Subnets[0].VpcId, "vpc_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.Vpcs[2].Subnets[0].Instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.Vpcs[0].Subnets[0].Instances[0].Id, "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[0].Instances[1].Id, "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[1].Instances[0].Id, "inst_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[1].Instances[1].Id, "inst_4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[0].Subnets[1].Instances[2].Id, "inst_5"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[2].Subnets[0].Instances[0].Id, "inst_6"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.Vpcs[2].Subnets[0].Instances[1].Id, "inst_7"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

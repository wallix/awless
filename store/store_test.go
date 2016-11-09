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

	if got, want := region.id, "us-trump"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.vpcs[0].id, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[0].subnets), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.vpcs[1].id, "vpc_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[1].subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.vpcs[2].id, "vpc_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[2].subnets), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.vpcs[0].subnets[0].id, "sub_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[0].vpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[0].subnets[0].instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.vpcs[0].subnets[1].id, "sub_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[1].vpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[0].subnets[1].instances), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := region.vpcs[2].subnets[0].id, "sub_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[2].subnets[0].vpcId, "vpc_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := len(region.vpcs[2].subnets[0].instances), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := region.vpcs[0].subnets[0].instances[0].id, "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[0].instances[1].id, "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[1].instances[0].id, "inst_3"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[1].instances[1].id, "inst_4"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[0].subnets[1].instances[2].id, "inst_5"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[2].subnets[0].instances[0].id, "inst_6"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := region.vpcs[2].subnets[0].instances[1].id, "inst_7"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

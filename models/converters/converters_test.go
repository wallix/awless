package converters

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/models"
)

func TestInstanceConversion(t *testing.T) {
	awsInst := &ec2.Instance{
		InstanceId:       aws.String("inst_1"),
		SubnetId:         aws.String("sub_1"),
		VpcId:            aws.String("vpc_1"),
		PublicIpAddress:  aws.String("127.0.0.1"),
		PrivateIpAddress: aws.String("172.0.0.1"),
	}
	inst := ConvertModel(awsInst).(*models.Instance)

	if got, want := inst.Id, "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := inst.SubnetId, "sub_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := inst.VpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := inst.PublicIp, "127.0.0.1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := inst.PrivateIp, "172.0.0.1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestVpcConversion(t *testing.T) {
	awsVpc := &ec2.Vpc{
		VpcId: aws.String("vpc_1"),
	}
	vpc := ConvertModel(awsVpc).(*models.Vpc)

	if got, want := vpc.Id, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestSubnetConversion(t *testing.T) {
	awsSubnet := &ec2.Subnet{
		SubnetId: aws.String("subnet_1"),
		VpcId:    aws.String("vpc_1"),
	}
	subnet := ConvertModel(awsSubnet).(*models.Subnet)

	if got, want := subnet.Id, "subnet_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := subnet.VpcId, "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

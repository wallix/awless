package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (inf *Infra) Vpc(id string) (interface{}, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{awssdk.String(id)},
	}

	return inf.DescribeVpcs(input)
}

func (inf *Infra) CreateInstance(ami string) (interface{}, error) {
	params := &ec2.RunInstancesInput{
		ImageId:      awssdk.String(ami),
		MaxCount:     awssdk.Int64(1),
		MinCount:     awssdk.Int64(1),
		InstanceType: awssdk.String("t2.micro"),
	}

	return inf.RunInstances(params)
}

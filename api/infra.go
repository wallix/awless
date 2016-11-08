package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Infra struct {
	*ec2.EC2
}

func NewInfra(region string) (*Infra, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &Infra{ec2.New(sess, aws.NewConfig().WithRegion(region))}, nil
}

func (inf *Infra) Instances() (interface{}, error) {
	return inf.DescribeInstances(&ec2.DescribeInstancesInput{})
}

func (inf *Infra) Regions() (interface{}, error) {
	return inf.DescribeRegions(&ec2.DescribeRegionsInput{})
}

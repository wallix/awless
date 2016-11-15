package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Infra struct {
	*ec2.EC2
}

func NewInfra(sess *session.Session) *Infra {
	return &Infra{ec2.New(sess)}
}

func (inf *Infra) Instances() (interface{}, error) {
	return inf.DescribeInstances(&ec2.DescribeInstancesInput{})
}

func (inf *Infra) Vpcs() (interface{}, error) {
	return inf.DescribeVpcs(&ec2.DescribeVpcsInput{})
}

func (inf *Infra) Subnets() (interface{}, error) {
	return inf.DescribeSubnets(&ec2.DescribeSubnetsInput{})
}

func (inf *Infra) Regions() (interface{}, error) {
	return inf.DescribeRegions(&ec2.DescribeRegionsInput{})
}

func (inf *Infra) Vpc(id string) (interface{}, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(id)},
	}

	return inf.DescribeVpcs(input)
}

func (inf *Infra) FetchInfra() ([]*ec2.Vpc, []*ec2.Subnet, []*ec2.Instance, error) {
	var fetchErr error
	var instances []*ec2.Instance
	var vpcs []*ec2.Vpc
	var subnets []*ec2.Subnet

	type fetchFn func() (interface{}, error)

	allFetch := []fetchFn{inf.Instances, inf.Subnets, inf.Vpcs}
	resultc := make(chan interface{})
	errc := make(chan error)

	for _, fetch := range allFetch {
		go func(fn fetchFn) {
			if r, err := fn(); err != nil {
				errc <- err
			} else {
				resultc <- r
			}
		}(fetch)
	}

	for range allFetch {
		select {
		case r := <-resultc:
			switch r.(type) {
			case *ec2.DescribeVpcsOutput:
				vpcs = r.(*ec2.DescribeVpcsOutput).Vpcs
			case *ec2.DescribeSubnetsOutput:
				subnets = r.(*ec2.DescribeSubnetsOutput).Subnets
			case *ec2.DescribeInstancesOutput:
				instances = r.(*ec2.DescribeInstancesOutput).Reservations[0].Instances
			}
		case fetchErr = <-errc:
			//
		}
	}

	return vpcs, subnets, instances, fetchErr
}

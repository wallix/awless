// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Infra struct {
	ec2iface.EC2API
}

func NewInfra(sess *session.Session) *Infra {
	return &Infra{ec2.New(sess)}
}

func (s *Infra) Instances() (interface{}, error) {
	return s.DescribeInstances(&ec2.DescribeInstancesInput{})
}
func (s *Infra) Subnets() (interface{}, error) {
	return s.DescribeSubnets(&ec2.DescribeSubnetsInput{})
}
func (s *Infra) Vpcs() (interface{}, error) {
	return s.DescribeVpcs(&ec2.DescribeVpcsInput{})
}
func (s *Infra) Keypairs() (interface{}, error) {
	return s.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
}
func (s *Infra) Securitygroups() (interface{}, error) {
	return s.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
}
func (s *Infra) Volumes() (interface{}, error) {
	return s.DescribeVolumes(&ec2.DescribeVolumesInput{})
}
func (s *Infra) Regions() (interface{}, error) {
	return s.DescribeRegions(&ec2.DescribeRegionsInput{})
}
func (s *Infra) Internetgateways() (interface{}, error) {
	return s.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
}
func (s *Infra) Routetables() (interface{}, error) {
	return s.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
}
func (s *Infra) Images() (interface{}, error) {
	return s.DescribeImages(&ec2.DescribeImagesInput{})
}

type AwsInfra struct {
	Instances        []*ec2.Instance
	Subnets          []*ec2.Subnet
	Vpcs             []*ec2.Vpc
	Keypairs         []*ec2.KeyPairInfo
	Securitygroups   []*ec2.SecurityGroup
	Volumes          []*ec2.Volume
	Regions          []*ec2.Region
	Internetgateways []*ec2.InternetGateway
	Routetables      []*ec2.RouteTable
	Images           []*ec2.Image
}

func (s *Infra) FetchAwsInfra() (*AwsInfra, error) {
	resultc, errc := multiFetch(
		s.Instances,
		s.Subnets,
		s.Vpcs,
		s.Keypairs,
		s.Securitygroups,
		s.Volumes,
		s.Regions,
		s.Internetgateways,
		s.Routetables,
		s.Images,
	)

	awsService := &AwsInfra{}

	for r := range resultc {
		switch rr := r.(type) {
		case *ec2.DescribeInstancesOutput:
			for _, c := range rr.Reservations {
				awsService.Instances = append(awsService.Instances, c.Instances...)
			}
		case *ec2.DescribeSubnetsOutput:
			awsService.Subnets = append(awsService.Subnets, rr.Subnets...)
		case *ec2.DescribeVpcsOutput:
			awsService.Vpcs = append(awsService.Vpcs, rr.Vpcs...)
		case *ec2.DescribeKeyPairsOutput:
			awsService.Keypairs = append(awsService.Keypairs, rr.KeyPairs...)
		case *ec2.DescribeSecurityGroupsOutput:
			awsService.Securitygroups = append(awsService.Securitygroups, rr.SecurityGroups...)
		case *ec2.DescribeVolumesOutput:
			awsService.Volumes = append(awsService.Volumes, rr.Volumes...)
		case *ec2.DescribeRegionsOutput:
			awsService.Regions = append(awsService.Regions, rr.Regions...)
		case *ec2.DescribeInternetGatewaysOutput:
			awsService.Internetgateways = append(awsService.Internetgateways, rr.InternetGateways...)
		case *ec2.DescribeRouteTablesOutput:
			awsService.Routetables = append(awsService.Routetables, rr.RouteTables...)
		case *ec2.DescribeImagesOutput:
			awsService.Images = append(awsService.Images, rr.Images...)
		}
	}

	return awsService, <-errc
}

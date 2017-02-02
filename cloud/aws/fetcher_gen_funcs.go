// DO NOT EDIT
// This file was automatically generated with go generate
package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/graph"
)

type Infra struct {
	region string
	ec2iface.EC2API
}

func NewInfra(sess *session.Session) *Infra {
	region := awssdk.StringValue(sess.Config.Region)
	return &Infra{EC2API: ec2.New(sess), region: region}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Provider() string {
	return "aws"
}

func (s *Infra) ProviderRunnableAPI() interface{} {
	return s.EC2API
}

func (s *Infra) ResourceTypes() (all []string) {

	all = append(all, "instance")
	all = append(all, "subnet")
	all = append(all, "vpc")
	all = append(all, "keypair")
	all = append(all, "securitygroup")
	all = append(all, "volume")
	all = append(all, "region")
	all = append(all, "internetgateway")
	all = append(all, "routetable")
	all = append(all, "image")
	return
}

func (s *Infra) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "instance":
		return s.fetch_all_instance_graph()
	case "subnet":
		return s.fetch_all_subnet_graph()
	case "vpc":
		return s.fetch_all_vpc_graph()
	case "keypair":
		return s.fetch_all_keypair_graph()
	case "securitygroup":
		return s.fetch_all_securitygroup_graph()
	case "volume":
		return s.fetch_all_volume_graph()
	case "region":
		return s.fetch_all_region_graph()
	case "internetgateway":
		return s.fetch_all_internetgateway_graph()
	case "routetable":
		return s.fetch_all_routetable_graph()
	case "image":
		return s.fetch_all_image_graph()
	default:
		return nil, fmt.Errorf("aws infra: unsupported fetch for type %s", t)
	}
}

func (s *Infra) fetch_all_instance() (interface{}, error) {
	return s.DescribeInstances(&ec2.DescribeInstancesInput{})
}
func (s *Infra) fetch_all_subnet() (interface{}, error) {
	return s.DescribeSubnets(&ec2.DescribeSubnetsInput{})
}
func (s *Infra) fetch_all_vpc() (interface{}, error) {
	return s.DescribeVpcs(&ec2.DescribeVpcsInput{})
}
func (s *Infra) fetch_all_keypair() (interface{}, error) {
	return s.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
}
func (s *Infra) fetch_all_securitygroup() (interface{}, error) {
	return s.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
}
func (s *Infra) fetch_all_volume() (interface{}, error) {
	return s.DescribeVolumes(&ec2.DescribeVolumesInput{})
}
func (s *Infra) fetch_all_region() (interface{}, error) {
	return s.DescribeRegions(&ec2.DescribeRegionsInput{})
}
func (s *Infra) fetch_all_internetgateway() (interface{}, error) {
	return s.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
}
func (s *Infra) fetch_all_routetable() (interface{}, error) {
	return s.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
}
func (s *Infra) fetch_all_image() (interface{}, error) {
	return s.DescribeImages(&ec2.DescribeImagesInput{})
}

func (s *Infra) fetch_all_instance_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_instance()
	if err != nil {
		return nil, err
	}

	for _, all := range out.(*ec2.DescribeInstancesOutput).Reservations {
		for _, output := range all.Instances {
			res, err := NewResource(output)
			if err != nil {
				return g, err
			}
			g.AddResource(res)
		}
	}

	return g, nil
}

func (s *Infra) fetch_all_subnet_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_subnet()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeSubnetsOutput).Subnets {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_vpc_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_vpc()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeVpcsOutput).Vpcs {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_keypair_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_keypair()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeKeyPairsOutput).KeyPairs {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_securitygroup_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_securitygroup()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeSecurityGroupsOutput).SecurityGroups {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_volume_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_volume()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeVolumesOutput).Volumes {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_region_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_region()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeRegionsOutput).Regions {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_internetgateway_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_internetgateway()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeInternetGatewaysOutput).InternetGateways {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_routetable_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_routetable()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeRouteTablesOutput).RouteTables {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Infra) fetch_all_image_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_image()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*ec2.DescribeImagesOutput).Images {
		res, err := NewResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

type AwsInfra struct {
	instanceList        []*ec2.Instance
	subnetList          []*ec2.Subnet
	vpcList             []*ec2.Vpc
	keypairList         []*ec2.KeyPairInfo
	securitygroupList   []*ec2.SecurityGroup
	volumeList          []*ec2.Volume
	regionList          []*ec2.Region
	internetgatewayList []*ec2.InternetGateway
	routetableList      []*ec2.RouteTable
	imageList           []*ec2.Image
}

func (s *Infra) fetch_ec2() (*AwsInfra, error) {
	resultc, errc := multiFetch(
		s.fetch_all_instance,
		s.fetch_all_subnet,
		s.fetch_all_vpc,
		s.fetch_all_keypair,
		s.fetch_all_securitygroup,
		s.fetch_all_volume,
		s.fetch_all_region,
		s.fetch_all_internetgateway,
		s.fetch_all_routetable,
		s.fetch_all_image,
	)

	awsService := &AwsInfra{}

	for r := range resultc {
		switch rr := r.(type) {
		case *ec2.DescribeInstancesOutput:
			for _, c := range rr.Reservations {
				awsService.instanceList = append(awsService.instanceList, c.Instances...)
			}
		case *ec2.DescribeSubnetsOutput:
			awsService.subnetList = append(awsService.subnetList, rr.Subnets...)
		case *ec2.DescribeVpcsOutput:
			awsService.vpcList = append(awsService.vpcList, rr.Vpcs...)
		case *ec2.DescribeKeyPairsOutput:
			awsService.keypairList = append(awsService.keypairList, rr.KeyPairs...)
		case *ec2.DescribeSecurityGroupsOutput:
			awsService.securitygroupList = append(awsService.securitygroupList, rr.SecurityGroups...)
		case *ec2.DescribeVolumesOutput:
			awsService.volumeList = append(awsService.volumeList, rr.Volumes...)
		case *ec2.DescribeRegionsOutput:
			awsService.regionList = append(awsService.regionList, rr.Regions...)
		case *ec2.DescribeInternetGatewaysOutput:
			awsService.internetgatewayList = append(awsService.internetgatewayList, rr.InternetGateways...)
		case *ec2.DescribeRouteTablesOutput:
			awsService.routetableList = append(awsService.routetableList, rr.RouteTables...)
		case *ec2.DescribeImagesOutput:
			awsService.imageList = append(awsService.imageList, rr.Images...)
		}
	}

	return awsService, <-errc
}

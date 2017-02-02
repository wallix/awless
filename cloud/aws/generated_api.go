// Auto generated implementation for the AWS cloud service
package aws

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
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

func (s *Infra) fetch_all_instance_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_instance()
	if err != nil {
		return nil, err
	}

	for _, all := range out.(*ec2.DescribeInstancesOutput).Reservations {
		for _, output := range all.Instances {
			res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
		res, err := newResource(output)
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
}

func (s *Infra) global_fetch() (*AwsInfra, error) {
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
		}
	}

	return awsService, <-errc
}

type Access struct {
	region string
	iamiface.IAMAPI
}

func NewAccess(sess *session.Session) *Access {
	region := awssdk.StringValue(sess.Config.Region)
	return &Access{IAMAPI: iam.New(sess), region: region}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Provider() string {
	return "aws"
}

func (s *Access) ProviderRunnableAPI() interface{} {
	return s.IAMAPI
}

func (s *Access) ResourceTypes() (all []string) {
	all = append(all, "user")
	all = append(all, "group")
	all = append(all, "role")
	all = append(all, "policy")
	return
}

func (s *Access) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "user":
		return s.fetch_all_user_graph()
	case "group":
		return s.fetch_all_group_graph()
	case "role":
		return s.fetch_all_role_graph()
	case "policy":
		return s.fetch_all_policy_graph()
	default:
		return nil, fmt.Errorf("aws access: unsupported fetch for type %s", t)
	}
}

func (s *Access) fetch_all_user() (interface{}, error) {
	return s.ListUsers(&iam.ListUsersInput{})
}
func (s *Access) fetch_all_group() (interface{}, error) {
	return s.ListGroups(&iam.ListGroupsInput{})
}
func (s *Access) fetch_all_role() (interface{}, error) {
	return s.ListRoles(&iam.ListRolesInput{})
}
func (s *Access) fetch_all_policy() (interface{}, error) {
	return s.ListPolicies(&iam.ListPoliciesInput{Scope: awssdk.String(iam.PolicyScopeTypeLocal)})
}

func (s *Access) fetch_all_user_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_user()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*iam.ListUsersOutput).Users {
		res, err := newResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Access) fetch_all_group_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_group()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*iam.ListGroupsOutput).Groups {
		res, err := newResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Access) fetch_all_role_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_role()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*iam.ListRolesOutput).Roles {
		res, err := newResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

func (s *Access) fetch_all_policy_graph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_policy()
	if err != nil {
		return nil, err
	}

	for _, output := range out.(*iam.ListPoliciesOutput).Policies {
		res, err := newResource(output)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}

	return g, nil
}

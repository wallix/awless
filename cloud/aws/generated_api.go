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

func init() {
	ServiceNames = append(ServiceNames, "infra")
	ServiceNames = append(ServiceNames, "access")
}

var ServiceNames = []string{}

var ResourceTypesPerAPI = map[string][]string{
	"ec2": []string{
		"instance",
		"subnet",
		"vpc",
		"keypair",
		"securitygroup",
		"volume",
		"region",
		"internetgateway",
		"routetable",
	},
	"iam": []string{
		"user",
		"group",
		"role",
		"policy",
	},
}

var ServicePerResourceType = map[string]string{
	"instance":        "infra",
	"subnet":          "infra",
	"vpc":             "infra",
	"keypair":         "infra",
	"securitygroup":   "infra",
	"volume":          "infra",
	"region":          "infra",
	"internetgateway": "infra",
	"routetable":      "infra",
	"user":            "access",
	"group":           "access",
	"role":            "access",
	"policy":          "access",
}

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

func (s *Infra) ProviderAPI() string {
	return "ec2"
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
		graph, _, err := s.fetch_all_instance_graph()
		return graph, err
	case "subnet":
		graph, _, err := s.fetch_all_subnet_graph()
		return graph, err
	case "vpc":
		graph, _, err := s.fetch_all_vpc_graph()
		return graph, err
	case "keypair":
		graph, _, err := s.fetch_all_keypair_graph()
		return graph, err
	case "securitygroup":
		graph, _, err := s.fetch_all_securitygroup_graph()
		return graph, err
	case "volume":
		graph, _, err := s.fetch_all_volume_graph()
		return graph, err
	case "region":
		graph, _, err := s.fetch_all_region_graph()
		return graph, err
	case "internetgateway":
		graph, _, err := s.fetch_all_internetgateway_graph()
		return graph, err
	case "routetable":
		graph, _, err := s.fetch_all_routetable_graph()
		return graph, err
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

func (s *Infra) fetch_all_instance_graph() (*graph.Graph, []*ec2.Instance, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Instance
	out, err := s.fetch_all_instance()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, all := range out.(*ec2.DescribeInstancesOutput).Reservations {
		for _, output := range all.Instances {
			cloudResources = append(cloudResources, output)
			res, err := newResource(output)
			if err != nil {
				return g, cloudResources, err
			}
			g.AddResource(res)
		}
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_subnet_graph() (*graph.Graph, []*ec2.Subnet, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Subnet
	out, err := s.fetch_all_subnet()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeSubnetsOutput).Subnets {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_vpc_graph() (*graph.Graph, []*ec2.Vpc, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Vpc
	out, err := s.fetch_all_vpc()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeVpcsOutput).Vpcs {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_keypair_graph() (*graph.Graph, []*ec2.KeyPairInfo, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.KeyPairInfo
	out, err := s.fetch_all_keypair()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeKeyPairsOutput).KeyPairs {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_securitygroup_graph() (*graph.Graph, []*ec2.SecurityGroup, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.SecurityGroup
	out, err := s.fetch_all_securitygroup()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeSecurityGroupsOutput).SecurityGroups {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_volume_graph() (*graph.Graph, []*ec2.Volume, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Volume
	out, err := s.fetch_all_volume()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeVolumesOutput).Volumes {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_region_graph() (*graph.Graph, []*ec2.Region, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Region
	out, err := s.fetch_all_region()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeRegionsOutput).Regions {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_internetgateway_graph() (*graph.Graph, []*ec2.InternetGateway, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.InternetGateway
	out, err := s.fetch_all_internetgateway()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeInternetGatewaysOutput).InternetGateways {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_routetable_graph() (*graph.Graph, []*ec2.RouteTable, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.RouteTable
	out, err := s.fetch_all_routetable()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeRouteTablesOutput).RouteTables {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
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

func (s *Access) ProviderAPI() string {
	return "iam"
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
		graph, _, err := s.fetch_all_user_graph()
		return graph, err
	case "group":
		graph, _, err := s.fetch_all_group_graph()
		return graph, err
	case "role":
		graph, _, err := s.fetch_all_role_graph()
		return graph, err
	case "policy":
		graph, _, err := s.fetch_all_policy_graph()
		return graph, err
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

func (s *Access) fetch_all_user_graph() (*graph.Graph, []*iam.User, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.User
	out, err := s.fetch_all_user()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.ListUsersOutput).Users {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Access) fetch_all_group_graph() (*graph.Graph, []*iam.Group, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.Group
	out, err := s.fetch_all_group()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.ListGroupsOutput).Groups {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Access) fetch_all_role_graph() (*graph.Graph, []*iam.Role, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.Role
	out, err := s.fetch_all_role()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.ListRolesOutput).Roles {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Access) fetch_all_policy_graph() (*graph.Graph, []*iam.Policy, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.Policy
	out, err := s.fetch_all_policy()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.ListPoliciesOutput).Policies {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

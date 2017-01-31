package aws

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

var ErrNoPublicIP = errors.New("This instance has no public IP address")
var ErrNoAccessKey = errors.New("This instance has no access key set")

func (inf *Infra) FetchRDFResources(resourceType graph.ResourceType) (*graph.Graph, error) {
	return cloud.FetchRDFResources(inf, resourceType)
}

func (access *Access) FetchRDFResources(resourceType graph.ResourceType) (*graph.Graph, error) {
	return cloud.FetchRDFResources(access, resourceType)
}

func (inf *Infra) InstancesGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	instances, err := inf.Instances()
	if err != nil {
		return nil, err
	}
	for _, res := range instances.(*ec2.DescribeInstancesOutput).Reservations {
		for _, inst := range res.Instances {
			res, err := NewResource(inst)
			if err != nil {
				return g, err
			}
			g.AddResource(res)
		}
	}
	return g, nil
}

func (inf *Infra) VpcsGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.Vpcs()
	if err != nil {
		return nil, err
	}
	for _, vpc := range out.(*ec2.DescribeVpcsOutput).Vpcs {
		res, err := NewResource(vpc)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) SubnetsGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.Subnets()
	if err != nil {
		return nil, err
	}
	for _, subnet := range out.(*ec2.DescribeSubnetsOutput).Subnets {
		res, err := NewResource(subnet)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) InternetgatewaysGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.InternetGateways()
	if err != nil {
		return nil, err
	}
	for _, gw := range out.(*ec2.DescribeInternetGatewaysOutput).InternetGateways {
		res, err := NewResource(gw)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) SecuritygroupsGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.SecurityGroups()
	if err != nil {
		return nil, err
	}
	for _, sec := range out.(*ec2.DescribeSecurityGroupsOutput).SecurityGroups {
		res, err := NewResource(sec)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) KeypairsGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.Keypairs()
	if err != nil {
		return nil, err
	}
	for _, keypair := range out.(*ec2.DescribeKeyPairsOutput).KeyPairs {
		res, err := NewResource(keypair)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) VolumesGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.Volumes()
	if err != nil {
		return nil, err
	}
	for _, vol := range out.(*ec2.DescribeVolumesOutput).Volumes {
		res, err := NewResource(vol)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (inf *Infra) RoutetablesGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := inf.RouteTables()
	if err != nil {
		return nil, err
	}
	for _, rt := range out.(*ec2.DescribeRouteTablesOutput).RouteTables {
		res, err := NewResource(rt)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (access *Access) UsersGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := access.Users()
	if err != nil {
		return nil, err
	}
	for _, user := range out.(*iam.ListUsersOutput).Users {
		res, err := NewResource(user)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (access *Access) RolesGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := access.Roles()
	if err != nil {
		return nil, err
	}
	for _, role := range out.(*iam.ListRolesOutput).Roles {
		res, err := NewResource(role)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (access *Access) GroupsGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := access.Groups()
	if err != nil {
		return nil, err
	}
	for _, group := range out.(*iam.ListGroupsOutput).Groups {
		res, err := NewResource(group)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func (access *Access) PoliciesGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	out, err := access.LocalPolicies()
	if err != nil {
		return nil, err
	}
	for _, pol := range out.(*iam.ListPoliciesOutput).Policies {
		res, err := NewResource(pol)
		if err != nil {
			return g, err
		}
		g.AddResource(res)
	}
	return g, nil
}

func BuildAwsAccessGraph(region string, access *AwsAccess) (*graph.Graph, error) {
	g := graph.NewGraph()

	regionN := graph.InitResource(region, graph.Region)
	g.AddResource(regionN)

	policiesIndex := make(map[string]*graph.Resource)
	for _, policy := range access.Policies {
		res, err := NewResource(policy)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		policiesIndex[awssdk.StringValue(policy.PolicyName)] = res
	}

	groupsIndex := make(map[string]*graph.Resource)
	for _, group := range access.GroupsDetail {
		res, err := NewResource(group)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		groupsIndex[res.Id()] = res

		if policies, ok := access.GroupPolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					g.AddParent(policyNode, res)
				}
			}
		}
	}

	for _, user := range access.Users {
		res, err := NewResource(user)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		if groupIds, ok := access.UserGroups[res.Id()]; ok {
			for _, groupId := range groupIds {
				if groupNode, present := groupsIndex[groupId]; present {
					g.AddParent(groupNode, res)
				}
			}
		}

		if policies, ok := access.UserPolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					g.AddParent(policyNode, res)
				}
			}
		}
	}

	for _, role := range access.RolesDetail {
		res, err := NewResource(role)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		if policies, ok := access.RolePolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					g.AddParent(policyNode, res)
				}
			}
		}
	}

	return g, nil
}

func BuildAwsInfraGraph(region string, awsInfra *AwsInfra) (g *graph.Graph, err error) {
	g = graph.NewGraph()
	var vpcNodes, subnetNodes, secGroupNodes []*graph.Resource

	regionN := graph.InitResource(region, graph.Region)
	g.AddResource(regionN)

	for _, vpc := range awsInfra.Vpcs {
		res, err := NewResource(vpc)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		vpcNodes = append(vpcNodes, res)
	}

	for _, subnet := range awsInfra.Subnets {
		res, err := NewResource(subnet)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)

		subnetNodes = append(subnetNodes, res)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(subnet.VpcId))
		if vpcN != nil {
			g.AddParent(vpcN, res)
		}
	}

	for _, secgroup := range awsInfra.SecurityGroups {
		res, err := NewResource(secgroup)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)

		secGroupNodes = append(secGroupNodes, res)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(secgroup.VpcId))
		if vpcN != nil {
			g.AddParent(vpcN, res)
		}
	}

	for _, keypair := range awsInfra.Keypairs {
		res, err := NewResource(keypair)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)
	}

	for _, gw := range awsInfra.InternetGateways {
		res, err := NewResource(gw)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		for _, att := range gw.Attachments {
			vpcN := findNodeById(vpcNodes, awssdk.StringValue(att.VpcId))
			if vpcN != nil {
				g.AddParent(vpcN, res)
			}
		}
	}

	for _, rt := range awsInfra.RouteTables {
		res, err := NewResource(rt)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(rt.VpcId))
		if vpcN != nil {
			g.AddParent(vpcN, res)
		}
		for _, assos := range rt.Associations {
			if awssdk.StringValue(assos.RouteTableId) == awssdk.StringValue(rt.RouteTableId) {
				subN := findNodeById(subnetNodes, awssdk.StringValue(assos.SubnetId))
				if subN != nil {
					g.AddParent(subN, res)
				}
			}
		}
	}

	for _, instance := range awsInfra.Instances {
		res, err := NewResource(instance)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)

		subnetN := findNodeById(subnetNodes, awssdk.StringValue(instance.SubnetId))
		if subnetN != nil {
			g.AddParent(subnetN, res)
		}

		for _, refSecGroup := range instance.SecurityGroups {
			secGroupN := findNodeById(secGroupNodes, awssdk.StringValue(refSecGroup.GroupId))

			if secGroupN != nil {
				g.AddParent(secGroupN, res)
			}
		}
	}

	return g, nil
}

func InstanceCredentialsFromGraph(g *graph.Graph, instanceID string) (*shell.Credentials, error) {
	inst, err := g.GetResource(graph.Instance, instanceID)
	if err != nil {
		return nil, err
	}

	ip, ok := inst.Properties["PublicIp"]
	if !ok {
		return nil, ErrNoPublicIP
	}

	key, ok := inst.Properties["KeyName"]
	if !ok {
		return nil, ErrNoAccessKey
	}
	return &shell.Credentials{IP: fmt.Sprint(ip), User: "", KeyName: fmt.Sprint(key)}, nil
}

func findNodeById(resources []*graph.Resource, id string) *graph.Resource {
	for _, r := range resources {
		if id == r.Id() {
			return r
		}
	}
	return nil
}

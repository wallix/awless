package aws

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

var ErrInstanceNotFound = errors.New("Unknown instance")
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
			if err := addCloudResourceToGraph(g, inst); err != nil {
				return g, err
			}
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
		if err := addCloudResourceToGraph(g, vpc); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, subnet); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, sec); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, user); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, role); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, group); err != nil {
			return g, err
		}
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
		if err := addCloudResourceToGraph(g, pol); err != nil {
			return g, err
		}
	}
	return g, nil
}

func BuildAwsAccessGraph(region string, access *AwsAccess) (*graph.Graph, error) {
	g := graph.NewGraph()

	regionN, err := node.NewNodeFromStrings(graph.Region.ToRDFString(), region)
	if err != nil {
		return g, err
	}

	t, err := graph.NewRegionTypeTriple(regionN)
	if err != nil {
		return g, err
	}
	g.Add(t)

	policiesIndex := make(map[string]*node.Node)
	for _, policy := range access.Policies {
		res, err := NewResource(policy)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = graph.NewParentOfTriple(regionN, n)
		if err != nil {
			return g, err
		}
		g.Add(t)

		policiesIndex[awssdk.StringValue(policy.PolicyName)] = n
	}

	groupsIndex := make(map[string]*node.Node)
	for _, group := range access.Groups {
		res, err := NewResource(group)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = graph.NewParentOfTriple(regionN, n)
		if err != nil {
			return g, err
		}
		g.Add(t)

		groupsIndex[res.Id()] = n

		if policies, ok := access.GroupPolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					t, err = graph.NewParentOfTriple(policyNode, n)
					if err != nil {
						return g, err
					}
					g.Add(t)
				}
			}
		}
	}

	for _, user := range access.Users {
		res, err := NewResource(user)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = graph.NewParentOfTriple(regionN, n)
		if err != nil {
			return g, err
		}
		g.Add(t)

		if groupIds, ok := access.UserGroups[res.Id()]; ok {
			for _, groupId := range groupIds {
				if groupNode, present := groupsIndex[groupId]; present {
					t, err = graph.NewParentOfTriple(groupNode, n)
					if err != nil {
						return g, err
					}
					g.Add(t)
				}
			}
		}

		if policies, ok := access.UserPolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					t, err = graph.NewParentOfTriple(policyNode, n)
					if err != nil {
						return g, err
					}
					g.Add(t)
				}
			}
		}
	}

	for _, role := range access.Roles {
		res, err := NewResource(role)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = graph.NewParentOfTriple(regionN, n)
		if err != nil {
			return g, err
		}
		g.Add(t)

		if policies, ok := access.RolePolicies[res.Id()]; ok {
			for _, policy := range policies {
				if policyNode, present := policiesIndex[policy]; present {
					t, err = graph.NewParentOfTriple(policyNode, n)
					if err != nil {
						return g, err
					}
					g.Add(t)
				}
			}
		}
	}

	return g, nil
}

func BuildAwsInfraGraph(region string, awsInfra *AwsInfra) (g *graph.Graph, err error) {
	g = graph.NewGraph()
	var vpcNodes, subnetNodes, secGroupNodes []*node.Node

	regionN, err := node.NewNodeFromStrings(graph.Region.ToRDFString(), region)
	if err != nil {
		return g, err
	}

	t, err := graph.NewRegionTypeTriple(regionN)
	if err != nil {
		return g, err
	}
	g.Add(t)

	for _, vpc := range awsInfra.Vpcs {
		res, err := NewResource(vpc)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}
		vpcNodes = append(vpcNodes, n)
		t, err := graph.NewParentOfTriple(regionN, n)
		if err != nil {
			return g, fmt.Errorf("region %s", err)
		}
		g.Add(t)
	}

	for _, subnet := range awsInfra.Subnets {
		res, err := NewResource(subnet)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}

		subnetNodes = append(subnetNodes, n)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(subnet.VpcId))
		if vpcN != nil {
			t, err := graph.NewParentOfTriple(vpcN, n)
			if err != nil {
				return g, fmt.Errorf("vpc %s", err)
			}
			g.Add(t)
		}
	}

	for _, secgroup := range awsInfra.SecurityGroups {
		res, err := NewResource(secgroup)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}

		secGroupNodes = append(secGroupNodes, n)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(secgroup.VpcId))
		if vpcN != nil {
			t, err := graph.NewParentOfTriple(vpcN, n)
			if err != nil {
				return g, fmt.Errorf("vpc %s", err)
			}
			g.Add(t)
		}
	}

	for _, instance := range awsInfra.Instances {
		res, err := NewResource(instance)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.BuildRdfSubject()
		if err != nil {
			return g, err
		}

		subnetN := findNodeById(subnetNodes, awssdk.StringValue(instance.SubnetId))

		if subnetN != nil {
			t, err := graph.NewParentOfTriple(subnetN, n)
			if err != nil {
				return g, fmt.Errorf("instances subnet %s", err)
			}
			g.Add(t)
		}

		for _, refSecGroup := range instance.SecurityGroups {
			secGroupN := findNodeById(secGroupNodes, awssdk.StringValue(refSecGroup.GroupId))

			if secGroupN != nil {
				t, err := graph.NewParentOfTriple(secGroupN, n)
				if err != nil {
					return g, fmt.Errorf("instances security groups %s", err)
				}
				g.Add(t)
			}
		}

	}

	return g, nil
}

func InstanceCredentialsFromGraph(g *graph.Graph, instanceID string) (*shell.Credentials, error) {
	inst := graph.InitResource(instanceID, graph.Instance)
	err := inst.UnmarshalFromGraph(g)
	if err != nil {
		return nil, err
	}

	if !inst.ExistsInGraph(g) {
		return nil, ErrInstanceNotFound
	}

	ip, ok := inst.Properties()["PublicIp"]
	if !ok {
		return nil, ErrNoPublicIP
	}

	key, ok := inst.Properties()["KeyName"]
	if !ok {
		return nil, ErrNoAccessKey
	}
	return &shell.Credentials{IP: fmt.Sprint(ip), User: "", KeyName: fmt.Sprint(key)}, nil
}

func findNodeById(nodes []*node.Node, id string) *node.Node {
	for _, n := range nodes {
		if id == n.ID().String() {
			return n
		}
	}
	return nil
}

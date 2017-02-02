package aws

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

var ErrNoPublicIP = errors.New("This instance has no public IP address")
var ErrNoAccessKey = errors.New("This instance has no access key set")

func (access *Access) FetchRDFResources(resourceType graph.ResourceType) (*graph.Graph, error) {
	fnName := fmt.Sprintf("%sGraph", strings.Title(resourceType.PluralString()))
	method := reflect.ValueOf(access).MethodByName(fnName)
	if method.IsValid() && !method.IsNil() {
		methodI := method.Interface()
		if graphFn, ok := methodI.(func() (*graph.Graph, error)); ok {
			return graphFn()
		}
	}
	return nil, (fmt.Errorf("Unknown type of resource: %s", resourceType.String()))
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

func (inf *Infra) FetchResources() (*graph.Graph, error) {
	infra, err := inf.fetch_ec2()
	if err != nil {
		return nil, err
	}

	return BuildAwsInfraGraph(inf.region, infra)
}

func (inf *Infra) FetchAwsInfra() (*AwsInfra, error) {
	return inf.fetch_ec2()
}

func BuildAwsInfraGraph(region string, awsInfra *AwsInfra) (g *graph.Graph, err error) {
	g = graph.NewGraph()
	var vpcNodes, subnetNodes, secGroupNodes []*graph.Resource

	regionN := graph.InitResource(region, graph.Region)
	g.AddResource(regionN)

	for _, vpc := range awsInfra.vpcList {
		res, err := NewResource(vpc)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		vpcNodes = append(vpcNodes, res)
	}

	for _, subnet := range awsInfra.subnetList {
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

	for _, secgroup := range awsInfra.securitygroupList {
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

	for _, keypair := range awsInfra.keypairList {
		res, err := NewResource(keypair)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)
	}

	for _, gw := range awsInfra.internetgatewayList {
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

	for _, rt := range awsInfra.routetableList {
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

	for _, instance := range awsInfra.instanceList {
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

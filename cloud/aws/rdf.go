package aws

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/shell"
)

var ErrNoPublicIP = errors.New("This instance has no public IP address")
var ErrNoAccessKey = errors.New("This instance has no access key set")

func (acc *Access) FetchResources() (*graph.Graph, error) {
	access, err := acc.global_fetch()
	if err != nil {
		return nil, err
	}

	return BuildAwsAccessGraph(acc.region, access)
}

func (acc *Access) FetchAwsAccess() (*AwsAccess, error) {
	return acc.global_fetch()
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
	infra, err := inf.global_fetch()
	if err != nil {
		return nil, err
	}

	return BuildAwsInfraGraph(inf.region, infra)
}

func (inf *Infra) FetchAwsInfra() (*AwsInfra, error) {
	return inf.global_fetch()
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

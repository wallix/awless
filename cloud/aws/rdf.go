package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/wallix/awless/graph"
)

func (acc *Access) FetchResources() (*graph.Graph, error) {
	access, err := acc.global_fetch()
	if err != nil {
		return nil, err
	}

	return buildAccessGraph(acc.region, access)
}

func buildAccessGraph(region string, access *AwsAccess) (*graph.Graph, error) {
	g := graph.NewGraph()

	regionN := graph.InitResource(region, graph.Region)
	g.AddResource(regionN)

	policiesIndex := make(map[string]*graph.Resource)
	for _, policy := range access.Policies {
		res, err := newResource(policy)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		policiesIndex[awssdk.StringValue(policy.PolicyName)] = res
	}

	groupsIndex := make(map[string]*graph.Resource)
	for _, group := range access.GroupsDetail {
		res, err := newResource(group)
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
		res, err := newResource(user)
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
		res, err := newResource(role)
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

	return buildInfraGraph(inf.region, infra)
}

func buildInfraGraph(region string, awsInfra *AwsInfra) (g *graph.Graph, err error) {
	g = graph.NewGraph()
	var vpcNodes, subnetNodes, secGroupNodes []*graph.Resource

	regionN := graph.InitResource(region, graph.Region)
	g.AddResource(regionN)

	for _, vpc := range awsInfra.vpcList {
		res, err := newResource(vpc)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)

		vpcNodes = append(vpcNodes, res)
	}

	for _, subnet := range awsInfra.subnetList {
		res, err := newResource(subnet)
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
		res, err := newResource(secgroup)
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
		res, err := newResource(keypair)
		if err != nil {
			return nil, err
		}
		g.AddResource(res)
		g.AddParent(regionN, res)
	}

	for _, gw := range awsInfra.internetgatewayList {
		res, err := newResource(gw)
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
		res, err := newResource(rt)
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
		res, err := newResource(instance)
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

func findNodeById(resources []*graph.Resource, id string) *graph.Resource {
	for _, r := range resources {
		if id == r.Id() {
			return r
		}
	}
	return nil
}

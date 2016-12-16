package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/rdf"
)

func (inf *Infra) InstancesGraph() (*rdf.Graph, error) {
	out, err := inf.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	var triples []*triple.Triple
	for _, res := range out.Reservations {
		for _, inst := range res.Instances {
			_, err := cloud.AddNodeWithPropertiesToTriples(rdf.INSTANCE, awssdk.StringValue(inst.InstanceId), inst, awsResourcesProperties, &triples)
			if err != nil {
				return nil, err
			}
		}
	}
	return rdf.NewGraphFromTriples(triples), nil
}

func (inf *Infra) VpcsGraph() (*rdf.Graph, error) {
	out, err := inf.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}
	var triples []*triple.Triple
	for _, vpc := range out.Vpcs {
		_, err := cloud.AddNodeWithPropertiesToTriples(rdf.VPC, awssdk.StringValue(vpc.VpcId), vpc, awsResourcesProperties, &triples)
		if err != nil {
			return nil, err
		}
	}
	return rdf.NewGraphFromTriples(triples), nil
}

func (inf *Infra) SubnetsGraph() (*rdf.Graph, error) {
	out, err := inf.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, err
	}
	var triples []*triple.Triple
	for _, subnet := range out.Subnets {
		_, err := cloud.AddNodeWithPropertiesToTriples(rdf.SUBNET, awssdk.StringValue(subnet.SubnetId), subnet, awsResourcesProperties, &triples)
		if err != nil {
			return nil, err
		}
	}
	return rdf.NewGraphFromTriples(triples), nil
}

func BuildAwsAccessGraph(region string, access *AwsAccess) (*rdf.Graph, error) {
	triples, err := buildAccessRdfTriples(region, access)
	if err != nil {
		return nil, err
	}

	return rdf.NewGraphFromTriples(triples), nil
}

func BuildAwsInfraGraph(region string, infra *AwsInfra) (*rdf.Graph, error) {
	triples, err := buildInfraRdfTriples(region, infra)
	if err != nil {
		return nil, err
	}

	return rdf.NewGraphFromTriples(triples), nil
}

func buildAccessRdfTriples(region string, access *AwsAccess) ([]*triple.Triple, error) {
	var triples []*triple.Triple

	regionN, err := node.NewNodeFromStrings(rdf.REGION, region)
	if err != nil {
		return triples, err
	}

	t, err := triple.New(regionN, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	usersIndex := make(map[string]*node.Node)
	for _, user := range access.Users {
		userId := awssdk.StringValue(user.UserId)
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.USER, userId, user, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		usersIndex[userId] = n
	}

	rolesIndex := make(map[string]*node.Node)
	for _, role := range access.Roles {
		roleId := awssdk.StringValue(role.RoleId)
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.ROLE, roleId, role, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		rolesIndex[roleId] = n
	}

	groupsIndex := make(map[string]*node.Node)
	for _, group := range access.Groups {
		groupId := awssdk.StringValue(group.GroupId)
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.GROUP, groupId, group, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		groupsIndex[groupId] = n

		for _, userId := range access.UsersByGroup[groupId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("group %s has user %s, but this user does not exist", groupId, userId)
			}
			t, err = triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	for _, policy := range access.LocalPolicies {
		policyId := awssdk.StringValue(policy.PolicyId)
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.POLICY, policyId, policy, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		for _, userId := range access.UsersByLocalPolicies[policyId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, userId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, groupId := range access.GroupsByLocalPolicies[policyId] {
			if groupsIndex[groupId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, groupId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(groupsIndex[groupId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, roleId := range access.RolesByLocalPolicies[policyId] {
			if rolesIndex[roleId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, roleId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(rolesIndex[roleId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	return triples, nil
}

func buildInfraRdfTriples(region string, awsInfra *AwsInfra) (triples []*triple.Triple, err error) {
	var vpcNodes, subnetNodes []*node.Node

	regionN, err := node.NewNodeFromStrings(rdf.REGION, region)
	if err != nil {
		return triples, err
	}

	t, err := triple.New(regionN, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	for _, vpc := range awsInfra.Vpcs {
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.VPC, awssdk.StringValue(vpc.VpcId), vpc, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}
		vpcNodes = append(vpcNodes, n)
		t, err := triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return triples, fmt.Errorf("region %s", err)
		}
		triples = append(triples, t)
	}

	for _, subnet := range awsInfra.Subnets {
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.SUBNET, awssdk.StringValue(subnet.SubnetId), subnet, awsResourcesProperties, &triples)
		if err != nil {
			return triples, fmt.Errorf("subnet %s", err)
		}

		subnetNodes = append(subnetNodes, n)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(subnet.VpcId))
		if vpcN != nil {
			t, err := triple.New(vpcN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
			if err != nil {
				return triples, fmt.Errorf("vpc %s", err)
			}
			triples = append(triples, t)
		}
	}

	for _, instance := range awsInfra.Instances {
		n, err := cloud.AddNodeWithPropertiesToTriples(rdf.INSTANCE, awssdk.StringValue(instance.InstanceId), instance, awsResourcesProperties, &triples)
		if err != nil {
			return triples, err
		}

		subnetN := findNodeById(subnetNodes, awssdk.StringValue(instance.SubnetId))

		if subnetN != nil {
			t, err := triple.New(subnetN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
			if err != nil {
				return triples, fmt.Errorf("instances subnet %s", err)
			}
			triples = append(triples, t)
		}
	}

	return triples, nil
}

func findNodeById(nodes []*node.Node, id string) *node.Node {
	for _, n := range nodes {
		if id == n.ID().String() {
			return n
		}
	}
	return nil
}

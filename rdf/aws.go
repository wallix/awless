package rdf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/api"
)

const (
	REGION   = "/region"
	VPC      = "/vpc"
	SUBNET   = "/subnet"
	INSTANCE = "/instance"
	USER     = "/user"
	ROLE     = "/role"
	GROUP    = "/group"
	POLICY   = "/policy"
)

var regionL, vpcL, subnetL, instanceL, userL, roleL, groupL, policyL *literal.Literal

func init() {
	var err error
	if regionL, err = literal.DefaultBuilder().Build(literal.Text, REGION); err != nil {
		panic(err)
	}
	if vpcL, err = literal.DefaultBuilder().Build(literal.Text, VPC); err != nil {
		panic(err)
	}
	if subnetL, err = literal.DefaultBuilder().Build(literal.Text, SUBNET); err != nil {
		panic(err)
	}
	if instanceL, err = literal.DefaultBuilder().Build(literal.Text, INSTANCE); err != nil {
		panic(err)
	}
	if userL, err = literal.DefaultBuilder().Build(literal.Text, USER); err != nil {
		panic(err)
	}
	if roleL, err = literal.DefaultBuilder().Build(literal.Text, ROLE); err != nil {
		panic(err)
	}
	if groupL, err = literal.DefaultBuilder().Build(literal.Text, GROUP); err != nil {
		panic(err)
	}
	if policyL, err = literal.DefaultBuilder().Build(literal.Text, POLICY); err != nil {
		panic(err)
	}
}

func BuildAwsAccessGraph(region string, access *api.AwsAccess) (*Graph, error) {
	triples, err := buildAccessRdfTriples(region, access)
	if err != nil {
		return nil, err
	}

	return NewGraphFromTriples(triples), nil
}

func BuildAwsInfraGraph(region string, infra *api.AwsInfra) (*Graph, error) {
	triples, err := buildInfraRdfTriples(region, infra)
	if err != nil {
		return nil, err
	}

	return NewGraphFromTriples(triples), nil
}

func buildAccessRdfTriples(region string, access *api.AwsAccess) ([]*triple.Triple, error) {
	var triples []*triple.Triple

	regionN, err := node.NewNodeFromStrings(REGION, region)
	if err != nil {
		return triples, err
	}

	t, err := triple.New(regionN, HasType, triple.NewLiteralObject(regionL))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	usersIndex := make(map[string]*node.Node)
	for _, user := range access.Users {
		userId := aws.StringValue(user.UserId)
		n, err := node.NewNodeFromStrings(USER, userId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(userL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, ParentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		usersIndex[userId] = n
	}

	rolesIndex := make(map[string]*node.Node)
	for _, role := range access.Roles {
		roleId := aws.StringValue(role.RoleId)
		n, err := node.NewNodeFromStrings(ROLE, roleId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(roleL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, ParentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		rolesIndex[roleId] = n
	}

	groupsIndex := make(map[string]*node.Node)
	for _, group := range access.Groups {
		groupId := aws.StringValue(group.GroupId)
		n, err := node.NewNodeFromStrings(GROUP, groupId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(groupL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, ParentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		groupsIndex[groupId] = n

		for _, userId := range access.UsersByGroup[groupId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("group %s has user %s, but this user does not exist", groupId, userId)
			}
			t, err := triple.New(n, ParentOf, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	for _, policy := range access.LocalPolicies {
		policyId := aws.StringValue(policy.PolicyId)
		n, err := node.NewNodeFromStrings(POLICY, policyId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(policyL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		t, err = triple.New(regionN, ParentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		for _, userId := range access.UsersByLocalPolicies[policyId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, userId)
			}
			t, err := triple.New(n, ParentOf, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, groupId := range access.GroupsByLocalPolicies[policyId] {
			if groupsIndex[groupId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, groupId)
			}
			t, err := triple.New(n, ParentOf, triple.NewNodeObject(groupsIndex[groupId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, roleId := range access.RolesByLocalPolicies[policyId] {
			if rolesIndex[roleId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, roleId)
			}
			t, err := triple.New(n, ParentOf, triple.NewNodeObject(rolesIndex[roleId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	return triples, nil
}

func buildInfraRdfTriples(region string, awsInfra *api.AwsInfra) ([]*triple.Triple, error) {
	var triples []*triple.Triple
	var vpcNodes, subnetNodes []*node.Node

	regionN, err := node.NewNodeFromStrings(REGION, region)
	if err != nil {
		return triples, err
	}
	t, err := triple.New(regionN, HasType, triple.NewLiteralObject(regionL))
	if err != nil {
		return triples, err
	}
	triples = append(triples, t)

	for _, vpc := range awsInfra.Vpcs {
		n, err := node.NewNodeFromStrings(VPC, aws.StringValue(vpc.VpcId))
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(vpcL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		vpcNodes = append(vpcNodes, n)
		t, err = triple.New(regionN, ParentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, fmt.Errorf("region %s", err)
		}
		triples = append(triples, t)
	}

	for _, subnet := range awsInfra.Subnets {
		n, err := node.NewNodeFromStrings(SUBNET, aws.StringValue(subnet.SubnetId))
		if err != nil {
			return triples, fmt.Errorf("subnet %s", err)
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(subnetL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		subnetNodes = append(subnetNodes, n)

		vpcN := findNodeById(vpcNodes, aws.StringValue(subnet.VpcId))
		if vpcN != nil {
			t, err := triple.New(vpcN, ParentOf, triple.NewNodeObject(n))
			if err != nil {
				return triples, fmt.Errorf("vpc %s", err)
			}
			triples = append(triples, t)
		}
	}

	for _, instance := range awsInfra.Instances {
		n, err := node.NewNodeFromStrings(INSTANCE, aws.StringValue(instance.InstanceId))
		if err != nil {
			return triples, err
		}
		t, err := triple.New(n, HasType, triple.NewLiteralObject(instanceL))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)
		subnetN := findNodeById(subnetNodes, aws.StringValue(instance.SubnetId))

		if subnetN != nil {
			t, err := triple.New(subnetN, ParentOf, triple.NewNodeObject(n))
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

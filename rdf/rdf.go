package rdf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"github.com/wallix/awless/api"
)

var parentOf *predicate.Predicate

func init() {
	var err error
	if parentOf, err = predicate.NewImmutable("parent_of"); err != nil {
		panic(err)
	}
}

func BuildAwsAccessGraph(graphname string, region string, access *api.AwsAccess) (*Graph, error) {
	triples, err := buildAccessRdfTriples(region, access)
	if err != nil {
		return nil, err
	}

	g, err := NewNamedGraphFromTriples(graphname, triples)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func BuildAwsInfraGraph(graphname string, region string, infra *api.AwsInfra) (*Graph, error) {
	triples, err := buildInfraRdfTriples(region, infra)
	if err != nil {
		return nil, err
	}

	g, err := NewNamedGraphFromTriples(graphname, triples)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func buildAccessRdfTriples(region string, access *api.AwsAccess) ([]*triple.Triple, error) {
	var triples []*triple.Triple

	regionN, err := node.NewNodeFromStrings("/region", region)
	if err != nil {
		return triples, err
	}

	usersIndex := make(map[string]*node.Node)
	for _, user := range access.Users {
		userId := aws.StringValue(user.UserId)
		n, err := node.NewNodeFromStrings("/user", userId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		usersIndex[userId] = n
	}

	rolesIndex := make(map[string]*node.Node)
	for _, role := range access.Roles {
		roleId := aws.StringValue(role.RoleId)
		n, err := node.NewNodeFromStrings("/role", roleId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		rolesIndex[roleId] = n
	}

	groupsIndex := make(map[string]*node.Node)
	for _, group := range access.Groups {
		groupId := aws.StringValue(group.GroupId)
		n, err := node.NewNodeFromStrings("/group", groupId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		groupsIndex[groupId] = n

		for _, userId := range access.UsersByGroup[groupId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("group %s has user %s, but this user does not exist", groupId, userId)
			}
			t, err := triple.New(n, parentOf, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	for _, policy := range access.LocalPolicies {
		policyId := aws.StringValue(policy.PolicyId)
		n, err := node.NewNodeFromStrings("/policy", policyId)
		if err != nil {
			return triples, err
		}
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		for _, userId := range access.UsersByLocalPolicies[policyId] {
			if usersIndex[userId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, userId)
			}
			t, err := triple.New(n, parentOf, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, groupId := range access.GroupsByLocalPolicies[policyId] {
			if groupsIndex[groupId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, groupId)
			}
			t, err := triple.New(n, parentOf, triple.NewNodeObject(groupsIndex[groupId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}

		for _, roleId := range access.RolesByLocalPolicies[policyId] {
			if rolesIndex[roleId] == nil {
				return triples, fmt.Errorf("policy %s has user %s, but this user does not exist", policyId, roleId)
			}
			t, err := triple.New(n, parentOf, triple.NewNodeObject(rolesIndex[roleId]))
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

	regionN, err := node.NewNodeFromStrings("/region", region)
	if err != nil {
		return triples, err
	}

	for _, vpc := range awsInfra.Vpcs {
		n, err := node.NewNodeFromStrings("/vpc", aws.StringValue(vpc.VpcId))
		if err != nil {
			return triples, err
		}

		vpcNodes = append(vpcNodes, n)
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, fmt.Errorf("region %s", err)
		}
		triples = append(triples, t)
	}

	for _, subnet := range awsInfra.Subnets {
		n, err := node.NewNodeFromStrings("/subnet", aws.StringValue(subnet.SubnetId))
		if err != nil {
			return triples, fmt.Errorf("subnet %s", err)
		}

		subnetNodes = append(subnetNodes, n)

		vpcN := findNodeById(vpcNodes, aws.StringValue(subnet.VpcId))
		if vpcN != nil {
			t, err := triple.New(vpcN, parentOf, triple.NewNodeObject(n))
			if err != nil {
				return triples, fmt.Errorf("vpc %s", err)
			}
			triples = append(triples, t)
		}
	}

	for _, instance := range awsInfra.Instances {
		n, err := node.NewNodeFromStrings("/instance", aws.StringValue(instance.InstanceId))
		if err != nil {
			return triples, err
		}
		subnetN := findNodeById(subnetNodes, aws.StringValue(instance.SubnetId))

		if subnetN != nil {
			t, err := triple.New(subnetN, parentOf, triple.NewNodeObject(n))
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

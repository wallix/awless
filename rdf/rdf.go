package rdf

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
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

func BuildAccessRdfTriples(region string, access *api.AwsAccess) ([]*triple.Triple, error) {
	var triples []*triple.Triple

	regionN, err := node.NewNodeFromStrings("/region", region)
	if err != nil {
		return triples, err
	}

	userNodes := make(map[string]*node.Node)
	for _, user := range access.Users {
		n, err := node.NewNodeFromStrings("/user", aws.StringValue(user.UserId))
		if err != nil {
			return triples, err
		}
		t, err := triple.New(regionN, parentOf, triple.NewNodeObject(n))
		if err != nil {
			return triples, err
		}
		triples = append(triples, t)

		userNodes[aws.StringValue(user.UserId)] = n
	}

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
	}

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

		for _, userId := range access.UsersByGroup[groupId] {
			if userNodes[userId] == nil {
				return triples, fmt.Errorf("group %s has user %s, but this user does not exist", groupId, userId)
			}
			t, err := triple.New(n, parentOf, triple.NewNodeObject(userNodes[userId]))
			if err != nil {
				return triples, err
			}
			triples = append(triples, t)
		}
	}

	return triples, nil
}

func BuildInfraRdfTriples(region string, awsVpcs []*ec2.Vpc, awsSubnets []*ec2.Subnet, awsInstances []*ec2.Instance) ([]*triple.Triple, error) {
	var triples []*triple.Triple
	var vpcNodes, subnetNodes []*node.Node

	regionN, err := node.NewNodeFromStrings("/region", region)
	if err != nil {
		return triples, err
	}

	for _, vpc := range awsVpcs {
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

	for _, subnet := range awsSubnets {
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

	for _, instance := range awsInstances {
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

func IntersectTriples(a, b []*triple.Triple) []*triple.Triple {
	var inter []*triple.Triple

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				inter = append(inter, a[i])
			}
		}
	}

	return inter
}

func SubstractTriples(a, b []*triple.Triple) []*triple.Triple {
	var sub []*triple.Triple

	for i := 0; i < len(a); i++ {
		var found bool
		for j := 0; j < len(b); j++ {
			if a[i].String() == b[j].String() {
				found = true
			}
		}
		if !found {
			sub = append(sub, a[i])
		}
	}

	return sub
}

func MarshalTriples(triples []*triple.Triple) string {
	var triplesString []string
	for _, triple := range triples {
		triplesString = append(triplesString, triple.String())
	}
	return strings.Join(triplesString, "\n")
}

func UnmarshalTriples(raw string) ([]*triple.Triple, error) {
	var triples []*triple.Triple
	for _, rawTriple := range strings.Split(raw, "\n") {
		if strings.TrimSpace(rawTriple) == "" {
			continue
		}
		triple, err := triple.Parse(rawTriple, literal.DefaultBuilder())
		if err != nil {
			return triples, err
		}
		triples = append(triples, triple)
	}
	return triples, nil
}

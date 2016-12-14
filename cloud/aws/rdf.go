package aws

import (
	"fmt"
	"reflect"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
	"github.com/wallix/awless/rdf"
)

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
		n, err := addNode(rdf.USER, userId, user, &triples)
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
		n, err := addNode(rdf.ROLE, roleId, role, &triples)
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
		n, err := addNode(rdf.GROUP, groupId, group, &triples)
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
		n, err := addNode(rdf.POLICY, policyId, policy, &triples)
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
		n, err := addNode(rdf.VPC, awssdk.StringValue(vpc.VpcId), vpc, &triples)
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
		n, err := addNode(rdf.SUBNET, awssdk.StringValue(subnet.SubnetId), subnet, &triples)
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
		n, err := addNode(rdf.INSTANCE, awssdk.StringValue(instance.InstanceId), instance, &triples)
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

type User struct {
	Id string `aws:"UserId"`
}

type Role struct {
	Id string `aws:"RoleId"`
}

type Group struct {
	Id string `aws:"GroupId"`
}

type Policy struct {
	Id string `aws:"PolicyId"`
}

type Vpc struct {
	Id string `aws:"VpcId"`
}

type Subnet struct {
	Id    string `aws:"SubnetId"`
	VpcId string `aws:"VpcId"`
}

type Instance struct {
	Id        string `aws:"InstanceId"`
	Type      string `aws:"InstanceType"`
	SubnetId  string `aws:"SubnetId"`
	VpcId     string `aws:"VpcId"`
	PublicIp  string `aws:"PublicIpAddress"`
	PrivateIp string `aws:"PrivateIpAddress"`
	ImageId   string `aws:"ImageId"`
}

func addNode(nodeType, id string, awsNode interface{}, triples *[]*triple.Triple) (*node.Node, error) {
	n, err := node.NewNodeFromStrings(nodeType, id)
	if err != nil {
		return nil, err
	}
	var lit *literal.Literal
	if lit, err = literal.DefaultBuilder().Build(literal.Text, nodeType); err != nil {
		return nil, err
	}
	t, err := triple.New(n, rdf.HasTypePredicate, triple.NewLiteralObject(lit))
	if err != nil {
		return nil, err
	}
	*triples = append(*triples, t)

	nodeV := reflect.ValueOf(awsNode).Elem()
	var propP *predicate.Predicate
	var destType reflect.Type
	switch nodeType {
	case rdf.VPC:
		destType = reflect.TypeOf(Vpc{})
	case rdf.SUBNET:
		destType = reflect.TypeOf(Subnet{})
	case rdf.INSTANCE:
		destType = reflect.TypeOf(Instance{})
	case rdf.USER:
		destType = reflect.TypeOf(User{})
	case rdf.ROLE:
		destType = reflect.TypeOf(Role{})
	case rdf.GROUP:
		destType = reflect.TypeOf(Group{})
	case rdf.POLICY:
		destType = reflect.TypeOf(Policy{})
	default:
		return nil, fmt.Errorf("type %s is not managed", nodeType)
	}

	for i := 0; i < destType.NumField(); i++ {
		if propP, err = predicate.NewImmutable(destType.Field(i).Name); err != nil {
			return nil, err
		}
		var propL *literal.Literal
		if awsTag, ok := destType.Field(i).Tag.Lookup("aws"); ok {
			sourceField := nodeV.FieldByName(awsTag)
			if sourceField.IsValid() {
				stringValue := awssdk.StringValue(sourceField.Interface().(*string))
				if stringValue != "" {
					if propL, err = literal.DefaultBuilder().Build(literal.Text, stringValue); err != nil {
						return nil, err
					}
					propT, err := triple.New(n, propP, triple.NewLiteralObject(propL))
					if err != nil {
						return nil, err
					}
					*triples = append(*triples, propT)
				}
			}
		}
	}

	return n, nil
}

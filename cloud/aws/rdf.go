package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

func (inf *Infra) InstancesGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
	out, err := inf.Instances()
	if err != nil {
		return nil, err
	}
	instances, ok := out.(*ec2.DescribeInstancesOutput)
	if !ok {
		return nil, fmt.Errorf("invalid instances type %T", out)
	}
	for _, res := range instances.Reservations {
		for _, inst := range res.Instances {
			res, err := NewResource(inst)
			if err != nil {
				return nil, err
			}
			triples, err := res.MarshalToTriples()
			if err != nil {
				return nil, err
			}
			g.Add(triples...)
		}
	}
	return g, nil
}

func (inf *Infra) VpcsGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
	out, err := inf.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}
	for _, vpc := range out.Vpcs {
		res, err := NewResource(vpc)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
	}
	return g, nil
}

func (inf *Infra) SubnetsGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
	out, err := inf.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, err
	}
	for _, subnet := range out.Subnets {
		res, err := NewResource(subnet)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
	}
	return g, nil
}

func BuildAwsAccessGraph(region string, access *AwsAccess) (*rdf.Graph, error) {
	g := rdf.NewGraph()

	regionN, err := node.NewNodeFromStrings(rdf.REGION, region)
	if err != nil {
		return g, err
	}

	t, err := triple.New(regionN, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
	if err != nil {
		return g, err
	}
	g.Add(t)

	usersIndex := make(map[string]*node.Node)
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return g, err
		}
		g.Add(t)

		usersIndex[res.id] = n
	}

	rolesIndex := make(map[string]*node.Node)
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return g, err
		}
		g.Add(t)

		rolesIndex[res.id] = n
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return g, err
		}
		g.Add(t)

		groupsIndex[res.id] = n

		for _, userId := range access.UsersByGroup[res.id] {
			if usersIndex[userId] == nil {
				return g, fmt.Errorf("group %s has user %s, but this user does not exist", res.id, userId)
			}
			t, err = triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return g, err
			}
			g.Add(t)
		}
	}

	for _, policy := range access.LocalPolicies {
		res, err := NewResource(policy)
		if err != nil {
			return nil, err
		}
		triples, err := res.MarshalToTriples()
		if err != nil {
			return nil, err
		}
		g.Add(triples...)
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}
		t, err = triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
		if err != nil {
			return g, err
		}
		g.Add(t)

		for _, userId := range access.UsersByLocalPolicies[res.id] {
			if usersIndex[userId] == nil {
				return g, fmt.Errorf("policy %s has user %s, but this user does not exist", res.id, userId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(usersIndex[userId]))
			if err != nil {
				return g, err
			}
			g.Add(t)
		}

		for _, groupId := range access.GroupsByLocalPolicies[res.id] {
			if groupsIndex[groupId] == nil {
				return g, fmt.Errorf("policy %s has user %s, but this user does not exist", res.id, groupId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(groupsIndex[groupId]))
			if err != nil {
				return g, err
			}
			g.Add(t)
		}

		for _, roleId := range access.RolesByLocalPolicies[res.id] {
			if rolesIndex[roleId] == nil {
				return g, fmt.Errorf("policy %s has user %s, but this user does not exist", res.id, roleId)
			}
			t, err := triple.New(n, rdf.ParentOfPredicate, triple.NewNodeObject(rolesIndex[roleId]))
			if err != nil {
				return g, err
			}
			g.Add(t)
		}
	}

	return g, nil
}

func BuildAwsInfraGraph(region string, awsInfra *AwsInfra) (g *rdf.Graph, err error) {
	g = rdf.NewGraph()
	var vpcNodes, subnetNodes []*node.Node

	regionN, err := node.NewNodeFromStrings(rdf.REGION, region)
	if err != nil {
		return g, err
	}

	t, err := triple.New(regionN, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}
		vpcNodes = append(vpcNodes, n)
		t, err := triple.New(regionN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}

		subnetNodes = append(subnetNodes, n)

		vpcN := findNodeById(vpcNodes, awssdk.StringValue(subnet.VpcId))
		if vpcN != nil {
			t, err := triple.New(vpcN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
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
		n, err := res.buildRdfSubject()
		if err != nil {
			return g, err
		}

		subnetN := findNodeById(subnetNodes, awssdk.StringValue(instance.SubnetId))

		if subnetN != nil {
			t, err := triple.New(subnetN, rdf.ParentOfPredicate, triple.NewNodeObject(n))
			if err != nil {
				return g, fmt.Errorf("instances subnet %s", err)
			}
			g.Add(t)
		}
	}

	return g, nil
}

func findNodeById(nodes []*node.Node, id string) *node.Node {
	for _, n := range nodes {
		if id == n.ID().String() {
			return n
		}
	}
	return nil
}

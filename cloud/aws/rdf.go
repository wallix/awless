package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/rdf"
)

func (inf *Infra) InstancesGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (inf *Infra) VpcsGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (inf *Infra) SubnetsGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (access *Access) UsersGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (access *Access) RolesGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (access *Access) GroupsGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

func (access *Access) PoliciesGraph() (*rdf.Graph, error) {
	g := rdf.NewGraph()
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

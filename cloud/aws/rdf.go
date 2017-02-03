package aws

import (
	"fmt"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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
	return inf.fetchAndBuildGraph()
}

type addParentFn func(*graph.Graph, interface{}) error

var addParentsFns = map[string][]addParentFn{
	graph.Vpc.String():             {addRegionParent},
	graph.Subnet.String():          {subnetAddVpcParent},
	graph.Instance.String():        {instanceAddSubnetParent, instanceAddSecurityGroupsParents},
	graph.SecurityGroup.String():   {secgroupAddVpcParent},
	graph.Keypair.String():         {addRegionParent},
	graph.InternetGateway.String(): {addRegionParent, gatewayAddVpcParents},
	graph.RouteTable.String():      {routeTableAddSubnetParents, routeTableAddVpcParent},
}

func (s *Infra) fetchAndBuildGraph() (*graph.Graph, error) {
	g := graph.NewGraph()
	regionN := graph.InitResource(s.region, graph.Region)
	g.AddResource(regionN)

	var vpcs []*ec2.Vpc

	var subnets []*ec2.Subnet
	var instances []*ec2.Instance
	var secgroups []*ec2.SecurityGroup
	var keypairs []*ec2.KeyPairInfo
	var igw []*ec2.InternetGateway
	var rt []*ec2.RouteTable
	errc := make(chan error)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		var vpcGraph *graph.Graph
		var err error
		vpcGraph, vpcs, err = s.fetch_all_vpc_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(vpcGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var subnetGraph *graph.Graph
		var err error
		subnetGraph, subnets, err = s.fetch_all_subnet_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(subnetGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var instGraph *graph.Graph
		var err error
		instGraph, instances, err = s.fetch_all_instance_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(instGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var secgroupGraph *graph.Graph
		var err error
		secgroupGraph, secgroups, err = s.fetch_all_securitygroup_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(secgroupGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var keypairGraph *graph.Graph
		var err error
		keypairGraph, keypairs, err = s.fetch_all_keypair_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(keypairGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var igwGraph *graph.Graph
		var err error
		igwGraph, igw, err = s.fetch_all_internetgateway_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(igwGraph)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var rtGraph *graph.Graph
		var err error
		rtGraph, rt, err = s.fetch_all_routetable_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(rtGraph)
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return g, err
		}
	}

	errc = make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range vpcs {
			for _, fn := range addParentsFns[graph.Vpc.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range subnets {
			for _, fn := range addParentsFns[graph.Subnet.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range instances {
			for _, fn := range addParentsFns[graph.Instance.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range secgroups {
			for _, fn := range addParentsFns[graph.SecurityGroup.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range keypairs {
			for _, fn := range addParentsFns[graph.Keypair.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range igw {
			for _, fn := range addParentsFns[graph.InternetGateway.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range rt {
			for _, fn := range addParentsFns[graph.RouteTable.String()] {
				err := fn(g, r)
				if err != nil {
					errc <- err
					return
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return g, err
		}
	}

	return g, nil
}

func addRegionParent(g *graph.Graph, i interface{}) error {
	resources, err := g.GetAllResources(graph.Region)
	if err != nil {
		return err
	}
	if len(resources) != 1 {
		return fmt.Errorf("aws fetch: expect exactly one region in graph, but got %d", len(resources))
	}
	regionN := resources[0]
	switch ii := i.(type) {
	case *ec2.Vpc:
		n, err := g.GetResource(graph.Vpc, awssdk.StringValue(ii.VpcId))
		if err != nil {
			return err
		}
		g.AddParent(regionN, n)
	case *ec2.KeyPairInfo:
		n, err := g.GetResource(graph.Keypair, awssdk.StringValue(ii.KeyName))
		if err != nil {
			return err
		}
		g.AddParent(regionN, n)
	case *ec2.InternetGateway:
		n, err := g.GetResource(graph.InternetGateway, awssdk.StringValue(ii.InternetGatewayId))
		if err != nil {
			return err
		}
		g.AddParent(regionN, n)
	default:
		return fmt.Errorf("aws fetch: unkown type of resource to add region: %T", i)
	}

	return nil
}

func instanceAddSubnetParent(g *graph.Graph, i interface{}) error {
	instance, ok := i.(*ec2.Instance)
	if !ok {
		return fmt.Errorf("aws fetch: not an instance, but a %T", i)
	}
	instanceN, err := g.GetResource(graph.Instance, awssdk.StringValue(instance.InstanceId))
	if err != nil {
		return err
	}
	if awssdk.StringValue(instance.SubnetId) == "" {
		return nil
	}
	subnetN, err := g.GetResource(graph.Subnet, awssdk.StringValue(instance.SubnetId))
	if err != nil {
		return err
	}
	g.AddParent(subnetN, instanceN)
	return nil
}

func subnetAddVpcParent(g *graph.Graph, i interface{}) error {
	subnet, ok := i.(*ec2.Subnet)
	if !ok {
		return fmt.Errorf("aws fetch: not an subnet, but a %T", i)
	}
	n, err := g.GetResource(graph.Subnet, awssdk.StringValue(subnet.SubnetId))
	if err != nil {
		return err
	}
	if awssdk.StringValue(subnet.VpcId) == "" {
		return nil
	}
	parent, err := g.GetResource(graph.Vpc, awssdk.StringValue(subnet.VpcId))
	if err != nil {
		return err
	}
	g.AddParent(parent, n)
	return nil
}

func secgroupAddVpcParent(g *graph.Graph, i interface{}) error {
	secgroup, ok := i.(*ec2.SecurityGroup)
	if !ok {
		return fmt.Errorf("aws fetch: not a security group, but a %T", i)
	}
	n, err := g.GetResource(graph.SecurityGroup, awssdk.StringValue(secgroup.GroupId))
	if err != nil {
		return err
	}
	if awssdk.StringValue(secgroup.VpcId) == "" {
		return nil
	}
	parent, err := g.GetResource(graph.Vpc, awssdk.StringValue(secgroup.VpcId))
	if err != nil {
		return err
	}
	g.AddParent(parent, n)
	return nil
}

func routeTableAddVpcParent(g *graph.Graph, i interface{}) error {
	rT, ok := i.(*ec2.RouteTable)
	if !ok {
		return fmt.Errorf("aws fetch: not a route table, but a %T", i)
	}
	n, err := g.GetResource(graph.RouteTable, awssdk.StringValue(rT.RouteTableId))
	if err != nil {
		return err
	}
	if awssdk.StringValue(rT.VpcId) == "" {
		return nil
	}
	parent, err := g.GetResource(graph.Vpc, awssdk.StringValue(rT.VpcId))
	if err != nil {
		return err
	}
	g.AddParent(parent, n)
	return nil
}

func instanceAddSecurityGroupsParents(g *graph.Graph, i interface{}) error {
	instance, ok := i.(*ec2.Instance)
	if !ok {
		return fmt.Errorf("aws fetch: not an instance, but a %T", i)
	}
	instanceN, err := g.GetResource(graph.Instance, awssdk.StringValue(instance.InstanceId))
	if err != nil {
		return err
	}

	for _, refSecGroup := range instance.SecurityGroups {
		if awssdk.StringValue(refSecGroup.GroupId) == "" {
			continue
		}
		secGroupN, err := g.GetResource(graph.SecurityGroup, awssdk.StringValue(refSecGroup.GroupId))
		if err != nil {
			return err
		}
		g.AddParent(secGroupN, instanceN)
	}
	return nil
}

func gatewayAddVpcParents(g *graph.Graph, i interface{}) error {
	igw, ok := i.(*ec2.InternetGateway)
	if !ok {
		return fmt.Errorf("aws fetch: not a gateway, but a %T", i)
	}
	n, err := g.GetResource(graph.InternetGateway, awssdk.StringValue(igw.InternetGatewayId))
	if err != nil {
		return err
	}

	for _, att := range igw.Attachments {
		if awssdk.StringValue(att.VpcId) == "" {
			continue
		}
		vpc, err := g.GetResource(graph.Vpc, awssdk.StringValue(att.VpcId))
		if err != nil {
			return err
		}
		g.AddParent(vpc, n)
	}
	return nil
}

func routeTableAddSubnetParents(g *graph.Graph, i interface{}) error {
	rt, ok := i.(*ec2.RouteTable)
	if !ok {
		return fmt.Errorf("aws fetch: not a route table, but a %T", i)
	}
	n, err := g.GetResource(graph.RouteTable, awssdk.StringValue(rt.RouteTableId))
	if err != nil {
		return err
	}

	for _, ass := range rt.Associations {
		if awssdk.StringValue(ass.RouteTableId) != awssdk.StringValue(rt.RouteTableId) {
			continue
		}
		if awssdk.StringValue(ass.SubnetId) == "" {
			continue
		}
		subnet, err := g.GetResource(graph.Subnet, awssdk.StringValue(ass.SubnetId))
		if err != nil {
			return err
		}
		g.AddParent(subnet, n)
	}
	return nil
}

func findNodeById(resources []*graph.Resource, id string) *graph.Resource {
	for _, r := range resources {
		if id == r.Id() {
			return r
		}
	}
	return nil
}

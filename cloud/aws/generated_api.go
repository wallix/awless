// Auto generated implementation for the AWS cloud service

/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"fmt"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

func init() {
	ServiceNames = append(ServiceNames, "infra")
	ServiceNames = append(ServiceNames, "access")
	ServiceNames = append(ServiceNames, "storage")
}

var ServiceNames = []string{}

var ResourceTypesPerAPI = map[string][]string{
	"ec2": []string{
		"instance",
		"subnet",
		"vpc",
		"keypair",
		"securitygroup",
		"volume",
		"internetgateway",
		"routetable",
	},
	"iam": []string{
		"user",
		"group",
		"role",
		"policy",
	},
	"s3": []string{
		"bucket",
		"storageobject",
	},
}

var ServicePerResourceType = map[string]string{
	"instance":        "infra",
	"subnet":          "infra",
	"vpc":             "infra",
	"keypair":         "infra",
	"securitygroup":   "infra",
	"volume":          "infra",
	"internetgateway": "infra",
	"routetable":      "infra",
	"user":            "access",
	"group":           "access",
	"role":            "access",
	"policy":          "access",
	"bucket":          "storage",
	"storageobject":   "storage",
}

type Infra struct {
	region string
	ec2iface.EC2API
}

func NewInfra(sess *session.Session) *Infra {
	region := awssdk.StringValue(sess.Config.Region)
	return &Infra{EC2API: ec2.New(sess), region: region}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Provider() string {
	return "aws"
}

func (s *Infra) ProviderAPI() string {
	return "ec2"
}

func (s *Infra) ProviderRunnableAPI() interface{} {
	return s.EC2API
}

func (s *Infra) ResourceTypes() (all []string) {
	all = append(all, "instance")
	all = append(all, "subnet")
	all = append(all, "vpc")
	all = append(all, "keypair")
	all = append(all, "securitygroup")
	all = append(all, "volume")
	all = append(all, "internetgateway")
	all = append(all, "routetable")
	return
}

func (s *Infra) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	regionN := graph.InitResource(s.region, graph.Region)
	g.AddResource(regionN)
	var instanceList []*ec2.Instance
	var subnetList []*ec2.Subnet
	var vpcList []*ec2.Vpc
	var keypairList []*ec2.KeyPairInfo
	var securitygroupList []*ec2.SecurityGroup
	var volumeList []*ec2.Volume
	var internetgatewayList []*ec2.InternetGateway
	var routetableList []*ec2.RouteTable

	errc := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, instanceList, err = s.fetch_all_instance_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, subnetList, err = s.fetch_all_subnet_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, vpcList, err = s.fetch_all_vpc_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, keypairList, err = s.fetch_all_keypair_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, securitygroupList, err = s.fetch_all_securitygroup_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, volumeList, err = s.fetch_all_volume_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, internetgatewayList, err = s.fetch_all_internetgateway_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, routetableList, err = s.fetch_all_routetable_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case "Access Denied":
				return g, cloud.ErrFetchAccessDenied
			default:
				return g, ee
			}
		case nil:
			continue
		default:
			return g, ee
		}
	}

	errc = make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range instanceList {
			for _, fn := range addParentsFns["instance"] {
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
		for _, r := range subnetList {
			for _, fn := range addParentsFns["subnet"] {
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
		for _, r := range vpcList {
			for _, fn := range addParentsFns["vpc"] {
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
		for _, r := range keypairList {
			for _, fn := range addParentsFns["keypair"] {
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
		for _, r := range securitygroupList {
			for _, fn := range addParentsFns["securitygroup"] {
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
		for _, r := range volumeList {
			for _, fn := range addParentsFns["volume"] {
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
		for _, r := range internetgatewayList {
			for _, fn := range addParentsFns["internetgateway"] {
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
		for _, r := range routetableList {
			for _, fn := range addParentsFns["routetable"] {
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

func (s *Infra) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "instance":
		graph, _, err := s.fetch_all_instance_graph()
		return graph, err
	case "subnet":
		graph, _, err := s.fetch_all_subnet_graph()
		return graph, err
	case "vpc":
		graph, _, err := s.fetch_all_vpc_graph()
		return graph, err
	case "keypair":
		graph, _, err := s.fetch_all_keypair_graph()
		return graph, err
	case "securitygroup":
		graph, _, err := s.fetch_all_securitygroup_graph()
		return graph, err
	case "volume":
		graph, _, err := s.fetch_all_volume_graph()
		return graph, err
	case "internetgateway":
		graph, _, err := s.fetch_all_internetgateway_graph()
		return graph, err
	case "routetable":
		graph, _, err := s.fetch_all_routetable_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws infra: unsupported fetch for type %s", t)
	}
}

func (s *Infra) fetch_all_instance() (interface{}, error) {
	return s.DescribeInstances(&ec2.DescribeInstancesInput{})
}
func (s *Infra) fetch_all_subnet() (interface{}, error) {
	return s.DescribeSubnets(&ec2.DescribeSubnetsInput{})
}
func (s *Infra) fetch_all_vpc() (interface{}, error) {
	return s.DescribeVpcs(&ec2.DescribeVpcsInput{})
}
func (s *Infra) fetch_all_keypair() (interface{}, error) {
	return s.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
}
func (s *Infra) fetch_all_securitygroup() (interface{}, error) {
	return s.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
}
func (s *Infra) fetch_all_volume() (interface{}, error) {
	return s.DescribeVolumes(&ec2.DescribeVolumesInput{})
}
func (s *Infra) fetch_all_internetgateway() (interface{}, error) {
	return s.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
}
func (s *Infra) fetch_all_routetable() (interface{}, error) {
	return s.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
}

func (s *Infra) fetch_all_instance_graph() (*graph.Graph, []*ec2.Instance, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Instance
	out, err := s.fetch_all_instance()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, all := range out.(*ec2.DescribeInstancesOutput).Reservations {
		for _, output := range all.Instances {
			cloudResources = append(cloudResources, output)
			res, err := newResource(output)
			if err != nil {
				return g, cloudResources, err
			}
			g.AddResource(res)
		}
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_subnet_graph() (*graph.Graph, []*ec2.Subnet, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Subnet
	out, err := s.fetch_all_subnet()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeSubnetsOutput).Subnets {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_vpc_graph() (*graph.Graph, []*ec2.Vpc, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Vpc
	out, err := s.fetch_all_vpc()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeVpcsOutput).Vpcs {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_keypair_graph() (*graph.Graph, []*ec2.KeyPairInfo, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.KeyPairInfo
	out, err := s.fetch_all_keypair()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeKeyPairsOutput).KeyPairs {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_securitygroup_graph() (*graph.Graph, []*ec2.SecurityGroup, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.SecurityGroup
	out, err := s.fetch_all_securitygroup()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeSecurityGroupsOutput).SecurityGroups {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_volume_graph() (*graph.Graph, []*ec2.Volume, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Volume
	out, err := s.fetch_all_volume()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeVolumesOutput).Volumes {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_internetgateway_graph() (*graph.Graph, []*ec2.InternetGateway, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.InternetGateway
	out, err := s.fetch_all_internetgateway()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeInternetGatewaysOutput).InternetGateways {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Infra) fetch_all_routetable_graph() (*graph.Graph, []*ec2.RouteTable, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.RouteTable
	out, err := s.fetch_all_routetable()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*ec2.DescribeRouteTablesOutput).RouteTables {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

type Access struct {
	region string
	iamiface.IAMAPI
}

func NewAccess(sess *session.Session) *Access {
	region := awssdk.StringValue(sess.Config.Region)
	return &Access{IAMAPI: iam.New(sess), region: region}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Provider() string {
	return "aws"
}

func (s *Access) ProviderAPI() string {
	return "iam"
}

func (s *Access) ProviderRunnableAPI() interface{} {
	return s.IAMAPI
}

func (s *Access) ResourceTypes() (all []string) {
	all = append(all, "user")
	all = append(all, "group")
	all = append(all, "role")
	all = append(all, "policy")
	return
}

func (s *Access) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	regionN := graph.InitResource(s.region, graph.Region)
	g.AddResource(regionN)
	var userList []*iam.UserDetail
	var groupList []*iam.GroupDetail
	var roleList []*iam.RoleDetail
	var policyList []*iam.Policy

	errc := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, userList, err = s.fetch_all_user_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, groupList, err = s.fetch_all_group_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, roleList, err = s.fetch_all_role_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, policyList, err = s.fetch_all_policy_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case "Access Denied":
				return g, cloud.ErrFetchAccessDenied
			default:
				return g, ee
			}
		case nil:
			continue
		default:
			return g, ee
		}
	}

	errc = make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range userList {
			for _, fn := range addParentsFns["user"] {
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
		for _, r := range groupList {
			for _, fn := range addParentsFns["group"] {
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
		for _, r := range roleList {
			for _, fn := range addParentsFns["role"] {
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
		for _, r := range policyList {
			for _, fn := range addParentsFns["policy"] {
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

func (s *Access) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "user":
		graph, _, err := s.fetch_all_user_graph()
		return graph, err
	case "group":
		graph, _, err := s.fetch_all_group_graph()
		return graph, err
	case "role":
		graph, _, err := s.fetch_all_role_graph()
		return graph, err
	case "policy":
		graph, _, err := s.fetch_all_policy_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws access: unsupported fetch for type %s", t)
	}
}

func (s *Access) fetch_all_group() (interface{}, error) {
	return s.GetAccountAuthorizationDetails(&iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}})
}
func (s *Access) fetch_all_role() (interface{}, error) {
	return s.GetAccountAuthorizationDetails(&iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}})
}
func (s *Access) fetch_all_policy() (interface{}, error) {
	return s.ListPolicies(&iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)})
}

func (s *Access) fetch_all_group_graph() (*graph.Graph, []*iam.GroupDetail, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.GroupDetail
	out, err := s.fetch_all_group()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.GetAccountAuthorizationDetailsOutput).GroupDetailList {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Access) fetch_all_role_graph() (*graph.Graph, []*iam.RoleDetail, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.RoleDetail
	out, err := s.fetch_all_role()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.GetAccountAuthorizationDetailsOutput).RoleDetailList {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

func (s *Access) fetch_all_policy_graph() (*graph.Graph, []*iam.Policy, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.Policy
	out, err := s.fetch_all_policy()
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.(*iam.ListPoliciesOutput).Policies {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil
}

type Storage struct {
	region string
	s3iface.S3API
}

func NewStorage(sess *session.Session) *Storage {
	region := awssdk.StringValue(sess.Config.Region)
	return &Storage{S3API: s3.New(sess), region: region}
}

func (s *Storage) Name() string {
	return "storage"
}

func (s *Storage) Provider() string {
	return "aws"
}

func (s *Storage) ProviderAPI() string {
	return "s3"
}

func (s *Storage) ProviderRunnableAPI() interface{} {
	return s.S3API
}

func (s *Storage) ResourceTypes() (all []string) {
	all = append(all, "bucket")
	all = append(all, "storageobject")
	return
}

func (s *Storage) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	regionN := graph.InitResource(s.region, graph.Region)
	g.AddResource(regionN)
	var bucketList []*s3.Bucket
	var storageobjectList []*s3.Object

	errc := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, bucketList, err = s.fetch_all_bucket_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var resGraph *graph.Graph
		var err error
		resGraph, storageobjectList, err = s.fetch_all_storageobject_graph()
		if err != nil {
			errc <- err
			return
		}
		g.AddGraph(resGraph)
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case "Access Denied":
				return g, cloud.ErrFetchAccessDenied
			default:
				return g, ee
			}
		case nil:
			continue
		default:
			return g, ee
		}
	}

	errc = make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, r := range bucketList {
			for _, fn := range addParentsFns["bucket"] {
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
		for _, r := range storageobjectList {
			for _, fn := range addParentsFns["storageobject"] {
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

func (s *Storage) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "bucket":
		graph, _, err := s.fetch_all_bucket_graph()
		return graph, err
	case "storageobject":
		graph, _, err := s.fetch_all_storageobject_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws storage: unsupported fetch for type %s", t)
	}
}

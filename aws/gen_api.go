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
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	awsdriver "github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

const accessDenied = "Access Denied"

func init() {
	ServiceNames = append(ServiceNames, "infra")
	ServiceNames = append(ServiceNames, "access")
	ServiceNames = append(ServiceNames, "storage")
	ServiceNames = append(ServiceNames, "notification")
	ServiceNames = append(ServiceNames, "queue")
	ServiceNames = append(ServiceNames, "dns")
}

var ServiceNames = []string{}

var ResourceTypes = []string{
	"instance",
	"subnet",
	"vpc",
	"keypair",
	"securitygroup",
	"volume",
	"internetgateway",
	"routetable",
	"availabilityzone",
	"loadbalancer",
	"targetgroup",
	"listener",
	"user",
	"group",
	"role",
	"policy",
	"bucket",
	"storageobject",
	"subscription",
	"topic",
	"queue",
	"zone",
	"record",
}

var ServicePerAPI = map[string]string{
	"ec2":     "infra",
	"elbv2":   "infra",
	"iam":     "access",
	"s3":      "storage",
	"sns":     "notification",
	"sqs":     "queue",
	"route53": "dns",
}

var ServicePerResourceType = map[string]string{
	"instance":         "infra",
	"subnet":           "infra",
	"vpc":              "infra",
	"keypair":          "infra",
	"securitygroup":    "infra",
	"volume":           "infra",
	"internetgateway":  "infra",
	"routetable":       "infra",
	"availabilityzone": "infra",
	"loadbalancer":     "infra",
	"targetgroup":      "infra",
	"listener":         "infra",
	"user":             "access",
	"group":            "access",
	"role":             "access",
	"policy":           "access",
	"bucket":           "storage",
	"storageobject":    "storage",
	"subscription":     "notification",
	"topic":            "notification",
	"queue":            "queue",
	"zone":             "dns",
	"record":           "dns",
}

type Infra struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	ec2iface.EC2API
	elbv2iface.ELBV2API
}

func NewInfra(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Infra{
		EC2API:   ec2.New(sess),
		ELBV2API: elbv2.New(sess),
		config:   awsconf,
		region:   region,
		log:      log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewEc2Driver(s.EC2API),
		awsdriver.NewElbv2Driver(s.ELBV2API),
	}
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
	all = append(all, "availabilityzone")
	all = append(all, "loadbalancer")
	all = append(all, "targetgroup")
	all = append(all, "listener")
	return
}

func (s *Infra) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var instanceList []*ec2.Instance
	var subnetList []*ec2.Subnet
	var vpcList []*ec2.Vpc
	var keypairList []*ec2.KeyPairInfo
	var securitygroupList []*ec2.SecurityGroup
	var volumeList []*ec2.Volume
	var internetgatewayList []*ec2.InternetGateway
	var routetableList []*ec2.RouteTable
	var availabilityzoneList []*ec2.AvailabilityZone
	var loadbalancerList []*elbv2.LoadBalancer
	var targetgroupList []*elbv2.TargetGroup
	var listenerList []*elbv2.Listener

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.infra.instance.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[instance]")
	}
	if s.config.getBool("aws.infra.subnet.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[subnet]")
	}
	if s.config.getBool("aws.infra.vpc.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[vpc]")
	}
	if s.config.getBool("aws.infra.keypair.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[keypair]")
	}
	if s.config.getBool("aws.infra.securitygroup.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[securitygroup]")
	}
	if s.config.getBool("aws.infra.volume.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[volume]")
	}
	if s.config.getBool("aws.infra.internetgateway.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[internetgateway]")
	}
	if s.config.getBool("aws.infra.routetable.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[routetable]")
	}
	if s.config.getBool("aws.infra.availabilityzone.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, availabilityzoneList, err = s.fetch_all_availabilityzone_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[availabilityzone]")
	}
	if s.config.getBool("aws.infra.loadbalancer.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, loadbalancerList, err = s.fetch_all_loadbalancer_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[loadbalancer]")
	}
	if s.config.getBool("aws.infra.targetgroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, targetgroupList, err = s.fetch_all_targetgroup_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[targetgroup]")
	}
	if s.config.getBool("aws.infra.listener.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, listenerList, err = s.fetch_all_listener_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[listener]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.infra.instance.sync", true) {
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
	}
	if s.config.getBool("aws.infra.subnet.sync", true) {
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
	}
	if s.config.getBool("aws.infra.vpc.sync", true) {
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
	}
	if s.config.getBool("aws.infra.keypair.sync", true) {
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
	}
	if s.config.getBool("aws.infra.securitygroup.sync", true) {
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
	}
	if s.config.getBool("aws.infra.volume.sync", true) {
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
	}
	if s.config.getBool("aws.infra.internetgateway.sync", true) {
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
	}
	if s.config.getBool("aws.infra.routetable.sync", true) {
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
	}
	if s.config.getBool("aws.infra.availabilityzone.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range availabilityzoneList {
				for _, fn := range addParentsFns["availabilityzone"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.loadbalancer.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range loadbalancerList {
				for _, fn := range addParentsFns["loadbalancer"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.targetgroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range targetgroupList {
				for _, fn := range addParentsFns["targetgroup"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.listener.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range listenerList {
				for _, fn := range addParentsFns["listener"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}

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
	case "availabilityzone":
		graph, _, err := s.fetch_all_availabilityzone_graph()
		return graph, err
	case "loadbalancer":
		graph, _, err := s.fetch_all_loadbalancer_graph()
		return graph, err
	case "targetgroup":
		graph, _, err := s.fetch_all_targetgroup_graph()
		return graph, err
	case "listener":
		graph, _, err := s.fetch_all_listener_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws infra: unsupported fetch for type %s", t)
	}
}

func (s *Infra) fetch_all_instance_graph() (*graph.Graph, []*ec2.Instance, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Instance
	var badResErr error
	err := s.DescribeInstancesPages(&ec2.DescribeInstancesInput{},
		func(out *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool) {
			for _, all := range out.Reservations {
				for _, output := range all.Instances {
					cloudResources = append(cloudResources, output)
					var res *graph.Resource
					res, badResErr = newResource(output)
					if badResErr != nil {
						return false
					}
					g.AddResource(res)
				}
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_subnet_graph() (*graph.Graph, []*ec2.Subnet, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Subnet
	out, err := s.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.Subnets {
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
	out, err := s.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.Vpcs {
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
	out, err := s.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.KeyPairs {
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
	out, err := s.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.SecurityGroups {
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
	var badResErr error
	err := s.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
		func(out *ec2.DescribeVolumesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Volumes {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_internetgateway_graph() (*graph.Graph, []*ec2.InternetGateway, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.InternetGateway
	out, err := s.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.InternetGateways {
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
	out, err := s.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.RouteTables {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil

}

func (s *Infra) fetch_all_availabilityzone_graph() (*graph.Graph, []*ec2.AvailabilityZone, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.AvailabilityZone
	out, err := s.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.AvailabilityZones {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil

}

func (s *Infra) fetch_all_loadbalancer_graph() (*graph.Graph, []*elbv2.LoadBalancer, error) {
	g := graph.NewGraph()
	var cloudResources []*elbv2.LoadBalancer
	var badResErr error
	err := s.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
		func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.LoadBalancers {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.NextMarker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_targetgroup_graph() (*graph.Graph, []*elbv2.TargetGroup, error) {
	g := graph.NewGraph()
	var cloudResources []*elbv2.TargetGroup
	out, err := s.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.TargetGroups {
		cloudResources = append(cloudResources, output)
		res, err := newResource(output)
		if err != nil {
			return g, cloudResources, err
		}
		g.AddResource(res)
	}

	return g, cloudResources, nil

}

func (s *Infra) IsSyncDisabled() bool {
	return !s.config.getBool("aws.infra.sync", true)
}

type Access struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	iamiface.IAMAPI
}

func NewAccess(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Access{
		IAMAPI: iam.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewIamDriver(s.IAMAPI),
	}
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
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var userList []*iam.UserDetail
	var groupList []*iam.GroupDetail
	var roleList []*iam.RoleDetail
	var policyList []*iam.Policy

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.access.user.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource access[user]")
	}
	if s.config.getBool("aws.access.group.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource access[group]")
	}
	if s.config.getBool("aws.access.role.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource access[role]")
	}
	if s.config.getBool("aws.access.policy.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource access[policy]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.access.user.sync", true) {
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
	}
	if s.config.getBool("aws.access.group.sync", true) {
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
	}
	if s.config.getBool("aws.access.role.sync", true) {
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
	}
	if s.config.getBool("aws.access.policy.sync", true) {
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
	}

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

func (s *Access) fetch_all_group_graph() (*graph.Graph, []*iam.GroupDetail, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.GroupDetail
	var badResErr error
	err := s.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeGroup)}},
		func(out *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.GroupDetailList {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Access) fetch_all_role_graph() (*graph.Graph, []*iam.RoleDetail, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.RoleDetail
	var badResErr error
	err := s.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{Filter: []*string{awssdk.String(iam.EntityTypeRole)}},
		func(out *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.RoleDetailList {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Access) fetch_all_policy_graph() (*graph.Graph, []*iam.Policy, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.Policy
	var badResErr error
	err := s.ListPoliciesPages(&iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)},
		func(out *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Policies {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Access) IsSyncDisabled() bool {
	return !s.config.getBool("aws.access.sync", true)
}

type Storage struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	s3iface.S3API
}

func NewStorage(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Storage{
		S3API:  s3.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Storage) Name() string {
	return "storage"
}

func (s *Storage) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewS3Driver(s.S3API),
	}
}

func (s *Storage) ResourceTypes() (all []string) {
	all = append(all, "bucket")
	all = append(all, "storageobject")
	return
}

func (s *Storage) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var bucketList []*s3.Bucket
	var storageobjectList []*s3.Object

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.storage.bucket.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource storage[bucket]")
	}
	if s.config.getBool("aws.storage.storageobject.sync", true) {
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
	} else {
		s.log.Verbose("sync: *disabled* for resource storage[storageobject]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.storage.bucket.sync", true) {
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
	}
	if s.config.getBool("aws.storage.storageobject.sync", true) {
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
	}

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

func (s *Storage) IsSyncDisabled() bool {
	return !s.config.getBool("aws.storage.sync", true)
}

type Notification struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	snsiface.SNSAPI
}

func NewNotification(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Notification{
		SNSAPI: sns.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Notification) Name() string {
	return "notification"
}

func (s *Notification) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewSnsDriver(s.SNSAPI),
	}
}

func (s *Notification) ResourceTypes() (all []string) {
	all = append(all, "subscription")
	all = append(all, "topic")
	return
}

func (s *Notification) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var subscriptionList []*sns.Subscription
	var topicList []*sns.Topic

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.notification.subscription.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, subscriptionList, err = s.fetch_all_subscription_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource notification[subscription]")
	}
	if s.config.getBool("aws.notification.topic.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, topicList, err = s.fetch_all_topic_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource notification[topic]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.notification.subscription.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range subscriptionList {
				for _, fn := range addParentsFns["subscription"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.notification.topic.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range topicList {
				for _, fn := range addParentsFns["topic"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}

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

func (s *Notification) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "subscription":
		graph, _, err := s.fetch_all_subscription_graph()
		return graph, err
	case "topic":
		graph, _, err := s.fetch_all_topic_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws notification: unsupported fetch for type %s", t)
	}
}

func (s *Notification) fetch_all_subscription_graph() (*graph.Graph, []*sns.Subscription, error) {
	g := graph.NewGraph()
	var cloudResources []*sns.Subscription
	var badResErr error
	err := s.ListSubscriptionsPages(&sns.ListSubscriptionsInput{},
		func(out *sns.ListSubscriptionsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Subscriptions {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Notification) fetch_all_topic_graph() (*graph.Graph, []*sns.Topic, error) {
	g := graph.NewGraph()
	var cloudResources []*sns.Topic
	var badResErr error
	err := s.ListTopicsPages(&sns.ListTopicsInput{},
		func(out *sns.ListTopicsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Topics {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Notification) IsSyncDisabled() bool {
	return !s.config.getBool("aws.notification.sync", true)
}

type Queue struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	sqsiface.SQSAPI
}

func NewQueue(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Queue{
		SQSAPI: sqs.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Queue) Name() string {
	return "queue"
}

func (s *Queue) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewSqsDriver(s.SQSAPI),
	}
}

func (s *Queue) ResourceTypes() (all []string) {
	all = append(all, "queue")
	return
}

func (s *Queue) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var queueList []*string

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.queue.queue.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, queueList, err = s.fetch_all_queue_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource queue[queue]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.queue.queue.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range queueList {
				for _, fn := range addParentsFns["queue"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}

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

func (s *Queue) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "queue":
		graph, _, err := s.fetch_all_queue_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws queue: unsupported fetch for type %s", t)
	}
}

func (s *Queue) IsSyncDisabled() bool {
	return !s.config.getBool("aws.queue.sync", true)
}

type Dns struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	route53iface.Route53API
}

func NewDns(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Dns{
		Route53API: route53.New(sess),
		config:     awsconf,
		region:     region,
		log:        log,
	}
}

func (s *Dns) Name() string {
	return "dns"
}

func (s *Dns) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewRoute53Driver(s.Route53API),
	}
}

func (s *Dns) ResourceTypes() (all []string) {
	all = append(all, "zone")
	all = append(all, "record")
	return
}

func (s *Dns) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(s.region, cloud.Region)
	g.AddResource(regionN)
	var zoneList []*route53.HostedZone
	var recordList []*route53.ResourceRecordSet

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.dns.zone.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, zoneList, err = s.fetch_all_zone_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource dns[zone]")
	}
	if s.config.getBool("aws.dns.record.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, recordList, err = s.fetch_all_record_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource dns[record]")
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		switch ee := err.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
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
	if s.config.getBool("aws.dns.zone.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range zoneList {
				for _, fn := range addParentsFns["zone"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.dns.record.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range recordList {
				for _, fn := range addParentsFns["record"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}

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

func (s *Dns) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "zone":
		graph, _, err := s.fetch_all_zone_graph()
		return graph, err
	case "record":
		graph, _, err := s.fetch_all_record_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws dns: unsupported fetch for type %s", t)
	}
}

func (s *Dns) fetch_all_zone_graph() (*graph.Graph, []*route53.HostedZone, error) {
	g := graph.NewGraph()
	var cloudResources []*route53.HostedZone
	var badResErr error
	err := s.ListHostedZonesPages(&route53.ListHostedZonesInput{},
		func(out *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.HostedZones {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				res, badResErr = newResource(output)
				if badResErr != nil {
					return false
				}
				g.AddResource(res)
			}
			return out.NextMarker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Dns) IsSyncDisabled() bool {
	return !s.config.getBool("aws.dns.sync", true)
}

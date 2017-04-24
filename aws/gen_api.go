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
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	"infra",
	"access",
	"storage",
	"notification",
	"queue",
	"dns",
	"lambda",
}

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
	"database",
	"dbsubnetgroup",
	"launchconfiguration",
	"autoscalinggroup",
	"user",
	"group",
	"role",
	"policy",
	"accesskey",
	"bucket",
	"s3object",
	"subscription",
	"topic",
	"queue",
	"zone",
	"record",
	"function",
}

var ServicePerAPI = map[string]string{
	"ec2":         "infra",
	"elbv2":       "infra",
	"rds":         "infra",
	"autoscaling": "infra",
	"iam":         "access",
	"sts":         "access",
	"s3":          "storage",
	"sns":         "notification",
	"sqs":         "queue",
	"route53":     "dns",
	"lambda":      "lambda",
}

var ServicePerResourceType = map[string]string{
	"instance":            "infra",
	"subnet":              "infra",
	"vpc":                 "infra",
	"keypair":             "infra",
	"securitygroup":       "infra",
	"volume":              "infra",
	"internetgateway":     "infra",
	"routetable":          "infra",
	"availabilityzone":    "infra",
	"loadbalancer":        "infra",
	"targetgroup":         "infra",
	"listener":            "infra",
	"database":            "infra",
	"dbsubnetgroup":       "infra",
	"launchconfiguration": "infra",
	"autoscalinggroup":    "infra",
	"user":                "access",
	"group":               "access",
	"role":                "access",
	"policy":              "access",
	"accesskey":           "access",
	"bucket":              "storage",
	"s3object":            "storage",
	"subscription":        "notification",
	"topic":               "notification",
	"queue":               "queue",
	"zone":                "dns",
	"record":              "dns",
	"function":            "lambda",
}

var APIPerResourceType = map[string]string{
	"instance":            "ec2",
	"subnet":              "ec2",
	"vpc":                 "ec2",
	"keypair":             "ec2",
	"securitygroup":       "ec2",
	"volume":              "ec2",
	"internetgateway":     "ec2",
	"routetable":          "ec2",
	"availabilityzone":    "ec2",
	"loadbalancer":        "elbv2",
	"targetgroup":         "elbv2",
	"listener":            "elbv2",
	"database":            "rds",
	"dbsubnetgroup":       "rds",
	"launchconfiguration": "autoscaling",
	"autoscalinggroup":    "autoscaling",
	"user":                "iam",
	"group":               "iam",
	"role":                "iam",
	"policy":              "iam",
	"accesskey":           "iam",
	"bucket":              "s3",
	"s3object":            "s3",
	"subscription":        "sns",
	"topic":               "sns",
	"queue":               "sqs",
	"zone":                "route53",
	"record":              "route53",
	"function":            "lambda",
}

type Infra struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	ec2iface.EC2API
	elbv2iface.ELBV2API
	rdsiface.RDSAPI
	autoscalingiface.AutoScalingAPI
}

func NewInfra(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Infra{
		EC2API:         ec2.New(sess),
		ELBV2API:       elbv2.New(sess),
		RDSAPI:         rds.New(sess),
		AutoScalingAPI: autoscaling.New(sess),
		config:         awsconf,
		region:         region,
		log:            log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewEc2Driver(s.EC2API),
		awsdriver.NewElbv2Driver(s.ELBV2API),
		awsdriver.NewRdsDriver(s.RDSAPI),
		awsdriver.NewAutoscalingDriver(s.AutoScalingAPI),
	}
}

func (s *Infra) ResourceTypes() []string {
	return []string{
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
		"database",
		"dbsubnetgroup",
		"launchconfiguration",
		"autoscalinggroup",
	}
}

func (s *Infra) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
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
	var databaseList []*rds.DBInstance
	var dbsubnetgroupList []*rds.DBSubnetGroup
	var launchconfigurationList []*autoscaling.LaunchConfiguration
	var autoscalinggroupList []*autoscaling.Group

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
	if s.config.getBool("aws.infra.database.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, databaseList, err = s.fetch_all_database_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[database]")
	}
	if s.config.getBool("aws.infra.dbsubnetgroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, dbsubnetgroupList, err = s.fetch_all_dbsubnetgroup_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[dbsubnetgroup]")
	}
	if s.config.getBool("aws.infra.launchconfiguration.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, launchconfigurationList, err = s.fetch_all_launchconfiguration_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[launchconfiguration]")
	}
	if s.config.getBool("aws.infra.autoscalinggroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, autoscalinggroupList, err = s.fetch_all_autoscalinggroup_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[autoscalinggroup]")
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
	if s.config.getBool("aws.infra.database.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range databaseList {
				for _, fn := range addParentsFns["database"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.dbsubnetgroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range dbsubnetgroupList {
				for _, fn := range addParentsFns["dbsubnetgroup"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.launchconfiguration.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range launchconfigurationList {
				for _, fn := range addParentsFns["launchconfiguration"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.autoscalinggroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range autoscalinggroupList {
				for _, fn := range addParentsFns["autoscalinggroup"] {
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
	case "database":
		graph, _, err := s.fetch_all_database_graph()
		return graph, err
	case "dbsubnetgroup":
		graph, _, err := s.fetch_all_dbsubnetgroup_graph()
		return graph, err
	case "launchconfiguration":
		graph, _, err := s.fetch_all_launchconfiguration_graph()
		return graph, err
	case "autoscalinggroup":
		graph, _, err := s.fetch_all_autoscalinggroup_graph()
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
					if res, badResErr = newResource(output); badResErr != nil {
						return false
					}
					if badResErr = g.AddResource(res); badResErr != nil {
						return false
					}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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
		if err = g.AddResource(res); err != nil {
			return g, cloudResources, err
		}
	}

	return g, cloudResources, nil

}

func (s *Infra) fetch_all_database_graph() (*graph.Graph, []*rds.DBInstance, error) {
	g := graph.NewGraph()
	var cloudResources []*rds.DBInstance
	var badResErr error
	err := s.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{},
		func(out *rds.DescribeDBInstancesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.DBInstances {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_dbsubnetgroup_graph() (*graph.Graph, []*rds.DBSubnetGroup, error) {
	g := graph.NewGraph()
	var cloudResources []*rds.DBSubnetGroup
	var badResErr error
	err := s.DescribeDBSubnetGroupsPages(&rds.DescribeDBSubnetGroupsInput{},
		func(out *rds.DescribeDBSubnetGroupsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.DBSubnetGroups {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_launchconfiguration_graph() (*graph.Graph, []*autoscaling.LaunchConfiguration, error) {
	g := graph.NewGraph()
	var cloudResources []*autoscaling.LaunchConfiguration
	var badResErr error
	err := s.DescribeLaunchConfigurationsPages(&autoscaling.DescribeLaunchConfigurationsInput{},
		func(out *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.LaunchConfigurations {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Infra) fetch_all_autoscalinggroup_graph() (*graph.Graph, []*autoscaling.Group, error) {
	g := graph.NewGraph()
	var cloudResources []*autoscaling.Group
	var badResErr error
	err := s.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{},
		func(out *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.AutoScalingGroups {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.NextToken != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
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
	stsiface.STSAPI
}

func NewAccess(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Access{
		IAMAPI: iam.New(sess),
		STSAPI: sts.New(sess),
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
		awsdriver.NewStsDriver(s.STSAPI),
	}
}

func (s *Access) ResourceTypes() []string {
	return []string{
		"user",
		"group",
		"role",
		"policy",
		"accesskey",
	}
}

func (s *Access) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var userList []*iam.UserDetail
	var groupList []*iam.GroupDetail
	var roleList []*iam.RoleDetail
	var policyList []*iam.Policy
	var accesskeyList []*iam.AccessKeyMetadata

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
	if s.config.getBool("aws.access.accesskey.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, accesskeyList, err = s.fetch_all_accesskey_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource access[accesskey]")
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
	if s.config.getBool("aws.access.accesskey.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range accesskeyList {
				for _, fn := range addParentsFns["accesskey"] {
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
	case "accesskey":
		graph, _, err := s.fetch_all_accesskey_graph()
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.Marker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Access) fetch_all_accesskey_graph() (*graph.Graph, []*iam.AccessKeyMetadata, error) {
	g := graph.NewGraph()
	var cloudResources []*iam.AccessKeyMetadata
	var badResErr error
	err := s.ListAccessKeysPages(&iam.ListAccessKeysInput{},
		func(out *iam.ListAccessKeysOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.AccessKeyMetadata {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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

func (s *Storage) ResourceTypes() []string {
	return []string{
		"bucket",
		"s3object",
	}
}

func (s *Storage) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var bucketList []*s3.Bucket
	var s3objectList []*s3.Object

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
	if s.config.getBool("aws.storage.s3object.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, s3objectList, err = s.fetch_all_s3object_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource storage[s3object]")
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
	if s.config.getBool("aws.storage.s3object.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range s3objectList {
				for _, fn := range addParentsFns["s3object"] {
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
	case "s3object":
		graph, _, err := s.fetch_all_s3object_graph()
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

func (s *Notification) ResourceTypes() []string {
	return []string{
		"subscription",
		"topic",
	}
}

func (s *Notification) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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

func (s *Queue) ResourceTypes() []string {
	return []string{
		"queue",
	}
}

func (s *Queue) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
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

func (s *Dns) ResourceTypes() []string {
	return []string{
		"zone",
		"record",
	}
}

func (s *Dns) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
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
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
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

type Lambda struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	lambdaiface.LambdaAPI
}

func NewLambda(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Lambda{
		LambdaAPI: lambda.New(sess),
		config:    awsconf,
		region:    region,
		log:       log,
	}
}

func (s *Lambda) Name() string {
	return "lambda"
}

func (s *Lambda) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewLambdaDriver(s.LambdaAPI),
	}
}

func (s *Lambda) ResourceTypes() []string {
	return []string{
		"function",
	}
}

func (s *Lambda) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var functionList []*lambda.FunctionConfiguration

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.lambda.function.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, functionList, err = s.fetch_all_function_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource lambda[function]")
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
	if s.config.getBool("aws.lambda.function.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range functionList {
				for _, fn := range addParentsFns["function"] {
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

func (s *Lambda) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "function":
		graph, _, err := s.fetch_all_function_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws lambda: unsupported fetch for type %s", t)
	}
}

func (s *Lambda) fetch_all_function_graph() (*graph.Graph, []*lambda.FunctionConfiguration, error) {
	g := graph.NewGraph()
	var cloudResources []*lambda.FunctionConfiguration
	var badResErr error
	err := s.ListFunctionsPages(&lambda.ListFunctionsInput{},
		func(out *lambda.ListFunctionsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Functions {
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.NextMarker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Lambda) IsSyncDisabled() bool {
	return !s.config.getBool("aws.lambda.sync", true)
}

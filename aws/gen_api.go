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
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
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
	"messaging",
	"dns",
	"lambda",
	"monitoring",
	"cdn",
	"cloudformation",
}

var ResourceTypes = []string{
	"instance",
	"subnet",
	"vpc",
	"keypair",
	"securitygroup",
	"volume",
	"internetgateway",
	"natgateway",
	"routetable",
	"availabilityzone",
	"image",
	"importimagetask",
	"elasticip",
	"snapshot",
	"loadbalancer",
	"targetgroup",
	"listener",
	"database",
	"dbsubnetgroup",
	"launchconfiguration",
	"scalinggroup",
	"scalingpolicy",
	"repository",
	"containercluster",
	"containerservice",
	"container",
	"containerinstance",
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
	"metric",
	"alarm",
	"distribution",
	"stack",
}

var ServicePerAPI = map[string]string{
	"ec2":         "infra",
	"elbv2":       "infra",
	"rds":         "infra",
	"autoscaling": "infra",
	"ecr":         "infra",
	"ecs":         "infra",
	"applicationautoscaling": "infra",
	"iam":            "access",
	"sts":            "access",
	"s3":             "storage",
	"sns":            "messaging",
	"sqs":            "messaging",
	"route53":        "dns",
	"lambda":         "lambda",
	"cloudwatch":     "monitoring",
	"cloudfront":     "cdn",
	"cloudformation": "cloudformation",
}

var ServicePerResourceType = map[string]string{
	"instance":            "infra",
	"subnet":              "infra",
	"vpc":                 "infra",
	"keypair":             "infra",
	"securitygroup":       "infra",
	"volume":              "infra",
	"internetgateway":     "infra",
	"natgateway":          "infra",
	"routetable":          "infra",
	"availabilityzone":    "infra",
	"image":               "infra",
	"importimagetask":     "infra",
	"elasticip":           "infra",
	"snapshot":            "infra",
	"loadbalancer":        "infra",
	"targetgroup":         "infra",
	"listener":            "infra",
	"database":            "infra",
	"dbsubnetgroup":       "infra",
	"launchconfiguration": "infra",
	"scalinggroup":        "infra",
	"scalingpolicy":       "infra",
	"repository":          "infra",
	"containercluster":    "infra",
	"containerservice":    "infra",
	"container":           "infra",
	"containerinstance":   "infra",
	"user":                "access",
	"group":               "access",
	"role":                "access",
	"policy":              "access",
	"accesskey":           "access",
	"bucket":              "storage",
	"s3object":            "storage",
	"subscription":        "messaging",
	"topic":               "messaging",
	"queue":               "messaging",
	"zone":                "dns",
	"record":              "dns",
	"function":            "lambda",
	"metric":              "monitoring",
	"alarm":               "monitoring",
	"distribution":        "cdn",
	"stack":               "cloudformation",
}

var APIPerResourceType = map[string]string{
	"instance":            "ec2",
	"subnet":              "ec2",
	"vpc":                 "ec2",
	"keypair":             "ec2",
	"securitygroup":       "ec2",
	"volume":              "ec2",
	"internetgateway":     "ec2",
	"natgateway":          "ec2",
	"routetable":          "ec2",
	"availabilityzone":    "ec2",
	"image":               "ec2",
	"importimagetask":     "ec2",
	"elasticip":           "ec2",
	"snapshot":            "ec2",
	"loadbalancer":        "elbv2",
	"targetgroup":         "elbv2",
	"listener":            "elbv2",
	"database":            "rds",
	"dbsubnetgroup":       "rds",
	"launchconfiguration": "autoscaling",
	"scalinggroup":        "autoscaling",
	"scalingpolicy":       "autoscaling",
	"repository":          "ecr",
	"containercluster":    "ecs",
	"containerservice":    "ecs",
	"container":           "ecs",
	"containerinstance":   "ecs",
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
	"metric":              "cloudwatch",
	"alarm":               "cloudwatch",
	"distribution":        "cloudfront",
	"stack":               "cloudformation",
}

var GlobalServices = []string{
	"access",
	"dns",
	"cdn",
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
	ecriface.ECRAPI
	ecsiface.ECSAPI
	applicationautoscalingiface.ApplicationAutoScalingAPI
}

func NewInfra(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Infra{
		EC2API:         ec2.New(sess),
		ELBV2API:       elbv2.New(sess),
		RDSAPI:         rds.New(sess),
		AutoScalingAPI: autoscaling.New(sess),
		ECRAPI:         ecr.New(sess),
		ECSAPI:         ecs.New(sess),
		ApplicationAutoScalingAPI: applicationautoscaling.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Region() string {
	return s.region
}

func (s *Infra) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewEc2Driver(s.EC2API),
		awsdriver.NewElbv2Driver(s.ELBV2API),
		awsdriver.NewRdsDriver(s.RDSAPI),
		awsdriver.NewAutoscalingDriver(s.AutoScalingAPI),
		awsdriver.NewEcrDriver(s.ECRAPI),
		awsdriver.NewEcsDriver(s.ECSAPI),
		awsdriver.NewApplicationautoscalingDriver(s.ApplicationAutoScalingAPI),
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
		"natgateway",
		"routetable",
		"availabilityzone",
		"image",
		"importimagetask",
		"elasticip",
		"snapshot",
		"loadbalancer",
		"targetgroup",
		"listener",
		"database",
		"dbsubnetgroup",
		"launchconfiguration",
		"scalinggroup",
		"scalingpolicy",
		"repository",
		"containercluster",
		"containerservice",
		"container",
		"containerinstance",
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
	var natgatewayList []*ec2.NatGateway
	var routetableList []*ec2.RouteTable
	var availabilityzoneList []*ec2.AvailabilityZone
	var imageList []*ec2.Image
	var importimagetaskList []*ec2.ImportImageTask
	var elasticipList []*ec2.Address
	var snapshotList []*ec2.Snapshot
	var loadbalancerList []*elbv2.LoadBalancer
	var targetgroupList []*elbv2.TargetGroup
	var listenerList []*elbv2.Listener
	var databaseList []*rds.DBInstance
	var dbsubnetgroupList []*rds.DBSubnetGroup
	var launchconfigurationList []*autoscaling.LaunchConfiguration
	var scalinggroupList []*autoscaling.Group
	var scalingpolicyList []*autoscaling.ScalingPolicy
	var repositoryList []*ecr.Repository
	var containerclusterList []*ecs.Cluster
	var containerserviceList []*ecs.TaskDefinition
	var containerList []*ecs.Container
	var containerinstanceList []*ecs.ContainerInstance

	fetchError := new(multiError)

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
	if s.config.getBool("aws.infra.natgateway.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, natgatewayList, err = s.fetch_all_natgateway_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[natgateway]")
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
	if s.config.getBool("aws.infra.image.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, imageList, err = s.fetch_all_image_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[image]")
	}
	if s.config.getBool("aws.infra.importimagetask.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, importimagetaskList, err = s.fetch_all_importimagetask_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[importimagetask]")
	}
	if s.config.getBool("aws.infra.elasticip.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, elasticipList, err = s.fetch_all_elasticip_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[elasticip]")
	}
	if s.config.getBool("aws.infra.snapshot.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, snapshotList, err = s.fetch_all_snapshot_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[snapshot]")
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
	if s.config.getBool("aws.infra.scalinggroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, scalinggroupList, err = s.fetch_all_scalinggroup_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[scalinggroup]")
	}
	if s.config.getBool("aws.infra.scalingpolicy.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, scalingpolicyList, err = s.fetch_all_scalingpolicy_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[scalingpolicy]")
	}
	if s.config.getBool("aws.infra.repository.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, repositoryList, err = s.fetch_all_repository_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[repository]")
	}
	if s.config.getBool("aws.infra.containercluster.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, containerclusterList, err = s.fetch_all_containercluster_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[containercluster]")
	}
	if s.config.getBool("aws.infra.containerservice.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, containerserviceList, err = s.fetch_all_containerservice_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[containerservice]")
	}
	if s.config.getBool("aws.infra.container.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, containerList, err = s.fetch_all_container_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[container]")
	}
	if s.config.getBool("aws.infra.containerinstance.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, containerinstanceList, err = s.fetch_all_containerinstance_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource infra[containerinstance]")
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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
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
	if s.config.getBool("aws.infra.natgateway.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range natgatewayList {
				for _, fn := range addParentsFns["natgateway"] {
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
	if s.config.getBool("aws.infra.image.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range imageList {
				for _, fn := range addParentsFns["image"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.importimagetask.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range importimagetaskList {
				for _, fn := range addParentsFns["importimagetask"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.elasticip.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range elasticipList {
				for _, fn := range addParentsFns["elasticip"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.snapshot.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range snapshotList {
				for _, fn := range addParentsFns["snapshot"] {
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
	if s.config.getBool("aws.infra.scalinggroup.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range scalinggroupList {
				for _, fn := range addParentsFns["scalinggroup"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.scalingpolicy.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range scalingpolicyList {
				for _, fn := range addParentsFns["scalingpolicy"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.repository.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range repositoryList {
				for _, fn := range addParentsFns["repository"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.containercluster.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range containerclusterList {
				for _, fn := range addParentsFns["containercluster"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.containerservice.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range containerserviceList {
				for _, fn := range addParentsFns["containerservice"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.container.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range containerList {
				for _, fn := range addParentsFns["container"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.infra.containerinstance.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range containerinstanceList {
				for _, fn := range addParentsFns["containerinstance"] {
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
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
	case "natgateway":
		graph, _, err := s.fetch_all_natgateway_graph()
		return graph, err
	case "routetable":
		graph, _, err := s.fetch_all_routetable_graph()
		return graph, err
	case "availabilityzone":
		graph, _, err := s.fetch_all_availabilityzone_graph()
		return graph, err
	case "image":
		graph, _, err := s.fetch_all_image_graph()
		return graph, err
	case "importimagetask":
		graph, _, err := s.fetch_all_importimagetask_graph()
		return graph, err
	case "elasticip":
		graph, _, err := s.fetch_all_elasticip_graph()
		return graph, err
	case "snapshot":
		graph, _, err := s.fetch_all_snapshot_graph()
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
	case "scalinggroup":
		graph, _, err := s.fetch_all_scalinggroup_graph()
		return graph, err
	case "scalingpolicy":
		graph, _, err := s.fetch_all_scalingpolicy_graph()
		return graph, err
	case "repository":
		graph, _, err := s.fetch_all_repository_graph()
		return graph, err
	case "containercluster":
		graph, _, err := s.fetch_all_containercluster_graph()
		return graph, err
	case "containerservice":
		graph, _, err := s.fetch_all_containerservice_graph()
		return graph, err
	case "container":
		graph, _, err := s.fetch_all_container_graph()
		return graph, err
	case "containerinstance":
		graph, _, err := s.fetch_all_containerinstance_graph()
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
					if badResErr != nil {
						return false
					}
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

	out, err := s.EC2API.DescribeSubnets(&ec2.DescribeSubnetsInput{})
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

	out, err := s.EC2API.DescribeVpcs(&ec2.DescribeVpcsInput{})
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

	out, err := s.EC2API.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
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

	out, err := s.EC2API.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
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
				if badResErr != nil {
					return false
				}
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

	out, err := s.EC2API.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
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

func (s *Infra) fetch_all_natgateway_graph() (*graph.Graph, []*ec2.NatGateway, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.NatGateway

	out, err := s.EC2API.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.NatGateways {
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

	out, err := s.EC2API.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
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

	out, err := s.EC2API.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
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

func (s *Infra) fetch_all_image_graph() (*graph.Graph, []*ec2.Image, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Image

	out, err := s.EC2API.DescribeImages(&ec2.DescribeImagesInput{Owners: []*string{awssdk.String("self")}})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.Images {
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

func (s *Infra) fetch_all_importimagetask_graph() (*graph.Graph, []*ec2.ImportImageTask, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.ImportImageTask

	out, err := s.EC2API.DescribeImportImageTasks(&ec2.DescribeImportImageTasksInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.ImportImageTasks {
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

func (s *Infra) fetch_all_elasticip_graph() (*graph.Graph, []*ec2.Address, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Address

	out, err := s.EC2API.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, cloudResources, err
	}

	for _, output := range out.Addresses {
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

func (s *Infra) fetch_all_snapshot_graph() (*graph.Graph, []*ec2.Snapshot, error) {
	g := graph.NewGraph()
	var cloudResources []*ec2.Snapshot
	var badResErr error
	err := s.DescribeSnapshotsPages(&ec2.DescribeSnapshotsInput{OwnerIds: []*string{awssdk.String("self")}},
		func(out *ec2.DescribeSnapshotsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Snapshots {
				if badResErr != nil {
					return false
				}
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

func (s *Infra) fetch_all_loadbalancer_graph() (*graph.Graph, []*elbv2.LoadBalancer, error) {
	g := graph.NewGraph()
	var cloudResources []*elbv2.LoadBalancer
	var badResErr error
	err := s.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
		func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.LoadBalancers {
				if badResErr != nil {
					return false
				}
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

	out, err := s.ELBV2API.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
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
				if badResErr != nil {
					return false
				}
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
				if badResErr != nil {
					return false
				}
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
				if badResErr != nil {
					return false
				}
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

func (s *Infra) fetch_all_scalinggroup_graph() (*graph.Graph, []*autoscaling.Group, error) {
	g := graph.NewGraph()
	var cloudResources []*autoscaling.Group
	var badResErr error
	err := s.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{},
		func(out *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.AutoScalingGroups {
				if badResErr != nil {
					return false
				}
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

func (s *Infra) fetch_all_scalingpolicy_graph() (*graph.Graph, []*autoscaling.ScalingPolicy, error) {
	g := graph.NewGraph()
	var cloudResources []*autoscaling.ScalingPolicy
	var badResErr error
	err := s.DescribePoliciesPages(&autoscaling.DescribePoliciesInput{},
		func(out *autoscaling.DescribePoliciesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.ScalingPolicies {
				if badResErr != nil {
					return false
				}
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

func (s *Infra) fetch_all_repository_graph() (*graph.Graph, []*ecr.Repository, error) {
	g := graph.NewGraph()
	var cloudResources []*ecr.Repository
	var badResErr error
	err := s.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{},
		func(out *ecr.DescribeRepositoriesOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Repositories {
				if badResErr != nil {
					return false
				}
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
	region := "global"
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

func (s *Access) Region() string {
	return s.region
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

	fetchError := new(multiError)

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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
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
				if badResErr != nil {
					return false
				}
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
				if badResErr != nil {
					return false
				}
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
				if badResErr != nil {
					return false
				}
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

func (s *Storage) Region() string {
	return s.region
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

	fetchError := new(multiError)

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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
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

type Messaging struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	snsiface.SNSAPI
	sqsiface.SQSAPI
}

func NewMessaging(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Messaging{
		SNSAPI: sns.New(sess),
		SQSAPI: sqs.New(sess),
		config: awsconf,
		region: region,
		log:    log,
	}
}

func (s *Messaging) Name() string {
	return "messaging"
}

func (s *Messaging) Region() string {
	return s.region
}

func (s *Messaging) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewSnsDriver(s.SNSAPI),
		awsdriver.NewSqsDriver(s.SQSAPI),
	}
}

func (s *Messaging) ResourceTypes() []string {
	return []string{
		"subscription",
		"topic",
		"queue",
	}
}

func (s *Messaging) FetchResources() (*graph.Graph, error) {
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
	var queueList []*string

	fetchError := new(multiError)

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.messaging.subscription.sync", true) {
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
		s.log.Verbose("sync: *disabled* for resource messaging[subscription]")
	}
	if s.config.getBool("aws.messaging.topic.sync", true) {
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
		s.log.Verbose("sync: *disabled* for resource messaging[topic]")
	}
	if s.config.getBool("aws.messaging.queue.sync", true) {
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
		s.log.Verbose("sync: *disabled* for resource messaging[queue]")
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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
		}
	}

	errc = make(chan error)
	if s.config.getBool("aws.messaging.subscription.sync", true) {
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
	if s.config.getBool("aws.messaging.topic.sync", true) {
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
	if s.config.getBool("aws.messaging.queue.sync", true) {
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
	}

	return g, nil
}

func (s *Messaging) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "subscription":
		graph, _, err := s.fetch_all_subscription_graph()
		return graph, err
	case "topic":
		graph, _, err := s.fetch_all_topic_graph()
		return graph, err
	case "queue":
		graph, _, err := s.fetch_all_queue_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws messaging: unsupported fetch for type %s", t)
	}
}

func (s *Messaging) fetch_all_subscription_graph() (*graph.Graph, []*sns.Subscription, error) {
	g := graph.NewGraph()
	var cloudResources []*sns.Subscription
	var badResErr error
	err := s.ListSubscriptionsPages(&sns.ListSubscriptionsInput{},
		func(out *sns.ListSubscriptionsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Subscriptions {
				if badResErr != nil {
					return false
				}
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

func (s *Messaging) fetch_all_topic_graph() (*graph.Graph, []*sns.Topic, error) {
	g := graph.NewGraph()
	var cloudResources []*sns.Topic
	var badResErr error
	err := s.ListTopicsPages(&sns.ListTopicsInput{},
		func(out *sns.ListTopicsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Topics {
				if badResErr != nil {
					return false
				}
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

func (s *Messaging) IsSyncDisabled() bool {
	return !s.config.getBool("aws.messaging.sync", true)
}

type Dns struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	route53iface.Route53API
}

func NewDns(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := "global"
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

func (s *Dns) Region() string {
	return s.region
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

	fetchError := new(multiError)

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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
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
				if badResErr != nil {
					return false
				}
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

func (s *Lambda) Region() string {
	return s.region
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

	fetchError := new(multiError)

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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
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
				if badResErr != nil {
					return false
				}
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

type Monitoring struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	cloudwatchiface.CloudWatchAPI
}

func NewMonitoring(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Monitoring{
		CloudWatchAPI: cloudwatch.New(sess),
		config:        awsconf,
		region:        region,
		log:           log,
	}
}

func (s *Monitoring) Name() string {
	return "monitoring"
}

func (s *Monitoring) Region() string {
	return s.region
}

func (s *Monitoring) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudwatchDriver(s.CloudWatchAPI),
	}
}

func (s *Monitoring) ResourceTypes() []string {
	return []string{
		"metric",
		"alarm",
	}
}

func (s *Monitoring) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var metricList []*cloudwatch.Metric
	var alarmList []*cloudwatch.MetricAlarm

	fetchError := new(multiError)

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.monitoring.metric.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, metricList, err = s.fetch_all_metric_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource monitoring[metric]")
	}
	if s.config.getBool("aws.monitoring.alarm.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, alarmList, err = s.fetch_all_alarm_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource monitoring[alarm]")
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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
		}
	}

	errc = make(chan error)
	if s.config.getBool("aws.monitoring.metric.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range metricList {
				for _, fn := range addParentsFns["metric"] {
					err := fn(g, r)
					if err != nil {
						errc <- err
						return
					}
				}
			}
		}()
	}
	if s.config.getBool("aws.monitoring.alarm.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range alarmList {
				for _, fn := range addParentsFns["alarm"] {
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
	}

	return g, nil
}

func (s *Monitoring) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "metric":
		graph, _, err := s.fetch_all_metric_graph()
		return graph, err
	case "alarm":
		graph, _, err := s.fetch_all_alarm_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws monitoring: unsupported fetch for type %s", t)
	}
}

func (s *Monitoring) fetch_all_metric_graph() (*graph.Graph, []*cloudwatch.Metric, error) {
	g := graph.NewGraph()
	var cloudResources []*cloudwatch.Metric
	var badResErr error
	err := s.ListMetricsPages(&cloudwatch.ListMetricsInput{},
		func(out *cloudwatch.ListMetricsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Metrics {
				if badResErr != nil {
					return false
				}
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

func (s *Monitoring) fetch_all_alarm_graph() (*graph.Graph, []*cloudwatch.MetricAlarm, error) {
	g := graph.NewGraph()
	var cloudResources []*cloudwatch.MetricAlarm
	var badResErr error
	err := s.DescribeAlarmsPages(&cloudwatch.DescribeAlarmsInput{},
		func(out *cloudwatch.DescribeAlarmsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.MetricAlarms {
				if badResErr != nil {
					return false
				}
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

func (s *Monitoring) IsSyncDisabled() bool {
	return !s.config.getBool("aws.monitoring.sync", true)
}

type Cdn struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	cloudfrontiface.CloudFrontAPI
}

func NewCdn(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := "global"
	return &Cdn{
		CloudFrontAPI: cloudfront.New(sess),
		config:        awsconf,
		region:        region,
		log:           log,
	}
}

func (s *Cdn) Name() string {
	return "cdn"
}

func (s *Cdn) Region() string {
	return s.region
}

func (s *Cdn) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudfrontDriver(s.CloudFrontAPI),
	}
}

func (s *Cdn) ResourceTypes() []string {
	return []string{
		"distribution",
	}
}

func (s *Cdn) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var distributionList []*cloudfront.DistributionSummary

	fetchError := new(multiError)

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.cdn.distribution.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, distributionList, err = s.fetch_all_distribution_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource cdn[distribution]")
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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
		}
	}

	errc = make(chan error)
	if s.config.getBool("aws.cdn.distribution.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range distributionList {
				for _, fn := range addParentsFns["distribution"] {
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
	}

	return g, nil
}

func (s *Cdn) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "distribution":
		graph, _, err := s.fetch_all_distribution_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws cdn: unsupported fetch for type %s", t)
	}
}

func (s *Cdn) fetch_all_distribution_graph() (*graph.Graph, []*cloudfront.DistributionSummary, error) {
	g := graph.NewGraph()
	var cloudResources []*cloudfront.DistributionSummary
	var badResErr error
	err := s.ListDistributionsPages(&cloudfront.ListDistributionsInput{},
		func(out *cloudfront.ListDistributionsOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.DistributionList.Items {
				if badResErr != nil {
					return false
				}
				cloudResources = append(cloudResources, output)
				var res *graph.Resource
				if res, badResErr = newResource(output); badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return out.DistributionList.NextMarker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}

func (s *Cdn) IsSyncDisabled() bool {
	return !s.config.getBool("aws.cdn.sync", true)
}

type Cloudformation struct {
	once   oncer
	region string
	config config
	log    *logger.Logger
	cloudformationiface.CloudFormationAPI
}

func NewCloudformation(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	return &Cloudformation{
		CloudFormationAPI: cloudformation.New(sess),
		config:            awsconf,
		region:            region,
		log:               log,
	}
}

func (s *Cloudformation) Name() string {
	return "cloudformation"
}

func (s *Cloudformation) Region() string {
	return s.region
}

func (s *Cloudformation) Drivers() []driver.Driver {
	return []driver.Driver{
		awsdriver.NewCloudformationDriver(s.CloudFormationAPI),
	}
}

func (s *Cloudformation) ResourceTypes() []string {
	return []string{
		"stack",
	}
}

func (s *Cloudformation) FetchResources() (*graph.Graph, error) {
	g := graph.NewGraph()
	if s.IsSyncDisabled() {
		return g, nil
	}

	regionN := graph.InitResource(cloud.Region, s.region)
	if err := g.AddResource(regionN); err != nil {
		return g, err
	}
	var stackList []*cloudformation.Stack

	fetchError := new(multiError)

	errc := make(chan error)
	var wg sync.WaitGroup

	if s.config.getBool("aws.cloudformation.stack.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resGraph *graph.Graph
			var err error
			resGraph, stackList, err = s.fetch_all_stack_graph()
			if err != nil {
				errc <- err
				return
			}
			g.AddGraph(resGraph)
		}()
	} else {
		s.log.Verbose("sync: *disabled* for resource cloudformation[stack]")
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
				fetchError.add(cloud.ErrFetchAccessDenied)
			default:
				fetchError.add(ee)
			}
		case nil:
			continue
		default:
			fetchError.add(ee)
		}
	}

	errc = make(chan error)
	if s.config.getBool("aws.cloudformation.stack.sync", true) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, r := range stackList {
				for _, fn := range addParentsFns["stack"] {
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
			fetchError.add(err)
		}
	}

	if fetchError.hasAny() {
		return g, fetchError
	}

	return g, nil
}

func (s *Cloudformation) FetchByType(t string) (*graph.Graph, error) {
	switch t {
	case "stack":
		graph, _, err := s.fetch_all_stack_graph()
		return graph, err
	default:
		return nil, fmt.Errorf("aws cloudformation: unsupported fetch for type %s", t)
	}
}

func (s *Cloudformation) fetch_all_stack_graph() (*graph.Graph, []*cloudformation.Stack, error) {
	g := graph.NewGraph()
	var cloudResources []*cloudformation.Stack
	var badResErr error
	err := s.DescribeStacksPages(&cloudformation.DescribeStacksInput{},
		func(out *cloudformation.DescribeStacksOutput, lastPage bool) (shouldContinue bool) {
			for _, output := range out.Stacks {
				if badResErr != nil {
					return false
				}
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

func (s *Cloudformation) IsSyncDisabled() bool {
	return !s.config.getBool("aws.cloudformation.sync", true)
}

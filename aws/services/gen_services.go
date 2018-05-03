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

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"
	"errors"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
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
	"github.com/wallix/awless/aws/fetch"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	tstore "github.com/wallix/triplestore"
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
	"networkinterface",
	"classicloadbalancer",
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
	"containertask",
	"container",
	"containerinstance",
	"certificate",
	"user",
	"group",
	"role",
	"policy",
	"accesskey",
	"instanceprofile",
	"mfadevice",
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
	"elb":         "infra",
	"rds":         "infra",
	"autoscaling": "infra",
	"ecr":         "infra",
	"ecs":         "infra",
	"applicationautoscaling": "infra",
	"acm":            "infra",
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
	"networkinterface":    "infra",
	"classicloadbalancer": "infra",
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
	"containertask":       "infra",
	"container":           "infra",
	"containerinstance":   "infra",
	"certificate":         "infra",
	"user":                "access",
	"group":               "access",
	"role":                "access",
	"policy":              "access",
	"accesskey":           "access",
	"instanceprofile":     "access",
	"mfadevice":           "access",
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
	"networkinterface":    "ec2",
	"classicloadbalancer": "elb",
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
	"containertask":       "ecs",
	"container":           "ecs",
	"containerinstance":   "ecs",
	"certificate":         "acm",
	"user":                "iam",
	"group":               "iam",
	"role":                "iam",
	"policy":              "iam",
	"accesskey":           "iam",
	"instanceprofile":     "iam",
	"mfadevice":           "iam",
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

type Infra struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	ec2iface.EC2API
	elbv2iface.ELBV2API
	elbiface.ELBAPI
	rdsiface.RDSAPI
	autoscalingiface.AutoScalingAPI
	ecriface.ECRAPI
	ecsiface.ECSAPI
	applicationautoscalingiface.ApplicationAutoScalingAPI
	acmiface.ACMAPI
}

func NewInfra(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	ec2API := ec2.New(sess)
	elbv2API := elbv2.New(sess)
	elbAPI := elb.New(sess)
	rdsAPI := rds.New(sess)
	autoscalingAPI := autoscaling.New(sess)
	ecrAPI := ecr.New(sess)
	ecsAPI := ecs.New(sess)
	applicationautoscalingAPI := applicationautoscaling.New(sess)
	acmAPI := acm.New(sess)

	fetchConfig := awsfetch.NewConfig(
		ec2API,
		elbv2API,
		elbAPI,
		rdsAPI,
		autoscalingAPI,
		ecrAPI,
		ecsAPI,
		applicationautoscalingAPI,
		acmAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Infra{
		EC2API:         ec2API,
		ELBV2API:       elbv2API,
		ELBAPI:         elbAPI,
		RDSAPI:         rdsAPI,
		AutoScalingAPI: autoscalingAPI,
		ECRAPI:         ecrAPI,
		ECSAPI:         ecsAPI,
		ApplicationAutoScalingAPI: applicationautoscalingAPI,
		ACMAPI:  acmAPI,
		fetcher: fetch.NewFetcher(awsfetch.BuildInfraFetchFuncs(fetchConfig)),
		config:  extraConf,
		region:  region,
		profile: profile,
		log:     log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Region() string {
	return s.region
}

func (s *Infra) Profile() string {
	return s.profile
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
		"networkinterface",
		"classicloadbalancer",
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
		"containertask",
		"container",
		"containerinstance",
		"certificate",
	}
}

func (s *Infra) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.infra.instance.sync", true) {
		list, err := s.fetcher.Get("instance_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Instance); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Instance' type from fetch context")
		}
		for _, r := range list.([]*ec2.Instance) {
			for _, fn := range addParentsFns["instance"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Instance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.subnet.sync", true) {
		list, err := s.fetcher.Get("subnet_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Subnet); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Subnet' type from fetch context")
		}
		for _, r := range list.([]*ec2.Subnet) {
			for _, fn := range addParentsFns["subnet"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Subnet) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.vpc.sync", true) {
		list, err := s.fetcher.Get("vpc_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Vpc); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Vpc' type from fetch context")
		}
		for _, r := range list.([]*ec2.Vpc) {
			for _, fn := range addParentsFns["vpc"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Vpc) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.keypair.sync", true) {
		list, err := s.fetcher.Get("keypair_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.KeyPairInfo); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.KeyPairInfo' type from fetch context")
		}
		for _, r := range list.([]*ec2.KeyPairInfo) {
			for _, fn := range addParentsFns["keypair"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.KeyPairInfo) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.securitygroup.sync", true) {
		list, err := s.fetcher.Get("securitygroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.SecurityGroup); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.SecurityGroup' type from fetch context")
		}
		for _, r := range list.([]*ec2.SecurityGroup) {
			for _, fn := range addParentsFns["securitygroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.SecurityGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.volume.sync", true) {
		list, err := s.fetcher.Get("volume_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Volume); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Volume' type from fetch context")
		}
		for _, r := range list.([]*ec2.Volume) {
			for _, fn := range addParentsFns["volume"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Volume) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.internetgateway.sync", true) {
		list, err := s.fetcher.Get("internetgateway_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.InternetGateway); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.InternetGateway' type from fetch context")
		}
		for _, r := range list.([]*ec2.InternetGateway) {
			for _, fn := range addParentsFns["internetgateway"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.InternetGateway) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.natgateway.sync", true) {
		list, err := s.fetcher.Get("natgateway_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.NatGateway); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.NatGateway' type from fetch context")
		}
		for _, r := range list.([]*ec2.NatGateway) {
			for _, fn := range addParentsFns["natgateway"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.NatGateway) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.routetable.sync", true) {
		list, err := s.fetcher.Get("routetable_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.RouteTable); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.RouteTable' type from fetch context")
		}
		for _, r := range list.([]*ec2.RouteTable) {
			for _, fn := range addParentsFns["routetable"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.RouteTable) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.availabilityzone.sync", true) {
		list, err := s.fetcher.Get("availabilityzone_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.AvailabilityZone); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.AvailabilityZone' type from fetch context")
		}
		for _, r := range list.([]*ec2.AvailabilityZone) {
			for _, fn := range addParentsFns["availabilityzone"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.AvailabilityZone) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.image.sync", true) {
		list, err := s.fetcher.Get("image_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Image); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Image' type from fetch context")
		}
		for _, r := range list.([]*ec2.Image) {
			for _, fn := range addParentsFns["image"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Image) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.importimagetask.sync", true) {
		list, err := s.fetcher.Get("importimagetask_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.ImportImageTask); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.ImportImageTask' type from fetch context")
		}
		for _, r := range list.([]*ec2.ImportImageTask) {
			for _, fn := range addParentsFns["importimagetask"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.ImportImageTask) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.elasticip.sync", true) {
		list, err := s.fetcher.Get("elasticip_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Address); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Address' type from fetch context")
		}
		for _, r := range list.([]*ec2.Address) {
			for _, fn := range addParentsFns["elasticip"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Address) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.snapshot.sync", true) {
		list, err := s.fetcher.Get("snapshot_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.Snapshot); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.Snapshot' type from fetch context")
		}
		for _, r := range list.([]*ec2.Snapshot) {
			for _, fn := range addParentsFns["snapshot"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.Snapshot) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.networkinterface.sync", true) {
		list, err := s.fetcher.Get("networkinterface_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ec2.NetworkInterface); !ok {
			return gph, errors.New("cannot cast to '[]*ec2.NetworkInterface' type from fetch context")
		}
		for _, r := range list.([]*ec2.NetworkInterface) {
			for _, fn := range addParentsFns["networkinterface"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2.NetworkInterface) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.classicloadbalancer.sync", true) {
		list, err := s.fetcher.Get("classicloadbalancer_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*elb.LoadBalancerDescription); !ok {
			return gph, errors.New("cannot cast to '[]*elb.LoadBalancerDescription' type from fetch context")
		}
		for _, r := range list.([]*elb.LoadBalancerDescription) {
			for _, fn := range addParentsFns["classicloadbalancer"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elb.LoadBalancerDescription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.loadbalancer.sync", true) {
		list, err := s.fetcher.Get("loadbalancer_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*elbv2.LoadBalancer); !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.LoadBalancer' type from fetch context")
		}
		for _, r := range list.([]*elbv2.LoadBalancer) {
			for _, fn := range addParentsFns["loadbalancer"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2.LoadBalancer) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.targetgroup.sync", true) {
		list, err := s.fetcher.Get("targetgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*elbv2.TargetGroup); !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.TargetGroup' type from fetch context")
		}
		for _, r := range list.([]*elbv2.TargetGroup) {
			for _, fn := range addParentsFns["targetgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2.TargetGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.listener.sync", true) {
		list, err := s.fetcher.Get("listener_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*elbv2.Listener); !ok {
			return gph, errors.New("cannot cast to '[]*elbv2.Listener' type from fetch context")
		}
		for _, r := range list.([]*elbv2.Listener) {
			for _, fn := range addParentsFns["listener"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2.Listener) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.database.sync", true) {
		list, err := s.fetcher.Get("database_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*rds.DBInstance); !ok {
			return gph, errors.New("cannot cast to '[]*rds.DBInstance' type from fetch context")
		}
		for _, r := range list.([]*rds.DBInstance) {
			for _, fn := range addParentsFns["database"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *rds.DBInstance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.dbsubnetgroup.sync", true) {
		list, err := s.fetcher.Get("dbsubnetgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*rds.DBSubnetGroup); !ok {
			return gph, errors.New("cannot cast to '[]*rds.DBSubnetGroup' type from fetch context")
		}
		for _, r := range list.([]*rds.DBSubnetGroup) {
			for _, fn := range addParentsFns["dbsubnetgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *rds.DBSubnetGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.launchconfiguration.sync", true) {
		list, err := s.fetcher.Get("launchconfiguration_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*autoscaling.LaunchConfiguration); !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.LaunchConfiguration' type from fetch context")
		}
		for _, r := range list.([]*autoscaling.LaunchConfiguration) {
			for _, fn := range addParentsFns["launchconfiguration"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscaling.LaunchConfiguration) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.scalinggroup.sync", true) {
		list, err := s.fetcher.Get("scalinggroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*autoscaling.Group); !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.Group' type from fetch context")
		}
		for _, r := range list.([]*autoscaling.Group) {
			for _, fn := range addParentsFns["scalinggroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscaling.Group) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.scalingpolicy.sync", true) {
		list, err := s.fetcher.Get("scalingpolicy_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*autoscaling.ScalingPolicy); !ok {
			return gph, errors.New("cannot cast to '[]*autoscaling.ScalingPolicy' type from fetch context")
		}
		for _, r := range list.([]*autoscaling.ScalingPolicy) {
			for _, fn := range addParentsFns["scalingpolicy"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscaling.ScalingPolicy) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.repository.sync", true) {
		list, err := s.fetcher.Get("repository_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ecr.Repository); !ok {
			return gph, errors.New("cannot cast to '[]*ecr.Repository' type from fetch context")
		}
		for _, r := range list.([]*ecr.Repository) {
			for _, fn := range addParentsFns["repository"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecr.Repository) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containercluster.sync", true) {
		list, err := s.fetcher.Get("containercluster_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ecs.Cluster); !ok {
			return gph, errors.New("cannot cast to '[]*ecs.Cluster' type from fetch context")
		}
		for _, r := range list.([]*ecs.Cluster) {
			for _, fn := range addParentsFns["containercluster"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecs.Cluster) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containertask.sync", true) {
		list, err := s.fetcher.Get("containertask_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ecs.TaskDefinition); !ok {
			return gph, errors.New("cannot cast to '[]*ecs.TaskDefinition' type from fetch context")
		}
		for _, r := range list.([]*ecs.TaskDefinition) {
			for _, fn := range addParentsFns["containertask"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecs.TaskDefinition) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.container.sync", true) {
		list, err := s.fetcher.Get("container_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ecs.Container); !ok {
			return gph, errors.New("cannot cast to '[]*ecs.Container' type from fetch context")
		}
		for _, r := range list.([]*ecs.Container) {
			for _, fn := range addParentsFns["container"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecs.Container) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containerinstance.sync", true) {
		list, err := s.fetcher.Get("containerinstance_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*ecs.ContainerInstance); !ok {
			return gph, errors.New("cannot cast to '[]*ecs.ContainerInstance' type from fetch context")
		}
		for _, r := range list.([]*ecs.ContainerInstance) {
			for _, fn := range addParentsFns["containerinstance"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecs.ContainerInstance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.infra.certificate.sync", true) {
		list, err := s.fetcher.Get("certificate_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*acm.CertificateSummary); !ok {
			return gph, errors.New("cannot cast to '[]*acm.CertificateSummary' type from fetch context")
		}
		for _, r := range list.([]*acm.CertificateSummary) {
			for _, fn := range addParentsFns["certificate"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *acm.CertificateSummary) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Infra) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Infra) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.infra.sync", true)
}

type Access struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	iamiface.IAMAPI
	stsiface.STSAPI
}

func NewAccess(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := "global"
	iamAPI := iam.New(sess)
	stsAPI := sts.New(sess)

	fetchConfig := awsfetch.NewConfig(
		iamAPI,
		stsAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Access{
		IAMAPI:  iamAPI,
		STSAPI:  stsAPI,
		fetcher: fetch.NewFetcher(awsfetch.BuildAccessFetchFuncs(fetchConfig)),
		config:  extraConf,
		region:  region,
		profile: profile,
		log:     log,
	}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Region() string {
	return s.region
}

func (s *Access) Profile() string {
	return s.profile
}

func (s *Access) ResourceTypes() []string {
	return []string{
		"user",
		"group",
		"role",
		"policy",
		"accesskey",
		"instanceprofile",
		"mfadevice",
	}
}

func (s *Access) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.access.user.sync", true) {
		list, err := s.fetcher.Get("user_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.UserDetail); !ok {
			return gph, errors.New("cannot cast to '[]*iam.UserDetail' type from fetch context")
		}
		for _, r := range list.([]*iam.UserDetail) {
			for _, fn := range addParentsFns["user"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.UserDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.group.sync", true) {
		list, err := s.fetcher.Get("group_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.GroupDetail); !ok {
			return gph, errors.New("cannot cast to '[]*iam.GroupDetail' type from fetch context")
		}
		for _, r := range list.([]*iam.GroupDetail) {
			for _, fn := range addParentsFns["group"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.GroupDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.role.sync", true) {
		list, err := s.fetcher.Get("role_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.RoleDetail); !ok {
			return gph, errors.New("cannot cast to '[]*iam.RoleDetail' type from fetch context")
		}
		for _, r := range list.([]*iam.RoleDetail) {
			for _, fn := range addParentsFns["role"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.RoleDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.policy.sync", true) {
		list, err := s.fetcher.Get("policy_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.Policy); !ok {
			return gph, errors.New("cannot cast to '[]*iam.Policy' type from fetch context")
		}
		for _, r := range list.([]*iam.Policy) {
			for _, fn := range addParentsFns["policy"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.Policy) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.accesskey.sync", true) {
		list, err := s.fetcher.Get("accesskey_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.AccessKeyMetadata); !ok {
			return gph, errors.New("cannot cast to '[]*iam.AccessKeyMetadata' type from fetch context")
		}
		for _, r := range list.([]*iam.AccessKeyMetadata) {
			for _, fn := range addParentsFns["accesskey"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.AccessKeyMetadata) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.instanceprofile.sync", true) {
		list, err := s.fetcher.Get("instanceprofile_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.InstanceProfile); !ok {
			return gph, errors.New("cannot cast to '[]*iam.InstanceProfile' type from fetch context")
		}
		for _, r := range list.([]*iam.InstanceProfile) {
			for _, fn := range addParentsFns["instanceprofile"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.InstanceProfile) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.access.mfadevice.sync", true) {
		list, err := s.fetcher.Get("mfadevice_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*iam.VirtualMFADevice); !ok {
			return gph, errors.New("cannot cast to '[]*iam.VirtualMFADevice' type from fetch context")
		}
		for _, r := range list.([]*iam.VirtualMFADevice) {
			for _, fn := range addParentsFns["mfadevice"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iam.VirtualMFADevice) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Access) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Access) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.access.sync", true)
}

type Storage struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	s3iface.S3API
}

func NewStorage(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	s3API := s3.New(sess)

	fetchConfig := awsfetch.NewConfig(
		s3API,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Storage{
		S3API:   s3API,
		fetcher: fetch.NewFetcher(awsfetch.BuildStorageFetchFuncs(fetchConfig)),
		config:  extraConf,
		region:  region,
		profile: profile,
		log:     log,
	}
}

func (s *Storage) Name() string {
	return "storage"
}

func (s *Storage) Region() string {
	return s.region
}

func (s *Storage) Profile() string {
	return s.profile
}

func (s *Storage) ResourceTypes() []string {
	return []string{
		"bucket",
		"s3object",
	}
}

func (s *Storage) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.storage.bucket.sync", true) {
		list, err := s.fetcher.Get("bucket_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*s3.Bucket); !ok {
			return gph, errors.New("cannot cast to '[]*s3.Bucket' type from fetch context")
		}
		for _, r := range list.([]*s3.Bucket) {
			for _, fn := range addParentsFns["bucket"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *s3.Bucket) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.storage.s3object.sync", true) {
		list, err := s.fetcher.Get("s3object_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*s3.Object); !ok {
			return gph, errors.New("cannot cast to '[]*s3.Object' type from fetch context")
		}
		for _, r := range list.([]*s3.Object) {
			for _, fn := range addParentsFns["s3object"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *s3.Object) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Storage) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Storage) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.storage.sync", true)
}

type Messaging struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	snsiface.SNSAPI
	sqsiface.SQSAPI
}

func NewMessaging(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	snsAPI := sns.New(sess)
	sqsAPI := sqs.New(sess)

	fetchConfig := awsfetch.NewConfig(
		snsAPI,
		sqsAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Messaging{
		SNSAPI:  snsAPI,
		SQSAPI:  sqsAPI,
		fetcher: fetch.NewFetcher(awsfetch.BuildMessagingFetchFuncs(fetchConfig)),
		config:  extraConf,
		region:  region,
		profile: profile,
		log:     log,
	}
}

func (s *Messaging) Name() string {
	return "messaging"
}

func (s *Messaging) Region() string {
	return s.region
}

func (s *Messaging) Profile() string {
	return s.profile
}

func (s *Messaging) ResourceTypes() []string {
	return []string{
		"subscription",
		"topic",
		"queue",
	}
}

func (s *Messaging) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.messaging.subscription.sync", true) {
		list, err := s.fetcher.Get("subscription_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*sns.Subscription); !ok {
			return gph, errors.New("cannot cast to '[]*sns.Subscription' type from fetch context")
		}
		for _, r := range list.([]*sns.Subscription) {
			for _, fn := range addParentsFns["subscription"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *sns.Subscription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.messaging.topic.sync", true) {
		list, err := s.fetcher.Get("topic_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*sns.Topic); !ok {
			return gph, errors.New("cannot cast to '[]*sns.Topic' type from fetch context")
		}
		for _, r := range list.([]*sns.Topic) {
			for _, fn := range addParentsFns["topic"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *sns.Topic) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.messaging.queue.sync", true) {
		list, err := s.fetcher.Get("queue_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*string); !ok {
			return gph, errors.New("cannot cast to '[]*string' type from fetch context")
		}
		for _, r := range list.([]*string) {
			for _, fn := range addParentsFns["queue"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *string) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Messaging) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Messaging) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.messaging.sync", true)
}

type Dns struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	route53iface.Route53API
}

func NewDns(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := "global"
	route53API := route53.New(sess)

	fetchConfig := awsfetch.NewConfig(
		route53API,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Dns{
		Route53API: route53API,
		fetcher:    fetch.NewFetcher(awsfetch.BuildDnsFetchFuncs(fetchConfig)),
		config:     extraConf,
		region:     region,
		profile:    profile,
		log:        log,
	}
}

func (s *Dns) Name() string {
	return "dns"
}

func (s *Dns) Region() string {
	return s.region
}

func (s *Dns) Profile() string {
	return s.profile
}

func (s *Dns) ResourceTypes() []string {
	return []string{
		"zone",
		"record",
	}
}

func (s *Dns) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.dns.zone.sync", true) {
		list, err := s.fetcher.Get("zone_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*route53.HostedZone); !ok {
			return gph, errors.New("cannot cast to '[]*route53.HostedZone' type from fetch context")
		}
		for _, r := range list.([]*route53.HostedZone) {
			for _, fn := range addParentsFns["zone"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *route53.HostedZone) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.dns.record.sync", true) {
		list, err := s.fetcher.Get("record_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*route53.ResourceRecordSet); !ok {
			return gph, errors.New("cannot cast to '[]*route53.ResourceRecordSet' type from fetch context")
		}
		for _, r := range list.([]*route53.ResourceRecordSet) {
			for _, fn := range addParentsFns["record"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *route53.ResourceRecordSet) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Dns) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Dns) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.dns.sync", true)
}

type Lambda struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	lambdaiface.LambdaAPI
}

func NewLambda(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	lambdaAPI := lambda.New(sess)

	fetchConfig := awsfetch.NewConfig(
		lambdaAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Lambda{
		LambdaAPI: lambdaAPI,
		fetcher:   fetch.NewFetcher(awsfetch.BuildLambdaFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *Lambda) Name() string {
	return "lambda"
}

func (s *Lambda) Region() string {
	return s.region
}

func (s *Lambda) Profile() string {
	return s.profile
}

func (s *Lambda) ResourceTypes() []string {
	return []string{
		"function",
	}
}

func (s *Lambda) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.lambda.function.sync", true) {
		list, err := s.fetcher.Get("function_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*lambda.FunctionConfiguration); !ok {
			return gph, errors.New("cannot cast to '[]*lambda.FunctionConfiguration' type from fetch context")
		}
		for _, r := range list.([]*lambda.FunctionConfiguration) {
			for _, fn := range addParentsFns["function"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *lambda.FunctionConfiguration) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Lambda) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Lambda) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.lambda.sync", true)
}

type Monitoring struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	cloudwatchiface.CloudWatchAPI
}

func NewMonitoring(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	cloudwatchAPI := cloudwatch.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudwatchAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Monitoring{
		CloudWatchAPI: cloudwatchAPI,
		fetcher:       fetch.NewFetcher(awsfetch.BuildMonitoringFetchFuncs(fetchConfig)),
		config:        extraConf,
		region:        region,
		profile:       profile,
		log:           log,
	}
}

func (s *Monitoring) Name() string {
	return "monitoring"
}

func (s *Monitoring) Region() string {
	return s.region
}

func (s *Monitoring) Profile() string {
	return s.profile
}

func (s *Monitoring) ResourceTypes() []string {
	return []string{
		"metric",
		"alarm",
	}
}

func (s *Monitoring) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.monitoring.metric.sync", true) {
		list, err := s.fetcher.Get("metric_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*cloudwatch.Metric); !ok {
			return gph, errors.New("cannot cast to '[]*cloudwatch.Metric' type from fetch context")
		}
		for _, r := range list.([]*cloudwatch.Metric) {
			for _, fn := range addParentsFns["metric"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudwatch.Metric) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}
	if getBool(s.config, "aws.monitoring.alarm.sync", true) {
		list, err := s.fetcher.Get("alarm_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*cloudwatch.MetricAlarm); !ok {
			return gph, errors.New("cannot cast to '[]*cloudwatch.MetricAlarm' type from fetch context")
		}
		for _, r := range list.([]*cloudwatch.MetricAlarm) {
			for _, fn := range addParentsFns["alarm"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudwatch.MetricAlarm) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Monitoring) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Monitoring) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.monitoring.sync", true)
}

type Cdn struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	cloudfrontiface.CloudFrontAPI
}

func NewCdn(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := "global"
	cloudfrontAPI := cloudfront.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudfrontAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Cdn{
		CloudFrontAPI: cloudfrontAPI,
		fetcher:       fetch.NewFetcher(awsfetch.BuildCdnFetchFuncs(fetchConfig)),
		config:        extraConf,
		region:        region,
		profile:       profile,
		log:           log,
	}
}

func (s *Cdn) Name() string {
	return "cdn"
}

func (s *Cdn) Region() string {
	return s.region
}

func (s *Cdn) Profile() string {
	return s.profile
}

func (s *Cdn) ResourceTypes() []string {
	return []string{
		"distribution",
	}
}

func (s *Cdn) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cdn.distribution.sync", true) {
		list, err := s.fetcher.Get("distribution_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*cloudfront.DistributionSummary); !ok {
			return gph, errors.New("cannot cast to '[]*cloudfront.DistributionSummary' type from fetch context")
		}
		for _, r := range list.([]*cloudfront.DistributionSummary) {
			for _, fn := range addParentsFns["distribution"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudfront.DistributionSummary) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cdn) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Cdn) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cdn.sync", true)
}

type Cloudformation struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]interface{}
	log             *logger.Logger
	cloudformationiface.CloudFormationAPI
}

func NewCloudformation(sess *session.Session, profile string, extraConf map[string]interface{}, log *logger.Logger) cloud.Service {
	region := awssdk.StringValue(sess.Config.Region)
	cloudformationAPI := cloudformation.New(sess)

	fetchConfig := awsfetch.NewConfig(
		cloudformationAPI,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Cloudformation{
		CloudFormationAPI: cloudformationAPI,
		fetcher:           fetch.NewFetcher(awsfetch.BuildCloudformationFetchFuncs(fetchConfig)),
		config:            extraConf,
		region:            region,
		profile:           profile,
		log:               log,
	}
}

func (s *Cloudformation) Name() string {
	return "cloudformation"
}

func (s *Cloudformation) Region() string {
	return s.region
}

func (s *Cloudformation) Profile() string {
	return s.profile
}

func (s *Cloudformation) ResourceTypes() []string {
	return []string{
		"stack",
	}
}

func (s *Cloudformation) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cloudformation.stack.sync", true) {
		list, err := s.fetcher.Get("stack_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]*cloudformation.Stack); !ok {
			return gph, errors.New("cannot cast to '[]*cloudformation.Stack' type from fetch context")
		}
		for _, r := range list.([]*cloudformation.Stack) {
			for _, fn := range addParentsFns["stack"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudformation.Stack) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cloudformation) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Cloudformation) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cloudformation.sync", true)
}

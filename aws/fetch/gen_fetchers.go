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

package awsfetch

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wallix/awless/aws/conv"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
)

func BuildInfraFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualInfraFetchFuncs(conf, funcs)

	funcs["instance"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Instance

		if !conf.getBoolDefaultTrue("aws.infra.instance.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[instance]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Ec2.DescribeInstancesPages(&ec2.DescribeInstancesInput{},
			func(out *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool) {
				for _, all := range out.Reservations {
					for _, output := range all.Instances {
						if badResErr != nil {
							return false
						}
						objects = append(objects, output)
						var res *graph.Resource
						if res, badResErr = awsconv.NewResource(output); badResErr != nil {
							return false
						}
						resources = append(resources, res)
					}
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["subnet"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Subnet

		if !conf.getBoolDefaultTrue("aws.infra.subnet.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[subnet]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeSubnets(&ec2.DescribeSubnetsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Subnets {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["vpc"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Vpc

		if !conf.getBoolDefaultTrue("aws.infra.vpc.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[vpc]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeVpcs(&ec2.DescribeVpcsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Vpcs {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["keypair"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.KeyPairInfo

		if !conf.getBoolDefaultTrue("aws.infra.keypair.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[keypair]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.KeyPairs {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["securitygroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.SecurityGroup

		if !conf.getBoolDefaultTrue("aws.infra.securitygroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[securitygroup]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.SecurityGroups {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["volume"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Volume

		if !conf.getBoolDefaultTrue("aws.infra.volume.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[volume]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Ec2.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
			func(out *ec2.DescribeVolumesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Volumes {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["internetgateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.InternetGateway

		if !conf.getBoolDefaultTrue("aws.infra.internetgateway.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[internetgateway]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.InternetGateways {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["natgateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.NatGateway

		if !conf.getBoolDefaultTrue("aws.infra.natgateway.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[natgateway]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.NatGateways {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["routetable"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.RouteTable

		if !conf.getBoolDefaultTrue("aws.infra.routetable.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[routetable]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeRouteTables(&ec2.DescribeRouteTablesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.RouteTables {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["availabilityzone"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.AvailabilityZone

		if !conf.getBoolDefaultTrue("aws.infra.availabilityzone.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[availabilityzone]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.AvailabilityZones {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["image"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Image

		if !conf.getBoolDefaultTrue("aws.infra.image.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[image]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeImages(&ec2.DescribeImagesInput{Owners: []*string{awssdk.String("self")}})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Images {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["importimagetask"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.ImportImageTask

		if !conf.getBoolDefaultTrue("aws.infra.importimagetask.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[importimagetask]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeImportImageTasks(&ec2.DescribeImportImageTasksInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.ImportImageTasks {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["elasticip"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Address

		if !conf.getBoolDefaultTrue("aws.infra.elasticip.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[elasticip]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeAddresses(&ec2.DescribeAddressesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.Addresses {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["snapshot"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.Snapshot

		if !conf.getBoolDefaultTrue("aws.infra.snapshot.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[snapshot]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Ec2.DescribeSnapshotsPages(&ec2.DescribeSnapshotsInput{OwnerIds: []*string{awssdk.String("self")}},
			func(out *ec2.DescribeSnapshotsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Snapshots {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["networkinterface"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ec2.NetworkInterface

		if !conf.getBoolDefaultTrue("aws.infra.networkinterface.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[networkinterface]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Ec2.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.NetworkInterfaces {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["classicloadbalancer"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*elb.LoadBalancerDescription

		if !conf.getBoolDefaultTrue("aws.infra.classicloadbalancer.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[classicloadbalancer]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Elb.DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{},
			func(out *elb.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.LoadBalancerDescriptions {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["loadbalancer"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*elbv2.LoadBalancer

		if !conf.getBoolDefaultTrue("aws.infra.loadbalancer.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[loadbalancer]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Elbv2.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
			func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.LoadBalancers {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["targetgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*elbv2.TargetGroup

		if !conf.getBoolDefaultTrue("aws.infra.targetgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[targetgroup]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Elbv2.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.TargetGroups {
			objects = append(objects, output)
			res, err := awsconv.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}

		return resources, objects, nil
	}

	funcs["database"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*rds.DBInstance

		if !conf.getBoolDefaultTrue("aws.infra.database.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[database]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Rds.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{},
			func(out *rds.DescribeDBInstancesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.DBInstances {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.Marker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["dbsubnetgroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*rds.DBSubnetGroup

		if !conf.getBoolDefaultTrue("aws.infra.dbsubnetgroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[dbsubnetgroup]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Rds.DescribeDBSubnetGroupsPages(&rds.DescribeDBSubnetGroupsInput{},
			func(out *rds.DescribeDBSubnetGroupsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.DBSubnetGroups {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.Marker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["launchconfiguration"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*autoscaling.LaunchConfiguration

		if !conf.getBoolDefaultTrue("aws.infra.launchconfiguration.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[launchconfiguration]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Autoscaling.DescribeLaunchConfigurationsPages(&autoscaling.DescribeLaunchConfigurationsInput{},
			func(out *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.LaunchConfigurations {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["scalinggroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*autoscaling.Group

		if !conf.getBoolDefaultTrue("aws.infra.scalinggroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[scalinggroup]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Autoscaling.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{},
			func(out *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.AutoScalingGroups {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["scalingpolicy"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*autoscaling.ScalingPolicy

		if !conf.getBoolDefaultTrue("aws.infra.scalingpolicy.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[scalingpolicy]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Autoscaling.DescribePoliciesPages(&autoscaling.DescribePoliciesInput{},
			func(out *autoscaling.DescribePoliciesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.ScalingPolicies {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["repository"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*ecr.Repository

		if !conf.getBoolDefaultTrue("aws.infra.repository.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[repository]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Ecr.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{},
			func(out *ecr.DescribeRepositoriesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Repositories {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["certificate"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*acm.CertificateSummary

		if !conf.getBoolDefaultTrue("aws.infra.certificate.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[certificate]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Acm.ListCertificatesPages(&acm.ListCertificatesInput{},
			func(out *acm.ListCertificatesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.CertificateSummaryList {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildAccessFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualAccessFetchFuncs(conf, funcs)

	funcs["instanceprofile"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.InstanceProfile

		if !conf.getBoolDefaultTrue("aws.access.instanceprofile.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[instanceprofile]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Iam.ListInstanceProfilesPages(&iam.ListInstanceProfilesInput{},
			func(out *iam.ListInstanceProfilesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.InstanceProfiles {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.Marker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["mfadevice"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*iam.VirtualMFADevice

		if !conf.getBoolDefaultTrue("aws.access.mfadevice.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[mfadevice]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Iam.ListVirtualMFADevicesPages(&iam.ListVirtualMFADevicesInput{},
			func(out *iam.ListVirtualMFADevicesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.VirtualMFADevices {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.Marker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildStorageFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualStorageFetchFuncs(conf, funcs)
	return funcs
}
func BuildMessagingFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualMessagingFetchFuncs(conf, funcs)

	funcs["subscription"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*sns.Subscription

		if !conf.getBoolDefaultTrue("aws.messaging.subscription.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[subscription]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Sns.ListSubscriptionsPages(&sns.ListSubscriptionsInput{},
			func(out *sns.ListSubscriptionsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Subscriptions {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["topic"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*sns.Topic

		if !conf.getBoolDefaultTrue("aws.messaging.topic.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[topic]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Sns.ListTopicsPages(&sns.ListTopicsInput{},
			func(out *sns.ListTopicsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Topics {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildDnsFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualDnsFetchFuncs(conf, funcs)

	funcs["zone"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*route53.HostedZone

		if !conf.getBoolDefaultTrue("aws.dns.zone.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource dns[zone]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Route53.ListHostedZonesPages(&route53.ListHostedZonesInput{},
			func(out *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.HostedZones {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildLambdaFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualLambdaFetchFuncs(conf, funcs)

	funcs["function"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*lambda.FunctionConfiguration

		if !conf.getBoolDefaultTrue("aws.lambda.function.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource lambda[function]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Lambda.ListFunctionsPages(&lambda.ListFunctionsInput{},
			func(out *lambda.ListFunctionsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Functions {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildMonitoringFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualMonitoringFetchFuncs(conf, funcs)

	funcs["metric"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*cloudwatch.Metric

		if !conf.getBoolDefaultTrue("aws.monitoring.metric.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource monitoring[metric]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Cloudwatch.ListMetricsPages(&cloudwatch.ListMetricsInput{},
			func(out *cloudwatch.ListMetricsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Metrics {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}

	funcs["alarm"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*cloudwatch.MetricAlarm

		if !conf.getBoolDefaultTrue("aws.monitoring.alarm.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource monitoring[alarm]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Cloudwatch.DescribeAlarmsPages(&cloudwatch.DescribeAlarmsInput{},
			func(out *cloudwatch.DescribeAlarmsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.MetricAlarms {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildCdnFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCdnFetchFuncs(conf, funcs)

	funcs["distribution"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*cloudfront.DistributionSummary

		if !conf.getBoolDefaultTrue("aws.cdn.distribution.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cdn[distribution]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Cloudfront.ListDistributionsPages(&cloudfront.ListDistributionsInput{},
			func(out *cloudfront.ListDistributionsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.DistributionList.Items {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.DistributionList.NextMarker != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}
func BuildCloudformationFetchFuncs(conf *Config) fetch.Funcs {
	funcs := make(map[string]fetch.Func)

	addManualCloudformationFetchFuncs(conf, funcs)

	funcs["stack"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*cloudformation.Stack

		if !conf.getBoolDefaultTrue("aws.cloudformation.stack.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource cloudformation[stack]")
			return resources, objects, nil
		}
		var badResErr error
		err := conf.APIs.Cloudformation.DescribeStacksPages(&cloudformation.DescribeStacksInput{},
			func(out *cloudformation.DescribeStacksOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.Stacks {
					if badResErr != nil {
						return false
					}
					objects = append(objects, output)
					var res *graph.Resource
					if res, badResErr = awsconv.NewResource(output); badResErr != nil {
						return false
					}
					resources = append(resources, res)
				}
				return out.NextToken != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
	}
	return funcs
}

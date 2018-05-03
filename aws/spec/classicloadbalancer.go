/* Copyright 2017 WALLIX

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

package awsspec

import (
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"
)

type CreateClassicLoadbalancer struct {
	_                     string `action:"create" entity:"classicloadbalancer" awsAPI:"elb" awsCall:"CreateLoadBalancer" awsInput:"elb.CreateLoadBalancerInput" awsOutput:"elb.CreateLoadBalancerOutput"`
	logger                *logger.Logger
	graph                 cloud.GraphAPI
	api                   elbiface.ELBAPI
	Name                  *string   `awsName:"LoadBalancerName" awsType:"awsstr" templateName:"name"`
	AvailabilityZones     []*string `awsName:"AvailabilityZones" awsType:"awsstringslice" templateName:"zones"`
	Listeners             []*string `awsName:"Listeners" awsType:"awsclassicloadblisteners" templateName:"listeners"`
	Subnets               []*string `awsName:"Subnets" awsType:"awsstringslice" templateName:"subnets"`
	Securitygroups        []*string `awsName:"SecurityGroups" awsType:"awsstringslice" templateName:"securitygroups"`
	Scheme                *string   `awsName:"Scheme" awsType:"awsstr" templateName:"scheme"`
	Tags                  []*string `awsName:"Tags" awsType:"awstagslice" templateName:"tags"`
	HealthCheckPathToPing string    `awsType:"awsstr" templateName:"healthcheck-path"`
}

func (cmd *CreateClassicLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("listeners"),
		params.AtLeastOneOf(params.Key("subnets"), params.Key("zones")),
		params.Opt("scheme", "securitygroups", "tags", "healthcheck-path"),
	))
}

func (cmd *CreateClassicLoadbalancer) ExtractResult(i interface{}) string {
	return awssdk.StringValue(cmd.Name)
}

func (cmd *CreateClassicLoadbalancer) AfterRun(renv env.Running, output interface{}) (err error) {
	var target string
	if len(cmd.Listeners) > 0 {
		splits := strings.SplitN(*cmd.Listeners[0], ":", 3) // format validated on previous command
		target = splits[len(splits)-1]
	}

	var path string
	if strings.Contains(target, "HTTP") {
		path = "/"
		if p := cmd.HealthCheckPathToPing; len(p) > 0 {
			if strings.HasPrefix(p, "/") {
				path = p
			} else {
				path = "/" + p
			}
		}
	}

	target = target + path

	updateClassic := CommandFactory.Build("updateclassicloadbalancer")().(*UpdateClassicLoadbalancer)
	entries := map[string]interface{}{
		"name":                *cmd.Name,
		"healthy-threshold":   10,
		"unhealthy-threshold": 2,
		"health-timeout":      5,
		"health-interval":     30,
		"health-target":       target,
	}
	if err := params.Validate(updateClassic.ParamsSpec().Validators(), entries); err != nil {
		return err
	}
	_, err = updateClassic.Run(renv, entries)
	return
}

type UpdateClassicLoadbalancer struct {
	_                             string `action:"update" entity:"classicloadbalancer" awsAPI:"elb" awsCall:"ConfigureHealthCheck" awsInput:"elb.ConfigureHealthCheckInput" awsOutput:"elb.ConfigureHealthCheckOutput"`
	logger                        *logger.Logger
	graph                         cloud.GraphAPI
	api                           elbiface.ELBAPI
	Name                          *string `awsName:"LoadBalancerName" awsType:"awsstr" templateName:"name"`
	HealthcheckHealthyThreshold   *int64  `awsName:"Healthcheck.HealthyThreshold" awsType:"awsint64" templateName:"healthy-threshold"`
	HealthcheckUnhealthyThreshold *int64  `awsName:"Healthcheck.UnhealthyThreshold" awsType:"awsint64" templateName:"unhealthy-threshold"`
	HealthcheckInterval           *int64  `awsName:"Healthcheck.Interval" awsType:"awsint64" templateName:"health-interval"`
	HealthcheckTimeout            *int64  `awsName:"Healthcheck.Timeout" awsType:"awsint64" templateName:"health-timeout"`
	HealthcheckTarget             *string `awsName:"Healthcheck.Target" awsType:"awsstr" templateName:"health-target"`
}

func (cmd *UpdateClassicLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("healthy-threshold"), params.Key("unhealthy-threshold"),
		params.Key("health-interval"), params.Key("health-timeout"), params.Key("health-target")),
	)
}

type DeleteClassicLoadbalancer struct {
	_      string `action:"delete" entity:"classicloadbalancer" awsAPI:"elb" awsCall:"DeleteLoadBalancer" awsInput:"elb.DeleteLoadBalancerInput" awsOutput:"elb.DeleteLoadBalancerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    elbiface.ELBAPI
	Name   *string `awsName:"LoadBalancerName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteClassicLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type AttachClassicLoadbalancer struct {
	_        string `action:"attach" entity:"classicloadbalancer" awsAPI:"elb" awsCall:"RegisterInstancesWithLoadBalancer" awsInput:"elb.RegisterInstancesWithLoadBalancerInput" awsOutput:"elb.RegisterInstancesWithLoadBalancerOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      elbiface.ELBAPI
	Name     *string `awsName:"LoadBalancerName" awsType:"awsstr" templateName:"name"`
	Instance *string `awsName:"Instances[0]InstanceId" awsType:"awsslicestruct" templateName:"instance"`
}

func (cmd *AttachClassicLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("instance")))
}

func (cmd *AttachClassicLoadbalancer) ExtractResult(interface{}) string {
	return awssdk.StringValue(cmd.Instance)
}

type DetachClassicLoadbalancer struct {
	_        string `action:"detach" entity:"classicloadbalancer" awsAPI:"elb" awsCall:"DeregisterInstancesFromLoadBalancer" awsInput:"elb.DeregisterInstancesFromLoadBalancerInput" awsOutput:"elb.DeregisterInstancesFromLoadBalancerOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      elbiface.ELBAPI
	Name     *string `awsName:"LoadBalancerName" awsType:"awsstr" templateName:"name"`
	Instance *string `awsName:"Instances[0]InstanceId" awsType:"awsslicestruct" templateName:"instance"`
}

func (cmd *DetachClassicLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name"), params.Key("instance")))
}

func (cmd *DetachClassicLoadbalancer) ExtractResult(interface{}) string {
	return awssdk.StringValue(cmd.Instance)
}

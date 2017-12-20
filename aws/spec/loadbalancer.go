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
	"fmt"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"
)

type CreateLoadbalancer struct {
	_              string `action:"create" entity:"loadbalancer" awsAPI:"elbv2" awsCall:"CreateLoadBalancer" awsInput:"elbv2.CreateLoadBalancerInput" awsOutput:"elbv2.CreateLoadBalancerOutput"`
	logger         *logger.Logger
	graph          cloud.GraphAPI
	api            elbv2iface.ELBV2API
	Name           *string   `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Subnets        []*string `awsName:"Subnets" awsType:"awsstringslice" templateName:"subnets"`
	SubnetMappings []*string `awsName:"SubnetMappings" awsType:"awssubnetmappings" templateName:"subnet-mappings"`
	Iptype         *string   `awsName:"IpAddressType" awsType:"awsstr" templateName:"iptype"`
	Scheme         *string   `awsName:"Scheme" awsType:"awsstr" templateName:"scheme"`
	Securitygroups []*string `awsName:"SecurityGroups" awsType:"awsstringslice" templateName:"securitygroups"`
	Type           *string   `awsName:"Type" awsType:"awsstr" templateName:"type"`
}

func (cmd *CreateLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("subnets"),
		params.Opt("subnet-mappings", "iptype", "scheme", "securitygroups", "type"),
	))
}

func (cmd *CreateLoadbalancer) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*elbv2.CreateLoadBalancerOutput).LoadBalancers[0].LoadBalancerArn)
}

type DeleteLoadbalancer struct {
	_      string `action:"delete" entity:"loadbalancer" awsAPI:"elbv2" awsCall:"DeleteLoadBalancer" awsInput:"elbv2.DeleteLoadBalancerInput" awsOutput:"elbv2.DeleteLoadBalancerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    elbv2iface.ELBV2API
	Id     *string `awsName:"LoadBalancerArn" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type CheckLoadbalancer struct {
	_       string `action:"check" entity:"loadbalancer" awsAPI:"elbv2"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     elbv2iface.ELBV2API
	Id      *string `templateName:"id"`
	State   *string `templateName:"state"`
	Timeout *int64  `templateName:"timeout"`
}

func (cmd *CheckLoadbalancer) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("id"), params.Key("state"), params.Key("timeout")),
		params.Validators{
			"state": params.IsInEnumIgnoreCase("provisioning", "active", "failed", notFoundState),
		})
}

func (cmd *CheckLoadbalancer) ManualRun(renv env.Running) (interface{}, error) {
	input := &elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []*string{cmd.Id},
	}

	c := &checker{
		description: fmt.Sprintf("loadbalancer %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeLoadBalancers(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "LoadBalancerNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				for _, lb := range output.LoadBalancers {
					if StringValue(lb.LoadBalancerArn) == StringValue(cmd.Id) {
						return StringValue(lb.State.Code), nil
					}
				}
			}
			return notFoundState, nil
		},
		expect: StringValue(cmd.State),
		logger: cmd.logger,
	}
	return nil, c.check()
}

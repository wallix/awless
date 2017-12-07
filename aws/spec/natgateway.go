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

	"github.com/wallix/awless/cloud/graph"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateNatgateway struct {
	_           string `action:"create" entity:"natgateway" awsAPI:"ec2" awsCall:"CreateNatGateway" awsInput:"ec2.CreateNatGatewayInput" awsOutput:"ec2.CreateNatGatewayOutput"`
	logger      *logger.Logger
	graph       cloudgraph.GraphAPI
	api         ec2iface.EC2API
	ElasticipId *string `awsName:"AllocationId" awsType:"awsstr" templateName:"elasticip-id" required:""`
	Subnet      *string `awsName:"SubnetId" awsType:"awsstr" templateName:"subnet" required:""`
}

func (cmd *CreateNatgateway) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateNatgateway) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateNatGatewayOutput).NatGateway.NatGatewayId)
}

type DeleteNatgateway struct {
	_      string `action:"delete" entity:"natgateway" awsAPI:"ec2" awsCall:"DeleteNatGateway" awsInput:"ec2.DeleteNatGatewayInput" awsOutput:"ec2.DeleteNatGatewayOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"NatGatewayId" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *DeleteNatgateway) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type CheckNatgateway struct {
	_       string `action:"check" entity:"natgateway" awsAPI:"ec2"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     ec2iface.EC2API
	Id      *string `templateName:"id" required:""`
	State   *string `templateName:"state" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckNatgateway) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CheckNatgateway) Validate_State() error {
	return NewEnumValidator("pending", "failed", "available", "deleting", "deleted", notFoundState).Validate(cmd.State)
}

func (cmd *CheckNatgateway) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{cmd.Id},
	}

	c := &checker{
		description: fmt.Sprintf("natgateway %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeNatGateways(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "NatGatewayNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				for _, nat := range output.NatGateways {
					if StringValue(nat.NatGatewayId) == StringValue(cmd.Id) {
						return StringValue(nat.State), nil
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

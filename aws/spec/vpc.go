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
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateVpc struct {
	_      string `action:"create" entity:"vpc" awsAPI:"ec2" awsCall:"CreateVpc" awsInput:"ec2.CreateVpcInput" awsOutput:"ec2.CreateVpcOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	CIDR   *string `awsName:"CidrBlock" awsType:"awsstr" templateName:"cidr"`
	Name   *string `awsName:"Name" templateName:"name"`
}

func (cmd *CreateVpc) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("cidr"), params.Opt(params.Suggested("name"))),
		params.Validators{"cidr": params.IsCIDR})
}

func (cmd *CreateVpc) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateVpcOutput).Vpc.VpcId)
}

func (cmd *CreateVpc) AfterRun(renv env.Running, output interface{}) error {
	return createNameTag(awssdk.String(cmd.ExtractResult(output)), cmd.Name, renv)
}

type DeleteVpc struct {
	_      string `action:"delete" entity:"vpc" awsAPI:"ec2" awsCall:"DeleteVpc" awsInput:"ec2.DeleteVpcInput" awsOutput:"ec2.DeleteVpcOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"VpcId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteVpc) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

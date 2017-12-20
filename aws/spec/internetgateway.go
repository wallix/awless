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
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateInternetgateway struct {
	_      string `action:"create" entity:"internetgateway" awsAPI:"ec2" awsCall:"CreateInternetGateway" awsInput:"ec2.CreateInternetGatewayInput" awsOutput:"ec2.CreateInternetGatewayOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
}

func (cmd *CreateInternetgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.None())
}

func (cmd *CreateInternetgateway) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateInternetGatewayOutput).InternetGateway.InternetGatewayId)
}

type DeleteInternetgateway struct {
	_      string `action:"delete" entity:"internetgateway" awsAPI:"ec2" awsCall:"DeleteInternetGateway" awsInput:"ec2.DeleteInternetGatewayInput" awsOutput:"ec2.DeleteInternetGatewayOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"InternetGatewayId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteInternetgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type AttachInternetgateway struct {
	_      string `action:"attach" entity:"internetgateway" awsAPI:"ec2" awsCall:"AttachInternetGateway" awsInput:"ec2.AttachInternetGatewayInput" awsOutput:"ec2.AttachInternetGatewayOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"InternetGatewayId" awsType:"awsstr" templateName:"id"`
	Vpc    *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
}

func (cmd *AttachInternetgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("vpc")))
}

type DetachInternetgateway struct {
	_      string `action:"detach" entity:"internetgateway" awsAPI:"ec2" awsCall:"DetachInternetGateway" awsInput:"ec2.DetachInternetGatewayInput" awsOutput:"ec2.DetachInternetGatewayOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"InternetGatewayId" awsType:"awsstr" templateName:"id"`
	Vpc    *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
}

func (cmd *DetachInternetgateway) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("vpc")))
}

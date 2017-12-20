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

type CreateRoutetable struct {
	_      string `action:"create" entity:"routetable" awsAPI:"ec2" awsCall:"CreateRouteTable" awsInput:"ec2.CreateRouteTableInput" awsOutput:"ec2.CreateRouteTableOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Vpc    *string `awsName:"VpcId" awsType:"awsstr" templateName:"vpc"`
}

func (cmd *CreateRoutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("vpc")))
}

func (cmd *CreateRoutetable) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.CreateRouteTableOutput).RouteTable.RouteTableId)
}

type DeleteRoutetable struct {
	_      string `action:"delete" entity:"routetable" awsAPI:"ec2" awsCall:"DeleteRouteTable" awsInput:"ec2.DeleteRouteTableInput" awsOutput:"ec2.DeleteRouteTableOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"RouteTableId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteRoutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

type AttachRoutetable struct {
	_      string `action:"attach" entity:"routetable" awsAPI:"ec2" awsCall:"AssociateRouteTable" awsInput:"ec2.AssociateRouteTableInput" awsOutput:"ec2.AssociateRouteTableOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"RouteTableId" awsType:"awsstr" templateName:"id"`
	Subnet *string `awsName:"SubnetId" awsType:"awsstr" templateName:"subnet"`
}

func (cmd *AttachRoutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("subnet")))
}

func (cmd *AttachRoutetable) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.AssociateRouteTableOutput).AssociationId)
}

type DetachRoutetable struct {
	_           string `action:"detach" entity:"routetable" awsAPI:"ec2" awsCall:"DisassociateRouteTable" awsInput:"ec2.DisassociateRouteTableInput" awsOutput:"ec2.DisassociateRouteTableOutput" awsDryRun:""`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         ec2iface.EC2API
	Association *string `awsName:"AssociationId" awsType:"awsstr" templateName:"association"`
}

func (cmd *DetachRoutetable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("association")))
}

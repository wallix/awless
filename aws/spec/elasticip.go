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
	"github.com/wallix/awless/cloud/graph"
	"github.com/wallix/awless/logger"
)

type CreateElasticip struct {
	_      string `action:"create" entity:"elasticip" awsAPI:"ec2" awsCall:"AllocateAddress" awsInput:"ec2.AllocateAddressInput" awsOutput:"ec2.AllocateAddressOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	Domain *string `awsName:"Domain" awsType:"awsstr" templateName:"domain" required:""`
}

func (cmd *CreateElasticip) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateElasticip) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.AllocateAddressOutput).AllocationId)
}

type DeleteElasticip struct {
	_      string `action:"delete" entity:"elasticip" awsAPI:"ec2" awsCall:"ReleaseAddress" awsInput:"ec2.ReleaseAddressInput" awsOutput:"ec2.ReleaseAddressOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	Id     *string `awsName:"AllocationId" awsType:"awsstr" templateName:"id"`
	Ip     *string `awsName:"PublicIp" awsType:"awsstr" templateName:"ip"`
}

func (cmd *DeleteElasticip) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type AttachElasticip struct {
	_                  string `action:"attach" entity:"elasticip" awsAPI:"ec2" awsCall:"AssociateAddress" awsInput:"ec2.AssociateAddressInput" awsOutput:"ec2.AssociateAddressOutput" awsDryRun:""`
	logger             *logger.Logger
	graph              cloudgraph.GraphAPI
	api                ec2iface.EC2API
	Id                 *string `awsName:"AllocationId" awsType:"awsstr" templateName:"id" required:""`
	Instance           *string `awsName:"InstanceId" awsType:"awsstr" templateName:"instance"`
	Networkinterface   *string `awsName:"NetworkInterfaceId" awsType:"awsstr" templateName:"networkinterface"`
	Privateip          *string `awsName:"PrivateIpAddress" awsType:"awsstr" templateName:"privateip"`
	AllowReassociation *bool   `awsName:"AllowReassociation" awsType:"awsbool" templateName:"allow-reassociation"`
}

func (cmd *AttachElasticip) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *AttachElasticip) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*ec2.AssociateAddressOutput).AssociationId)
}

type DetachElasticip struct {
	_           string `action:"detach" entity:"elasticip" awsAPI:"ec2" awsCall:"DisassociateAddress" awsInput:"ec2.DisassociateAddressInput" awsOutput:"ec2.DisassociateAddressOutput" awsDryRun:""`
	logger      *logger.Logger
	graph       cloudgraph.GraphAPI
	api         ec2iface.EC2API
	Association *string `awsName:"AssociationId" awsType:"awsstr" templateName:"association" required:""`
}

func (cmd *DetachElasticip) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

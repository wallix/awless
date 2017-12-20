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
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/params"
)

type CreateZone struct {
	_               string `action:"create" entity:"zone" awsAPI:"route53" awsCall:"CreateHostedZone" awsInput:"route53.CreateHostedZoneInput" awsOutput:"route53.CreateHostedZoneOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             route53iface.Route53API
	Callerreference *string `awsName:"CallerReference" awsType:"awsstr" templateName:"callerreference"`
	Name            *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Delegationsetid *string `awsName:"DelegationSetId" awsType:"awsstr" templateName:"delegationsetid"`
	Comment         *string `awsName:"HostedZoneConfig.Comment" awsType:"awsstr" templateName:"comment"`
	Isprivate       *bool   `awsName:"HostedZoneConfig.PrivateZone" awsType:"awsbool" templateName:"isprivate"`
	Vpcid           *string `awsName:"VPC.VPCId" awsType:"awsstr" templateName:"vpcid"`
	Vpcregion       *string `awsName:"VPC.VPCRegion" awsType:"awsstr" templateName:"vpcregion"`
}

func (cmd *CreateZone) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("callerreference"), params.Key("name"),
		params.Opt("comment", "delegationsetid", "isprivate", "vpcid", "vpcregion"),
	))
}

func (cmd *CreateZone) ExtractResult(i interface{}) string {
	return awssdk.StringValue(i.(*route53.CreateHostedZoneOutput).HostedZone.Id)
}

type DeleteZone struct {
	_      string `action:"delete" entity:"zone" awsAPI:"route53" awsCall:"DeleteHostedZone" awsInput:"route53.DeleteHostedZoneInput" awsOutput:"route53.DeleteHostedZoneOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    route53iface.Route53API
	Id     *string `awsName:"Id" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteZone) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

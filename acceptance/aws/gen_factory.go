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

// DO NOT EDIT
// This file was automatically generated with go generate
package awsat

import (
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling/applicationautoscalingiface"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/logger"
)

type AcceptanceFactory struct {
	Mock   interface{}
	Logger *logger.Logger
}

func NewAcceptanceFactory(mock interface{}, l ...*logger.Logger) *AcceptanceFactory {
	logger := logger.DiscardLogger
	if len(l) > 0 {
		logger = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Logger: logger}
}

func (f *AcceptanceFactory) Build(key string) func() interface{} {
	switch key {
	case "attachalarm":
		return func() interface{} {
			cmd := awsspec.NewAttachAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "attachcontainertask":
		return func() interface{} {
			cmd := awsspec.NewAttachContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "attachelasticip":
		return func() interface{} {
			cmd := awsspec.NewAttachElasticip(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachinstance":
		return func() interface{} {
			cmd := awsspec.NewAttachInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "attachinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewAttachInstanceprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewAttachInternetgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachmfadevice":
		return func() interface{} {
			cmd := awsspec.NewAttachMfadevice(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "attachnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewAttachNetworkinterface(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachpolicy":
		return func() interface{} {
			cmd := awsspec.NewAttachPolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "attachrole":
		return func() interface{} {
			cmd := awsspec.NewAttachRole(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "attachroutetable":
		return func() interface{} {
			cmd := awsspec.NewAttachRoutetable(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewAttachSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "attachuser":
		return func() interface{} {
			cmd := awsspec.NewAttachUser(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "attachvolume":
		return func() interface{} {
			cmd := awsspec.NewAttachVolume(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "authenticateregistry":
		return func() interface{} {
			cmd := awsspec.NewAuthenticateRegistry(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecriface.ECRAPI))
			return cmd
		}
	case "checkcertificate":
		return func() interface{} {
			cmd := awsspec.NewCheckCertificate(nil, f.Logger)
			cmd.SetApi(f.Mock.(acmiface.ACMAPI))
			return cmd
		}
	case "checkdatabase":
		return func() interface{} {
			cmd := awsspec.NewCheckDatabase(nil, f.Logger)
			cmd.SetApi(f.Mock.(rdsiface.RDSAPI))
			return cmd
		}
	case "checkdistribution":
		return func() interface{} {
			cmd := awsspec.NewCheckDistribution(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudfrontiface.CloudFrontAPI))
			return cmd
		}
	case "checkinstance":
		return func() interface{} {
			cmd := awsspec.NewCheckInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "checkloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewCheckLoadbalancer(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "checknatgateway":
		return func() interface{} {
			cmd := awsspec.NewCheckNatgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "checknetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewCheckNetworkinterface(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "checkscalinggroup":
		return func() interface{} {
			cmd := awsspec.NewCheckScalinggroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "checksecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewCheckSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "checkvolume":
		return func() interface{} {
			cmd := awsspec.NewCheckVolume(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "copyimage":
		return func() interface{} {
			cmd := awsspec.NewCopyImage(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "copysnapshot":
		return func() interface{} {
			cmd := awsspec.NewCopySnapshot(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createaccesskey":
		return func() interface{} {
			cmd := awsspec.NewCreateAccesskey(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createalarm":
		return func() interface{} {
			cmd := awsspec.NewCreateAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "createappscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreateAppscalingpolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(applicationautoscalingiface.ApplicationAutoScalingAPI))
			return cmd
		}
	case "createappscalingtarget":
		return func() interface{} {
			cmd := awsspec.NewCreateAppscalingtarget(nil, f.Logger)
			cmd.SetApi(f.Mock.(applicationautoscalingiface.ApplicationAutoScalingAPI))
			return cmd
		}
	case "createbucket":
		return func() interface{} {
			cmd := awsspec.NewCreateBucket(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "createcertificate":
		return func() interface{} {
			cmd := awsspec.NewCreateCertificate(nil, f.Logger)
			cmd.SetApi(f.Mock.(acmiface.ACMAPI))
			return cmd
		}
	case "createcontainercluster":
		return func() interface{} {
			cmd := awsspec.NewCreateContainercluster(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "createdatabase":
		return func() interface{} {
			cmd := awsspec.NewCreateDatabase(nil, f.Logger)
			cmd.SetApi(f.Mock.(rdsiface.RDSAPI))
			return cmd
		}
	case "createdbsubnetgroup":
		return func() interface{} {
			cmd := awsspec.NewCreateDbsubnetgroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(rdsiface.RDSAPI))
			return cmd
		}
	case "createdistribution":
		return func() interface{} {
			cmd := awsspec.NewCreateDistribution(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudfrontiface.CloudFrontAPI))
			return cmd
		}
	case "createelasticip":
		return func() interface{} {
			cmd := awsspec.NewCreateElasticip(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createfunction":
		return func() interface{} {
			cmd := awsspec.NewCreateFunction(nil, f.Logger)
			cmd.SetApi(f.Mock.(lambdaiface.LambdaAPI))
			return cmd
		}
	case "creategroup":
		return func() interface{} {
			cmd := awsspec.NewCreateGroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createimage":
		return func() interface{} {
			cmd := awsspec.NewCreateImage(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createinstance":
		return func() interface{} {
			cmd := awsspec.NewCreateInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewCreateInstanceprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewCreateInternetgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createkeypair":
		return func() interface{} {
			cmd := awsspec.NewCreateKeypair(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createlaunchconfiguration":
		return func() interface{} {
			cmd := awsspec.NewCreateLaunchconfiguration(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "createlistener":
		return func() interface{} {
			cmd := awsspec.NewCreateListener(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "createloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewCreateLoadbalancer(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "createloginprofile":
		return func() interface{} {
			cmd := awsspec.NewCreateLoginprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createmfadevice":
		return func() interface{} {
			cmd := awsspec.NewCreateMfadevice(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createnatgateway":
		return func() interface{} {
			cmd := awsspec.NewCreateNatgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewCreateNetworkinterface(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreatePolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createqueue":
		return func() interface{} {
			cmd := awsspec.NewCreateQueue(nil, f.Logger)
			cmd.SetApi(f.Mock.(sqsiface.SQSAPI))
			return cmd
		}
	case "createrecord":
		return func() interface{} {
			cmd := awsspec.NewCreateRecord(nil, f.Logger)
			cmd.SetApi(f.Mock.(route53iface.Route53API))
			return cmd
		}
	case "createrepository":
		return func() interface{} {
			cmd := awsspec.NewCreateRepository(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecriface.ECRAPI))
			return cmd
		}
	case "createrole":
		return func() interface{} {
			cmd := awsspec.NewCreateRole(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createroute":
		return func() interface{} {
			cmd := awsspec.NewCreateRoute(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createroutetable":
		return func() interface{} {
			cmd := awsspec.NewCreateRoutetable(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "creates3object":
		return func() interface{} {
			cmd := awsspec.NewCreateS3object(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "createscalinggroup":
		return func() interface{} {
			cmd := awsspec.NewCreateScalinggroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "createscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreateScalingpolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "createsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewCreateSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createsnapshot":
		return func() interface{} {
			cmd := awsspec.NewCreateSnapshot(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createstack":
		return func() interface{} {
			cmd := awsspec.NewCreateStack(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudformationiface.CloudFormationAPI))
			return cmd
		}
	case "createsubnet":
		return func() interface{} {
			cmd := awsspec.NewCreateSubnet(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createsubscription":
		return func() interface{} {
			cmd := awsspec.NewCreateSubscription(nil, f.Logger)
			cmd.SetApi(f.Mock.(snsiface.SNSAPI))
			return cmd
		}
	case "createtag":
		return func() interface{} {
			cmd := awsspec.NewCreateTag(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createtargetgroup":
		return func() interface{} {
			cmd := awsspec.NewCreateTargetgroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "createtopic":
		return func() interface{} {
			cmd := awsspec.NewCreateTopic(nil, f.Logger)
			cmd.SetApi(f.Mock.(snsiface.SNSAPI))
			return cmd
		}
	case "createuser":
		return func() interface{} {
			cmd := awsspec.NewCreateUser(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "createvolume":
		return func() interface{} {
			cmd := awsspec.NewCreateVolume(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createvpc":
		return func() interface{} {
			cmd := awsspec.NewCreateVpc(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "createzone":
		return func() interface{} {
			cmd := awsspec.NewCreateZone(nil, f.Logger)
			cmd.SetApi(f.Mock.(route53iface.Route53API))
			return cmd
		}
	case "deleteaccesskey":
		return func() interface{} {
			cmd := awsspec.NewDeleteAccesskey(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deletealarm":
		return func() interface{} {
			cmd := awsspec.NewDeleteAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "deleteappscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewDeleteAppscalingpolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(applicationautoscalingiface.ApplicationAutoScalingAPI))
			return cmd
		}
	case "deleteappscalingtarget":
		return func() interface{} {
			cmd := awsspec.NewDeleteAppscalingtarget(nil, f.Logger)
			cmd.SetApi(f.Mock.(applicationautoscalingiface.ApplicationAutoScalingAPI))
			return cmd
		}
	case "deletebucket":
		return func() interface{} {
			cmd := awsspec.NewDeleteBucket(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "deletecertificate":
		return func() interface{} {
			cmd := awsspec.NewDeleteCertificate(nil, f.Logger)
			cmd.SetApi(f.Mock.(acmiface.ACMAPI))
			return cmd
		}
	case "deletecontainercluster":
		return func() interface{} {
			cmd := awsspec.NewDeleteContainercluster(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "deletecontainertask":
		return func() interface{} {
			cmd := awsspec.NewDeleteContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "deletedatabase":
		return func() interface{} {
			cmd := awsspec.NewDeleteDatabase(nil, f.Logger)
			cmd.SetApi(f.Mock.(rdsiface.RDSAPI))
			return cmd
		}
	case "deletedbsubnetgroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteDbsubnetgroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(rdsiface.RDSAPI))
			return cmd
		}
	case "deletedistribution":
		return func() interface{} {
			cmd := awsspec.NewDeleteDistribution(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudfrontiface.CloudFrontAPI))
			return cmd
		}
	case "deleteelasticip":
		return func() interface{} {
			cmd := awsspec.NewDeleteElasticip(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletefunction":
		return func() interface{} {
			cmd := awsspec.NewDeleteFunction(nil, f.Logger)
			cmd.SetApi(f.Mock.(lambdaiface.LambdaAPI))
			return cmd
		}
	case "deletegroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteGroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deleteimage":
		return func() interface{} {
			cmd := awsspec.NewDeleteImage(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deleteinstance":
		return func() interface{} {
			cmd := awsspec.NewDeleteInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deleteinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewDeleteInstanceprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deleteinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewDeleteInternetgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletekeypair":
		return func() interface{} {
			cmd := awsspec.NewDeleteKeypair(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletelaunchconfiguration":
		return func() interface{} {
			cmd := awsspec.NewDeleteLaunchconfiguration(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "deletelistener":
		return func() interface{} {
			cmd := awsspec.NewDeleteListener(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "deleteloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewDeleteLoadbalancer(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "deleteloginprofile":
		return func() interface{} {
			cmd := awsspec.NewDeleteLoginprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deletemfadevice":
		return func() interface{} {
			cmd := awsspec.NewDeleteMfadevice(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deletenatgateway":
		return func() interface{} {
			cmd := awsspec.NewDeleteNatgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletenetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewDeleteNetworkinterface(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletepolicy":
		return func() interface{} {
			cmd := awsspec.NewDeletePolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deletequeue":
		return func() interface{} {
			cmd := awsspec.NewDeleteQueue(nil, f.Logger)
			cmd.SetApi(f.Mock.(sqsiface.SQSAPI))
			return cmd
		}
	case "deleterecord":
		return func() interface{} {
			cmd := awsspec.NewDeleteRecord(nil, f.Logger)
			cmd.SetApi(f.Mock.(route53iface.Route53API))
			return cmd
		}
	case "deleterepository":
		return func() interface{} {
			cmd := awsspec.NewDeleteRepository(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecriface.ECRAPI))
			return cmd
		}
	case "deleterole":
		return func() interface{} {
			cmd := awsspec.NewDeleteRole(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deleteroute":
		return func() interface{} {
			cmd := awsspec.NewDeleteRoute(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deleteroutetable":
		return func() interface{} {
			cmd := awsspec.NewDeleteRoutetable(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletes3object":
		return func() interface{} {
			cmd := awsspec.NewDeleteS3object(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "deletescalinggroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteScalinggroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "deletescalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewDeleteScalingpolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "deletesecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletesnapshot":
		return func() interface{} {
			cmd := awsspec.NewDeleteSnapshot(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletestack":
		return func() interface{} {
			cmd := awsspec.NewDeleteStack(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudformationiface.CloudFormationAPI))
			return cmd
		}
	case "deletesubnet":
		return func() interface{} {
			cmd := awsspec.NewDeleteSubnet(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletesubscription":
		return func() interface{} {
			cmd := awsspec.NewDeleteSubscription(nil, f.Logger)
			cmd.SetApi(f.Mock.(snsiface.SNSAPI))
			return cmd
		}
	case "deletetag":
		return func() interface{} {
			cmd := awsspec.NewDeleteTag(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletetargetgroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteTargetgroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "deletetopic":
		return func() interface{} {
			cmd := awsspec.NewDeleteTopic(nil, f.Logger)
			cmd.SetApi(f.Mock.(snsiface.SNSAPI))
			return cmd
		}
	case "deleteuser":
		return func() interface{} {
			cmd := awsspec.NewDeleteUser(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "deletevolume":
		return func() interface{} {
			cmd := awsspec.NewDeleteVolume(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletevpc":
		return func() interface{} {
			cmd := awsspec.NewDeleteVpc(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "deletezone":
		return func() interface{} {
			cmd := awsspec.NewDeleteZone(nil, f.Logger)
			cmd.SetApi(f.Mock.(route53iface.Route53API))
			return cmd
		}
	case "detachalarm":
		return func() interface{} {
			cmd := awsspec.NewDetachAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "detachcontainertask":
		return func() interface{} {
			cmd := awsspec.NewDetachContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "detachelasticip":
		return func() interface{} {
			cmd := awsspec.NewDetachElasticip(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachinstance":
		return func() interface{} {
			cmd := awsspec.NewDetachInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	case "detachinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewDetachInstanceprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewDetachInternetgateway(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachmfadevice":
		return func() interface{} {
			cmd := awsspec.NewDetachMfadevice(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "detachnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewDetachNetworkinterface(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachpolicy":
		return func() interface{} {
			cmd := awsspec.NewDetachPolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "detachrole":
		return func() interface{} {
			cmd := awsspec.NewDetachRole(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "detachroutetable":
		return func() interface{} {
			cmd := awsspec.NewDetachRoutetable(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewDetachSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "detachuser":
		return func() interface{} {
			cmd := awsspec.NewDetachUser(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "detachvolume":
		return func() interface{} {
			cmd := awsspec.NewDetachVolume(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "importimage":
		return func() interface{} {
			cmd := awsspec.NewImportImage(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "startalarm":
		return func() interface{} {
			cmd := awsspec.NewStartAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "startcontainertask":
		return func() interface{} {
			cmd := awsspec.NewStartContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "startinstance":
		return func() interface{} {
			cmd := awsspec.NewStartInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "stopalarm":
		return func() interface{} {
			cmd := awsspec.NewStopAlarm(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudwatchiface.CloudWatchAPI))
			return cmd
		}
	case "stopcontainertask":
		return func() interface{} {
			cmd := awsspec.NewStopContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "stopinstance":
		return func() interface{} {
			cmd := awsspec.NewStopInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "updatebucket":
		return func() interface{} {
			cmd := awsspec.NewUpdateBucket(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "updatecontainertask":
		return func() interface{} {
			cmd := awsspec.NewUpdateContainertask(nil, f.Logger)
			cmd.SetApi(f.Mock.(ecsiface.ECSAPI))
			return cmd
		}
	case "updatedistribution":
		return func() interface{} {
			cmd := awsspec.NewUpdateDistribution(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudfrontiface.CloudFrontAPI))
			return cmd
		}
	case "updateimage":
		return func() interface{} {
			cmd := awsspec.NewUpdateImage(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "updateinstance":
		return func() interface{} {
			cmd := awsspec.NewUpdateInstance(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "updateloginprofile":
		return func() interface{} {
			cmd := awsspec.NewUpdateLoginprofile(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "updatepolicy":
		return func() interface{} {
			cmd := awsspec.NewUpdatePolicy(nil, f.Logger)
			cmd.SetApi(f.Mock.(iamiface.IAMAPI))
			return cmd
		}
	case "updaterecord":
		return func() interface{} {
			cmd := awsspec.NewUpdateRecord(nil, f.Logger)
			cmd.SetApi(f.Mock.(route53iface.Route53API))
			return cmd
		}
	case "updates3object":
		return func() interface{} {
			cmd := awsspec.NewUpdateS3object(nil, f.Logger)
			cmd.SetApi(f.Mock.(s3iface.S3API))
			return cmd
		}
	case "updatescalinggroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateScalinggroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(autoscalingiface.AutoScalingAPI))
			return cmd
		}
	case "updatesecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateSecuritygroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "updatestack":
		return func() interface{} {
			cmd := awsspec.NewUpdateStack(nil, f.Logger)
			cmd.SetApi(f.Mock.(cloudformationiface.CloudFormationAPI))
			return cmd
		}
	case "updatesubnet":
		return func() interface{} {
			cmd := awsspec.NewUpdateSubnet(nil, f.Logger)
			cmd.SetApi(f.Mock.(ec2iface.EC2API))
			return cmd
		}
	case "updatetargetgroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateTargetgroup(nil, f.Logger)
			cmd.SetApi(f.Mock.(elbv2iface.ELBV2API))
			return cmd
		}
	}
	return nil
}

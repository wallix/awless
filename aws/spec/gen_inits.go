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
package awsspec

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
)

type Factory interface {
	Build(key string) func() interface{}
}

var CommandFactory Factory

var MockAWSSessionFactory = &AWSFactory{
	Log:  logger.DiscardLogger,
	Sess: mock.Session,
}

type AWSFactory struct {
	Log   *logger.Logger
	Sess  *session.Session
	Graph cloud.GraphAPI
}

func (f *AWSFactory) Build(key string) func() interface{} {
	switch key {
	case "attachalarm":
		return func() interface{} { return NewAttachAlarm(f.Sess, f.Graph, f.Log) }
	case "attachclassicloadbalancer":
		return func() interface{} { return NewAttachClassicLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "attachcontainertask":
		return func() interface{} { return NewAttachContainertask(f.Sess, f.Graph, f.Log) }
	case "attachelasticip":
		return func() interface{} { return NewAttachElasticip(f.Sess, f.Graph, f.Log) }
	case "attachinstance":
		return func() interface{} { return NewAttachInstance(f.Sess, f.Graph, f.Log) }
	case "attachinstanceprofile":
		return func() interface{} { return NewAttachInstanceprofile(f.Sess, f.Graph, f.Log) }
	case "attachinternetgateway":
		return func() interface{} { return NewAttachInternetgateway(f.Sess, f.Graph, f.Log) }
	case "attachlistener":
		return func() interface{} { return NewAttachListener(f.Sess, f.Graph, f.Log) }
	case "attachmfadevice":
		return func() interface{} { return NewAttachMfadevice(f.Sess, f.Graph, f.Log) }
	case "attachnetworkinterface":
		return func() interface{} { return NewAttachNetworkinterface(f.Sess, f.Graph, f.Log) }
	case "attachpolicy":
		return func() interface{} { return NewAttachPolicy(f.Sess, f.Graph, f.Log) }
	case "attachrole":
		return func() interface{} { return NewAttachRole(f.Sess, f.Graph, f.Log) }
	case "attachroutetable":
		return func() interface{} { return NewAttachRoutetable(f.Sess, f.Graph, f.Log) }
	case "attachsecuritygroup":
		return func() interface{} { return NewAttachSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "attachuser":
		return func() interface{} { return NewAttachUser(f.Sess, f.Graph, f.Log) }
	case "attachvolume":
		return func() interface{} { return NewAttachVolume(f.Sess, f.Graph, f.Log) }
	case "authenticateregistry":
		return func() interface{} { return NewAuthenticateRegistry(f.Sess, f.Graph, f.Log) }
	case "checkcertificate":
		return func() interface{} { return NewCheckCertificate(f.Sess, f.Graph, f.Log) }
	case "checkdatabase":
		return func() interface{} { return NewCheckDatabase(f.Sess, f.Graph, f.Log) }
	case "checkdistribution":
		return func() interface{} { return NewCheckDistribution(f.Sess, f.Graph, f.Log) }
	case "checkinstance":
		return func() interface{} { return NewCheckInstance(f.Sess, f.Graph, f.Log) }
	case "checkloadbalancer":
		return func() interface{} { return NewCheckLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "checknatgateway":
		return func() interface{} { return NewCheckNatgateway(f.Sess, f.Graph, f.Log) }
	case "checknetworkinterface":
		return func() interface{} { return NewCheckNetworkinterface(f.Sess, f.Graph, f.Log) }
	case "checkscalinggroup":
		return func() interface{} { return NewCheckScalinggroup(f.Sess, f.Graph, f.Log) }
	case "checksecuritygroup":
		return func() interface{} { return NewCheckSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "checkvolume":
		return func() interface{} { return NewCheckVolume(f.Sess, f.Graph, f.Log) }
	case "copyimage":
		return func() interface{} { return NewCopyImage(f.Sess, f.Graph, f.Log) }
	case "copysnapshot":
		return func() interface{} { return NewCopySnapshot(f.Sess, f.Graph, f.Log) }
	case "createaccesskey":
		return func() interface{} { return NewCreateAccesskey(f.Sess, f.Graph, f.Log) }
	case "createalarm":
		return func() interface{} { return NewCreateAlarm(f.Sess, f.Graph, f.Log) }
	case "createappscalingpolicy":
		return func() interface{} { return NewCreateAppscalingpolicy(f.Sess, f.Graph, f.Log) }
	case "createappscalingtarget":
		return func() interface{} { return NewCreateAppscalingtarget(f.Sess, f.Graph, f.Log) }
	case "createbucket":
		return func() interface{} { return NewCreateBucket(f.Sess, f.Graph, f.Log) }
	case "createcertificate":
		return func() interface{} { return NewCreateCertificate(f.Sess, f.Graph, f.Log) }
	case "createclassicloadbalancer":
		return func() interface{} { return NewCreateClassicLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "createcontainercluster":
		return func() interface{} { return NewCreateContainercluster(f.Sess, f.Graph, f.Log) }
	case "createdatabase":
		return func() interface{} { return NewCreateDatabase(f.Sess, f.Graph, f.Log) }
	case "createdbsubnetgroup":
		return func() interface{} { return NewCreateDbsubnetgroup(f.Sess, f.Graph, f.Log) }
	case "createdistribution":
		return func() interface{} { return NewCreateDistribution(f.Sess, f.Graph, f.Log) }
	case "createelasticip":
		return func() interface{} { return NewCreateElasticip(f.Sess, f.Graph, f.Log) }
	case "createfunction":
		return func() interface{} { return NewCreateFunction(f.Sess, f.Graph, f.Log) }
	case "creategroup":
		return func() interface{} { return NewCreateGroup(f.Sess, f.Graph, f.Log) }
	case "createimage":
		return func() interface{} { return NewCreateImage(f.Sess, f.Graph, f.Log) }
	case "createinstance":
		return func() interface{} { return NewCreateInstance(f.Sess, f.Graph, f.Log) }
	case "createinstanceprofile":
		return func() interface{} { return NewCreateInstanceprofile(f.Sess, f.Graph, f.Log) }
	case "createinternetgateway":
		return func() interface{} { return NewCreateInternetgateway(f.Sess, f.Graph, f.Log) }
	case "createkeypair":
		return func() interface{} { return NewCreateKeypair(f.Sess, f.Graph, f.Log) }
	case "createlaunchconfiguration":
		return func() interface{} { return NewCreateLaunchconfiguration(f.Sess, f.Graph, f.Log) }
	case "createlistener":
		return func() interface{} { return NewCreateListener(f.Sess, f.Graph, f.Log) }
	case "createloadbalancer":
		return func() interface{} { return NewCreateLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "createloginprofile":
		return func() interface{} { return NewCreateLoginprofile(f.Sess, f.Graph, f.Log) }
	case "createmfadevice":
		return func() interface{} { return NewCreateMfadevice(f.Sess, f.Graph, f.Log) }
	case "createnatgateway":
		return func() interface{} { return NewCreateNatgateway(f.Sess, f.Graph, f.Log) }
	case "createnetworkinterface":
		return func() interface{} { return NewCreateNetworkinterface(f.Sess, f.Graph, f.Log) }
	case "createpolicy":
		return func() interface{} { return NewCreatePolicy(f.Sess, f.Graph, f.Log) }
	case "createqueue":
		return func() interface{} { return NewCreateQueue(f.Sess, f.Graph, f.Log) }
	case "createrecord":
		return func() interface{} { return NewCreateRecord(f.Sess, f.Graph, f.Log) }
	case "createrepository":
		return func() interface{} { return NewCreateRepository(f.Sess, f.Graph, f.Log) }
	case "createrole":
		return func() interface{} { return NewCreateRole(f.Sess, f.Graph, f.Log) }
	case "createroute":
		return func() interface{} { return NewCreateRoute(f.Sess, f.Graph, f.Log) }
	case "createroutetable":
		return func() interface{} { return NewCreateRoutetable(f.Sess, f.Graph, f.Log) }
	case "creates3object":
		return func() interface{} { return NewCreateS3object(f.Sess, f.Graph, f.Log) }
	case "createscalinggroup":
		return func() interface{} { return NewCreateScalinggroup(f.Sess, f.Graph, f.Log) }
	case "createscalingpolicy":
		return func() interface{} { return NewCreateScalingpolicy(f.Sess, f.Graph, f.Log) }
	case "createsecuritygroup":
		return func() interface{} { return NewCreateSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "createsnapshot":
		return func() interface{} { return NewCreateSnapshot(f.Sess, f.Graph, f.Log) }
	case "createstack":
		return func() interface{} { return NewCreateStack(f.Sess, f.Graph, f.Log) }
	case "createsubnet":
		return func() interface{} { return NewCreateSubnet(f.Sess, f.Graph, f.Log) }
	case "createsubscription":
		return func() interface{} { return NewCreateSubscription(f.Sess, f.Graph, f.Log) }
	case "createtag":
		return func() interface{} { return NewCreateTag(f.Sess, f.Graph, f.Log) }
	case "createtargetgroup":
		return func() interface{} { return NewCreateTargetgroup(f.Sess, f.Graph, f.Log) }
	case "createtopic":
		return func() interface{} { return NewCreateTopic(f.Sess, f.Graph, f.Log) }
	case "createuser":
		return func() interface{} { return NewCreateUser(f.Sess, f.Graph, f.Log) }
	case "createvolume":
		return func() interface{} { return NewCreateVolume(f.Sess, f.Graph, f.Log) }
	case "createvpc":
		return func() interface{} { return NewCreateVpc(f.Sess, f.Graph, f.Log) }
	case "createzone":
		return func() interface{} { return NewCreateZone(f.Sess, f.Graph, f.Log) }
	case "deleteaccesskey":
		return func() interface{} { return NewDeleteAccesskey(f.Sess, f.Graph, f.Log) }
	case "deletealarm":
		return func() interface{} { return NewDeleteAlarm(f.Sess, f.Graph, f.Log) }
	case "deleteappscalingpolicy":
		return func() interface{} { return NewDeleteAppscalingpolicy(f.Sess, f.Graph, f.Log) }
	case "deleteappscalingtarget":
		return func() interface{} { return NewDeleteAppscalingtarget(f.Sess, f.Graph, f.Log) }
	case "deletebucket":
		return func() interface{} { return NewDeleteBucket(f.Sess, f.Graph, f.Log) }
	case "deletecertificate":
		return func() interface{} { return NewDeleteCertificate(f.Sess, f.Graph, f.Log) }
	case "deleteclassicloadbalancer":
		return func() interface{} { return NewDeleteClassicLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "deletecontainercluster":
		return func() interface{} { return NewDeleteContainercluster(f.Sess, f.Graph, f.Log) }
	case "deletecontainertask":
		return func() interface{} { return NewDeleteContainertask(f.Sess, f.Graph, f.Log) }
	case "deletedatabase":
		return func() interface{} { return NewDeleteDatabase(f.Sess, f.Graph, f.Log) }
	case "deletedbsubnetgroup":
		return func() interface{} { return NewDeleteDbsubnetgroup(f.Sess, f.Graph, f.Log) }
	case "deletedistribution":
		return func() interface{} { return NewDeleteDistribution(f.Sess, f.Graph, f.Log) }
	case "deleteelasticip":
		return func() interface{} { return NewDeleteElasticip(f.Sess, f.Graph, f.Log) }
	case "deletefunction":
		return func() interface{} { return NewDeleteFunction(f.Sess, f.Graph, f.Log) }
	case "deletegroup":
		return func() interface{} { return NewDeleteGroup(f.Sess, f.Graph, f.Log) }
	case "deleteimage":
		return func() interface{} { return NewDeleteImage(f.Sess, f.Graph, f.Log) }
	case "deleteinstance":
		return func() interface{} { return NewDeleteInstance(f.Sess, f.Graph, f.Log) }
	case "deleteinstanceprofile":
		return func() interface{} { return NewDeleteInstanceprofile(f.Sess, f.Graph, f.Log) }
	case "deleteinternetgateway":
		return func() interface{} { return NewDeleteInternetgateway(f.Sess, f.Graph, f.Log) }
	case "deletekeypair":
		return func() interface{} { return NewDeleteKeypair(f.Sess, f.Graph, f.Log) }
	case "deletelaunchconfiguration":
		return func() interface{} { return NewDeleteLaunchconfiguration(f.Sess, f.Graph, f.Log) }
	case "deletelistener":
		return func() interface{} { return NewDeleteListener(f.Sess, f.Graph, f.Log) }
	case "deleteloadbalancer":
		return func() interface{} { return NewDeleteLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "deleteloginprofile":
		return func() interface{} { return NewDeleteLoginprofile(f.Sess, f.Graph, f.Log) }
	case "deletemfadevice":
		return func() interface{} { return NewDeleteMfadevice(f.Sess, f.Graph, f.Log) }
	case "deletenatgateway":
		return func() interface{} { return NewDeleteNatgateway(f.Sess, f.Graph, f.Log) }
	case "deletenetworkinterface":
		return func() interface{} { return NewDeleteNetworkinterface(f.Sess, f.Graph, f.Log) }
	case "deletepolicy":
		return func() interface{} { return NewDeletePolicy(f.Sess, f.Graph, f.Log) }
	case "deletequeue":
		return func() interface{} { return NewDeleteQueue(f.Sess, f.Graph, f.Log) }
	case "deleterecord":
		return func() interface{} { return NewDeleteRecord(f.Sess, f.Graph, f.Log) }
	case "deleterepository":
		return func() interface{} { return NewDeleteRepository(f.Sess, f.Graph, f.Log) }
	case "deleterole":
		return func() interface{} { return NewDeleteRole(f.Sess, f.Graph, f.Log) }
	case "deleteroute":
		return func() interface{} { return NewDeleteRoute(f.Sess, f.Graph, f.Log) }
	case "deleteroutetable":
		return func() interface{} { return NewDeleteRoutetable(f.Sess, f.Graph, f.Log) }
	case "deletes3object":
		return func() interface{} { return NewDeleteS3object(f.Sess, f.Graph, f.Log) }
	case "deletescalinggroup":
		return func() interface{} { return NewDeleteScalinggroup(f.Sess, f.Graph, f.Log) }
	case "deletescalingpolicy":
		return func() interface{} { return NewDeleteScalingpolicy(f.Sess, f.Graph, f.Log) }
	case "deletesecuritygroup":
		return func() interface{} { return NewDeleteSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "deletesnapshot":
		return func() interface{} { return NewDeleteSnapshot(f.Sess, f.Graph, f.Log) }
	case "deletestack":
		return func() interface{} { return NewDeleteStack(f.Sess, f.Graph, f.Log) }
	case "deletesubnet":
		return func() interface{} { return NewDeleteSubnet(f.Sess, f.Graph, f.Log) }
	case "deletesubscription":
		return func() interface{} { return NewDeleteSubscription(f.Sess, f.Graph, f.Log) }
	case "deletetag":
		return func() interface{} { return NewDeleteTag(f.Sess, f.Graph, f.Log) }
	case "deletetargetgroup":
		return func() interface{} { return NewDeleteTargetgroup(f.Sess, f.Graph, f.Log) }
	case "deletetopic":
		return func() interface{} { return NewDeleteTopic(f.Sess, f.Graph, f.Log) }
	case "deleteuser":
		return func() interface{} { return NewDeleteUser(f.Sess, f.Graph, f.Log) }
	case "deletevolume":
		return func() interface{} { return NewDeleteVolume(f.Sess, f.Graph, f.Log) }
	case "deletevpc":
		return func() interface{} { return NewDeleteVpc(f.Sess, f.Graph, f.Log) }
	case "deletezone":
		return func() interface{} { return NewDeleteZone(f.Sess, f.Graph, f.Log) }
	case "detachalarm":
		return func() interface{} { return NewDetachAlarm(f.Sess, f.Graph, f.Log) }
	case "detachclassicloadbalancer":
		return func() interface{} { return NewDetachClassicLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "detachcontainertask":
		return func() interface{} { return NewDetachContainertask(f.Sess, f.Graph, f.Log) }
	case "detachelasticip":
		return func() interface{} { return NewDetachElasticip(f.Sess, f.Graph, f.Log) }
	case "detachinstance":
		return func() interface{} { return NewDetachInstance(f.Sess, f.Graph, f.Log) }
	case "detachinstanceprofile":
		return func() interface{} { return NewDetachInstanceprofile(f.Sess, f.Graph, f.Log) }
	case "detachinternetgateway":
		return func() interface{} { return NewDetachInternetgateway(f.Sess, f.Graph, f.Log) }
	case "detachmfadevice":
		return func() interface{} { return NewDetachMfadevice(f.Sess, f.Graph, f.Log) }
	case "detachnetworkinterface":
		return func() interface{} { return NewDetachNetworkinterface(f.Sess, f.Graph, f.Log) }
	case "detachpolicy":
		return func() interface{} { return NewDetachPolicy(f.Sess, f.Graph, f.Log) }
	case "detachrole":
		return func() interface{} { return NewDetachRole(f.Sess, f.Graph, f.Log) }
	case "detachroutetable":
		return func() interface{} { return NewDetachRoutetable(f.Sess, f.Graph, f.Log) }
	case "detachsecuritygroup":
		return func() interface{} { return NewDetachSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "detachuser":
		return func() interface{} { return NewDetachUser(f.Sess, f.Graph, f.Log) }
	case "detachvolume":
		return func() interface{} { return NewDetachVolume(f.Sess, f.Graph, f.Log) }
	case "importimage":
		return func() interface{} { return NewImportImage(f.Sess, f.Graph, f.Log) }
	case "restartdatabase":
		return func() interface{} { return NewRestartDatabase(f.Sess, f.Graph, f.Log) }
	case "restartinstance":
		return func() interface{} { return NewRestartInstance(f.Sess, f.Graph, f.Log) }
	case "startalarm":
		return func() interface{} { return NewStartAlarm(f.Sess, f.Graph, f.Log) }
	case "startcontainertask":
		return func() interface{} { return NewStartContainertask(f.Sess, f.Graph, f.Log) }
	case "startdatabase":
		return func() interface{} { return NewStartDatabase(f.Sess, f.Graph, f.Log) }
	case "startinstance":
		return func() interface{} { return NewStartInstance(f.Sess, f.Graph, f.Log) }
	case "stopalarm":
		return func() interface{} { return NewStopAlarm(f.Sess, f.Graph, f.Log) }
	case "stopcontainertask":
		return func() interface{} { return NewStopContainertask(f.Sess, f.Graph, f.Log) }
	case "stopdatabase":
		return func() interface{} { return NewStopDatabase(f.Sess, f.Graph, f.Log) }
	case "stopinstance":
		return func() interface{} { return NewStopInstance(f.Sess, f.Graph, f.Log) }
	case "updatebucket":
		return func() interface{} { return NewUpdateBucket(f.Sess, f.Graph, f.Log) }
	case "updateclassicloadbalancer":
		return func() interface{} { return NewUpdateClassicLoadbalancer(f.Sess, f.Graph, f.Log) }
	case "updatecontainertask":
		return func() interface{} { return NewUpdateContainertask(f.Sess, f.Graph, f.Log) }
	case "updatedistribution":
		return func() interface{} { return NewUpdateDistribution(f.Sess, f.Graph, f.Log) }
	case "updateimage":
		return func() interface{} { return NewUpdateImage(f.Sess, f.Graph, f.Log) }
	case "updateinstance":
		return func() interface{} { return NewUpdateInstance(f.Sess, f.Graph, f.Log) }
	case "updateloginprofile":
		return func() interface{} { return NewUpdateLoginprofile(f.Sess, f.Graph, f.Log) }
	case "updatepolicy":
		return func() interface{} { return NewUpdatePolicy(f.Sess, f.Graph, f.Log) }
	case "updaterecord":
		return func() interface{} { return NewUpdateRecord(f.Sess, f.Graph, f.Log) }
	case "updates3object":
		return func() interface{} { return NewUpdateS3object(f.Sess, f.Graph, f.Log) }
	case "updatescalinggroup":
		return func() interface{} { return NewUpdateScalinggroup(f.Sess, f.Graph, f.Log) }
	case "updatesecuritygroup":
		return func() interface{} { return NewUpdateSecuritygroup(f.Sess, f.Graph, f.Log) }
	case "updatestack":
		return func() interface{} { return NewUpdateStack(f.Sess, f.Graph, f.Log) }
	case "updatesubnet":
		return func() interface{} { return NewUpdateSubnet(f.Sess, f.Graph, f.Log) }
	case "updatetargetgroup":
		return func() interface{} { return NewUpdateTargetgroup(f.Sess, f.Graph, f.Log) }
	}
	return nil
}

var (
	_ command = &AttachAlarm{}
	_ command = &AttachClassicLoadbalancer{}
	_ command = &AttachContainertask{}
	_ command = &AttachElasticip{}
	_ command = &AttachInstance{}
	_ command = &AttachInstanceprofile{}
	_ command = &AttachInternetgateway{}
	_ command = &AttachListener{}
	_ command = &AttachMfadevice{}
	_ command = &AttachNetworkinterface{}
	_ command = &AttachPolicy{}
	_ command = &AttachRole{}
	_ command = &AttachRoutetable{}
	_ command = &AttachSecuritygroup{}
	_ command = &AttachUser{}
	_ command = &AttachVolume{}
	_ command = &AuthenticateRegistry{}
	_ command = &CheckCertificate{}
	_ command = &CheckDatabase{}
	_ command = &CheckDistribution{}
	_ command = &CheckInstance{}
	_ command = &CheckLoadbalancer{}
	_ command = &CheckNatgateway{}
	_ command = &CheckNetworkinterface{}
	_ command = &CheckScalinggroup{}
	_ command = &CheckSecuritygroup{}
	_ command = &CheckVolume{}
	_ command = &CopyImage{}
	_ command = &CopySnapshot{}
	_ command = &CreateAccesskey{}
	_ command = &CreateAlarm{}
	_ command = &CreateAppscalingpolicy{}
	_ command = &CreateAppscalingtarget{}
	_ command = &CreateBucket{}
	_ command = &CreateCertificate{}
	_ command = &CreateClassicLoadbalancer{}
	_ command = &CreateContainercluster{}
	_ command = &CreateDatabase{}
	_ command = &CreateDbsubnetgroup{}
	_ command = &CreateDistribution{}
	_ command = &CreateElasticip{}
	_ command = &CreateFunction{}
	_ command = &CreateGroup{}
	_ command = &CreateImage{}
	_ command = &CreateInstance{}
	_ command = &CreateInstanceprofile{}
	_ command = &CreateInternetgateway{}
	_ command = &CreateKeypair{}
	_ command = &CreateLaunchconfiguration{}
	_ command = &CreateListener{}
	_ command = &CreateLoadbalancer{}
	_ command = &CreateLoginprofile{}
	_ command = &CreateMfadevice{}
	_ command = &CreateNatgateway{}
	_ command = &CreateNetworkinterface{}
	_ command = &CreatePolicy{}
	_ command = &CreateQueue{}
	_ command = &CreateRecord{}
	_ command = &CreateRepository{}
	_ command = &CreateRole{}
	_ command = &CreateRoute{}
	_ command = &CreateRoutetable{}
	_ command = &CreateS3object{}
	_ command = &CreateScalinggroup{}
	_ command = &CreateScalingpolicy{}
	_ command = &CreateSecuritygroup{}
	_ command = &CreateSnapshot{}
	_ command = &CreateStack{}
	_ command = &CreateSubnet{}
	_ command = &CreateSubscription{}
	_ command = &CreateTag{}
	_ command = &CreateTargetgroup{}
	_ command = &CreateTopic{}
	_ command = &CreateUser{}
	_ command = &CreateVolume{}
	_ command = &CreateVpc{}
	_ command = &CreateZone{}
	_ command = &DeleteAccesskey{}
	_ command = &DeleteAlarm{}
	_ command = &DeleteAppscalingpolicy{}
	_ command = &DeleteAppscalingtarget{}
	_ command = &DeleteBucket{}
	_ command = &DeleteCertificate{}
	_ command = &DeleteClassicLoadbalancer{}
	_ command = &DeleteContainercluster{}
	_ command = &DeleteContainertask{}
	_ command = &DeleteDatabase{}
	_ command = &DeleteDbsubnetgroup{}
	_ command = &DeleteDistribution{}
	_ command = &DeleteElasticip{}
	_ command = &DeleteFunction{}
	_ command = &DeleteGroup{}
	_ command = &DeleteImage{}
	_ command = &DeleteInstance{}
	_ command = &DeleteInstanceprofile{}
	_ command = &DeleteInternetgateway{}
	_ command = &DeleteKeypair{}
	_ command = &DeleteLaunchconfiguration{}
	_ command = &DeleteListener{}
	_ command = &DeleteLoadbalancer{}
	_ command = &DeleteLoginprofile{}
	_ command = &DeleteMfadevice{}
	_ command = &DeleteNatgateway{}
	_ command = &DeleteNetworkinterface{}
	_ command = &DeletePolicy{}
	_ command = &DeleteQueue{}
	_ command = &DeleteRecord{}
	_ command = &DeleteRepository{}
	_ command = &DeleteRole{}
	_ command = &DeleteRoute{}
	_ command = &DeleteRoutetable{}
	_ command = &DeleteS3object{}
	_ command = &DeleteScalinggroup{}
	_ command = &DeleteScalingpolicy{}
	_ command = &DeleteSecuritygroup{}
	_ command = &DeleteSnapshot{}
	_ command = &DeleteStack{}
	_ command = &DeleteSubnet{}
	_ command = &DeleteSubscription{}
	_ command = &DeleteTag{}
	_ command = &DeleteTargetgroup{}
	_ command = &DeleteTopic{}
	_ command = &DeleteUser{}
	_ command = &DeleteVolume{}
	_ command = &DeleteVpc{}
	_ command = &DeleteZone{}
	_ command = &DetachAlarm{}
	_ command = &DetachClassicLoadbalancer{}
	_ command = &DetachContainertask{}
	_ command = &DetachElasticip{}
	_ command = &DetachInstance{}
	_ command = &DetachInstanceprofile{}
	_ command = &DetachInternetgateway{}
	_ command = &DetachMfadevice{}
	_ command = &DetachNetworkinterface{}
	_ command = &DetachPolicy{}
	_ command = &DetachRole{}
	_ command = &DetachRoutetable{}
	_ command = &DetachSecuritygroup{}
	_ command = &DetachUser{}
	_ command = &DetachVolume{}
	_ command = &ImportImage{}
	_ command = &RestartDatabase{}
	_ command = &RestartInstance{}
	_ command = &StartAlarm{}
	_ command = &StartContainertask{}
	_ command = &StartDatabase{}
	_ command = &StartInstance{}
	_ command = &StopAlarm{}
	_ command = &StopContainertask{}
	_ command = &StopDatabase{}
	_ command = &StopInstance{}
	_ command = &UpdateBucket{}
	_ command = &UpdateClassicLoadbalancer{}
	_ command = &UpdateContainertask{}
	_ command = &UpdateDistribution{}
	_ command = &UpdateImage{}
	_ command = &UpdateInstance{}
	_ command = &UpdateLoginprofile{}
	_ command = &UpdatePolicy{}
	_ command = &UpdateRecord{}
	_ command = &UpdateS3object{}
	_ command = &UpdateScalinggroup{}
	_ command = &UpdateSecuritygroup{}
	_ command = &UpdateStack{}
	_ command = &UpdateSubnet{}
	_ command = &UpdateTargetgroup{}
)

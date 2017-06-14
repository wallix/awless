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
package awsdriver

import (
	"strings"

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
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

type Ec2Driver struct {
	dryRun bool
	logger *logger.Logger
	ec2iface.EC2API
}

func (d *Ec2Driver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *Ec2Driver) SetLogger(l *logger.Logger) { d.logger = l }
func NewEc2Driver(api ec2iface.EC2API) driver.Driver {
	return &Ec2Driver{false, logger.DiscardLogger, api}
}

func (d *Ec2Driver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createvpc":
		if d.dryRun {
			return d.Create_Vpc_DryRun, nil
		}
		return d.Create_Vpc, nil

	case "deletevpc":
		if d.dryRun {
			return d.Delete_Vpc_DryRun, nil
		}
		return d.Delete_Vpc, nil

	case "createsubnet":
		if d.dryRun {
			return d.Create_Subnet_DryRun, nil
		}
		return d.Create_Subnet, nil

	case "updatesubnet":
		if d.dryRun {
			return d.Update_Subnet_DryRun, nil
		}
		return d.Update_Subnet, nil

	case "deletesubnet":
		if d.dryRun {
			return d.Delete_Subnet_DryRun, nil
		}
		return d.Delete_Subnet, nil

	case "createinstance":
		if d.dryRun {
			return d.Create_Instance_DryRun, nil
		}
		return d.Create_Instance, nil

	case "updateinstance":
		if d.dryRun {
			return d.Update_Instance_DryRun, nil
		}
		return d.Update_Instance, nil

	case "deleteinstance":
		if d.dryRun {
			return d.Delete_Instance_DryRun, nil
		}
		return d.Delete_Instance, nil

	case "startinstance":
		if d.dryRun {
			return d.Start_Instance_DryRun, nil
		}
		return d.Start_Instance, nil

	case "stopinstance":
		if d.dryRun {
			return d.Stop_Instance_DryRun, nil
		}
		return d.Stop_Instance, nil

	case "checkinstance":
		if d.dryRun {
			return d.Check_Instance_DryRun, nil
		}
		return d.Check_Instance, nil

	case "createsecuritygroup":
		if d.dryRun {
			return d.Create_Securitygroup_DryRun, nil
		}
		return d.Create_Securitygroup, nil

	case "updatesecuritygroup":
		if d.dryRun {
			return d.Update_Securitygroup_DryRun, nil
		}
		return d.Update_Securitygroup, nil

	case "deletesecuritygroup":
		if d.dryRun {
			return d.Delete_Securitygroup_DryRun, nil
		}
		return d.Delete_Securitygroup, nil

	case "checksecuritygroup":
		if d.dryRun {
			return d.Check_Securitygroup_DryRun, nil
		}
		return d.Check_Securitygroup, nil

	case "attachsecuritygroup":
		if d.dryRun {
			return d.Attach_Securitygroup_DryRun, nil
		}
		return d.Attach_Securitygroup, nil

	case "detachsecuritygroup":
		if d.dryRun {
			return d.Detach_Securitygroup_DryRun, nil
		}
		return d.Detach_Securitygroup, nil

	case "copyimage":
		if d.dryRun {
			return d.Copy_Image_DryRun, nil
		}
		return d.Copy_Image, nil

	case "importimage":
		if d.dryRun {
			return d.Import_Image_DryRun, nil
		}
		return d.Import_Image, nil

	case "deleteimage":
		if d.dryRun {
			return d.Delete_Image_DryRun, nil
		}
		return d.Delete_Image, nil

	case "createvolume":
		if d.dryRun {
			return d.Create_Volume_DryRun, nil
		}
		return d.Create_Volume, nil

	case "checkvolume":
		if d.dryRun {
			return d.Check_Volume_DryRun, nil
		}
		return d.Check_Volume, nil

	case "deletevolume":
		if d.dryRun {
			return d.Delete_Volume_DryRun, nil
		}
		return d.Delete_Volume, nil

	case "attachvolume":
		if d.dryRun {
			return d.Attach_Volume_DryRun, nil
		}
		return d.Attach_Volume, nil

	case "detachvolume":
		if d.dryRun {
			return d.Detach_Volume_DryRun, nil
		}
		return d.Detach_Volume, nil

	case "createsnapshot":
		if d.dryRun {
			return d.Create_Snapshot_DryRun, nil
		}
		return d.Create_Snapshot, nil

	case "deletesnapshot":
		if d.dryRun {
			return d.Delete_Snapshot_DryRun, nil
		}
		return d.Delete_Snapshot, nil

	case "copysnapshot":
		if d.dryRun {
			return d.Copy_Snapshot_DryRun, nil
		}
		return d.Copy_Snapshot, nil

	case "createinternetgateway":
		if d.dryRun {
			return d.Create_Internetgateway_DryRun, nil
		}
		return d.Create_Internetgateway, nil

	case "deleteinternetgateway":
		if d.dryRun {
			return d.Delete_Internetgateway_DryRun, nil
		}
		return d.Delete_Internetgateway, nil

	case "attachinternetgateway":
		if d.dryRun {
			return d.Attach_Internetgateway_DryRun, nil
		}
		return d.Attach_Internetgateway, nil

	case "detachinternetgateway":
		if d.dryRun {
			return d.Detach_Internetgateway_DryRun, nil
		}
		return d.Detach_Internetgateway, nil

	case "createroutetable":
		if d.dryRun {
			return d.Create_Routetable_DryRun, nil
		}
		return d.Create_Routetable, nil

	case "deleteroutetable":
		if d.dryRun {
			return d.Delete_Routetable_DryRun, nil
		}
		return d.Delete_Routetable, nil

	case "attachroutetable":
		if d.dryRun {
			return d.Attach_Routetable_DryRun, nil
		}
		return d.Attach_Routetable, nil

	case "detachroutetable":
		if d.dryRun {
			return d.Detach_Routetable_DryRun, nil
		}
		return d.Detach_Routetable, nil

	case "createroute":
		if d.dryRun {
			return d.Create_Route_DryRun, nil
		}
		return d.Create_Route, nil

	case "deleteroute":
		if d.dryRun {
			return d.Delete_Route_DryRun, nil
		}
		return d.Delete_Route, nil

	case "createtag":
		if d.dryRun {
			return d.Create_Tag_DryRun, nil
		}
		return d.Create_Tag, nil

	case "deletetag":
		if d.dryRun {
			return d.Delete_Tag_DryRun, nil
		}
		return d.Delete_Tag, nil

	case "createkeypair":
		if d.dryRun {
			return d.Create_Keypair_DryRun, nil
		}
		return d.Create_Keypair, nil

	case "deletekeypair":
		if d.dryRun {
			return d.Delete_Keypair_DryRun, nil
		}
		return d.Delete_Keypair, nil

	case "createelasticip":
		if d.dryRun {
			return d.Create_Elasticip_DryRun, nil
		}
		return d.Create_Elasticip, nil

	case "deleteelasticip":
		if d.dryRun {
			return d.Delete_Elasticip_DryRun, nil
		}
		return d.Delete_Elasticip, nil

	case "attachelasticip":
		if d.dryRun {
			return d.Attach_Elasticip_DryRun, nil
		}
		return d.Attach_Elasticip, nil

	case "detachelasticip":
		if d.dryRun {
			return d.Detach_Elasticip_DryRun, nil
		}
		return d.Detach_Elasticip, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type Elbv2Driver struct {
	dryRun bool
	logger *logger.Logger
	elbv2iface.ELBV2API
}

func (d *Elbv2Driver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *Elbv2Driver) SetLogger(l *logger.Logger) { d.logger = l }
func NewElbv2Driver(api elbv2iface.ELBV2API) driver.Driver {
	return &Elbv2Driver{false, logger.DiscardLogger, api}
}

func (d *Elbv2Driver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createloadbalancer":
		if d.dryRun {
			return d.Create_Loadbalancer_DryRun, nil
		}
		return d.Create_Loadbalancer, nil

	case "deleteloadbalancer":
		if d.dryRun {
			return d.Delete_Loadbalancer_DryRun, nil
		}
		return d.Delete_Loadbalancer, nil

	case "checkloadbalancer":
		if d.dryRun {
			return d.Check_Loadbalancer_DryRun, nil
		}
		return d.Check_Loadbalancer, nil

	case "createlistener":
		if d.dryRun {
			return d.Create_Listener_DryRun, nil
		}
		return d.Create_Listener, nil

	case "deletelistener":
		if d.dryRun {
			return d.Delete_Listener_DryRun, nil
		}
		return d.Delete_Listener, nil

	case "createtargetgroup":
		if d.dryRun {
			return d.Create_Targetgroup_DryRun, nil
		}
		return d.Create_Targetgroup, nil

	case "deletetargetgroup":
		if d.dryRun {
			return d.Delete_Targetgroup_DryRun, nil
		}
		return d.Delete_Targetgroup, nil

	case "attachinstance":
		if d.dryRun {
			return d.Attach_Instance_DryRun, nil
		}
		return d.Attach_Instance, nil

	case "detachinstance":
		if d.dryRun {
			return d.Detach_Instance_DryRun, nil
		}
		return d.Detach_Instance, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type AutoscalingDriver struct {
	dryRun bool
	logger *logger.Logger
	autoscalingiface.AutoScalingAPI
}

func (d *AutoscalingDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *AutoscalingDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewAutoscalingDriver(api autoscalingiface.AutoScalingAPI) driver.Driver {
	return &AutoscalingDriver{false, logger.DiscardLogger, api}
}

func (d *AutoscalingDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createlaunchconfiguration":
		if d.dryRun {
			return d.Create_Launchconfiguration_DryRun, nil
		}
		return d.Create_Launchconfiguration, nil

	case "deletelaunchconfiguration":
		if d.dryRun {
			return d.Delete_Launchconfiguration_DryRun, nil
		}
		return d.Delete_Launchconfiguration, nil

	case "createscalinggroup":
		if d.dryRun {
			return d.Create_Scalinggroup_DryRun, nil
		}
		return d.Create_Scalinggroup, nil

	case "updatescalinggroup":
		if d.dryRun {
			return d.Update_Scalinggroup_DryRun, nil
		}
		return d.Update_Scalinggroup, nil

	case "deletescalinggroup":
		if d.dryRun {
			return d.Delete_Scalinggroup_DryRun, nil
		}
		return d.Delete_Scalinggroup, nil

	case "checkscalinggroup":
		if d.dryRun {
			return d.Check_Scalinggroup_DryRun, nil
		}
		return d.Check_Scalinggroup, nil

	case "createscalingpolicy":
		if d.dryRun {
			return d.Create_Scalingpolicy_DryRun, nil
		}
		return d.Create_Scalingpolicy, nil

	case "deletescalingpolicy":
		if d.dryRun {
			return d.Delete_Scalingpolicy_DryRun, nil
		}
		return d.Delete_Scalingpolicy, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type RdsDriver struct {
	dryRun bool
	logger *logger.Logger
	rdsiface.RDSAPI
}

func (d *RdsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *RdsDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewRdsDriver(api rdsiface.RDSAPI) driver.Driver {
	return &RdsDriver{false, logger.DiscardLogger, api}
}

func (d *RdsDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createdatabase":
		if d.dryRun {
			return d.Create_Database_DryRun, nil
		}
		return d.Create_Database, nil

	case "deletedatabase":
		if d.dryRun {
			return d.Delete_Database_DryRun, nil
		}
		return d.Delete_Database, nil

	case "checkdatabase":
		if d.dryRun {
			return d.Check_Database_DryRun, nil
		}
		return d.Check_Database, nil

	case "createdbsubnetgroup":
		if d.dryRun {
			return d.Create_Dbsubnetgroup_DryRun, nil
		}
		return d.Create_Dbsubnetgroup, nil

	case "deletedbsubnetgroup":
		if d.dryRun {
			return d.Delete_Dbsubnetgroup_DryRun, nil
		}
		return d.Delete_Dbsubnetgroup, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type EcrDriver struct {
	dryRun bool
	logger *logger.Logger
	ecriface.ECRAPI
}

func (d *EcrDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *EcrDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewEcrDriver(api ecriface.ECRAPI) driver.Driver {
	return &EcrDriver{false, logger.DiscardLogger, api}
}

func (d *EcrDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createrepository":
		if d.dryRun {
			return d.Create_Repository_DryRun, nil
		}
		return d.Create_Repository, nil

	case "deleterepository":
		if d.dryRun {
			return d.Delete_Repository_DryRun, nil
		}
		return d.Delete_Repository, nil

	case "authenticateregistry":
		if d.dryRun {
			return d.Authenticate_Registry_DryRun, nil
		}
		return d.Authenticate_Registry, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type EcsDriver struct {
	dryRun bool
	logger *logger.Logger
	ecsiface.ECSAPI
}

func (d *EcsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *EcsDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewEcsDriver(api ecsiface.ECSAPI) driver.Driver {
	return &EcsDriver{false, logger.DiscardLogger, api}
}

func (d *EcsDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createcontainercluster":
		if d.dryRun {
			return d.Create_Containercluster_DryRun, nil
		}
		return d.Create_Containercluster, nil

	case "deletecontainercluster":
		if d.dryRun {
			return d.Delete_Containercluster_DryRun, nil
		}
		return d.Delete_Containercluster, nil

	case "startcontainerservice":
		if d.dryRun {
			return d.Start_Containerservice_DryRun, nil
		}
		return d.Start_Containerservice, nil

	case "stopcontainerservice":
		if d.dryRun {
			return d.Stop_Containerservice_DryRun, nil
		}
		return d.Stop_Containerservice, nil

	case "updatecontainerservice":
		if d.dryRun {
			return d.Update_Containerservice_DryRun, nil
		}
		return d.Update_Containerservice, nil

	case "createcontainer":
		if d.dryRun {
			return d.Create_Container_DryRun, nil
		}
		return d.Create_Container, nil

	case "deletecontainer":
		if d.dryRun {
			return d.Delete_Container_DryRun, nil
		}
		return d.Delete_Container, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type StsDriver struct {
	dryRun bool
	logger *logger.Logger
	stsiface.STSAPI
}

func (d *StsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *StsDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewStsDriver(api stsiface.STSAPI) driver.Driver {
	return &StsDriver{false, logger.DiscardLogger, api}
}

func (d *StsDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type IamDriver struct {
	dryRun bool
	logger *logger.Logger
	iamiface.IAMAPI
}

func (d *IamDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *IamDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewIamDriver(api iamiface.IAMAPI) driver.Driver {
	return &IamDriver{false, logger.DiscardLogger, api}
}

func (d *IamDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createuser":
		if d.dryRun {
			return d.Create_User_DryRun, nil
		}
		return d.Create_User, nil

	case "deleteuser":
		if d.dryRun {
			return d.Delete_User_DryRun, nil
		}
		return d.Delete_User, nil

	case "attachuser":
		if d.dryRun {
			return d.Attach_User_DryRun, nil
		}
		return d.Attach_User, nil

	case "detachuser":
		if d.dryRun {
			return d.Detach_User_DryRun, nil
		}
		return d.Detach_User, nil

	case "createaccesskey":
		if d.dryRun {
			return d.Create_Accesskey_DryRun, nil
		}
		return d.Create_Accesskey, nil

	case "deleteaccesskey":
		if d.dryRun {
			return d.Delete_Accesskey_DryRun, nil
		}
		return d.Delete_Accesskey, nil

	case "createloginprofile":
		if d.dryRun {
			return d.Create_Loginprofile_DryRun, nil
		}
		return d.Create_Loginprofile, nil

	case "updateloginprofile":
		if d.dryRun {
			return d.Update_Loginprofile_DryRun, nil
		}
		return d.Update_Loginprofile, nil

	case "deleteloginprofile":
		if d.dryRun {
			return d.Delete_Loginprofile_DryRun, nil
		}
		return d.Delete_Loginprofile, nil

	case "creategroup":
		if d.dryRun {
			return d.Create_Group_DryRun, nil
		}
		return d.Create_Group, nil

	case "deletegroup":
		if d.dryRun {
			return d.Delete_Group_DryRun, nil
		}
		return d.Delete_Group, nil

	case "createrole":
		if d.dryRun {
			return d.Create_Role_DryRun, nil
		}
		return d.Create_Role, nil

	case "deleterole":
		if d.dryRun {
			return d.Delete_Role_DryRun, nil
		}
		return d.Delete_Role, nil

	case "attachrole":
		if d.dryRun {
			return d.Attach_Role_DryRun, nil
		}
		return d.Attach_Role, nil

	case "detachrole":
		if d.dryRun {
			return d.Detach_Role_DryRun, nil
		}
		return d.Detach_Role, nil

	case "createinstanceprofile":
		if d.dryRun {
			return d.Create_Instanceprofile_DryRun, nil
		}
		return d.Create_Instanceprofile, nil

	case "deleteinstanceprofile":
		if d.dryRun {
			return d.Delete_Instanceprofile_DryRun, nil
		}
		return d.Delete_Instanceprofile, nil

	case "createpolicy":
		if d.dryRun {
			return d.Create_Policy_DryRun, nil
		}
		return d.Create_Policy, nil

	case "deletepolicy":
		if d.dryRun {
			return d.Delete_Policy_DryRun, nil
		}
		return d.Delete_Policy, nil

	case "attachpolicy":
		if d.dryRun {
			return d.Attach_Policy_DryRun, nil
		}
		return d.Attach_Policy, nil

	case "detachpolicy":
		if d.dryRun {
			return d.Detach_Policy_DryRun, nil
		}
		return d.Detach_Policy, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type S3Driver struct {
	dryRun bool
	logger *logger.Logger
	s3iface.S3API
}

func (d *S3Driver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *S3Driver) SetLogger(l *logger.Logger) { d.logger = l }
func NewS3Driver(api s3iface.S3API) driver.Driver {
	return &S3Driver{false, logger.DiscardLogger, api}
}

func (d *S3Driver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createbucket":
		if d.dryRun {
			return d.Create_Bucket_DryRun, nil
		}
		return d.Create_Bucket, nil

	case "updatebucket":
		if d.dryRun {
			return d.Update_Bucket_DryRun, nil
		}
		return d.Update_Bucket, nil

	case "deletebucket":
		if d.dryRun {
			return d.Delete_Bucket_DryRun, nil
		}
		return d.Delete_Bucket, nil

	case "creates3object":
		if d.dryRun {
			return d.Create_S3object_DryRun, nil
		}
		return d.Create_S3object, nil

	case "updates3object":
		if d.dryRun {
			return d.Update_S3object_DryRun, nil
		}
		return d.Update_S3object, nil

	case "deletes3object":
		if d.dryRun {
			return d.Delete_S3object_DryRun, nil
		}
		return d.Delete_S3object, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type SnsDriver struct {
	dryRun bool
	logger *logger.Logger
	snsiface.SNSAPI
}

func (d *SnsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *SnsDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewSnsDriver(api snsiface.SNSAPI) driver.Driver {
	return &SnsDriver{false, logger.DiscardLogger, api}
}

func (d *SnsDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createtopic":
		if d.dryRun {
			return d.Create_Topic_DryRun, nil
		}
		return d.Create_Topic, nil

	case "deletetopic":
		if d.dryRun {
			return d.Delete_Topic_DryRun, nil
		}
		return d.Delete_Topic, nil

	case "createsubscription":
		if d.dryRun {
			return d.Create_Subscription_DryRun, nil
		}
		return d.Create_Subscription, nil

	case "deletesubscription":
		if d.dryRun {
			return d.Delete_Subscription_DryRun, nil
		}
		return d.Delete_Subscription, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type SqsDriver struct {
	dryRun bool
	logger *logger.Logger
	sqsiface.SQSAPI
}

func (d *SqsDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *SqsDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewSqsDriver(api sqsiface.SQSAPI) driver.Driver {
	return &SqsDriver{false, logger.DiscardLogger, api}
}

func (d *SqsDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createqueue":
		if d.dryRun {
			return d.Create_Queue_DryRun, nil
		}
		return d.Create_Queue, nil

	case "deletequeue":
		if d.dryRun {
			return d.Delete_Queue_DryRun, nil
		}
		return d.Delete_Queue, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type Route53Driver struct {
	dryRun bool
	logger *logger.Logger
	route53iface.Route53API
}

func (d *Route53Driver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *Route53Driver) SetLogger(l *logger.Logger) { d.logger = l }
func NewRoute53Driver(api route53iface.Route53API) driver.Driver {
	return &Route53Driver{false, logger.DiscardLogger, api}
}

func (d *Route53Driver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createzone":
		if d.dryRun {
			return d.Create_Zone_DryRun, nil
		}
		return d.Create_Zone, nil

	case "deletezone":
		if d.dryRun {
			return d.Delete_Zone_DryRun, nil
		}
		return d.Delete_Zone, nil

	case "createrecord":
		if d.dryRun {
			return d.Create_Record_DryRun, nil
		}
		return d.Create_Record, nil

	case "deleterecord":
		if d.dryRun {
			return d.Delete_Record_DryRun, nil
		}
		return d.Delete_Record, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type LambdaDriver struct {
	dryRun bool
	logger *logger.Logger
	lambdaiface.LambdaAPI
}

func (d *LambdaDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *LambdaDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewLambdaDriver(api lambdaiface.LambdaAPI) driver.Driver {
	return &LambdaDriver{false, logger.DiscardLogger, api}
}

func (d *LambdaDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createfunction":
		if d.dryRun {
			return d.Create_Function_DryRun, nil
		}
		return d.Create_Function, nil

	case "deletefunction":
		if d.dryRun {
			return d.Delete_Function_DryRun, nil
		}
		return d.Delete_Function, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type CloudwatchDriver struct {
	dryRun bool
	logger *logger.Logger
	cloudwatchiface.CloudWatchAPI
}

func (d *CloudwatchDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *CloudwatchDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewCloudwatchDriver(api cloudwatchiface.CloudWatchAPI) driver.Driver {
	return &CloudwatchDriver{false, logger.DiscardLogger, api}
}

func (d *CloudwatchDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createalarm":
		if d.dryRun {
			return d.Create_Alarm_DryRun, nil
		}
		return d.Create_Alarm, nil

	case "deletealarm":
		if d.dryRun {
			return d.Delete_Alarm_DryRun, nil
		}
		return d.Delete_Alarm, nil

	case "startalarm":
		if d.dryRun {
			return d.Start_Alarm_DryRun, nil
		}
		return d.Start_Alarm, nil

	case "stopalarm":
		if d.dryRun {
			return d.Stop_Alarm_DryRun, nil
		}
		return d.Stop_Alarm, nil

	case "attachalarm":
		if d.dryRun {
			return d.Attach_Alarm_DryRun, nil
		}
		return d.Attach_Alarm, nil

	case "detachalarm":
		if d.dryRun {
			return d.Detach_Alarm_DryRun, nil
		}
		return d.Detach_Alarm, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type CloudfrontDriver struct {
	dryRun bool
	logger *logger.Logger
	cloudfrontiface.CloudFrontAPI
}

func (d *CloudfrontDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *CloudfrontDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewCloudfrontDriver(api cloudfrontiface.CloudFrontAPI) driver.Driver {
	return &CloudfrontDriver{false, logger.DiscardLogger, api}
}

func (d *CloudfrontDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createdistribution":
		if d.dryRun {
			return d.Create_Distribution_DryRun, nil
		}
		return d.Create_Distribution, nil

	case "checkdistribution":
		if d.dryRun {
			return d.Check_Distribution_DryRun, nil
		}
		return d.Check_Distribution, nil

	case "updatedistribution":
		if d.dryRun {
			return d.Update_Distribution_DryRun, nil
		}
		return d.Update_Distribution, nil

	case "deletedistribution":
		if d.dryRun {
			return d.Delete_Distribution_DryRun, nil
		}
		return d.Delete_Distribution, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

type CloudformationDriver struct {
	dryRun bool
	logger *logger.Logger
	cloudformationiface.CloudFormationAPI
}

func (d *CloudformationDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *CloudformationDriver) SetLogger(l *logger.Logger) { d.logger = l }
func NewCloudformationDriver(api cloudformationiface.CloudFormationAPI) driver.Driver {
	return &CloudformationDriver{false, logger.DiscardLogger, api}
}

func (d *CloudformationDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	switch strings.Join(lookups, "") {

	case "createstack":
		if d.dryRun {
			return d.Create_Stack_DryRun, nil
		}
		return d.Create_Stack, nil

	case "updatestack":
		if d.dryRun {
			return d.Update_Stack_DryRun, nil
		}
		return d.Update_Stack, nil

	case "deletestack":
		if d.dryRun {
			return d.Delete_Stack_DryRun, nil
		}
		return d.Delete_Stack, nil

	default:
		return nil, driver.ErrDriverFnNotFound
	}
}

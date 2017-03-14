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
package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
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

	case "createvolume":
		if d.dryRun {
			return d.Create_Volume_DryRun, nil
		}
		return d.Create_Volume, nil

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

	case "deletebucket":
		if d.dryRun {
			return d.Delete_Bucket_DryRun, nil
		}
		return d.Delete_Bucket, nil

	case "createstorageobject":
		if d.dryRun {
			return d.Create_Storageobject_DryRun, nil
		}
		return d.Create_Storageobject, nil

	case "deletestorageobject":
		if d.dryRun {
			return d.Delete_Storageobject_DryRun, nil
		}
		return d.Delete_Storageobject, nil

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

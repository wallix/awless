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
	"errors"
	"fmt"
	"time"

	"github.com/wallix/awless/cloud/graph"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/wallix/awless/logger"
)

type CreateDatabase struct {
	_      string `action:"create" entity:"database" awsAPI:"rds"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    rdsiface.RDSAPI

	// Required for DB
	Type     *string `awsName:"DBInstanceClass" awsType:"awsstr" templateName:"type"`
	Id       *string `awsName:"DBInstanceIdentifier" awsType:"awsstr" templateName:"id"`
	Engine   *string `awsName:"Engine" awsType:"awsstr" templateName:"engine"`
	Password *string `awsName:"MasterUserPassword" awsType:"awsstr" templateName:"password"`
	Username *string `awsName:"MasterUsername" awsType:"awsstr" templateName:"username"`
	Size     *int64  `awsName:"AllocatedStorage" awsType:"awsint64" templateName:"size"`

	// Required for read replica DB
	ReadReplicaSourceDB   *string `awsName:"SourceDBInstanceIdentifier" awsType:"awsstr" templateName:"replica-source"`
	ReadReplicaIdentifier *string `awsName:"DBInstanceIdentifier" awsType:"awsstr" templateName:"replica"`

	// Extras common to both replica DB and source DB
	Autoupgrade      *bool   `awsName:"AutoMinorVersionUpgrade" awsType:"awsbool" templateName:"autoupgrade"`
	Availabilityzone *string `awsName:"AvailabilityZone" awsType:"awsstr" templateName:"availabilityzone"`
	Subnetgroup      *string `awsName:"DBSubnetGroupName" awsType:"awsstr" templateName:"subnetgroup"`
	Iops             *int64  `awsName:"Iops" awsType:"awsint64" templateName:"iops"`
	Optiongroup      *string `awsName:"OptionGroupName" awsType:"awsstr" templateName:"optiongroup"`
	Port             *int64  `awsName:"Port" awsType:"awsint64" templateName:"port"`
	Public           *bool   `awsName:"PubliclyAccessible" awsType:"awsbool" templateName:"public"`
	Storagetype      *string `awsName:"StorageType" awsType:"awsstr" templateName:"storagetype"`

	// Extra only for DB
	Backupretention   *int64    `awsName:"BackupRetentionPeriod" awsType:"awsint64" templateName:"backupretention"`
	Backupwindow      *string   `awsName:"PreferredBackupWindow" awsType:"awsstr" templateName:"backupwindow"`
	Cluster           *string   `awsName:"DBClusterIdentifier" awsType:"awsstr" templateName:"cluster"`
	Dbname            *string   `awsName:"DBName" awsType:"awsstr" templateName:"dbname"`
	Dbsecuritygroups  []*string `awsName:"DBSecurityGroups" awsType:"awsstringslice" templateName:"dbsecuritygroups"`
	Domain            *string   `awsName:"Domain" awsType:"awsstr" templateName:"domain"`
	Encrypted         *bool     `awsName:"StorageEncrypted" awsType:"awsbool" templateName:"encrypted"`
	Iamrole           *string   `awsName:"DomainIAMRoleName" awsType:"awsstr" templateName:"iamrole"`
	License           *string   `awsName:"LicenseModel" awsType:"awsstr" templateName:"license"`
	Maintenancewindow *string   `awsName:"PreferredMaintenanceWindow" awsType:"awsstr" templateName:"maintenancewindow"`
	Multiaz           *bool     `awsName:"MultiAZ" awsType:"awsbool" templateName:"multiaz"`
	Parametergroup    *string   `awsName:"DBParameterGroupName" awsType:"awsstr" templateName:"parametergroup"`
	Timezone          *string   `awsName:"Timezone" awsType:"awsstr" templateName:"timezone"`
	Vpcsecuritygroups []*string `awsName:"VpcSecurityGroupIds" awsType:"awsstringslice" templateName:"vpcsecuritygroups"`
	Version           *string   `awsName:"EngineVersion" awsType:"awsstr" templateName:"version"`

	// Extra only for replica DB
	CopyTagsToSnapshot *string `awsName:"CopyTagsToSnapshot" awsType:"awsbool" templateName:"copytagstosnapshot"`
}

func (cmd *CreateDatabase) ValidateParams(params []string) ([]string, error) {
	return paramRule{
		tree: oneOf(allOf(node("type"), node("id"), node("engine"), node("password"), node("username"), node("size")), allOf(node("replica"), node("replica-source"))),
		extras: []string{"autoupgrade", "availabilityzone", "backupretention", "cluster", "dbname", "parametergroup",
			"dbsecuritygroups", "subnetgroup", "domain", "iamrole", "version", "iops", "license", "multiaz", "optiongroup",
			"port", "backupwindow", "maintenancewindow", "public", "encrypted", "storagetype", "timezone", "vpcsecuritygroups"},
	}.verify(params)
}

func (cmd *CreateDatabase) Validate_Password() (err error) {
	if pass := cmd.Password; pass != nil {
		if len(*pass) < 8 {
			err = errors.New("should at least be 8 characters")
		}
	}
	return
}

func (cmd *CreateDatabase) Validate_ReadReplicaIdentifier() error {
	if cmd.ReadReplicaIdentifier != nil {
		msg := "param not allowed in replica (either not applicable or directly inherited from the source DB)"
		if cmd.Backupretention != nil {
			return fmt.Errorf("'backupretention' %s", msg)
		}
		if cmd.Backupwindow != nil {
			return fmt.Errorf("'backupwindow' %s", msg)
		}
		if cmd.Cluster != nil {
			return fmt.Errorf("'cluster' %s", msg)
		}
		if cmd.Dbname != nil {
			return fmt.Errorf("'dbname' %s", msg)
		}
		if cmd.Dbsecuritygroups != nil {
			return fmt.Errorf("'dbsecuritygroups' %s", msg)
		}
		if cmd.Domain != nil {
			return fmt.Errorf("'domain' %s", msg)
		}
		if cmd.Encrypted != nil {
			return fmt.Errorf("'encrypted' %s", msg)
		}
		if cmd.Iamrole != nil {
			return fmt.Errorf("'iamrole' %s", msg)
		}
		if cmd.License != nil {
			return fmt.Errorf("'license' %s", msg)
		}
		if cmd.Maintenancewindow != nil {
			return fmt.Errorf("'maintenancewindow' %s", msg)
		}
		if cmd.Multiaz != nil {
			return fmt.Errorf("'multiaz' %s", msg)
		}
		if cmd.Parametergroup != nil {
			return fmt.Errorf("'parametergroup' %s", msg)
		}
		if cmd.Timezone != nil {
			return fmt.Errorf("'timezone' %s", msg)
		}
		if cmd.Vpcsecuritygroups != nil {
			return fmt.Errorf("'vpcsecuritygroups' %s", msg)
		}
		if cmd.Version != nil {
			return fmt.Errorf("'version' %s", msg)
		}
	}
	return nil
}

func (cmd *CreateDatabase) ManualRun(ctx map[string]interface{}) (output interface{}, err error) {
	if replica := cmd.ReadReplicaIdentifier; replica != nil {
		input := &rds.CreateDBInstanceReadReplicaInput{}
		if ierr := structInjector(cmd, input, ctx); ierr != nil {
			return nil, fmt.Errorf("cannot inject in rds.CreateDBInstanceReadReplicaInput: %s", ierr)
		}
		start := time.Now()
		output, err = cmd.api.CreateDBInstanceReadReplica(input)
		cmd.logger.ExtraVerbosef("rds.CreateDBInstanceReadReplica call took %s", time.Since(start))
	} else {
		input := &rds.CreateDBInstanceInput{}
		if ierr := structInjector(cmd, input, ctx); ierr != nil {
			return nil, fmt.Errorf("cannot inject in rds.CreateDBInstanceInput: %s", ierr)
		}
		start := time.Now()
		output, err = cmd.api.CreateDBInstance(input)
		cmd.logger.ExtraVerbosef("rds.CreateDBInstance call took %s", time.Since(start))
	}
	if err != nil {
		return output, err
	}
	return output, nil
}

func (cmd *CreateDatabase) ExtractResult(i interface{}) string {
	switch i.(type) {
	case *rds.CreateDBInstanceOutput:
		return awssdk.StringValue(i.(*rds.CreateDBInstanceOutput).DBInstance.DBInstanceIdentifier)
	case *rds.CreateDBInstanceReadReplicaOutput:
		return awssdk.StringValue(i.(*rds.CreateDBInstanceReadReplicaOutput).DBInstance.DBInstanceIdentifier)
	default:
		logger.Errorf("unexpected interface type %T", i)
		return ""
	}
}

type DeleteDatabase struct {
	_            string `action:"delete" entity:"database" awsAPI:"rds" awsCall:"DeleteDBInstance" awsInput:"rds.DeleteDBInstanceInput" awsOutput:"rds.DeleteDBInstanceOutput"`
	logger       *logger.Logger
	graph        cloudgraph.GraphAPI
	api          rdsiface.RDSAPI
	Id           *string `awsName:"DBInstanceIdentifier" awsType:"awsstr" templateName:"id" required:""`
	SkipSnapshot *bool   `awsName:"SkipFinalSnapshot" awsType:"awsbool" templateName:"skip-snapshot"`
	Snapshot     *string `awsName:"FinalDBSnapshotIdentifier" awsType:"awsstr" templateName:"snapshot"`
}

func (cmd *DeleteDatabase) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type CheckDatabase struct {
	_       string `action:"check" entity:"database" awsAPI:"rds"`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     rdsiface.RDSAPI
	Id      *string `templateName:"id" required:""`
	State   *string `templateName:"state" required:""`
	Timeout *int64  `templateName:"timeout" required:""`
}

func (cmd *CheckDatabase) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CheckDatabase) Validate_State() error {
	return NewEnumValidator(
		"available",
		"backing-up",
		"creating",
		"deleting",
		"failed",
		"maintenance",
		"modifying",
		"rebooting",
		"renaming",
		"resetting-master-credentials",
		"restore-error",
		"storage-full",
		"upgrading",
		notFoundState).Validate(cmd.State)
}

func (cmd *CheckDatabase) ManualRun(ctx map[string]interface{}) (interface{}, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: cmd.Id,
	}

	c := &checker{
		description: fmt.Sprintf("database %s", StringValue(cmd.Id)),
		timeout:     time.Duration(Int64AsIntValue(cmd.Timeout)) * time.Second,
		frequency:   5 * time.Second,
		fetchFunc: func() (string, error) {
			output, err := cmd.api.DescribeDBInstances(input)
			if err != nil {
				if awserr, ok := err.(awserr.Error); ok {
					if awserr.Code() == "DatabaseNotFound" {
						return notFoundState, nil
					}
				} else {
					return "", err
				}
			} else {
				if res := output.DBInstances; len(res) > 0 {
					for _, dbinst := range res {
						if StringValue(dbinst.DBInstanceIdentifier) == StringValue(cmd.Id) {
							return StringValue(dbinst.DBInstanceStatus), nil
						}
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

type StartDatabase struct {
	_      string `action:"start" entity:"database" awsAPI:"rds" awsCall:"StartDBInstance" awsInput:"rds.StartDBInstanceInput" awsOutput:"rds.StartDBInstanceOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    rdsiface.RDSAPI
	Id     *string `awsName:"DBInstanceIdentifier" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *StartDatabase) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

type StopDatabase struct {
	_      string `action:"stop" entity:"database" awsAPI:"rds" awsCall:"StopDBInstance" awsInput:"rds.StopDBInstanceInput" awsOutput:"rds.StopDBInstanceOutput"`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    rdsiface.RDSAPI
	Id     *string `awsName:"DBInstanceIdentifier" awsType:"awsstr" templateName:"id" required:""`
}

func (cmd *StopDatabase) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

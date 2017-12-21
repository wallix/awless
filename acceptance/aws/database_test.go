package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
)

func TestDatabase(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("db", func(t *testing.T) {
			Template("create database type=my-db-type id=my-db-id engine=my-db-engine password=my-db-password username=my-db-username size=12 "+
				"autoupgrade=true availabilityzone=my-db-availabilityzone backupretention=10 cluster=my-db-cluster "+
				"dbname=my-db-dbname parametergroup=my-db-parametergroup dbsecuritygroups=my-db-dbsecuritygroup-1,my-db-dbsecuritygroup-2 subnetgroup=my-db-subnetgroup "+
				"domain=my-db-domain iamrole=my-db-iamrole version=my-db-version iops=1024 license=my-db-license multiaz=true "+
				"optiongroup=my-db-optiongroup port=3306 backupwindow=my-db-backupwindow maintenancewindow=my-db-maintenancewindow "+
				"public=true encrypted=true storagetype=my-db-storagetype timezone=my-db-timezone vpcsecuritygroups=my-db-vpcsecuritygroup-1,my-db-vpcsecuritygroup-2").
				Mock(&rdsMock{
					CreateDBInstanceFunc: func(param0 *rds.CreateDBInstanceInput) (*rds.CreateDBInstanceOutput, error) {
						return &rds.CreateDBInstanceOutput{DBInstance: &rds.DBInstance{DBInstanceIdentifier: String("new-database-id")}}, nil
					},
				}).ExpectInput("CreateDBInstance", &rds.CreateDBInstanceInput{
				DBInstanceClass:         String("my-db-type"),
				DBInstanceIdentifier:    String("my-db-id"),
				Engine:                  String("my-db-engine"),
				MasterUserPassword:      String("my-db-password"),
				MasterUsername:          String("my-db-username"),
				AllocatedStorage:        Int64(12),
				AutoMinorVersionUpgrade: Bool(true),
				AvailabilityZone:        String("my-db-availabilityzone"),
				BackupRetentionPeriod:   Int64(10),
				DBClusterIdentifier:     String("my-db-cluster"),
				DBName:                  String("my-db-dbname"),
				DBParameterGroupName:    String("my-db-parametergroup"),
				DBSecurityGroups:        []*string{String("my-db-dbsecuritygroup-1"), String("my-db-dbsecuritygroup-2")},
				DBSubnetGroupName:       String("my-db-subnetgroup"),
				Domain:                  String("my-db-domain"),
				DomainIAMRoleName:       String("my-db-iamrole"),
				EngineVersion:           String("my-db-version"),
				Iops:                    Int64(1024),
				LicenseModel:            String("my-db-license"),
				MultiAZ:                 Bool(true),
				OptionGroupName:         String("my-db-optiongroup"),
				Port:                    Int64(3306),
				PreferredBackupWindow:      String("my-db-backupwindow"),
				PreferredMaintenanceWindow: String("my-db-maintenancewindow"),
				PubliclyAccessible:         Bool(true),
				StorageEncrypted:           Bool(true),
				StorageType:                String("my-db-storagetype"),
				Timezone:                   String("my-db-timezone"),
				VpcSecurityGroupIds:        []*string{String("my-db-vpcsecuritygroup-1"), String("my-db-vpcsecuritygroup-2")},
			}).ExpectCommandResult("new-database-id").ExpectCalls("CreateDBInstance").Run(t)
		})
		t.Run("read replica db", func(t *testing.T) {
			Template("create database replica=my-replica-id replica-source=my-source-id").
				Mock(&rdsMock{
					CreateDBInstanceReadReplicaFunc: func(param0 *rds.CreateDBInstanceReadReplicaInput) (*rds.CreateDBInstanceReadReplicaOutput, error) {
						return &rds.CreateDBInstanceReadReplicaOutput{DBInstance: &rds.DBInstance{DBInstanceIdentifier: String("new-replica-id")}}, nil
					},
				}).ExpectInput("CreateDBInstanceReadReplica", &rds.CreateDBInstanceReadReplicaInput{
				DBInstanceIdentifier:       String("my-replica-id"),
				SourceDBInstanceIdentifier: String("my-source-id"),
			}).ExpectCommandResult("new-replica-id").ExpectCalls("CreateDBInstanceReadReplica").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete database id=db-1234 skip-snapshot=false snapshot=my-snapshot-id").
			Mock(&rdsMock{
				DeleteDBInstanceFunc: func(param0 *rds.DeleteDBInstanceInput) (*rds.DeleteDBInstanceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteDBInstance", &rds.DeleteDBInstanceInput{
			DBInstanceIdentifier:      String("db-1234"),
			SkipFinalSnapshot:         Bool(false),
			FinalDBSnapshotIdentifier: String("my-snapshot-id"),
		}).ExpectCalls("DeleteDBInstance").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check database id=db-1234 state=creating timeout=1").
			Mock(&rdsMock{
				DescribeDBInstancesFunc: func(param0 *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
					return &rds.DescribeDBInstancesOutput{
						DBInstances: []*rds.DBInstance{
							{
								DBInstanceIdentifier: String("db-1234"),
								DBInstanceStatus:     String("creating"),
							},
						},
					}, nil
				},
			}).ExpectInput("DescribeDBInstances", &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: String("db-1234"),
		}).ExpectCalls("DescribeDBInstances").Run(t)
	})

	t.Run("start", func(t *testing.T) {
		Template("start database id=db-1234").
			Mock(&rdsMock{
				StartDBInstanceFunc: func(param0 *rds.StartDBInstanceInput) (*rds.StartDBInstanceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("StartDBInstance", &rds.StartDBInstanceInput{
			DBInstanceIdentifier: String("db-1234"),
		}).ExpectCalls("StartDBInstance").Run(t)
	})

	t.Run("stop", func(t *testing.T) {
		Template("stop database id=db-1234").
			Mock(&rdsMock{
				StopDBInstanceFunc: func(param0 *rds.StopDBInstanceInput) (*rds.StopDBInstanceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("StopDBInstance", &rds.StopDBInstanceInput{
			DBInstanceIdentifier: String("db-1234"),
		}).ExpectCalls("StopDBInstance").Run(t)
	})

	t.Run("restart", func(t *testing.T) {
		Template("restart database id=db-1234 with-failover=true").
			Mock(&rdsMock{
				RebootDBInstanceFunc: func(param0 *rds.RebootDBInstanceInput) (*rds.RebootDBInstanceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("RebootDBInstance", &rds.RebootDBInstanceInput{
			DBInstanceIdentifier: String("db-1234"),
			ForceFailover:        Bool(true),
		}).ExpectCalls("RebootDBInstance").Run(t)
	})
}

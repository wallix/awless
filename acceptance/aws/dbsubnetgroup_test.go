package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
)

func TestDbsubnetgroup(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create dbsubnetgroup name=my-dbsubnetgroup description=db-subnet-description subnets=sub-1234,sub-2345").
			Mock(&rdsMock{
				CreateDBSubnetGroupFunc: func(param0 *rds.CreateDBSubnetGroupInput) (*rds.CreateDBSubnetGroupOutput, error) {
					return &rds.CreateDBSubnetGroupOutput{
						DBSubnetGroup: &rds.DBSubnetGroup{DBSubnetGroupName: String("new-dbsubnetgroup-name")},
					}, nil
				},
			}).ExpectInput("CreateDBSubnetGroup", &rds.CreateDBSubnetGroupInput{
			DBSubnetGroupName:        String("my-dbsubnetgroup"),
			DBSubnetGroupDescription: String("db-subnet-description"),
			SubnetIds:                []*string{String("sub-1234"), String("sub-2345")},
		}).ExpectCommandResult("new-dbsubnetgroup-name").ExpectCalls("CreateDBSubnetGroup").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete dbsubnetgroup name=dbsubnetgroup-to-delete").
			Mock(&rdsMock{
				DeleteDBSubnetGroupFunc: func(param0 *rds.DeleteDBSubnetGroupInput) (*rds.DeleteDBSubnetGroupOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteDBSubnetGroup", &rds.DeleteDBSubnetGroupInput{DBSubnetGroupName: String("dbsubnetgroup-to-delete")}).
			ExpectCalls("DeleteDBSubnetGroup").Run(t)
	})
}

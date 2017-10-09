package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestGroup(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create group name=my-group").
			Mock(&iamMock{
				CreateGroupFunc: func(param0 *iam.CreateGroupInput) (*iam.CreateGroupOutput, error) {
					return &iam.CreateGroupOutput{Group: &iam.Group{GroupId: String("new-group-id")}}, nil
				},
			}).ExpectInput("CreateGroup", &iam.CreateGroupInput{
			GroupName: String("my-group"),
		}).ExpectCommandResult("new-group-id").ExpectCalls("CreateGroup").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete group name=my-group").Mock(&iamMock{
			DeleteGroupFunc: func(param0 *iam.DeleteGroupInput) (*iam.DeleteGroupOutput, error) { return nil, nil },
		}).ExpectInput("DeleteGroup", &iam.DeleteGroupInput{GroupName: String("my-group")}).
			ExpectCalls("DeleteGroup").Run(t)
	})
}

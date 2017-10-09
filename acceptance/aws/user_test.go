package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestUser(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create user name=donald").Mock(&iamMock{
			CreateUserFunc: func(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
				return &iam.CreateUserOutput{User: &iam.User{UserId: String("new-user-id")}}, nil
			}}).ExpectInput("CreateUser", &iam.CreateUserInput{
			UserName: String("donald"),
		}).ExpectCommandResult("new-user-id").ExpectCalls("CreateUser").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete user name=donald").Mock(&iamMock{
			DeleteUserFunc: func(input *iam.DeleteUserInput) (*iam.DeleteUserOutput, error) {
				return nil, nil
			}}).ExpectInput("DeleteUser", &iam.DeleteUserInput{
			UserName: String("donald"),
		}).ExpectCalls("DeleteUser").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach user name=donald group=trolls").Mock(&iamMock{
			AddUserToGroupFunc: func(input *iam.AddUserToGroupInput) (*iam.AddUserToGroupOutput, error) {
				return nil, nil
			}}).ExpectInput("AddUserToGroup", &iam.AddUserToGroupInput{
			UserName:  String("donald"),
			GroupName: String("trolls"),
		}).ExpectCalls("AddUserToGroup").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach user name=donald group=trolls").Mock(&iamMock{
			RemoveUserFromGroupFunc: func(input *iam.RemoveUserFromGroupInput) (*iam.RemoveUserFromGroupOutput, error) {
				return nil, nil
			}}).ExpectInput("RemoveUserFromGroup", &iam.RemoveUserFromGroupInput{
			UserName:  String("donald"),
			GroupName: String("trolls"),
		}).ExpectCalls("RemoveUserFromGroup").Run(t)
	})
}

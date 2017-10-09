package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestLoginProfile(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create loginprofile username=jdoe password=temporary-password password-reset=true").
			Mock(&iamMock{
				CreateLoginProfileFunc: func(*iam.CreateLoginProfileInput) (*iam.CreateLoginProfileOutput, error) {
					return &iam.CreateLoginProfileOutput{LoginProfile: &iam.LoginProfile{UserName: String("jdoe")}}, nil
				},
			}).ExpectInput("CreateLoginProfile", &iam.CreateLoginProfileInput{
			UserName:              String("jdoe"),
			Password:              String("temporary-password"),
			PasswordResetRequired: Bool(true),
		}).
			ExpectCommandResult("jdoe").ExpectCalls("CreateLoginProfile").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template("update loginprofile username=jdoe password=temporary-password password-reset=true").
			Mock(&iamMock{
				UpdateLoginProfileFunc: func(param0 *iam.UpdateLoginProfileInput) (*iam.UpdateLoginProfileOutput, error) {
					return nil, nil
				},
			}).ExpectInput("UpdateLoginProfile", &iam.UpdateLoginProfileInput{
			UserName:              String("jdoe"),
			Password:              String("temporary-password"),
			PasswordResetRequired: Bool(true),
		}).ExpectCalls("UpdateLoginProfile").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete loginprofile username=jdoe").
			Mock(&iamMock{
				DeleteLoginProfileFunc: func(param0 *iam.DeleteLoginProfileInput) (*iam.DeleteLoginProfileOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteLoginProfile", &iam.DeleteLoginProfileInput{
			UserName: String("jdoe"),
		}).ExpectCalls("DeleteLoginProfile").Run(t)
	})
}

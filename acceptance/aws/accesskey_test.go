package awsat

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestAccesskey(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		defer redirectStdErrToDevNull()()

		Template("create accesskey user=jdoe no-prompt=true").
			Mock(&iamMock{
				CreateAccessKeyFunc: func(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
					return &iam.CreateAccessKeyOutput{AccessKey: &iam.AccessKey{AccessKeyId: String("new-keypair-id")}}, nil
				},
			}).ExpectInput("CreateAccessKey", &iam.CreateAccessKeyInput{UserName: String("jdoe")}).
			ExpectCommandResult("new-keypair-id").ExpectCalls("CreateAccessKey").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete accesskey id=ACCESSKEYID user=jdoe").
			Mock(&iamMock{
				DeleteAccessKeyFunc: func(param0 *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteAccessKey", &iam.DeleteAccessKeyInput{
			UserName:    String("jdoe"),
			AccessKeyId: String("ACCESSKEYID"),
		}).ExpectCalls("DeleteAccessKey").Run(t)
	})
}

func redirectStdErrToDevNull() func() {
	originalStdErr := os.Stderr
	toDefer := func() {
		os.Stderr = originalStdErr
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	os.Stderr = devNull
	return toDefer
}

package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestMFADevice(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		defer redirectStdErrToDevNull()()
		Template("create mfadevice name=my-new-mfadevice").
			Mock(&iamMock{
				CreateVirtualMFADeviceFunc: func(param0 *iam.CreateVirtualMFADeviceInput) (*iam.CreateVirtualMFADeviceOutput, error) {
					return &iam.CreateVirtualMFADeviceOutput{VirtualMFADevice: &iam.VirtualMFADevice{SerialNumber: String("arn:new:mfadevice")}}, nil
				},
			}).ExpectInput("CreateVirtualMFADevice", &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: String("my-new-mfadevice"),
		}).
			ExpectCommandResult("arn:new:mfadevice").ExpectCalls("CreateVirtualMFADevice").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete mfadevice id=arn:mfadevice:to:delete").
			Mock(&iamMock{
				DeleteVirtualMFADeviceFunc: func(param0 *iam.DeleteVirtualMFADeviceInput) (*iam.DeleteVirtualMFADeviceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteVirtualMFADevice", &iam.DeleteVirtualMFADeviceInput{SerialNumber: String("arn:mfadevice:to:delete")}).
			ExpectCalls("DeleteVirtualMFADevice").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach mfadevice id=arn:mfadevice:to:attach user=my-username mfa-code-1=012345 mfa-code-2=123456 no-prompt=true").
			Mock(&iamMock{
				EnableMFADeviceFunc: func(param0 *iam.EnableMFADeviceInput) (*iam.EnableMFADeviceOutput, error) { return nil, nil },
			}).ExpectInput("EnableMFADevice", &iam.EnableMFADeviceInput{
			SerialNumber:        String("arn:mfadevice:to:attach"),
			UserName:            String("my-username"),
			AuthenticationCode1: String("012345"),
			AuthenticationCode2: String("123456"),
		}).
			ExpectCalls("EnableMFADevice").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach mfadevice id=arn:mfadevice:to:detach user=my-username").
			Mock(&iamMock{
				DeactivateMFADeviceFunc: func(param0 *iam.DeactivateMFADeviceInput) (*iam.DeactivateMFADeviceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeactivateMFADevice", &iam.DeactivateMFADeviceInput{
			SerialNumber: String("arn:mfadevice:to:detach"),
			UserName:     String("my-username"),
		}).
			ExpectCalls("DeactivateMFADevice").Run(t)
	})

}

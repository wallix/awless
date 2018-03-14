package awsat

import (
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func TestLaunchConfiguration(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Skip("need support of multiple API in mocks")
		_, userdataFile, cleanup := generateTmpFile("this is my content with {{ .AWLESS.oneRef }} content")
		defer cleanup()
		t.Run("with image", func(t *testing.T) {
			Template("oneRef=awesome\n"+
				"create launchconfiguration name=new-launchconfiguration type=t2.nano image=ami-1234 public=true "+
				"keypair=an-existing-kp userdata="+userdataFile+" securitygroups=sg-1234,sg-2345 role=my-role spotprice=12.5").
				Mock(&autoscalingMock{
					CreateLaunchConfigurationFunc: func(param0 *autoscaling.CreateLaunchConfigurationInput) (*autoscaling.CreateLaunchConfigurationOutput, error) {
						return &autoscaling.CreateLaunchConfigurationOutput{}, nil
					},
				}).ExpectInput("CreateLaunchConfiguration", &autoscaling.CreateLaunchConfigurationInput{
				LaunchConfigurationName:  String("new-launchconfiguration"),
				InstanceType:             String("t2.nano"),
				ImageId:                  String("ami-1234"),
				AssociatePublicIpAddress: Bool(true),
				KeyName:                  String("an-existing-kp"),
				UserData:                 String(base64.StdEncoding.EncodeToString([]byte("this is my content with awesome content"))),
				SecurityGroups:           []*string{String("sg-1234"), String("sg-2345")},
				IamInstanceProfile:       String("my-role"),
				SpotPrice:                String("12.5"),
			}).ExpectCommandResult("new-launchconfiguration").ExpectCalls("CreateLaunchConfiguration").Run(t)
		})
		t.Run("with distro", func(t *testing.T) {
			Template("oneRef=awesome\n"+
				"create launchconfiguration name=new-launchconfiguration type=t2.nano distro=debian public=true "+
				"keypair=an-existing-kp userdata="+userdataFile+" securitygroups=sg-1234,sg-2345 role=my-role spotprice=12.5").
				Mock(&autoscalingMock{
					CreateLaunchConfigurationFunc: func(param0 *autoscaling.CreateLaunchConfigurationInput) (*autoscaling.CreateLaunchConfigurationOutput, error) {
						return &autoscaling.CreateLaunchConfigurationOutput{}, nil
					},
				}).ExpectInput("CreateLaunchConfiguration", &autoscaling.CreateLaunchConfigurationInput{
				LaunchConfigurationName:  String("new-launchconfiguration"),
				InstanceType:             String("t2.nano"),
				ImageId:                  String("ami-1234"),
				AssociatePublicIpAddress: Bool(true),
				KeyName:                  String("an-existing-kp"),
				UserData:                 String(base64.StdEncoding.EncodeToString([]byte("this is my content with awesome content"))),
				SecurityGroups:           []*string{String("sg-1234"), String("sg-2345")},
				IamInstanceProfile:       String("my-role"),
				SpotPrice:                String("12.5"),
			}).ExpectCommandResult("new-launchconfiguration").ExpectCalls("CreateLaunchConfiguration").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete launchconfiguration name=my-launchconfig-to-delete").
			Mock(&autoscalingMock{
				DeleteLaunchConfigurationFunc: func(param0 *autoscaling.DeleteLaunchConfigurationInput) (*autoscaling.DeleteLaunchConfigurationOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteLaunchConfiguration", &autoscaling.DeleteLaunchConfigurationInput{LaunchConfigurationName: String("my-launchconfig-to-delete")}).
			ExpectCalls("DeleteLaunchConfiguration").Run(t)
	})
}

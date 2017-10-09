package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func TestScalingGroup(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create scalinggroup name=new-autoscaling launchconfiguration=config max-size=12 min-size=10 subnets=sub_1,sub_2 cooldown=3 desired-capacity=12 healthcheck-grace-period=4 healthcheck-type=healthy new-instances-protected=true targetgroups=tg_1,tg_2").Mock(&autoscalingMock{
			CreateAutoScalingGroupFunc: func(input *autoscaling.CreateAutoScalingGroupInput) (*autoscaling.CreateAutoScalingGroupOutput, error) {
				return &autoscaling.CreateAutoScalingGroupOutput{}, nil
			}}).
			ExpectInput("CreateAutoScalingGroup", &autoscaling.CreateAutoScalingGroupInput{
				AutoScalingGroupName:             String("new-autoscaling"),
				LaunchConfigurationName:          String("config"),
				MaxSize:                          Int64(12),
				MinSize:                          Int64(10),
				DefaultCooldown:                  Int64(3),
				DesiredCapacity:                  Int64(12),
				HealthCheckGracePeriod:           Int64(4),
				HealthCheckType:                  String("healthy"),
				NewInstancesProtectedFromScaleIn: Bool(true),
				VPCZoneIdentifier:                String("sub_1,sub_2"),
				TargetGroupARNs:                  []*string{String("tg_1"), String("tg_2")},
			}).ExpectCommandResult("new-autoscaling").ExpectCalls("CreateAutoScalingGroup").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template("update scalinggroup name=new-autoscaling launchconfiguration=config max-size=12 min-size=10 subnets=sub_1,sub_2 cooldown=3 desired-capacity=12 healthcheck-grace-period=4 healthcheck-type=healthy new-instances-protected=true").Mock(&autoscalingMock{
			UpdateAutoScalingGroupFunc: func(input *autoscaling.UpdateAutoScalingGroupInput) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
				return nil, nil
			}}).
			ExpectInput("UpdateAutoScalingGroup", &autoscaling.UpdateAutoScalingGroupInput{
				AutoScalingGroupName:             String("new-autoscaling"),
				LaunchConfigurationName:          String("config"),
				MaxSize:                          Int64(12),
				MinSize:                          Int64(10),
				DefaultCooldown:                  Int64(3),
				DesiredCapacity:                  Int64(12),
				HealthCheckGracePeriod:           Int64(4),
				HealthCheckType:                  String("healthy"),
				NewInstancesProtectedFromScaleIn: Bool(true),
				VPCZoneIdentifier:                String("sub_1,sub_2"),
			}).ExpectCalls("UpdateAutoScalingGroup").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete scalinggroup name=any-sg force=true").Mock(&autoscalingMock{
			DeleteAutoScalingGroupFunc: func(input *autoscaling.DeleteAutoScalingGroupInput) (*autoscaling.DeleteAutoScalingGroupOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteAutoScalingGroup", &autoscaling.DeleteAutoScalingGroupInput{
				AutoScalingGroupName: String("any-sg"),
				ForceDelete:          Bool(true),
			}).ExpectCalls("DeleteAutoScalingGroup").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check scalinggroup name=any-sg count=1 timeout=0").Mock(&autoscalingMock{
			DescribeAutoScalingGroupsFunc: func(input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
				return &autoscaling.DescribeAutoScalingGroupsOutput{
					AutoScalingGroups: []*autoscaling.Group{
						{AutoScalingGroupName: String("any-sg"), Instances: []*autoscaling.Instance{{InstanceId: String("one")}}},
					}}, nil
			}}).
			ExpectInput("DescribeAutoScalingGroups", &autoscaling.DescribeAutoScalingGroupsInput{
				AutoScalingGroupNames: []*string{String("any-sg")},
			}).ExpectCalls("DescribeAutoScalingGroups").Run(t)
	})
}

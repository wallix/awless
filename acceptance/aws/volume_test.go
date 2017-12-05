package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestVolume(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create volume availabilityzone=eu-west-1 size=1").Mock(&ec2Mock{
			CreateVolumeFunc: func(input *ec2.CreateVolumeInput) (*ec2.Volume, error) {
				return &ec2.Volume{VolumeId: String("new-volume-id")}, nil
			}}).
			ExpectInput("CreateVolume", &ec2.CreateVolumeInput{
				AvailabilityZone: String("eu-west-1"),
				Size:             Int64(1),
			}).ExpectCommandResult("new-volume-id").ExpectCalls("CreateVolume").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete volume id=any-volume-id").Mock(&ec2Mock{
			DeleteVolumeFunc: func(*ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteVolume", &ec2.DeleteVolumeInput{
				VolumeId: String("any-volume-id"),
			}).ExpectCalls("DeleteVolume").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check volume id=my-volume-id state=available timeout=2").Mock(&ec2Mock{
			DescribeVolumesFunc: func(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
				return &ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{
					{VolumeId: String("my-volume-id"), State: String("available")},
				}}, nil
			}}).ExpectInput("DescribeVolumes", &ec2.DescribeVolumesInput{
			VolumeIds: []*string{String("my-volume-id")},
		}).ExpectCalls("DescribeVolumes").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach volume id=my-volume-id device=dev instance=my-instance-id").Mock(&ec2Mock{
			AttachVolumeFunc: func(param0 *ec2.AttachVolumeInput) (*ec2.VolumeAttachment, error) {
				return &ec2.VolumeAttachment{VolumeId: String("my-volume-id")}, nil
			}}).ExpectInput("AttachVolume", &ec2.AttachVolumeInput{
			Device:     String("dev"),
			InstanceId: String("my-instance-id"),
			VolumeId:   String("my-volume-id"),
		}).ExpectCalls("AttachVolume").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach volume id=my-volume-id device=dev instance=my-instance-id force=true").Mock(&ec2Mock{
			DetachVolumeFunc: func(input *ec2.DetachVolumeInput) (*ec2.VolumeAttachment, error) {
				return &ec2.VolumeAttachment{VolumeId: String("my-volume-id")}, nil
			}}).ExpectInput("DetachVolume", &ec2.DetachVolumeInput{
			Device:     String("dev"),
			Force:      Bool(true),
			InstanceId: String("my-instance-id"),
			VolumeId:   String("my-volume-id"),
		}).ExpectCalls("DetachVolume").Run(t)
	})
}

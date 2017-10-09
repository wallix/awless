package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSnapshot(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create snapshot volume=my-volume-id description='this is the description of my snapshot'").
			Mock(&ec2Mock{
				CreateSnapshotFunc: func(param0 *ec2.CreateSnapshotInput) (*ec2.Snapshot, error) {
					return &ec2.Snapshot{SnapshotId: String("new-snapshot-id")}, nil
				},
			}).ExpectInput("CreateSnapshot", &ec2.CreateSnapshotInput{
			VolumeId:    String("my-volume-id"),
			Description: String("this is the description of my snapshot"),
		}).
			ExpectCommandResult("new-snapshot-id").ExpectCalls("CreateSnapshot").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete snapshot id=snap-1234").
			Mock(&ec2Mock{
				DeleteSnapshotFunc: func(param0 *ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteSnapshot", &ec2.DeleteSnapshotInput{SnapshotId: String("snap-1234")}).
			ExpectCalls("DeleteSnapshot").Run(t)
	})

	t.Run("copy", func(t *testing.T) {
		Template("copy snapshot source-id=my-origin-id source-region=my-origin-region encrypted=true description='an encrypted snapshot'").
			Mock(&ec2Mock{
				CopySnapshotFunc: func(param0 *ec2.CopySnapshotInput) (*ec2.CopySnapshotOutput, error) {
					return &ec2.CopySnapshotOutput{SnapshotId: String("my-snapshotcopy-id")}, nil
				},
			}).ExpectInput("CopySnapshot", &ec2.CopySnapshotInput{
			SourceSnapshotId: String("my-origin-id"),
			SourceRegion:     String("my-origin-region"),
			Encrypted:        Bool(true),
			Description:      String("an encrypted snapshot"),
		}).ExpectCommandResult("my-snapshotcopy-id").ExpectCalls("CopySnapshot").Run(t)
	})

}

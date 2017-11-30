package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestImage(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("with reboot", func(t *testing.T) {
			Template("create image name=my-image-name instance=my-instance-id reboot=true description='an new image'").
				Mock(&ec2Mock{
					CreateImageFunc: func(param0 *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
						return &ec2.CreateImageOutput{ImageId: String("new-image-id")}, nil
					},
				}).ExpectInput("CreateImage", &ec2.CreateImageInput{
				Name:        String("my-image-name"),
				InstanceId:  String("my-instance-id"),
				Description: String("an new image"),
			}).ExpectCommandResult("new-image-id").ExpectCalls("CreateImage").Run(t)
		})

		t.Run("with no reboot", func(t *testing.T) {
			Template("create image name=my-image-name instance=my-instance-id description='an new image'").
				Mock(&ec2Mock{
					CreateImageFunc: func(param0 *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
						return &ec2.CreateImageOutput{ImageId: String("new-image-id")}, nil
					},
				}).ExpectInput("CreateImage", &ec2.CreateImageInput{
				Name:        String("my-image-name"),
				InstanceId:  String("my-instance-id"),
				Description: String("an new image"),
				NoReboot:    Bool(true),
			}).ExpectCommandResult("new-image-id").ExpectCalls("CreateImage").Run(t)
		})
	})

	t.Run("copy", func(t *testing.T) {
		Template("copy image name=my-image-name source-id=my-origin-id source-region=my-origin-region encrypted=true description='an encrypted image'").
			Mock(&ec2Mock{
				CopyImageFunc: func(param0 *ec2.CopyImageInput) (*ec2.CopyImageOutput, error) {
					return &ec2.CopyImageOutput{ImageId: String("my-imagecopy-id")}, nil
				},
			}).ExpectInput("CopyImage", &ec2.CopyImageInput{
			Name:          String("my-image-name"),
			SourceImageId: String("my-origin-id"),
			SourceRegion:  String("my-origin-region"),
			Encrypted:     Bool(true),
			Description:   String("an encrypted image"),
		}).ExpectCommandResult("my-imagecopy-id").ExpectCalls("CopyImage").Run(t)
	})

	t.Run("import", func(t *testing.T) {
		t.Run("from ebs snapshot", func(t *testing.T) {
			Template("import image architecture=x86_64 description='my image desc' license=BYOL platform=Linux role=vmimport snapshot=my-ebs-snapshot").
				Mock(&ec2Mock{
					ImportImageFunc: func(param0 *ec2.ImportImageInput) (*ec2.ImportImageOutput, error) {
						return &ec2.ImportImageOutput{ImportTaskId: String("my-import-task-id")}, nil
					},
				}).ExpectInput("ImportImage", &ec2.ImportImageInput{
				Architecture: String("x86_64"),
				Description:  String("my image desc"),
				LicenseType:  String("BYOL"),
				Platform:     String("Linux"),
				RoleName:     String("vmimport"),
				DiskContainers: []*ec2.ImageDiskContainer{
					{SnapshotId: String("my-ebs-snapshot")},
				},
			}).ExpectCommandResult("my-import-task-id").ExpectCalls("ImportImage").Run(t)
		})
		t.Run("from url", func(t *testing.T) {
			Template("import image url=http://download.image.from.here").
				Mock(&ec2Mock{
					ImportImageFunc: func(param0 *ec2.ImportImageInput) (*ec2.ImportImageOutput, error) {
						return &ec2.ImportImageOutput{ImportTaskId: String("my-import-task-id")}, nil
					},
				}).ExpectInput("ImportImage", &ec2.ImportImageInput{
				DiskContainers: []*ec2.ImageDiskContainer{
					{Url: String("http://download.image.from.here")},
				},
			}).ExpectCommandResult("my-import-task-id").ExpectCalls("ImportImage").Run(t)
		})
		t.Run("from s3", func(t *testing.T) {
			Template("import image bucket=my-bucket s3object=my/s3/image/file.img").
				Mock(&ec2Mock{
					ImportImageFunc: func(param0 *ec2.ImportImageInput) (*ec2.ImportImageOutput, error) {
						return &ec2.ImportImageOutput{ImportTaskId: String("my-import-task-id")}, nil
					},
				}).ExpectInput("ImportImage", &ec2.ImportImageInput{
				DiskContainers: []*ec2.ImageDiskContainer{
					{UserBucket: &ec2.UserBucket{
						S3Bucket: String("my-bucket"),
						S3Key:    String("my/s3/image/file.img"),
					}},
				},
			}).ExpectCommandResult("my-import-task-id").ExpectCalls("ImportImage").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete image id=ami-to-delete delete-snapshots=true").
			Mock(&ec2Mock{
				DescribeImagesFunc: func(param0 *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
					return &ec2.DescribeImagesOutput{
						Images: []*ec2.Image{{BlockDeviceMappings: []*ec2.BlockDeviceMapping{
							{Ebs: &ec2.EbsBlockDevice{SnapshotId: String("snapshot-of-ami")}},
						}},
						}}, nil
				},
				DeregisterImageFunc: func(param0 *ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error) {
					return nil, nil
				},
				DeleteSnapshotFunc: func(param0 *ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DescribeImages", &ec2.DescribeImagesInput{
			ImageIds: []*string{String("ami-to-delete")},
		}).ExpectInput("DeregisterImage", &ec2.DeregisterImageInput{
			ImageId: String("ami-to-delete"),
		}).ExpectInput("DeleteSnapshot", &ec2.DeleteSnapshotInput{SnapshotId: String("snapshot-of-ami")}).
			ExpectCalls("DescribeImages", "DeregisterImage", "DeleteSnapshot").Run(t)
	})
}

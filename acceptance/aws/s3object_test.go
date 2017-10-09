package awsat

import (
	"testing"

	"os"

	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wallix/awless/aws/spec"
)

func TestS3object(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		f, filePath, cleanup := generateTmpFile("body content")
		defer cleanup()

		readSeeker, err := awsspec.NewProgressReader(f)
		if err != nil {
			t.Fatal(err)
		}
		awsspec.ProgressBarFactory = func(*os.File) (*awsspec.ProgressReadSeeker, error) {
			return readSeeker, nil
		}

		t.Run("with filename", func(t *testing.T) {
			Template("create s3object name=my-s3object file="+filePath+" bucket=any-bucket acl=public-read").Mock(&s3Mock{
				PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return &s3.PutObjectOutput{}, nil
				}}).
				ExpectInput("PutObject", &s3.PutObjectInput{
					ACL:    String("public-read"),
					Bucket: String("any-bucket"),
					Key:    String("my-s3object"),
					Body:   readSeeker,
				}).ExpectCommandResult("my-s3object").ExpectCalls("PutObject").Run(t)
		})

		t.Run("no filename", func(t *testing.T) {
			filename := filepath.Base(filePath)
			Template("create s3object file="+filePath+" bucket=any-bucket acl=public-read").Mock(&s3Mock{
				PutObjectFunc: func(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return &s3.PutObjectOutput{}, nil
				}}).
				ExpectInput("PutObject", &s3.PutObjectInput{
					ACL:    String("public-read"),
					Bucket: String("any-bucket"),
					Key:    String(filename),
					Body:   readSeeker,
				}).ExpectCommandResult(filename).ExpectCalls("PutObject").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		Template("update s3object name=any-file bucket=other-bucket acl=public-read version=2").Mock(&s3Mock{
			PutObjectAclFunc: func(input *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
				return &s3.PutObjectAclOutput{}, nil
			}}).
			ExpectInput("PutObjectAcl", &s3.PutObjectAclInput{
				ACL:       String("public-read"),
				Key:       String("any-file"),
				VersionId: String("2"),
				Bucket:    String("other-bucket"),
			}).ExpectCalls("PutObjectAcl").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete s3object name=any-file bucket=any-bucket").Mock(&s3Mock{
			DeleteObjectFunc: func(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			}}).
			ExpectInput("DeleteObject", &s3.DeleteObjectInput{
				Key:    String("any-file"),
				Bucket: String("any-bucket"),
			}).ExpectCalls("DeleteObject").Run(t)
	})
}

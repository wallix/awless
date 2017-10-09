package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestBucket(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create bucket name=my-new-bucket acl=public-read").
			Mock(&s3Mock{
				CreateBucketFunc: func(param0 *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
					return &s3.CreateBucketOutput{}, nil
				},
			}).ExpectInput("CreateBucket", &s3.CreateBucketInput{
			Bucket: String("my-new-bucket"),
			ACL:    String("public-read"),
		}).ExpectCommandResult("my-new-bucket").ExpectCalls("CreateBucket").Run(t)
	})

	t.Run("update", func(t *testing.T) {
		Template("update bucket name=my-bucket-to-update acl=public-read").
			Mock(&s3Mock{
				PutBucketAclFunc: func(param0 *s3.PutBucketAclInput) (*s3.PutBucketAclOutput, error) {
					return nil, nil
				},
			}).ExpectInput("PutBucketAcl", &s3.PutBucketAclInput{
			Bucket: String("my-bucket-to-update"),
			ACL:    String("public-read"),
		}).ExpectCalls("PutBucketAcl").Run(t)

		Template("update bucket name=my-bucket-to-update public-website=true redirect-hostname='http://myhostname.com' enforce-https=true").
			Mock(&s3Mock{
				PutBucketWebsiteFunc: func(param0 *s3.PutBucketWebsiteInput) (*s3.PutBucketWebsiteOutput, error) {
					return nil, nil
				},
			}).ExpectInput("PutBucketWebsite", &s3.PutBucketWebsiteInput{
			Bucket: String("my-bucket-to-update"),
			WebsiteConfiguration: &s3.WebsiteConfiguration{
				RedirectAllRequestsTo: &s3.RedirectAllRequestsTo{HostName: String("http://myhostname.com"), Protocol: String("https")},
			},
		}).ExpectCalls("PutBucketWebsite").Run(t)

		Template("update bucket name=my-bucket-to-update public-website=true index-suffix='index.go'").
			Mock(&s3Mock{
				PutBucketWebsiteFunc: func(param0 *s3.PutBucketWebsiteInput) (*s3.PutBucketWebsiteOutput, error) {
					return nil, nil
				},
			}).ExpectInput("PutBucketWebsite", &s3.PutBucketWebsiteInput{
			Bucket: String("my-bucket-to-update"),
			WebsiteConfiguration: &s3.WebsiteConfiguration{
				IndexDocument: &s3.IndexDocument{Suffix: String("index.go")},
			},
		}).ExpectCalls("PutBucketWebsite").Run(t)

		Template("update bucket name=my-bucket-to-update public-website=false").
			Mock(&s3Mock{
				DeleteBucketWebsiteFunc: func(param0 *s3.DeleteBucketWebsiteInput) (*s3.DeleteBucketWebsiteOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteBucketWebsite", &s3.DeleteBucketWebsiteInput{
			Bucket: String("my-bucket-to-update"),
		}).ExpectCalls("DeleteBucketWebsite").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete bucket name=my-bucket-to-delete").
			Mock(&s3Mock{
				DeleteBucketFunc: func(param0 *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteBucket", &s3.DeleteBucketInput{Bucket: String("my-bucket-to-delete")}).
			ExpectCalls("DeleteBucket").Run(t)
	})
}

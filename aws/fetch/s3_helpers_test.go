package awsfetch

import (
	"context"
	"fmt"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/graph"
)

func TestFetchFunctions(t *testing.T) {
	t.Parallel()
	t.Run("fetchAndExtractGrants", func(t *testing.T) {
		bucketsACL := map[string][]*s3.Grant{
			"bucket_1": {
				{Permission: awssdk.String("Read"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1"), Type: awssdk.String("my_type_1")}},
				{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2"), DisplayName: awssdk.String("my_user_2"), Type: awssdk.String("my_type_2")}},
				{Permission: awssdk.String("Execute"), Grantee: &s3.Grantee{ID: awssdk.String("usr_3"), DisplayName: awssdk.String("my_user_3"), EmailAddress: awssdk.String("user@domain"), Type: awssdk.String("my_type_3")}},
			},
			"bucket_2": {
				{Permission: awssdk.String("Read"), Grantee: &s3.Grantee{URI: awssdk.String("group_uri"), Type: awssdk.String("Group")}},
				{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1"), Type: awssdk.String("my_type_2")}},
			},
		}
		bucket1 := &s3.Bucket{Name: awssdk.String("bucket_1")}
		mock := &mockS3{grants: bucketsACL}
		res, err := fetchAndExtractGrantsFn(context.Background(), mock, awssdk.StringValue(bucket1.Name))
		if err != nil {
			t.Fatal(err)
		}
		expected := []*graph.Grant{
			{Permission: "Read", Grantee: graph.Grantee{GranteeID: "usr_1", GranteeType: "my_type_1"}},
			{Permission: "Write", Grantee: graph.Grantee{GranteeID: "usr_2", GranteeDisplayName: "my_user_2", GranteeType: "my_type_2"}},
			{Permission: "Execute", Grantee: graph.Grantee{GranteeID: "usr_3", GranteeDisplayName: "my_user_3<user@domain>", GranteeType: "my_type_3"}},
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}

		bucket2 := &s3.Bucket{Name: awssdk.String("bucket_2")}
		res, err = fetchAndExtractGrantsFn(context.Background(), mock, awssdk.StringValue(bucket2.Name))
		if err != nil {
			t.Fatal(err)
		}
		expected = []*graph.Grant{
			{Permission: "Read", Grantee: graph.Grantee{GranteeID: "group_uri", GranteeType: "Group"}},
			{Permission: "Write", Grantee: graph.Grantee{GranteeID: "usr_1", GranteeType: "my_type_2"}},
		}
		for i := range expected {
			if got, want := res[i].String(), expected[i].String(); got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		}
	})
}

type mockS3 struct {
	s3iface.S3API
	buckets map[string][]*s3.Bucket
	objects map[string][]*s3.Object
	grants  map[string][]*s3.Grant
}

func (m *mockS3) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	return &s3.GetBucketAclOutput{Grants: m.grants[awssdk.StringValue(input.Bucket)]}, nil
}

func (m *mockS3) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	var buckets []*s3.Bucket
	for _, b := range m.buckets {
		buckets = append(buckets, b...)
	}
	return &s3.ListBucketsOutput{Buckets: buckets}, nil
}
func (m *mockS3) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return &s3.ListObjectsOutput{Contents: m.objects[awssdk.StringValue(input.Bucket)]}, nil
}
func (m *mockS3) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	for region, buckets := range m.buckets {
		for _, bucket := range buckets {
			if awssdk.StringValue(bucket.Name) == awssdk.StringValue(input.Bucket) {
				return &s3.GetBucketLocationOutput{LocationConstraint: awssdk.String(region)}, nil
			}
		}
	}
	return nil, fmt.Errorf("bucket location mock: bucket %s not found", awssdk.StringValue(input.Bucket))
}

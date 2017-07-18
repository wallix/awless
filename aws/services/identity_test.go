package awsservices

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

type mockSTS struct {
	stsiface.STSAPI
	output *sts.GetCallerIdentityOutput
}

func (m *mockSTS) GetCallerIdentity(in *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return m.output, nil
}

func TestGetIdentityParseAllTypesOfUsername(t *testing.T) {
	tcases := []struct {
		arn, expResource, expResourceType, expResourcePath string
	}{
		{arn: "", expResource: "", expResourceType: "", expResourcePath: ""},
		{arn: "arn:", expResource: "", expResourceType: "", expResourcePath: ""},
		{arn: "arn:aws:iam::123456789012:root", expResource: "root", expResourceType: "user", expResourcePath: "root"},
		{arn: "arn:aws:iam::123456789012:user/Bob", expResource: "Bob", expResourceType: "user", expResourcePath: "user/Bob"},
		{arn: "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Donald", expResource: "division_abc/subdivision_xyz/Donald", expResourceType: "user", expResourcePath: "user/division_abc/subdivision_xyz/Donald"},
	}

	for _, tcase := range tcases {
		out := &sts.GetCallerIdentityOutput{Arn: awssdk.String(tcase.arn)}
		access := Access{STSAPI: &mockSTS{output: out}}
		id, err := access.GetIdentity()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := id.Resource, tcase.expResource; got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
		if got, want := id.ResourceType, tcase.expResourceType; got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
		if got, want := id.ResourcePath, tcase.expResourcePath; got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	}
}

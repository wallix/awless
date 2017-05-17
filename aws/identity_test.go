package aws

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
		arn, expResource, expResourceType string
	}{
		{arn: "arn:aws:iam::123456789012:root", expResource: "root", expResourceType: "user"},
		{arn: "arn:aws:iam::123456789012:user/Bob", expResource: "Bob", expResourceType: "user"},
		{arn: "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/Donald", expResource: "division_abc/subdivision_xyz/Donald", expResourceType: "user"},
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
	}
}

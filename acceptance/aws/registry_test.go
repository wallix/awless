package awsat

import (
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecr"
)

func TestRegistry(t *testing.T) {
	t.Run("authenticate", func(t *testing.T) {
		Template("authenticate registry accounts=my-registry-id-1,my-registry-id-2 no-docker-login=true").Mock(&ecrMock{
			GetAuthorizationTokenFunc: func(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
				return &ecr.GetAuthorizationTokenOutput{
					AuthorizationData: []*ecr.AuthorizationData{{AuthorizationToken: String(base64.StdEncoding.EncodeToString([]byte("user:my-authorization-token")))}},
				}, nil
			}}).
			ExpectInput("GetAuthorizationToken", &ecr.GetAuthorizationTokenInput{
				RegistryIds: []*string{String("my-registry-id-1"), String("my-registry-id-2")},
			}).ExpectCalls("GetAuthorizationToken").Run(t)
	})
}

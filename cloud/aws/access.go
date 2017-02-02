package aws

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Secu struct {
	*sts.STS
}

func NewSecu(sess *session.Session) *Secu {
	return &Secu{sts.New(sess)}
}

func (s *Secu) CallerIdentity() (interface{}, error) {
	return s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (s *Secu) GetUserId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Arn), nil
}

func (s *Secu) GetAccountId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Account), nil
}

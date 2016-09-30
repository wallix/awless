package awsprovider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Config is the awless config needs for aws
type Config struct {
	Credentials credentials.Value
}

// CreateSession creates a new session
func (c *Config) CreateSession() *session.Session {
	provider := credentials.StaticProvider{Value: c.Credentials}
	return session.New(&aws.Config{Credentials: credentials.NewCredentials(&provider)})
}

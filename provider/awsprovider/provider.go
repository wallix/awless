package awsprovider

import (
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/users"
)

// Provider for aws
type Provider struct {
	config *config.Config
	log    *log.Entry
	sess   *session.Session
}

// NewProvider creates a new aws provider
func NewProvider(config *config.Config, sess *session.Session) *Provider {
	return &Provider{
		config: config,
		log:    config.Log.WithField("provider", "aws"),
		sess:   sess,
	}
}

// ListUserNames part of provider.Interface
func (p *Provider) ListUserNames() ([]string, error) {
	p.log.Debugf("ListUserNames")
	c := iam.New(p.sess)
	l, err := c.ListUsers(nil)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(l.Users), len(l.Users))
	for i, u := range l.Users {
		names[i] = *u.UserName
	}
	return names, nil
}

// CreateUser part of provider.Interface
func (p *Provider) CreateUser(u users.User) error {
	c := iam.New(p.sess)
	out, err := c.CreateUser(&iam.CreateUserInput{UserName: &u.UserName})
	if err != nil {
		p.log.WithField("arn", *out.User.Arn).Debugf("user created")
	}
	return err
}

// DeleteUser part of provider.Interface
func (p *Provider) DeleteUser(u users.User) error {
	c := iam.New(p.sess)
	_, err := c.DeleteUser(&iam.DeleteUserInput{UserName: &u.UserName})
	return err
}

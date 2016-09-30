package awsprovider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/provider"
)

// ProviderBuilder implements provider.Builder interface
type ProviderBuilder struct {
}

// ReadConfig part of provider.Builder interface
func (p *ProviderBuilder) ReadConfig() (interface{}, error) {
	c := &Config{}
	fmt.Printf("AWS Access Key ID: ")
	_, err := fmt.Scanln(&c.Credentials.AccessKeyID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("AWS Secret Access Key: ")
	_, err = fmt.Scanln(&c.Credentials.SecretAccessKey)
	return c, err
}

// CheckConfig part of provider.Builder interface
func (p *ProviderBuilder) CheckConfig(i interface{}) error {
	c := i.(*Config)
	iamclient := iam.New(c.CreateSession())
	out, err := iamclient.GetUser(nil)
	fmt.Printf("AWS user::arn(%v)\n", *out.User.Arn)
	return err
}

// EmptyConfig part of provider.Builder interface
func (p *ProviderBuilder) EmptyConfig() interface{} {
	return &Config{}
}

// NewProvider part of provider.Builder interface
func (p *ProviderBuilder) NewProvider(config *config.Config, custom interface{}) provider.Interface {
	return NewProvider(config, custom.(*Config).CreateSession())
}

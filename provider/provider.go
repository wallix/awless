package provider

import (
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/users"
)

// Interface of a provider
type Interface interface {
	// ListUserNames list usernames
	ListUserNames() ([]string, error)

	// CreateUser creates a new user
	CreateUser(users.User) error

	// DeleteUser deletes an existing user
	DeleteUser(users.User) error
}

// LoadProvider load the provider name or the default if name is empty.
func LoadProvider(c *config.Config, name string) (Interface, error) {
	cp := &ConfigureParam{}
	var err error
	if name == "" {
		err = c.Load(prefix+".default", cp)
	} else {
		err = c.Load(prefix+name+"/config", cp)
	}
	if err != nil {
		return nil, err
	}
	builder, err := GetBuilder(cp.Kind)
	if err != nil {
		return nil, err
	}
	config := builder.EmptyConfig()
	err = c.Load(prefix+cp.Name+"/custom", config)
	if err != nil {
		return nil, err
	}
	return builder.NewProvider(c, config), nil
}

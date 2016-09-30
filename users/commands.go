package users

import "github.com/wallix/awless/config"

const (
	prefix = "users/"
)

// Commands contains all users commands
type Commands struct {
	config *config.Config
	api    *API
}

// NewCommands creates a new commands for the given config
func NewCommands(config *config.Config) *Commands {
	return &Commands{config: config, api: NewAPI(config)}
}

// Create a new awless user
func (c *Commands) Create(u User) {
	l := c.config.Log.WithField("user-name", u.UserName)
	l.Infof("adding user")
	if u.UserName == "" {
		c.config.Log.Errorf("name cannot be empty")
	}
	err := c.api.Create(u)
	if err == config.ErrExist {
		c.config.Log.Errorf("user already exists")
	}
}

// List all awless users
func (c *Commands) List(_ struct{}) {
	users := c.api.List()
	for _, user := range users {
		c.config.Log.Infof(user.UserName)
	}
}

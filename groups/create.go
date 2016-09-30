package groups

import "github.com/wallix/awless/config"

// CreateParam input parameters of Create
type CreateParam struct {
	GroupName string
}

// Create a users to the infrastructure
func (c *Commands) Create(p CreateParam) {
	l := c.config.Log.WithField("group-name", p.GroupName)
	l.Infof("adding group")
	if p.GroupName == "" {
		c.config.Log.Errorf("name cannot be empty")
	}
	group := Group{
		Name:  p.GroupName,
		Users: []string{},
	}
	err := c.config.Create(prefix+p.GroupName, group)
	if err == config.ErrExist {
		c.config.Log.Errorf("user already exists")
	}
}

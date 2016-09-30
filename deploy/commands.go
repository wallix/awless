package deploy

import (
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/provider"
)

// Commands contains all explain commands
type Commands struct {
	config *config.Config
}

// NewCommands creates a new commands for the given config
func NewCommands(config *config.Config) *Commands {
	return &Commands{config: config}
}

// ExplainParam params of explain
type ExplainParam struct {
	Provider string
	Noop     bool
	NoDelete bool
}

// Explain which operations is needed to apply on the provider
func (c *Commands) Explain(p ExplainParam) {
	pro, err := provider.LoadProvider(c.config, p.Provider)
	if err != nil {
		c.config.Log.WithError(err).Errorf("cannot load provider")
		return
	}
	api := NewAPI(c.config, pro)
	ops, err := api.UsersPatch(p.Noop)
	if err != nil {
		c.config.Log.WithError(err).Errorf("cannot get users patch")
		return
	}
	for _, op := range ops {
		if op.Kind == Del && p.NoDelete {
			continue
		}
		c.config.Log.Infof("%s user %s", op.Kind.Symbol(), op.User.UserName)
	}
}

// ApplyParam params of explain
type ApplyParam struct {
	Provider string
	NoDelete bool
}

// Apply the awless configuration on the provider
func (c *Commands) Apply(p ApplyParam) {
	pro, err := provider.LoadProvider(c.config, p.Provider)
	if err != nil {
		c.config.Log.WithError(err).Errorf("cannot load provider")
		return
	}
	api := NewAPI(c.config, pro)
	ops, err := api.UsersPatch(false)
	if err != nil {
		c.config.Log.WithError(err).Errorf("cannot get users patch")
		return
	}
	for _, op := range ops {
		switch op.Kind {
		case Noop:
			continue // should not occurs
		case Add:
			pro.CreateUser(op.User)
		case Del:
			if p.NoDelete {
				continue
			}
			pro.DeleteUser(op.User)
		}
		c.config.Log.Infof("%s user %s", op.Kind.Symbol(), op.User.UserName)
	}
}

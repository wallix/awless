package config

// Commands contains all explain commands
type Commands struct {
	config *Config
}

// NewCommands creates a new commands for the given config
func NewCommands(config *Config) *Commands {
	return &Commands{config: config}
}

// CommitParam params of commit
type CommitParam struct {
	message string
}

// Commit configuration change
func (c *Commands) Commit(p CommitParam) {
	c.config.Commit(p.message, true)
}

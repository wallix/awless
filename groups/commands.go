package groups

import "github.com/wallix/awless/config"

const (
	prefix = "groups/"
)

// Commands contains all users commands
type Commands struct {
	config *config.Config
}

// NewCommands creates a new commands for the given config
func NewCommands(config *config.Config) *Commands {
	return &Commands{config: config}
}

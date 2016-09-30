package provider

import "github.com/wallix/awless/config"

// Commands contains all provider commands
type Commands struct {
	config *config.Config
}

// NewCommands creates a new commands for the given config
func NewCommands(config *config.Config) *Commands {
	return &Commands{config: config}
}

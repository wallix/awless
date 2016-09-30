package provider

import (
	"fmt"

	"github.com/wallix/awless/config"
)

var (
	cprov = make(map[string]Builder)
)

// Builder is a tools for reading and check a custom provider
// configuration
type Builder interface {
	// ReadConfig reads config from the environnement
	ReadConfig() (interface{}, error)

	// CheckConfig checks if the config is valid
	CheckConfig(interface{}) error

	// EmptyConfig returns an empty config
	EmptyConfig() interface{}

	// NewProvider returns a new provider Interface
	NewProvider(c *config.Config, config interface{}) Interface
}

// RegisterBuilder register a new kind of builder
func RegisterBuilder(kind string, builder Builder) {
	cprov[kind] = builder
}

// GetBuilder returns a builder for the given provider kind, or an error if
// kind is unknown
func GetBuilder(kind string) (Builder, error) {
	b, o := cprov[kind]
	if o {
		return b, nil
	}
	return nil, fmt.Errorf("provider kind[%s] is unknown", kind)
}

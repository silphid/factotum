package cfg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/silphid/factotum/cli/src/internal/ctx"
	"github.com/silphid/factotum/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = "2021.04"
	SharedFileName = "shared.yaml"
	UserFileName   = "user.yaml"
	ContextBase    = "base"
	ContextNone    = "none"
)

// Config represents factotum's current config persisted to disk
type Config struct {
	Version  string
	Base     ctx.Context
	Contexts map[string]ctx.Context
}

// Load loads the config file from given directory
func Load(dir, fileName string) (*Config, error) {
	var config Config
	path := filepath.Join(dir, fileName)
	if !helpers.PathExists(path) {
		return &config, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading config file: %w", err)
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling config file yaml: %w", err)
	}

	if config.Version != currentVersion {
		return nil, fmt.Errorf("unsupported config file %q version %s (expected %s)", path, config.Version, currentVersion)
	}

	return &config, nil
}

// Clone returns a deep-copy of this config
func (c Config) Clone() Config {
	var config Config
	config.Version = c.Version
	config.Base = c.Base.Clone()
	config.Contexts = make(map[string]ctx.Context, len(c.Contexts))
	for key, value := range c.Contexts {
		config.Contexts[key] = value.Clone()
	}
	return config
}

// Merge creates a deep-copy of this config and copies values from given source config on top of it
func (c Config) Merge(source Config) Config {
	config := c.Clone()
	if source.Version != "" {
		config.Version = source.Version
	}

	// Base context
	config.Base = c.Base.Merge(source.Base)

	// Named contexts
	for key, value := range source.Contexts {
		targetContext := config.Contexts[key]
		config.Contexts[key] = targetContext.Merge(value)
	}

	return config
}

// GetContext returns context with given name merged on top of base context.
// If name is "base", only the base context is returned.
// If name is "none", an empty context is returned.
func (c Config) GetContext(name string) (ctx.Context, error) {
	// No context
	var context ctx.Context
	if name == ContextNone {
		return context, nil
	}

	// Base context
	context = context.Merge(c.Base)
	if name == ContextBase {
		return context, nil
	}

	// Named context
	namedContext, ok := c.Contexts[name]
	if !ok {
		return ctx.Context{}, fmt.Errorf("context %q not found", name)
	}
	return context.Merge(namedContext), nil
}

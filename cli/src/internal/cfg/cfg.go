package cfg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/silphid/factotum/cli/src/internal/ctx"
	"github.com/silphid/factotum/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = "2021.04"
	sharedFileName = "shared.yaml"
	userFileName   = "user.yaml"
	ContextBase    = "base"
	ContextNone    = "none"
)

// Config represents factotum's current config persisted to disk
type Config struct {
	Version  string
	Base     ctx.Context
	Contexts map[string]ctx.Context
}

// Load loads the config file from given file
func Load(file string) (Config, error) {
	var config Config
	if !helpers.PathExists(file) {
		return config, nil
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("loading config file: %w", err)
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return config, fmt.Errorf("unmarshalling yaml of config file %q: %w", file, err)
	}

	if config.Version != currentVersion {
		return config, fmt.Errorf("unsupported version %s (expected %s) in config file %q", config.Version, currentVersion, file)
	}

	return config, nil
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

// LoadConfigs load both the shared and user config files
func LoadConfigs(sharedDir, userDir string) (Config, Config, error) {
	sharedConfig, err := Load(filepath.Join(sharedDir, sharedFileName))
	if err != nil {
		return Config{}, Config{}, fmt.Errorf("loading shared config: %w", err)
	}
	userConfig, err := Load(filepath.Join(userDir, userFileName))
	if err != nil {
		return Config{}, Config{}, fmt.Errorf("loading user config: %w", err)
	}
	return sharedConfig, userConfig, nil
}

// GetContextNames returns the list of all context names user can choose from including
// the special "base" and "none" contexts.
func GetContextNames(sharedDir, userDir string) ([]string, error) {
	sharedConfig, userConfig, err := LoadConfigs(sharedDir, userDir)
	if err != nil {
		return nil, err
	}
	return getContextNames(sharedConfig, userConfig)
}

// getContextNames returns the list of all context names user can choose from including
// the special "base" and "none" contexts.
func getContextNames(sharedConfig, userConfig Config) ([]string, error) {
	// Extract unique names
	namesMap := make(map[string]bool)
	for name := range sharedConfig.Contexts {
		namesMap[name] = true
	}
	for name := range userConfig.Contexts {
		namesMap[name] = true
	}

	// Sort
	sortedNames := make([]string, 0, len(namesMap))
	for name := range namesMap {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	// Prepend special contexts
	names := make([]string, 0, len(sortedNames)+2)
	names = append(names, "none", "base")
	names = append(names, sortedNames...)

	return names, nil
}

// GetContext finds shared/user base/named contexts and returns their merged result.
// If name is "base", only the merged base context is returned.
// If name is "none", an empty context is returned.
func GetContext(sharedDir, userDir, name string) (ctx.Context, error) {
	// No context
	if name == ContextNone {
		return ctx.Context{}, nil
	}

	sharedConfig, userConfig, err := LoadConfigs(sharedDir, userDir)
	if err != nil {
		return ctx.Context{}, err
	}
	return getContext(sharedConfig, userConfig, name)
}

// getContext finds shared/user base/named contexts and returns their merged result.
// If name is "base", only the merged base context is returned.
func getContext(sharedConfig, userConfig Config, name string) (ctx.Context, error) {
	// Base contexts
	baseContext := sharedConfig.Base.Merge(userConfig.Base)
	if name == ContextBase {
		return baseContext, nil
	}

	// Named contexts
	sharedNamedContext, sharedOK := sharedConfig.Contexts[name]
	userNamedContext, userOK := userConfig.Contexts[name]
	if !sharedOK && !userOK {
		return ctx.Context{}, fmt.Errorf("context not found %q", name)
	}

	// Merge contexts
	return baseContext.Merge(sharedNamedContext).Merge(userNamedContext), nil
}

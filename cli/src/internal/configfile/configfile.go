package configfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

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

// RegistryType represents the type of docker registry factotum image should be retrieved from
type RegistryType string

const (
	RegistryGCR       RegistryType = "gcr"
	RegistryECR       RegistryType = "ecr"
	RegistryDockerHub RegistryType = "dockerhub"
)

// Container represents the information on how to retrieve factotum docker image
type Container struct {
	Registry RegistryType
	Image    string
}

// Context represents an execution context for factotum (env vars and volumes)
type Context struct {
	Env     map[string]string
	Volumes map[string]string
}

// Config represents factotum's current config persisted to disk
type Config struct {
	Version   string
	Container Container
	Base      Context
	Contexts  map[string]Context
}

// Save saves config file to given directory
func (s Config) Save(dir, fileName string) error {
	s.Version = currentVersion

	doc, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, fileName)
	return ioutil.WriteFile(path, doc, 0644)
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

func copyContext(source, target *Context) {
	if target.Env == nil {
		target.Env = make(map[string]string)
	}
	for key, value := range source.Env {
		target.Env[key] = value
	}

	if target.Volumes == nil {
		target.Volumes = make(map[string]string)
	}
	for key, value := range source.Volumes {
		target.Volumes[key] = value
	}
}

func (c Config) GetContext(name string) (Context, error) {
	// No context
	if name == ContextNone {
		return Context{}, nil
	}

	// Base context
	var context Context
	copyContext(&c.Base, &context)
	if name == ContextBase {
		return context, nil
	}

	// Named context
	namedContext, ok := c.Contexts[name]
	if !ok {
		return Context{}, fmt.Errorf("context %q not found", name)
	}
	copyContext(&namedContext, &context)

	return context, nil
}

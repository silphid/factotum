package statefile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/silphid/factotum/cli/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

const (
	currentVersion = "2021.04"
	fileName       = "state.yaml"
)

// State represents factotum's current state persisted to disk
type State struct {
	Version        string `yaml:"version"`
	CloneDir       string `yaml:"cloneDir"`
	ImageVersion   string `yaml:"imageVersion"`
	CurrentContext string `yaml:"currentContext"`
}

// Save saves state file to given directory
func (s State) Save(dir string) error {
	s.Version = currentVersion

	doc, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, fileName)
	return ioutil.WriteFile(path, doc, 0644)
}

// Load loads the state file from given directory
func Load(dir string) (*State, error) {
	var state State
	path := filepath.Join(dir, fileName)
	if !helpers.PathExists(path) {
		return &state, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading state file: %w", err)
	}
	err = yaml.Unmarshal(buf, &state)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling state file yaml: %w", err)
	}

	if state.Version != currentVersion {
		return nil, fmt.Errorf("unsupported state file %q version %s (expected %s)", path, state.Version, currentVersion)
	}

	return &state, nil
}

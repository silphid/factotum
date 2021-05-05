package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/silphid/factotum/cli/src/internal/cfg"
	"github.com/silphid/factotum/cli/src/internal/ctx"
	"github.com/silphid/factotum/cli/src/internal/logging"
	"github.com/silphid/factotum/cli/src/internal/statefile"
)

const (
	homeVar       = "FACTOTUM_HOME"
	homeDirName   = ".factotum"
	configDirName = "config"
)

type Core struct {
	homeDir   string
	sharedDir string
}

func New() (Core, error) {
	homeDir, err := getFactotumHomeDir()
	if err != nil {
		return Core{}, err
	}
	cloneDir, err := getFactotumCloneDir(homeDir)
	if err != nil {
		return Core{}, err
	}
	sharedDir := filepath.Join(cloneDir, configDirName)
	return Core{
		homeDir:   homeDir,
		sharedDir: sharedDir,
	}, nil
}

// GetContextNames returns the list of all context names user can choose from including
// the special "base" and "none" contexts.
func (c Core) GetContextNames() ([]string, error) {
	return cfg.GetContextNames(c.sharedDir, c.homeDir)
}

// GetContext finds shared/user base/named contexts and returns their merged result.
// If name is "base", only the merged base context is returned.
// If name is "none", an empty context is returned.
func (c Core) GetContext(name string) (ctx.Context, error) {
	return cfg.GetContext(c.sharedDir, c.homeDir, name)
}

// UseContext validates that given context is valid and saves it as default context
// to state file. If an empty string is passed, prompts user to select from list of
// available contexts.
func (c Core) UseContext(name string) error {
	if name == "" {
		var err error
		name, err = c.promptContext()
		if err != nil {
			return err
		}
	} else {
		err := c.validateContextName(name)
		if err != nil {
			return err
		}
	}
	state, err := statefile.Load(c.homeDir)
	if err != nil {
		return err
	}
	state.CurrentContext = name
	return state.Save()
}

func (c Core) promptContext() (string, error) {
	// Get context names
	names, err := c.GetContextNames()
	if err != nil {
		return "", err
	}

	// Determine default context index
	defaultIndex := 0
	state, err := statefile.Load(c.homeDir)
	if err != nil {
		return "", err
	}
	for i, n := range names {
		if n == state.CurrentContext {
			defaultIndex = i
			break
		}
	}

	// Show selection prompt
	prompt := &survey.Select{
		Message: "Select context:",
		Options: names,
		Default: defaultIndex,
	}
	selectedIndex := defaultIndex
	if err := survey.AskOne(prompt, &selectedIndex); err != nil {
		return "", err
	}

	return names[selectedIndex], nil
}

// validateContextName returns an error if given context names is invalid
func (c Core) validateContextName(name string) error {
	names, err := c.GetContextNames()
	if err != nil {
		return err
	}
	found := false
	for _, n := range names {
		if n == name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("context %q invalid (expecting one of: %s)", name, strings.Join(names, ", "))
	}
	return nil
}

// GetCloneDir returns currently configured clone directory
func (c Core) GetState() (statefile.State, error) {
	return statefile.Load(c.homeDir)
}

// getFactotumCloneDir returns the path to the factotum clone directory, specified by required FACTOTUM_CLONE env var
func getFactotumCloneDir(homeDir string) (string, error) {
	state, err := statefile.Load(homeDir)
	if err != nil {
		return "", err
	}
	if state.CloneDir == "" {
		return "", fmt.Errorf("clone dir not defined in state.yaml, please reinstall")
	}
	return state.CloneDir, nil
}

// getFactotumHomeDir returns the path to the factotum home directory, optionally specified by FACTOTUM_HOME env var
// (defaults to ~/.factotum)
func getFactotumHomeDir() (homeDir string, err error) {
	defer func() {
		if err == nil {
			logging.Log("Using factotum home dir: %s", homeDir)
		}
	}()

	homeDir, ok := os.LookupEnv(homeVar)
	if ok && homeDir != "" {
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		err = fmt.Errorf("failed to detect user home directory: %w", err)
		return
	}
	homeDir = filepath.Join(home, homeDirName)
	return
}

package home

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/silphid/factotum/cli/src/internal/logging"
)

const (
	homeVar     = "FACTOTUM_HOME"
	homeDirName = ".factotum"
)

// GetFactotumHomeDir returns the path to the factotum home directory, optionally specified by FACTOTUM_HOME env var
// (defaults to ~/.factotum)
func GetFactotumHomeDir() (homeDir string, err error) {
	defer func() {
		if err == nil {
			logging.Log("Using clone dir: %s", homeDir)
		}
	}()

	homeDir, ok := os.LookupEnv(homeVar)
	if ok && homeDir != "" {
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		err = fmt.Errorf("failed to detect home directory: %w", err)
		return
	}
	homeDir = filepath.Join(home, homeDirName)
	return
}

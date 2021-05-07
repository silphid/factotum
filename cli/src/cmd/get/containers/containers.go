package containers

import (
	"github.com/silphid/factotum/cli/src/internal/docker"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "containers",
		Short: "Lists running containers",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	return docker.ListContainers()
}
package cmd

import (
	"github.com/silphid/factotum/cli/src/internal/logging"
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot(version string) *cobra.Command {
	c := &cobra.Command{
		Use:          "factotum",
		Short:        "A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher",
		Long:         `A DevOps & CI/CD & Kubernetes-oriented general purpose Docker container with CLI launcher`,
		SilenceUsage: true,
	}

	// var options internal.Options
	c.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "display verbose messages")
	// c.AddCommand(require.New(&options))
	return c
}

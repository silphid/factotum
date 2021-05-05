package cmd

import (
	"github.com/silphid/factotum/cli/src/cmd/get"
	"github.com/silphid/factotum/cli/src/cmd/set"
	"github.com/silphid/factotum/cli/src/cmd/use"
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
	c.AddCommand(get.New())
	c.AddCommand(set.New())
	c.AddCommand(use.New())
	return c
}

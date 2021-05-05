package set

import (
	"github.com/silphid/factotum/cli/src/cmd/set/clone"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "set",
		Short: "Sets a factotum state variable to given value",
	}
	c.AddCommand(clone.New())
	return c
}

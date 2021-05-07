package cli

import (
	"fmt"

	"github.com/silphid/factotum/cli/src/internal/ctx"
)

type CLI struct{}

func (c CLI) Start(ct ctx.Context, imageTag, containerName string) error {
	return fmt.Errorf("not implemented")
}

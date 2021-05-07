package docker

import "github.com/silphid/factotum/cli/src/internal/ctx"

type Docker interface {
	Start(c ctx.Context, imageTag, containerName string) error
}

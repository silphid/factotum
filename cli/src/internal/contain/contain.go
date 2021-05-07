package contain

import (
	"fmt"

	"github.com/silphid/factotum/cli/src/internal/ctx"
	"github.com/silphid/factotum/cli/src/internal/docker/api"
)

const (
	factotumContainerPrefix = "factotum"
)

func Start(c ctx.Context, imageTag string) error {
	if c.Container == "" {
		c.Container = "factotum"
	}

	containerName := fmt.Sprintf("%s-%s-%s-%s", factotumContainerPrefix, c.Container, c.Name, imageTag)

	docker := api.API{}
	// docker := cli.CLI{}
	return docker.Start(c, imageTag, containerName)
}

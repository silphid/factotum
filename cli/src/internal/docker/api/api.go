package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/silphid/factotum/cli/src/internal/ctx"
)

type API struct{}

func (a API) Start(ct ctx.Context, imageTag, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	container, err := getContainer(cli, containerName)
	if err != nil {
		return err
	}

	if container == nil {
		fmt.Printf("Creating container\n")
		container, err = createContainer(cli, ct, containerName)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Status: %s\n", container.Status)
	fmt.Printf("State: %s\n", container.State)

	err = cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Status: %s\n", container.Status)
	fmt.Printf("State: %s\n", container.State)

	return nil
}

func getContainer(cli *client.Client, name string) (*types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", name)),
	})
	if err != nil {
		return nil, err
	}

	if len(containers) > 0 {
		return &containers[0], nil
	}
	return nil, nil
}

func createContainer(cli *client.Client, c ctx.Context, name string) (*types.Container, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to detect user home directory: %w", err)
	}

	mounts := make([]mount.Mount, 0)
	for source, target := range c.Mounts {
		source = strings.ReplaceAll(source, "$HOME", home)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		})
	}

	config := container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Entrypoint:   []string{"zsh"},
	}
	hostConfig := container.HostConfig{
		Mounts: mounts,
	}
	networkingConfig := network.NetworkingConfig{}

	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networkingConfig, nil, name)
	if err != nil {
		return nil, err
	}

	return getContainer(cli, name)
}

// func ListContainers() error {
// 	cli, err := client.NewClientWithOpts()
// 	if err != nil {
// 		return err
// 	}

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
// 		All: true,
// 		Filters: filters.NewArgs(
// 			filters.Arg("name", "factotum-*")),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for _, container := range containers {
// 		fmt.Println(strings.TrimPrefix(container.Names[0], "/"))
// 	}
// 	return nil
// }

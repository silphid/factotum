package docker

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
	"github.com/silphid/factotum/cli/src/internal/ctx"
)

const (
	factotumContainerPrefix = "factotum"
)

func Start(c ctx.Context, imageTag string) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	containerName := fmt.Sprintf("%s-%s-%s-%s", factotumContainerPrefix, c.Container, c.Name, imageTag)

	// Is container created?
	isCreated, err := isContainerCreated(cli, containerName)
	if err != nil {
		return err
	}
	if isCreated {
		isRunning, err := isContainerRunning(cli, containerName)
		if err != nil {
			return err
		}
		if isRunning {
			// Running
			fmt.Printf("Running")
		} else {
			// Created
			fmt.Printf("Created")
		}
	} else {
		// Absent
		fmt.Printf("Absent")
		err := createContainer(cli, c, containerName)
		if err != nil {
			return err
		}
	}

	return nil
}

func isContainerRunning(cli *client.Client, name string) (bool, error) {
	return isContainerListed(cli, name, false)
}

func isContainerCreated(cli *client.Client, name string) (bool, error) {
	return isContainerListed(cli, name, true)
}

func isContainerListed(cli *client.Client, name string, all bool) (bool, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: all,
		Filters: filters.NewArgs(
			filters.Arg("name", name)),
	})
	if err != nil {
		return false, err
	}
	return len(containers) > 0, nil
}

func createContainer(cli *client.Client, c ctx.Context, name string) error {
	mounts := make([]mount.Mount, 0)
	for k, v := range c.Mounts {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: k,
			Target: v,
		})
	}

	config := container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}
	hostConfig := container.HostConfig{
		Mounts: mounts,
	}
	networkingConfig := network.NetworkingConfig{}

	body, err := cli.ContainerCreate(context.Background(), &config, &hostConfig, &networkingConfig, nil, name)
	if err != nil {
		return err
	}

	fmt.Printf("Created container ID: %s\n", body.ID)
	return nil
}

func ListContainers() error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", "factotum-*")),
	})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println(strings.TrimPrefix(container.Names[0], "/"))
	}
	return nil
}

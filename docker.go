package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type Docker struct {
	client *client.Client
}

func NewDocker() (res Docker, err error) {

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.24"}
	client, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		return
	}

	res = Docker{
		client: client,
	}
	return

}

func (d Docker) DeleteContainer(ctx context.Context, name string) error {

	filter := filters.NewArgs()
	filter.Add("name", name)

	res, err := d.client.ContainerList(ctx, types.ContainerListOptions{
		All:    true,
		Quiet:  true,
		Filter: filter,
	})
	if err != nil {
		return err
	}

	for _, container := range res {

		err := d.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		})
		if err != nil {
			return err
		}

	}

	return nil

}

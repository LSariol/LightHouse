package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Handler struct {
	cli *client.Client
}

func NewHandler() (*Handler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Handler{cli: cli}, nil
}

func (h *Handler) ListContainers(ctx context.Context) ([]types.Container, error) {
	f := filters.NewArgs()

	return h.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: f,
	})
}

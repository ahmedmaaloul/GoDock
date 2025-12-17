package dockerwrapper

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Container represents a simplified view of a Docker container for the UI.
type Container struct {
	ID     string
	Image  string
	Status string // "Up 2 hours", "Exited (0) 5 seconds ago"
	State  string // "running", "exited"
}

// Client wraps the Docker SDK client.
type Client struct {
	cli *client.Client
}

// NewClient initializes a new Docker client.
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

// ListContainers returns a list of all containers (running and stopped).
func (c *Client) ListContainers() ([]Container, error) {
	containers, err := c.cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []Container
	for _, cnt := range containers {
		// ID is usually long, let's take the first 12 chars
		shortID := cnt.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		result = append(result, Container{
			ID:     shortID,
			Image:  cnt.Image,
			Status: cnt.Status,
			State:  cnt.State,
		})
	}
	return result, nil
}

// StartContainer starts a container by ID.
func (c *Client) StartContainer(containerID string) error {
	return c.cli.ContainerStart(context.Background(), containerID, container.StartOptions{})
}

// StopContainer stops a container by ID.
func (c *Client) StopContainer(containerID string) error {
	return c.cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
}

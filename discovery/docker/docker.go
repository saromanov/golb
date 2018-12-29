package docker

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/saromanov/golb/server"
)

const defaultEndpoint = "unix:///var/run/docker.sock"

// Discovery provides discoverying of the servers
// via docker
type Discovery struct {
	client  *docker.Client
	servers []*server.Server
}

// New provides initialization of docker and Discovery
func New() (Discovery, error) {
	client, err := docker.NewClient(defaultEndpoint)
	if err != nil {
		return Discovery{}, err
	}

	return Discovery{
		client: client,
	}, nil
}

// Search provides getting of containers and add their host
// on backend representation of server
func (d *Discovery) Search() error {
	containers, err := d.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return err
	}

	for _, c := range containers {
		fmt.Println(c)
	}
	return nil
}

// GetServers retruns list of servers
func (d *Discovery) GetServers() []*server.Server {
	return d.servers
}

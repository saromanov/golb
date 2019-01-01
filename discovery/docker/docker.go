package docker

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/saromanov/golb/discovery"
	"github.com/saromanov/golb/server"
)

const defaultEndpoint = "unix:///var/run/docker.sock"

// Discovery provides discoverying of the servers
// via docker
type Discovery struct {
	client  *docker.Client
	servers []*server.Server
	cfg     *discovery.Config
}

// New provides initialization of docker and Discovery
func New(cfg *discovery.Config) (Discovery, error) {
	endpoint := defaultEndpoint
	if cfg.DockerEndpoint != "" {
		endpoint = cfg.DockerEndpoint
	}
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return Discovery{}, err
	}

	return Discovery{
		client: client,
		cfg:    cfg,
	}, nil
}

// Search provides getting of containers and add their host
// on backend representation of server
func (d Discovery) Search() error {
	var filters map[string][]string
	if d.cfg.Filters != "" {
		filters = map[string][]string{
			"label": []string{d.cfg.Filters},
		}
	}
	containers, err := d.client.ListContainers(docker.ListContainersOptions{
		Filters: filters,
	})
	if err != nil {
		return err
	}

	for _, c := range containers {
		fmt.Println(c)
	}
	return nil
}

// GetServers retruns list of servers
func (d Discovery) GetServers() []*server.Server {
	return d.servers
}

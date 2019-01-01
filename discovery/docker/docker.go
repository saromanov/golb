package docker

import (
	"regexp"

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
		for _, p := range c.Ports {
			d.servers = append(d.servers, &server.Server{
				Host: d.getContainerHost(c.ID, p.IP),
				Port: uint32(p.PublicPort),
			})
		}
	}
	return nil
}

// GetServers retruns list of servers
func (d Discovery) GetServers() []*server.Server {
	return d.servers
}

func (d *Discovery) getContainerHost(id, portHost string) string {
	if portHost != "0.0.0.0" {
		return portHost
	}
	var reg = regexp.MustCompile("(.*?)://(?P<host>[-.A-Za-z0-9]+)/?(.*)")
	match := reg.FindStringSubmatch(d.cfg.DockerEndpoint)

	if len(match) == 0 {
		return portHost
	}

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if name != "" {
			result[name] = match[i]
		}
	}

	h, ok := result["host"]
	if !ok {
		return portHost
	}

	return h
}

package json

import (
	"errors"

	"github.com/saromanov/golb/config"
	"github.com/saromanov/golb/discovery"
	"github.com/saromanov/golb/server"
)

var errEmptyConfig = errors.New("config is not defined")

// Discovery provides importing of servers based on json
type Discovery struct {
	servers []*server.Server
}

// New provides getting of server definition from config file
func New(c *config.Config) (Discovery, error) {
	if c == nil {
		return Discovery{}, errEmptyConfig
	}

	srvs := getServers(c)
	return Discovery{
		servers: srvs,
	}, nil

}

// Search provides searching of servers
func (d Discovery) Search() error {
	servers := getServers(nil)
	d.servers = servers
	return nil
}

// getServers returns servers from config
func getServers(c *config.Config) []*server.Server {
	result := make([]*server.Server, len(c.Servers))
	for i, s := range c.Servers {
		result[i] = &server.Server{
			Host: s.Host,
			Port: s.Port,
		}
		result[i].ID = discovery.GenID(result[i])
	}
	return result
}

// GetServers retruns list of servers
func (d Discovery) GetServers() []*server.Server {
	return d.servers
}

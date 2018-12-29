package json

import (
	"errors"
	"github.com/saromanov/golb/server"
	"github.com/saromanov/golb/config"
)

var errEmptyConfig = errors.New("config is not defined")

// Discovery provides importing of servers based on json
type Discovery struct {
	servers []server.Server
}

// New provides getting of server definition from config file
func New(c* config.Config)(*Discovery, error) {
	if c == nil {
		return nil, errEmptyConfig
	}

}

// getServers returns servers from config
func getServers(c *config.Config)[]server.Server {

}
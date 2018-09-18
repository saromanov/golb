package golb

import (
	"time"

	"github.com/pkg/errors"
)

var errNoServers = errors.New("server is not defined")

type GoLB struct {
	Servers           []Server
	MaxConnections    uint32
	ClientTimeout     time.Time
	ConnectionTimeout time.Time
}

// SelectServers return server by the index
func (g *GoLB) SelectServer() (*Server, error) {
	if len(g.Servers) == 0 {
		return nil, errNoServers
	}

	return nil, nil
}

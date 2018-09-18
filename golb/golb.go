package golb

import (
	"time"

	"github.com/pkg/errors"
	"github.com/saromanov/golb/balancer"
)

var (
	errNoServers = errors.New("server is not defined")
	errUnknownBalanceType = errors.New("unknown balance type")
)

type GoLB struct {
	Servers           []Server
	MaxConnections    uint32
	ClientTimeout     time.Duration
	ConnectionTimeout time.Duration
	balance           *balancer.Balancer
	Balancer          string
}

// Build provides building of the GoLB
func (g *GoLB) Build() error {
	if balance == nil {
		g.balance = balancer.RoundRobin{Servers: g.Servers}
	}
	switch g.Balancer {
	case "rr":
		g.balance = balancer.RoundRobin{Servers: g.Servers}
	case "lc":
		g.balance = balancer.LeastConnect{Servers: g.Servers}
	default:
		return errUnknownBalanceType
	}
	return nil
}

// SelectServers return server by the index
func (g *GoLB) SelectServer() (*Server, error) {
	if len(g.Servers) == 0 {
		return nil, errNoServers
	}

	return nil, nil
}

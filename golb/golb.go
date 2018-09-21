package golb

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/saromanov/golb/balancer"
	"github.com/saromanov/golb/server"
)

var (
	errNoServers          = errors.New("server is not defined")
	errUnknownBalanceType = errors.New("unknown balance type")
	errNoBalancer         = errors.New("balancer is not defined")
)

type GoLB struct {
	Servers           []*server.Server
	MaxConnections    uint32
	ClientTimeout     time.Duration
	ConnectionTimeout time.Duration
	balance           balancer.Balancer
	Balancer          string
	Protocol          string
	Port              uint32
	Scheme            string
	Connections       uint32
}

// Build provides building of the GoLB
func (g *GoLB) Build() error {
	if g.balance == nil {
		g.balance = &balancer.RoundRobin{Servers: g.Servers}
	}
	switch g.Balancer {
	case "rr":
		g.balance = &balancer.RoundRobin{Servers: g.Servers}
	case "lc":
		g.balance = &balancer.LeastConnect{Servers: g.Servers}
	default:
		return errUnknownBalanceType
	}
	if g.Scheme == "" {
		g.Scheme = "http"
	}
	if g.Protocol == "" {
		g.Protocol = "tcp"
	}
	return nil
}

// AddServer adds a new server to the GoLB
func (g *GoLB) AddServer(s *server.Server) error {
	if s == nil || s.Host == "" {
		return fmt.Errorf("unable to add server")
	}
	g.Servers = append(g.Servers, s)
	return nil
}

// SelectServer return server by the index
func (g *GoLB) SelectServer() (*server.Server, error) {
	if len(g.Servers) == 0 {
		return nil, errNoServers
	}
	if g.balance == nil {
		return nil, errNoBalancer
	}
	serv, err := g.balance.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to apply balancing: %v", err)
	}
	return serv, nil
}

// HandleHTTP implements middleware for http requests
func (g *GoLB) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	serv, err := g.SelectServer()
	switch err {
	case errNoServers:
		log.Printf(err.Error())
		return
	case errUnknownBalanceType:
		log.Printf(err.Error())
		return
	default:
		// TODO: try again
	}
	g.Connections++
	proxy := &HTTPProxy{serv: serv, Scheme: g.Scheme}
	err = proxy.Do(w, r)
	switch err := err.(type) {
	case urlParseError:
		log.Printf("Err Parse: %v", err)
	case httpRequestError:
		log.Printf("HandleHTTP error: %v", err)
	}
	serv.IncSuccessRequests()
}

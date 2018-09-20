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

	fmt.Println(g.balance)
	serv, err := g.balance.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to apply balancing: %v", err)
	}
	return serv, nil
}

// HandleHTTP implements middleware for http requests
func (g *GoLB) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	serv, err := g.SelectServer()
	if err != nil {
		log.Printf("Err: %v", err)
	}

	fmt.Println(serv)
	proxy := &HTTPProxy{serv: serv}
	err = proxy.Do(w, r)
	if err != nil {
		log.Printf("Err Proxy: %v", err)
	}
}

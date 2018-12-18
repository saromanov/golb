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
	Servers             []*server.Server
	MaxConnections      uint32
	ClientTimeout       time.Duration
	ConnectionTimeout   time.Duration
	balance             balancer.Balancer
	Balancer            string
	Protocol            string
	Port                uint32
	Scheme              string
	Connections         uint32
	ProxyHeaders        map[string]string
	FailedRequestsLimit uint32
	Stats               *Stats
	ServerScheme        string
	KeyFilePath         string
	CertFilePath        string
}

//New returns golb object after reading of config
func New(conf *Config) golb.GoLB {
	g := golb.GoLB{
		MaxConnections:      conf.MaxConnections,
		Balancer:            conf.Balancer,
		Protocol:            conf.Protocol,
		Port:                conf.Port,
		Scheme:              conf.Scheme,
		FailedRequestsLimit: conf.FailedRequestsLimit,
		ServerScheme:        conf.ServerScheme,
		CertFilePath:        conf.CertFilePath,
		KeyFilePath:         conf.KeyFilePath,
	}

	servers := []*server.Server{}
	for _, s := range conf.Servers {
		servers = append(servers, &server.Server{
			Host: s.Host,
			Port: s.Port,
		})
	}

	g.Servers = servers
	return g
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
	case "wrr":
		g.balance = &balancer.WeightedRoundRobin{Servers: g.Servers}
	default:
		return errUnknownBalanceType
	}
	if g.Scheme == "" {
		g.Scheme = "http"
	}
	if g.Protocol == "" {
		g.Protocol = "tcp"
	}
	if len(g.ProxyHeaders) == 0 {
		g.ProxyHeaders = map[string]string{}
	}
	g.Stats = &Stats{StatusCodes: map[int]uint32{}}
	return nil
}

// AddServer adds a new server to the GoLB
func (g *GoLB) AddServer(s *server.Server) error {
	if s == nil || s.Host == "" {
		return fmt.Errorf("unable to add server")
	}
	g.Servers = append(g.Servers, s)
	g.Stats.Servers++
	return nil
}

// SelectServer return server by the index
func (g *GoLB) SelectServer() (*server.Server, error) {
	g.Stats.Requests++
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
	g.Stats.CompleteRequests++
	return serv, nil
}

// GetStats returns stats for GoLB
func (g *GoLB) GetStats() *Stats {
	return g.Stats
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
	resp, err := proxy.Do(w, r)
	switch err := err.(type) {
	case urlParseError:
		log.Printf("Err Parse: %v", err)
	case httpRequestError:
		log.Printf("HandleHTTP error: %v", err)
		g.checkFailedRequests(serv)
	}
	g.updateStats(serv, resp)
}

func (g *GoLB) updateStats(s *server.Server, r *HTTPProxyResponse) {
	g.Stats.StatusCodes[r.statusCode]++
	s.IncSuccessRequests()
}

// checkfailedRequests increments number of failed requests
// for the server and if this is reached limit,
// then its removing server from the list
func (g *GoLB) checkFailedRequests(serv *server.Server) {
	serv.FailedRequests++
	if g.FailedRequestsLimit != 0 && serv.FailedRequests > g.FailedRequestsLimit {
		log.Printf("Remove server: %s:%d from the list ", serv.Host, serv.Port)
		serv.RemovedFromList = true
		for i, s := range g.Servers {
			if s.RemovedFromList {
				g.Servers = append(g.Servers[:i], g.Servers[i+1:]...)
			}
		}
	}
}

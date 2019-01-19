package golb

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/saromanov/golb/balancer"
	"github.com/saromanov/golb/config"
	"github.com/saromanov/golb/discovery"
	"github.com/saromanov/golb/discovery/docker"
	"github.com/saromanov/golb/discovery/json"
	"github.com/saromanov/golb/logger"
	"github.com/saromanov/golb/server"
)

var (
	errNoServers           = errors.New("server is not defined")
	errUnknownBalancerType = errors.New("unknown balancer type")
	errNoBalancer          = errors.New("balancer is not defined")
	errServerNotFound      = errors.New("server not found")
)

// GoLB defines main struct of app
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
	mu                  *sync.RWMutex
	conf                *config.Config
	disc                discovery.Discovery
}

//New returns golb object after reading of config
func New(conf *config.Config) *GoLB {
	g := &GoLB{
		MaxConnections:      conf.MaxConnections,
		Balancer:            conf.Balancer,
		Protocol:            conf.Protocol,
		Port:                conf.Port,
		Scheme:              conf.Scheme,
		FailedRequestsLimit: conf.FailedRequestsLimit,
		ServerScheme:        conf.ServerScheme,
		CertFilePath:        conf.CertFilePath,
		KeyFilePath:         conf.KeyFilePath,
		mu:                  &sync.RWMutex{},
		conf:                conf,
	}
	conf.Discovery = "docker"

	switch conf.Discovery {
	case "docker":
		d, err := docker.New(&discovery.Config{})
		if err != nil {
			log.Fatalf("unable to discover servers: %v", err)
		}
		g.disc = d
		servers := d.GetServers()
		logger.Infof("Docker: discovered %d servers", len(servers))
		g.Servers = servers
	default:
		d, err := json.New(conf)
		if err != nil {
			log.Fatalf("unable to discover servers: %v", err)
		}
		g.disc = d
		servers := d.GetServers()
		logger.Infof("Json: discovered %d servers", len(servers))
		g.Servers = servers
	}
	return g
}

// Build provides building of the GoLB
func (g *GoLB) Build() error {
	if g.balance == nil {
		g.balance = &balancer.RoundRobin{Servers: g.Servers}
	}
	fmt.Println(g.conf)
	switch g.conf.Balancer {
	case "rr":
		g.balance = &balancer.RoundRobin{Servers: g.Servers}
	case "lc":
		g.balance = &balancer.LeastConnect{Servers: g.Servers}
	case "wrr":
		g.balance = &balancer.WeightedRoundRobin{Servers: g.Servers}
	default:
		return errUnknownBalancerType
	}
	if g.conf.Scheme == "" {
		g.conf.Scheme = "http"
	}
	if g.conf.Protocol == "" {
		g.conf.Protocol = "tcp"
	}
	if len(g.ProxyHeaders) == 0 {
		g.ProxyHeaders = map[string]string{}
	}

	if g.Scheme == "https" {
		go makeHTTPSMetricsServer(g.conf)
	} else {
		go makeHTTPMetricsServer()
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

// RemoveServer provides removing of server from list
func (g *GoLB) RemoveServer(ID string) error {
	for i, x := range g.Servers {
		if x.ID == ID {
			g.Servers = append(g.Servers[:i], g.Servers[i+1:]...)
			return nil
		}
	}
	return errServerNotFound
}

// SelectServer return server by the index
func (g *GoLB) SelectServer() (*server.Server, error) {
	servers := g.disc.GetServers()
	g.Stats.Requests++
	if len(servers) == 0 {
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
	case errUnknownBalancerType:
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
	g.mu.RLock()
	defer g.mu.RUnlock()
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

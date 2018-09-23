package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/saromanov/golb/golb"
	"github.com/saromanov/golb/server"
)

// Config provides definition of the config
type Config struct {
	MaxConnections      uint32            `json:"max_connections"`
	ClientTimeout       string            `json:"client_timeout"`
	ConnectionTimeout   string            `json:"connection_timeout"`
	Balancer            string            `json:"balancer"`
	Protocol            string            `json:"protocol"`
	Port                uint32            `json:"port"`
	Scheme              string            `json:"scheme"`
	ProxyHeaders        map[string]string `json:"proxy_headers"`
	Servers             []Server          `json:"servers"`
	FailedRequestsLimit uint32            `json:"failed_requests_limit"`
}

// Server defines config for the server
type Server struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

// ReadConfig provides redaing of the config
func ReadConfig(path string) (*Config, error) {
	var conf *Config
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %v", err)
	}
	return conf, nil
}

//MakeGoLBObject returns golb object after reading of config
func MakeGoLBObject(conf *Config) golb.GoLB {
	g := golb.GoLB{
		MaxConnections:      conf.MaxConnections,
		Balancer:            conf.Balancer,
		Protocol:            conf.Protocol,
		Port:                conf.Port,
		Scheme:              conf.Scheme,
		FailedRequestsLimit: conf.FailedRequestsLimit,
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

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/saromanov/golb/golb"
)

// Config provides definition of the config
type Config struct {
	MaxConnections    uint32            `json:"max_connections"`
	ClientTimeout     string            `json:"client_timeout"`
	ConnectionTimeout string            `json:"connection_timeout"`
	Balancer          string            `json:"balancer"`
	Protocol          string            `json:"protocol"`
	Port              uint32            `json:"port"`
	Scheme            string            `json:"scheme"`
	ProxyHeaders      map[string]string `json:"proxy_headers"`
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
	return golb.GoLB{
		MaxConnections: conf.MaxConnections,
		Balancer:       conf.Balancer,
		Protocol:       conf.Protocol,
		Port:           conf.Port,
		Scheme:         conf.Scheme,
	}
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
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
	ServerScheme        string            `json:"server_scheme"`
	CertFilePath        string            `json:"cert_file_path"`
	KeyFilePath         string            `json:"key_file_path"`
	Discovery           string            `json:"discovery"`
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

// MakeDefaultConfig provides loading of default conig
func MakeDefaultConfig() *Config {
	return &Config{
		MaxConnections:    10,
		ClientTimeout:     "5s",
		ConnectionTimeout: "5s",
		Balancer:          "rr",
	}
}

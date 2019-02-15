package config

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/consul/api"
)

// ReadConsulConfig provides reading of config from Comsul
func ReadConsulConfig(key string) (*Config, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("consul: unable to create client: %v", err)
	}
	kv := client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("consul: unable to get value by the key: %v", err)
	}

	var c *Config
	err = json.Unmarshal(pair.Value, &c)
	if err != nil {
		return nil, fmt.Errorf("consul: unable to unmarshal response: %v", err)
	}
	return c, nil
}

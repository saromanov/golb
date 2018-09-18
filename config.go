package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
}

// ReadConfig provides redaing of the config
func ReadConfig(path string) (*Config, error) {
	var conf *Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

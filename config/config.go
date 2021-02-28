package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Config struct {
	Instances []*InstancesConfig `yaml:"instances,omitempty"`
	Port string `yaml:"port,omitempty"`
}

type InstancesConfig struct {
	Name     string `yaml:"name,omitempty"`
	Address  string   `yaml:"address"`
}

func ConfigParse(r io.Reader) (*Config, error) {
	// read everything from io.Reader
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Create config instance from yaml
	config := &Config{}
	err = yaml.Unmarshal(buffer, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ConfigLoadFromFile() (*Config, error) {
	// Load from file
	buffer, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// Parse config
	config, err := ConfigParse(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}

	if config.Port == "" {
		config.Port = "9879"
	}
	return config, nil
}
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token     string  `yaml:"token"`
	Proxy     string  `yaml:"proxy"`
	Whitelist []int64 `yaml:"whitelist"`
}

func Load() (*Config, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	data := &Config{}
	return data, yaml.Unmarshal(file, data)
}

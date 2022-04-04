package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type HTTPConf struct {
	Host string
	Port string
}

type GRPCConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level string
	File  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

func NewConfig() Config {
	return Config{}
}

func LoadConfiguration(configFile string) (*Config, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := NewConfig()
	err = yaml.Unmarshal(content, &newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file  %s: %w", configFile, err)
	}

	return &newConfig, nil
}

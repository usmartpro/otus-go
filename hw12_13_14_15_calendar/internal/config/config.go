package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type CalendarConf struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type SchedulerConf struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
}

type SenderConf struct {
	Logger  LoggerConf
	Storage StorageConf
	Rabbit  RabbitConf
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

type RabbitConf struct {
	Dsn      string
	Exchange string
	Queue    string
}

func NewCalendarConfig() CalendarConf {
	return CalendarConf{}
}

func NewSchedulerConfig() SchedulerConf {
	return SchedulerConf{}
}

func NewSenderConfig() SenderConf {
	return SenderConf{}
}

func LoadCalendarConfiguration(configFile string) (*CalendarConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := NewCalendarConfig()
	err = yaml.Unmarshal(content, &newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return &newConfig, nil
}

func LoadSchedulerConfig(configFile string) (*SchedulerConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := NewSchedulerConfig()
	err = yaml.Unmarshal(content, &newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return &newConfig, nil
}

func LoadSenderConfig(configFile string) (*SenderConf, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("wrong configuration file %s: %w", configFile, err)
	}

	newConfig := NewSenderConfig()
	err = yaml.Unmarshal(content, &newConfig)
	if err != nil {
		return nil, fmt.Errorf("wrong params in configuration file %s: %w", configFile, err)
	}

	return &newConfig, nil
}

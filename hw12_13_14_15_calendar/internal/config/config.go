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
}

type (
	Level   string
	Storage string
)

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

type LoggerConf struct {
	Level    Level
	Filename string
}

const (
	SQL       Storage = "sql"
	InMemmory Storage = "in-memory"
)

type StorageConf struct {
	Type Storage
	URL  Storage
}

type HTTPConf struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(path string) (*Config, error) {
	resultConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("invalid config %s: %w", path, err)
	}

	config := NewConfig()
	err = yaml.Unmarshal(resultConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid unmarshal config %s:%w", path, err)
	}

	return &config, nil
}

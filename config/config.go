package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseConfig struct {
	Name                  string `json:"name"`
	Host                  string `json:"host"`
	Port                  uint16 `json:"port"`
	User                  string `json:"user"`
	MaxConns              int32  `json:"pool_max_conns"`
	MinConns              int32  `json:"pool_min_conns"`
	MaxConnLifetime       string `json:"pool_max_conn_lifetime"`
	MaxConnIdleTime       string `json:"pool_max_conn_idle_time"`
	HealthCheckPeriod     string `json:"pool_health_check_period"`
	MaxConnLifetimeJitter string `json:"pool_max_conn_lifetime_jitter"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

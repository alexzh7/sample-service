package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Application configuration
type Config struct {
	Postgres PostgresConfig
	GRPC     GRPCConfig
}

// Postgresql config
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GRPC config
type GRPCConfig struct {
	Port string
}

// NewConfig parses config file and returns app config
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.ReadInConfig: %v", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("config.Unmarshal: %v", err)
	}

	return &config, nil
}

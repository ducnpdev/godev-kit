package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		GRPC    GRPC
		RMQ     RMQ
		Metrics Metrics
		Swagger Swagger
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME" mapstructure:"APP_NAME"`
		Version string `env:"APP_VERSION" mapstructure:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT" mapstructure:"HTTP_PORT"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false" mapstructure:"HTTP_USE_PREFORK_MODE"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL" mapstructure:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX" mapstructure:"PG_POOL_MAX"`
		URL     string `env:"PG_URL" mapstructure:"PG_URL"`
	}

	// GRPC -.
	GRPC struct {
		Port string `env:"GRPC_PORT" mapstructure:"GRPC_PORT"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER" mapstructure:"RMQ_RPC_SERVER"`
		ClientExchange string `env:"RMQ_RPC_CLIENT" mapstructure:"RMQ_RPC_CLIENT"`
		URL            string `env:"RMQ_URL" mapstructure:"RMQ_URL"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" mapstructure:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" mapstructure:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// EnvValue configuration -.
	EnvValue struct {
		Type string `env:"TYPE" mapstructure:"TYPE" envDefault:"env"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

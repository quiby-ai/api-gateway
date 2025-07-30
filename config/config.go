package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Host         string
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv(EnvGatewayPort, DefaultPort),
			Host:         getEnv(EnvGatewayHost, DefaultHost),
			ReadTimeout:  getEnvDuration(EnvGatewayReadTimeout, DefaultReadTimeout),
			WriteTimeout: getEnvDuration(EnvGatewayWriteTimeout, DefaultWriteTimeout),
			IdleTimeout:  getEnvDuration(EnvGatewayIdleTimeout, DefaultIdleTimeout),
		},
		Logging: LoggingConfig{
			Level:  getEnv(EnvLogLevel, DefaultLogLevel),
			Format: getEnv(EnvLogFormat, DefaultLogFormat),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	var errors []ValidationError

	if err := c.validateServer(); err != nil {
		errors = append(errors, *err)
	}

	if err := c.validateLogging(); err != nil {
		errors = append(errors, *err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %v", errors)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if duration, err := time.ParseDuration(v); err == nil {
			return duration
		}
	}
	return fallback
}

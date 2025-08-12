package config

import (
	"time"

	"github.com/quiby-ai/api-gateway/internal/gateway"
	"github.com/spf13/viper"
)

type GatewayConfig struct {
	Routes []gateway.Route `mapstructure:"routes"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
	AllowedMethods []string `mapstructure:"allowedMethods"`
	AllowedHeaders []string `mapstructure:"allowedHeaders"`
	MaxAge         int      `mapstructure:"maxAge"`
}

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Logging LoggingConfig `mapstructure:"logging"`
	Gateway GatewayConfig `mapstructure:"gateway"`
	CORS    CORSConfig    `mapstructure:"cors"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

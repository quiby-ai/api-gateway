package config

import "time"

const (
	DefaultPort         = "8080"
	DefaultHost         = "0.0.0.0"
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultIdleTimeout  = 60 * time.Second

	DefaultLogLevel  = "info"
	DefaultLogFormat = "json"
)

const (
	EnvGatewayPort         = "GATEWAY_PORT"
	EnvGatewayHost         = "GATEWAY_HOST"
	EnvGatewayReadTimeout  = "GATEWAY_READ_TIMEOUT"
	EnvGatewayWriteTimeout = "GATEWAY_WRITE_TIMEOUT"
	EnvGatewayIdleTimeout  = "GATEWAY_IDLE_TIMEOUT"

	EnvLogLevel  = "LOG_LEVEL"
	EnvLogFormat = "LOG_FORMAT"
)

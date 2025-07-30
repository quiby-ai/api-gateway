package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// validateServer validates server configuration
func (c *Config) validateServer() *ValidationError {
	// Validate port
	if c.Server.Port == "" {
		return &ValidationError{Field: "Server.Port", Message: "port is required"}
	}

	if port, err := strconv.Atoi(c.Server.Port); err != nil || port <= 0 || port > 65535 {
		return &ValidationError{Field: "Server.Port", Message: "port must be a valid number between 1 and 65535"}
	}

	// Validate host
	if c.Server.Host == "" {
		return &ValidationError{Field: "Server.Host", Message: "host is required"}
	}

	// Validate timeouts
	if c.Server.ReadTimeout <= 0 {
		return &ValidationError{Field: "Server.ReadTimeout", Message: "read timeout must be positive"}
	}

	if c.Server.WriteTimeout <= 0 {
		return &ValidationError{Field: "Server.WriteTimeout", Message: "write timeout must be positive"}
	}

	if c.Server.IdleTimeout <= 0 {
		return &ValidationError{Field: "Server.IdleTimeout", Message: "idle timeout must be positive"}
	}

	return nil
}

// validateLogging validates logging configuration
func (c *Config) validateLogging() *ValidationError {
	// Validate log level
	validLogLevels := []string{"debug", "info", "warn", "warning", "error", "fatal", "panic"}
	if !slices.Contains(validLogLevels, strings.ToLower(c.Logging.Level)) {
		return &ValidationError{Field: "Logging.Level", Message: fmt.Sprintf("log level must be one of: %s", strings.Join(validLogLevels, ", "))}
	}

	// Validate log format
	validLogFormats := []string{"json", "text", "console"}
	if !slices.Contains(validLogFormats, strings.ToLower(c.Logging.Format)) {
		return &ValidationError{Field: "Logging.Format", Message: fmt.Sprintf("log format must be one of: %s", strings.Join(validLogFormats, ", "))}
	}

	return nil
}

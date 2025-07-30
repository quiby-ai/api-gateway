# API Gateway

Production-ready Go API Gateway with pluggable middleware for authentication, rate limiting, caching, and monitoring, providing a unified reverse-proxy interface for microservices.

## Configuration

The API Gateway uses environment variables for configuration. Copy `env.example` to `.env` and update the values as needed.

### Environment Variables

#### Server Configuration
- `GATEWAY_PORT` - Server port (default: 8080)
- `GATEWAY_HOST` - Server host (default: 0.0.0.0)
- `GATEWAY_READ_TIMEOUT` - Read timeout (default: 5s)
- `GATEWAY_WRITE_TIMEOUT` - Write timeout (default: 10s)
- `GATEWAY_IDLE_TIMEOUT` - Idle timeout (default: 60s)

#### Logging Configuration
- `LOG_LEVEL` - Log level (default: info)
- `LOG_FORMAT` - Log format (default: json)

### Quick Start

1. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

2. Run the gateway:
```bash
go run cmd/main.go
```

3. Test the health endpoint:
```bash
curl http://localhost:8080/healthz
```

## Development

### Project Structure
```
api-gateway/
├── cmd/
│   └── main.go          # Application entry point
├── config/
│   ├── config.go        # Configuration structs and loading
│   ├── defaults.go      # Default values and env var names
│   └── validation.go    # Configuration validation
├── internal/            # Private application code
├── pkg/                 # Public libraries
└── env.example         # Example environment variables
```

### Configuration Validation

The configuration system includes comprehensive validation:
- Port numbers must be between 1-65535
- Timeouts must be positive
- Log levels and formats are validated against allowed values

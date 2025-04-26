# Simple Redis-based Load Balancer

A lightweight HTTP load balancer implementation using Go, Gin framework, and Redis for tracking server loads.

## Overview

This project implements a basic load balancing service that distributes incoming requests across multiple backend servers. It uses Redis sorted sets to keep track of the current load on each server and directs new requests to the server with the lowest current load.

## Features

- Server selection based on current load using Redis sorted sets
- Dynamic scoring system that increments when a server receives a request and decrements when processing is complete
- Basic server health checking through a heartbeat endpoint
- Simple HTTP request forwarding

## Requirements

- Go 1.16+
- Redis server (running on localhost:6379)
- Backend servers configured with a `/heartbeat` endpoint

## Installation

1. Clone this repository
2. Install dependencies:
```
go mod tidy
```

## Usage

1. Ensure Redis is running on localhost:6379
2. Start your backend servers on ports 8080, 8081, 8082, and 8083
3. Run the load balancer:
```
go run main.go
```
4. The load balancer will be available at http://localhost:9090
5. Access the load balancer by sending requests to http://localhost:9090/gethandler

## How It Works

1. The load balancer initializes a Redis sorted set with all backend server ports, each with an initial score of 0
2. When a request comes in:
   - The server with the lowest score is selected
   - Its score is incremented by 1
   - The request is forwarded to that server
   - Upon completion, the server's score is decremented by 1

## Areas for Improvement

### Error Handling
- Replace `panic()` calls with proper error handling to prevent service disruption
- Add retry mechanisms for failed Redis operations
- Implement graceful degradation when backend servers are unavailable

### Health Checking
- Add proper health check logic that validates server responses
- Implement periodic health checks instead of checking only during requests
- Add a mechanism to temporarily remove unhealthy servers from the pool

### Connection Management
- Properly close Redis connections when the service shuts down
- Implement connection pooling for better performance
- Add timeouts for Redis operations and HTTP requests

### Load Balancing Algorithm
- Improve server selection strategy to avoid the "thundering herd" problem
- Consider weighted round-robin or other advanced algorithms
- Add support for server weights or priority levels

### Concurrency Handling
- Make score increment/decrement operations atomic to prevent race conditions
- Implement proper locking or use Redis transactions where appropriate
- Add request rate limiting

### Resilience
- Implement circuit breaking for consistently failing servers
- Add exponential backoff for retry attempts
- Create a mechanism to dynamically add/remove servers from the pool

### Configuration
- Move hardcoded values to configuration files or environment variables
- Make Redis connection parameters configurable
- Allow dynamic configuration of backend servers

### Logging and Monitoring
- Replace fmt.Println with a proper logging framework
- Add structured logging for better filtering and analysis
- Implement metrics collection for monitoring load balancer performance

### Context Handling
- Use consistent context objects throughout the application
- Properly propagate context cancellations and timeouts
- Add request tracing for better debugging

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

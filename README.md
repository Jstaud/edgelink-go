# EdgeLink-Go

A no-hardware SCADA gateway built in Go that simulates industrial devices, polls them concurrently, and publishes data via MQTT while serving a REST API.

## Features

- **Concurrent polling**: One goroutine per device with context-based timeouts
- **MQTT publishing**: QoS 1 message delivery to MQTT broker
- **REST API**: Query latest readings via HTTP
- **Prometheus metrics**: Built-in observability 
- **Graceful shutdown**: Proper signal handling and cleanup
- **Fake devices**: Realistic sensor simulation with jitter

## Quick Start

1. **Run with Docker Compose:**
   ```bash
   docker compose up --build
   ```

2. **Test the API:**
   ```bash
   # Get latest reading from a device
   curl localhost:8080/api/readings/press-1 | jq
   
   # Check metrics
   curl localhost:8080/metrics
   
   # Health check
   curl localhost:8080/healthz
   ```

3. **Monitor MQTT** (if you have mosquitto_sub):
   ```bash
   mosquitto_sub -h localhost -t 'plant/readings/#' -v
   ```

## Configuration

Edit `configs/example.yaml` to configure:
- HTTP server address
- MQTT broker settings  
- Device specifications (poll intervals, metrics, jitter)

## Project Structure

```
cmd/edgelink/           # Application entry point
internal/
  api/                  # REST API handlers
  config/               # Configuration loading
  device/               # Device driver interfaces
    fake/               # Fake device implementation
  log/                  # Logging setup
  observability/        # Prometheus metrics
  pipeline/             # Data processing pipeline
    publisher/          # MQTT publishing
pkg/models/             # Domain models
configs/                # Configuration files
```

## Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Run locally:**
   ```bash
   go run ./cmd/edgelink --config configs/example.yaml
   ```

3. **Build:**
   ```bash
   go build -o bin/edgelink ./cmd/edgelink
   ```

## Next Steps

- Add real device drivers (Modbus, OPC UA)
- Implement circuit breaker pattern
- Add persistent storage for readings
- Create Grafana dashboards for visualization

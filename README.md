# OpenTelemetry Demo

A hands-on demonstration of distributed tracing with OpenTelemetry in a Go microservice architecture.

## Author

Elda Mahaindra ([faith030@gmail.com](mailto:faith030@gmail.com))

## Overview

This project demonstrates how to implement **distributed tracing** with **OpenTelemetry** in a Go application. It showcases the complete observability pipeline from trace generation to visualization, helping you understand how OpenTelemetry works in practice.

### Architecture

```
┌─────────────┐    ┌──────────────────┐    ┌──────────────┐    ┌─────────────┐
│   Client    │───▶│    service-a     │───▶│ OpenTelemetry│───▶│   Jaeger    │
│   (curl)    │    │   (Go/Fiber)     │    │  Collector   │    │    UI       │
└─────────────┘    └──────────────────┘    └──────────────┘    └─────────────┘
                           │                      │                    │
                    ┌──────┴──────┐               │              ┌─────┴─────────┐
                    │   3 Layers  │               │              │  Traces       │
                    │ • API       │               │              │ Visualization │
                    │ • Service   │               │              └───────────────┘
                    │ • Store     │               │
                    └─────────────┘               │
                                          ┌───────┴──────────┐
                                          │ Trace Processing │
                                          │ • OTLP Receiver  │
                                          │ • Zipkin Export  │
                                          │ • Debug Logging  │
                                          └──────────────────┘
```

## Features

- ✅ **3-Layer Architecture**: API → Service → Store with tracing at each level
- ✅ **OpenTelemetry Integration**: Complete OTLP implementation
- ✅ **Distributed Context Propagation**: Traces flow through all layers
- ✅ **Rich Span Attributes**: Input/output data, operation metadata
- ✅ **Span Events**: Database operation timing markers
- ✅ **Error Handling**: Proper error recording in traces
- ✅ **Performance Monitoring**: Processing time visualization
- ✅ **Jaeger Integration**: Beautiful trace visualization UI

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)

### Running the Demo

1. **Start all services:**

   ```bash
   docker compose up -d
   ```

2. **Make a test request:**

   ```bash
   curl "http://localhost:4000/ping?message=Hello%20World"
   ```

3. **View traces in Jaeger:**
   - Open http://localhost:16686
   - Select "service-a" from the service dropdown
   - Click "Find Traces"
   - Explore the trace hierarchy and span details

### Test Different Scenarios

```bash
# Normal request
curl "http://localhost:4000/ping?message=test"

# Trigger an error (to see error tracing)
curl "http://localhost:4000/ping?message=error"

# Test with special characters
curl "http://localhost:4000/ping?message=special%20test"
```

## OpenTelemetry Implementation

### Dependencies Added

Added to `service-a/go.mod`:

- `go.opentelemetry.io/otel` - Core OpenTelemetry library
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` - OTLP gRPC exporter
- `go.opentelemetry.io/otel/sdk` - OpenTelemetry SDK
- `go.opentelemetry.io/otel/trace` - Tracing interface

### Tracing Components

#### 1. Tracer Initialization (`service-a/util/tracing/tracing.go`)

- `InitTracer()` - Sets up OpenTelemetry with OTLP exporter
- `GetTracer()` - Provides tracer instances for span creation
- Configured to send traces to OpenTelemetry Collector

#### 2. API Layer Tracing (`service-a/api/ping.go`)

- **Root span creation** for incoming HTTP requests
- **Request attributes**: endpoint, method, query parameters
- **Context propagation** to downstream layers
- **Error and success status** recording

#### 3. Service Layer Tracing (`service-a/service/ping.go`)

- **Child span creation** for business logic
- **Operation metadata**: input/output parameters
- **Processing simulation** (500ms delay)
- **Context forwarding** to data layer

#### 4. Store Layer Tracing (`service-a/store/ping.go`)

- **Data access spans** for storage operations
- **Database simulation** (1 second delay)
- **Span events**: query start/end timestamps
- **Operation metrics**: duration, query type

### Observability Features

#### Span Attributes (Visible in Jaeger "Tags" section)

```go
span.SetAttributes(
    attribute.String("api.endpoint", "/ping"),
    attribute.String("api.method", "GET"),
    attribute.String("ping.message", message),
    attribute.String("response.pong_message", result.PongMessage),
)
```

#### Span Events (Visible in Jaeger "Events" section)

```go
span.AddEvent("database_query_start")
time.Sleep(1 * time.Second)
span.AddEvent("database_query_end")
```

## Infrastructure Components

### OpenTelemetry Collector (`otel-collector-config.yaml`)

- **OTLP Receiver**: Accepts traces via gRPC (4317) and HTTP (4318)
- **Batch Processor**: Optimizes trace export performance
- **Zipkin Exporter**: Sends traces to Jaeger in Zipkin format
- **Debug Exporter**: Console logging for development

### Jaeger All-in-One

- **Trace Storage**: In-memory storage for demo purposes
- **Zipkin Compatibility**: Receives traces via Zipkin API (port 9411)
- **Web UI**: Trace visualization and analysis (port 16686)
- **OTLP Support**: Additional OTLP gRPC endpoint (port 14250)

### Service Ports

- **service-a**: 4000 (HTTP API)
- **OpenTelemetry Collector**: 4317 (OTLP gRPC), 4318 (OTLP HTTP)
- **Jaeger UI**: 16686 (Web interface)
- **Jaeger Zipkin**: 9411 (Trace ingestion)

## Understanding Traces

### Trace Hierarchy

When you make a request, you'll see a trace with three spans:

1. **`api.ping`** - HTTP request handling (API layer)

   - Duration: ~1.5 seconds total
   - Attributes: HTTP method, endpoint, input message
   - Status: Success/Error indication

2. **`service.ping`** - Business logic processing (Service layer)

   - Duration: ~1.5 seconds (includes store call)
   - Attributes: Operation type, input/output data
   - Parent: api.ping

3. **`store.ping`** - Data access simulation (Store layer)
   - Duration: ~1 second (simulated database query)
   - Attributes: Query type, database operation details
   - Events: Query start/end timestamps
   - Parent: service.ping

### Trace Analysis Benefits

- **Performance Monitoring**: Identify bottlenecks in specific layers
- **Error Tracking**: See exactly where and why errors occur
- **Request Correlation**: Follow a single request through all systems
- **Dependency Mapping**: Understand service communication patterns
- **Resource Utilization**: Monitor processing time across components

## Development

### Project Structure

```
.
├── docker-compose.yml           # Infrastructure orchestration
├── otel-collector-config.yaml   # OpenTelemetry collector configuration
├── README.md                    # This file
└── service-a/                   # Go microservice
    ├── api/                     # HTTP handlers with tracing
    ├── service/                 # Business logic with tracing
    ├── store/                   # Data layer with tracing
    ├── util/tracing/            # OpenTelemetry utilities
    ├── cmd/                     # Application entry points
    ├── config.json              # Service configuration
    ├── Dockerfile               # Container build instructions
    └── go.mod                   # Go dependencies
```

### Building service-a

```bash
cd service-a
go build -o service-a ./cmd/
./service-a start
```

### Running Tests

```bash
# Test normal flow
curl "http://localhost:4000/ping?message=test"

# Test error handling
curl "http://localhost:4000/ping?message=error"
```

## Troubleshooting

### Common Issues

1. **"No traces found in Jaeger"**

   - Ensure all services are running: `docker compose ps`
   - Check service-a logs: `docker compose logs service-a`
   - Verify OpenTelemetry collector logs: `docker compose logs otel-collector`

2. **"Connection refused" errors**

   - Restart services: `docker compose down && docker compose up -d`
   - Check network connectivity between containers

3. **"Service-a not in dropdown"**
   - Make a request first to generate traces
   - Wait 10-15 seconds for trace processing
   - Refresh Jaeger UI

### Logs

```bash
# View all logs
docker compose logs

# View specific service logs
docker compose logs service-a
docker compose logs otel-collector
docker compose logs jaeger
```

## Learning Resources

This demo helps you understand:

- **OpenTelemetry fundamentals**: Traces, spans, attributes, events
- **Distributed tracing patterns**: Context propagation, span hierarchy
- **Observability best practices**: Error handling, performance monitoring
- **Infrastructure setup**: Collectors, exporters, backends

## License

This project is for educational purposes. Feel free to use it as a reference for implementing OpenTelemetry in your own applications.

---

_Built with ❤️ to help developers understand distributed tracing with OpenTelemetry_

# OpenTelemetry Demo

A hands-on demonstration of distributed tracing with OpenTelemetry in a Go microservice architecture.

## Author

Elda Mahaindra ([faith030@gmail.com](mailto:faith030@gmail.com))

## Overview

This project demonstrates how to implement **distributed tracing** with **OpenTelemetry** in a Go microservice architecture. It showcases the complete observability pipeline from trace generation to visualization, including **cross-service context propagation** via gRPC, helping you understand how OpenTelemetry works in practice across multiple services.

### Architecture

```
┌─────────────┐    ┌──────────────────┐    ┌──────────────────┐    ┌──────────────┐    ┌─────────────┐
│   Client    │───▶│    service-a     │───▶│    service-b     │───▶│ OpenTelemetry│───▶│   Jaeger    │
│   (curl)    │    │   (Go/Fiber)     │    │   (Go/gRPC)      │    │  Collector   │    │    UI       │
└─────────────┘    └──────────────────┘    └──────────────────┘    └──────────────┘    └─────────────┘
                           │                         │                       │                    │
                    ┌──────┴──────┐          ┌───────┴───────┐               │              ┌─────┴─────────┐
                    │   3 Layers  │          │   3 Layers    │               │              │  Traces       │
                    │ • API       │   gRPC   │ • API         │               │              │ Visualization │
                    │ • Service   │◄────────▶│ • Service     │               │              └───────────────┘
                    │ • Adapter   │          │ • Store       │               │
                    └─────────────┘          └───────────────┘               │
                                                     │                       │
                                            ┌────────┴────────┐              │
                                            │  Database Ops   │              │
                                            │ • Query Sim     │              │
                                            │ • Events        │              │
                                            └─────────────────┘              │
                                                                    ┌────────┴──────────┐
                                                                    │ Trace Processing  │
                                                                    │ • OTLP Receiver   │
                                                                    │ • Zipkin Export   │
                                                                    │ • Debug Logging   │
                                                                    └───────────────────┘
```

## Features

- ✅ **Multi-Service Architecture**: service-a (HTTP) → service-b (gRPC) with tracing
- ✅ **OpenTelemetry Integration**: Complete OTLP implementation across services
- ✅ **Cross-Service Context Propagation**: Traces flow seamlessly between services via gRPC
- ✅ **Rich Span Attributes**: Input/output data, operation metadata in both services
- ✅ **Span Events**: Database operation timing markers
- ✅ **Error Handling**: Proper error recording in traces across service boundaries
- ✅ **Performance Monitoring**: End-to-end processing time visualization
- ✅ **gRPC Instrumentation**: Automatic trace propagation with OpenTelemetry gRPC interceptors
- ✅ **Jaeger Integration**: Beautiful trace visualization UI showing complete request journey

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
   - Select "demo-app" from the service dropdown
   - Click "Find Traces"
   - Explore the complete trace hierarchy across both services

### Test Different Scenarios

```bash
# Normal request (flows through both services)
curl "http://localhost:4000/ping?message=test"

# Trigger an error in service-b (to see error tracing across services)
curl "http://localhost:4000/ping?message=error"

# Test with special characters
curl "http://localhost:4000/ping?message=special%20test"
```

## OpenTelemetry Implementation

### Dependencies Added

Added to both `service-a/go.mod` and `service-b/go.mod`:

- `go.opentelemetry.io/otel` - Core OpenTelemetry library
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` - OTLP gRPC exporter
- `go.opentelemetry.io/otel/sdk` - OpenTelemetry SDK
- `go.opentelemetry.io/otel/trace` - Tracing interface
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` - gRPC tracing

### Cross-Service Tracing Components

#### 1. Tracer Initialization (Both services)

- `InitTracer()` - Sets up OpenTelemetry with OTLP exporter
- `GetTracer()` - Provides tracer instances for span creation
- **Same service name** (`demo-app`) for unified traces
- Configured to send traces to OpenTelemetry Collector

#### 2. service-a (HTTP → gRPC Client)

**API Layer** (`service-a/api/ping.go`):

- **Root span creation** for incoming HTTP requests
- **Request attributes**: endpoint, method, query parameters
- **Context propagation** to service layer

**Service Layer** (`service-a/service/ping.go`):

- **Child span creation** for business logic
- **Operation metadata**: input/output parameters
- **Context forwarding** to gRPC adapter

**gRPC Adapter** (`service-a/adapter/service_b_adapter/ping.go`):

- **gRPC client spans** for outbound calls
- **Context propagation** via gRPC metadata
- **Automatic trace headers** injection

#### 3. service-b (gRPC Server)

**API Layer** (`service-b/api/ping.go`):

- **gRPC span creation** for incoming requests
- **Automatic context extraction** from gRPC metadata
- **Request attributes**: gRPC method, input data

**Service Layer** (`service-b/service/ping.go`):

- **Child span creation** for business logic
- **Operation metadata**: input/output parameters
- **Context forwarding** to store layer

**Store Layer** (`service-b/store/ping.go`):

- **Data access spans** for storage operations
- **Database simulation** (500ms delay)
- **Span events**: query start/end timestamps
- **Operation metrics**: duration, query type

### gRPC Context Propagation Configuration

#### Client Configuration (service-a)

```go
grpc.WithStatsHandler(otelgrpc.NewClientHandler(
    otelgrpc.WithPropagators(propagation.TraceContext{}),
))
```

#### Server Configuration (service-b)

```go
grpc.StatsHandler(otelgrpc.NewServerHandler(
    otelgrpc.WithPropagators(propagation.TraceContext{}),
))
```

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
- **service-b**: 50051 (gRPC API)
- **OpenTelemetry Collector**: 4317 (OTLP gRPC), 4318 (OTLP HTTP)
- **Jaeger UI**: 16686 (Web interface)
- **Jaeger Zipkin**: 9411 (Trace ingestion)

## Understanding Traces

### Complete Trace Hierarchy

When you make a request, you'll see **one unified trace** with spans across both services:

#### service-a spans:

1. **`api.Api.Ping`** - HTTP request handling (API layer)

   - Duration: ~2 seconds total (includes service-b call)
   - Attributes: HTTP method, endpoint, input message
   - Status: Success/Error indication

2. **`service.Service.Ping`** - Business logic processing (Service layer)

   - Duration: ~1.5 seconds (includes gRPC call)
   - Attributes: Operation type, input/output data
   - Parent: api.Api.Ping

3. **`service_b_adapter.Adapter.Ping`** - gRPC client call (Adapter layer)
   - Duration: ~1 second (gRPC call time)
   - Attributes: gRPC endpoint, request/response data
   - Parent: service.Service.Ping

#### service-b spans:

4. **`api.Api.Ping`** - gRPC request handling (service-b API layer)

   - Duration: ~750ms (service + store processing)
   - Attributes: gRPC method, input message
   - Parent: service_b_adapter.Adapter.Ping

5. **`service.Service.Ping`** - Business logic processing (service-b Service layer)

   - Duration: ~500ms (includes store call)
   - Attributes: Operation type, input/output data
   - Parent: service-b api.Api.Ping

6. **`store.Store.Ping`** - Data access simulation (service-b Store layer)
   - Duration: ~500ms (simulated database query)
   - Attributes: Query type, database operation details
   - Events: Query start/end timestamps
   - Parent: service-b service.Service.Ping

### Trace Analysis Benefits

- **End-to-End Performance**: See complete request journey across services
- **Cross-Service Error Tracking**: Identify which service and layer errors occur in
- **Request Correlation**: Follow a single request through multiple services
- **Service Dependency Mapping**: Understand inter-service communication patterns
- **gRPC Performance Monitoring**: Monitor network latency between services
- **Resource Utilization**: Monitor processing time across all components

## Development

### Project Structure

```
.
├── docker-compose.yml           # Infrastructure orchestration
├── otel-collector-config.yaml   # OpenTelemetry collector configuration
├── README.md                    # This file
├── service-a/                   # HTTP microservice
│   ├── api/                     # HTTP handlers with tracing
│   ├── service/                 # Business logic with tracing
│   ├── adapter/                 # gRPC client adapter with tracing
│   ├── util/tracing/            # OpenTelemetry utilities
│   ├── cmd/                     # Application entry points
│   ├── config.json              # Service configuration
│   ├── Dockerfile               # Container build instructions
│   └── go.mod                   # Go dependencies
└── service-b/                   # gRPC microservice
    ├── api/                     # gRPC handlers with tracing
    ├── service/                 # Business logic with tracing
    ├── store/                   # Data layer with tracing
    ├── util/tracing/            # OpenTelemetry utilities
    ├── cmd/                     # Application entry points
    ├── config.json              # Service configuration
    ├── Dockerfile               # Container build instructions
    └── go.mod                   # Go dependencies
```

### Building Services

```bash
# Build service-a
cd service-a
go build -o service-a ./cmd/
./service-a start

# Build service-b
cd service-b
go build -o service-b ./cmd/
./service-b start
```

### Running Tests

```bash
# Test normal cross-service flow
curl "http://localhost:4000/ping?message=test"

# Test error handling across services
curl "http://localhost:4000/ping?message=error"
```

## Troubleshooting

### Common Issues

1. **"No traces found in Jaeger"**

   - Ensure all services are running: `docker compose ps`
   - Check service-a logs: `docker compose logs service-a`
   - Check service-b logs: `docker compose logs service-b`
   - Verify OpenTelemetry collector logs: `docker compose logs otel-collector`

2. **"Connection refused" errors**

   - Restart services: `docker compose down && docker compose up -d`
   - Check network connectivity between containers
   - Verify service-b is listening on port 50051

3. **"demo-app not in dropdown"**

   - Make a request first to generate traces
   - Wait 10-15 seconds for trace processing
   - Refresh Jaeger UI

4. **"Separate traces instead of one unified trace"**
   - Ensure both services use the same tracer name in config.json
   - Verify gRPC propagators are configured on both client and server
   - Check that context is properly passed through all layers

### Logs

```bash
# View all logs
docker compose logs

# View specific service logs
docker compose logs service-a
docker compose logs service-b
docker compose logs otel-collector
docker compose logs jaeger

# Follow logs in real-time
docker compose logs -f service-a service-b
```

### Context Propagation Debugging

If traces aren't connecting across services:

1. **Check propagator configuration** in both services
2. **Verify same service name** in both config.json files
3. **Confirm context flow** through all layers using log trace IDs
4. **Test gRPC connectivity** directly between services

## Learning Resources

This demo helps you understand:

- **OpenTelemetry fundamentals**: Traces, spans, attributes, events
- **Cross-service tracing**: Context propagation via gRPC
- **Distributed tracing patterns**: Span hierarchy across microservices
- **Observability best practices**: Error handling, performance monitoring
- **Infrastructure setup**: Collectors, exporters, backends
- **gRPC instrumentation**: Automatic trace propagation configuration

## License

This project is for educational purposes. Feel free to use it as a reference for implementing OpenTelemetry in your own microservice applications.

---

_Built with ❤️ to help developers understand distributed tracing with OpenTelemetry across microservices_

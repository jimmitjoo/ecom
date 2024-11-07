# E-commerce Product API

A robust and scalable API for product management in e-commerce systems, built with Go following clean architecture principles.

## Technical Highlights

### Performance Optimizations
- Lock-free algorithms for high-throughput operations
- Connection pooling with configurable limits
- Efficient memory management with object pooling
- Optimized concurrent batch processing

### Reliability Features
- Circuit breakers for external dependencies
- Graceful degradation strategies
- Health check endpoints with detailed diagnostics
- Automated recovery mechanisms

### Monitoring & Observability
- Prometheus metrics for real-time monitoring
- Structured logging with correlation IDs
- Distributed tracing with OpenTelemetry
- Performance analytics dashboard

### Security
- Input validation with custom rules
- Advanced rate limiting with sliding window algorithm
  - Per-IP tracking
  - Configurable limits per endpoint
  - Automatic recovery
- CORS with configurable origins
- Secure WebSocket connections

## Project Structure

```
src/
├── domain/           # Business logic and domain models
│   ├── models/
│   ├── repositories/
│   └── events/
├── application/      # Application logic
│   ├── services/
│   └── interfaces/
├── infrastructure/   # External implementations
│   ├── handlers/
│   ├── repositories/
│   └── events/
└── main.go
```

## Prerequisites

- Go 1.20 or later
- gorilla/mux for routing
- gorilla/websocket for WebSocket support
- validator/v10 for data validation

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/jimmitjoo/ecom.git
cd ecom
```

2. Install dependencies:
```bash
go mod download
```

3. Install Air for hot-reloading (optional):
```bash
go install github.com/air-verse/air@latest
```

4. Start the server:
```bash
# Without hot-reloading
go run src/main.go

# With hot-reloading
air
```

The server will start on `http://localhost:8080`

## API Documentation

API documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`
- OpenAPI YAML: `http://localhost:8080/swagger/doc.yaml`

The documentation is automatically generated from code comments. To update:


```bash
swag init -g src/main.go
```

This will:
- Automatically generate OpenAPI/Swagger documentation from code comments
- Provide an interactive Swagger UI to test the API
- Document both REST endpoints and WebSocket connections
- Automatically update when you run `swag init`

## Development

### Architecture

The project follows clean architecture principles:

1. **Domain Layer** (innermost)
   - Contains business logic and domain models
   - No external dependencies

2. **Application Layer**
   - Contains application-specific logic
   - Orchestrates domain objects
   - Implements use cases

3. **Infrastructure Layer** (outermost)
   - Contains all external implementations
   - Database, HTTP handlers, WebSocket etc.

### Technical Implementation

#### Concurrent Processing
- Goroutine pools for batch operations
- Mutex-protected shared resources
- Channel-based event processing
- Lock-free data structures where possible

#### Event Sourcing
- Versioned events with hash chains
- Optimistic concurrency control
- Event replay capabilities
- Snapshot support

#### State Management
- In-memory state with versioning
- Event-based state reconstruction
- Conflict resolution strategies
- Consistency guarantees

#### Monitoring & Observability

##### Metrics (Prometheus)
- Request latency histograms
- Operation counters by type
- WebSocket connection metrics
- Event processing metrics
- Repository operation latency

##### Structured Logging (Zap)
- Request/response logging
- Error tracking
- Performance metrics
- Audit trail

##### Development Mode Logging
I utvecklingsläge aktiveras detaljerad loggning automatiskt med:
- Färgkodade loggnivåer för bättre läsbarhet
- Automatisk stack trace vid fel
- Request ID för varje anrop
- Detaljerad timing för alla operationer
- Källkodsinformation (fil och rad)
- Utökad kontextuell information

För att aktivera development logging:
```bash
export GO_ENV=development
```

##### Distributed Tracing (OpenTelemetry)
- End-to-end request tracing
- Cross-service correlation
- Performance bottleneck analysis
- Error propagation tracking

### Development Tools

#### Test Data Generator
The project includes a built-in test data generator accessible through the WebSocket test interface (`test/websocket.html`). This tool provides:

- Quick generation of test products with realistic data
- Support for batch generation (1-50 products at once)
- Multi-market and multi-currency test data
- Randomized but valid product attributes

To use the generator:
1. Start the server
2. Open `test/websocket.html` in your browser
3. Use the generator controls to create test data:
   - Generate single products
   - Generate batch products (5, 20, or random amount)
   - Automatic WebSocket updates for generated products

The generator automatically creates products with:
- Unique SKU numbers
- Randomized prices in multiple currencies (SEK, NOK, DKK, EUR)
- Market-specific metadata (SE, NO, DK, FI)
- Realistic titles and descriptions

#### Hot Reloading
The project supports hot-reloading using Air, which automatically rebuilds and restarts the application when file changes are detected. This significantly improves the development experience.

To use hot-reloading:
1. Install Air: `go install github.com/air-verse/air@latest`
2. Run the server with Air: `air`

Air will monitor your source files and automatically rebuild when changes are detected.

#### Performance Profiling
In development mode, the following profiling endpoints are available:

- `/debug/pprof/cpu` - CPU profiling data
- `/debug/pprof/heap` - Heap memory profiling
- `/debug/pprof/goroutine` - Goroutine profiling

Use with Go's built-in tools:
```bash
# CPU Profiling
go tool pprof http://localhost:8080/debug/pprof/cpu

# Heap Profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine Profiling
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

For visualization, use the `-http` flag:
```bash
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap
```

This provides:
1. Built-in profiling with Go's pprof
2. Secure access (development mode only)
3. Real-time analysis of:
   - CPU usage
   - Memory allocation
   - Goroutine management
4. Visual profiling data tools

### Testing

1. Unit Tests
```bash
go test ./...
```

2. Integration Testing with WebSocket Interface
```bash
open test/websocket.html
```
Features:
- Real-time WebSocket connection testing
- Built-in test data generator
- Visual feedback for all operations
- Automatic updates for product changes

3. Performance Testing
```bash
go test -bench=.
```

4. Load Testing & Benchmarking
```bash
# Run all benchmark tests
go test -bench=. ./src/testing/load/...

# For detailed output with memory information
go test -bench=. -benchmem ./src/testing/load/...

# To run a specific benchmark
go test -bench=BenchmarkBatchOperations ./src/testing/load/...
```

The benchmark tests include:
- Batch operations with various sizes (10, 100, 1000 products)
- Sequential and parallel operations
- Performance measurements for:
  - Small batches (10 products)
  - Medium batches (100 products)
  - Large batches (1000 products)
- Automatic measurement of:
  - Operations per second
  - Memory allocation
  - Allocation frequency

5. Chaos Testing
```bash
go test ./src/testing/chaos -v -timeout 10m
```

The project includes comprehensive chaos testing to ensure system resilience:

#### Network Chaos
- Simulates network latency (100ms-500ms)
- Packet loss simulation (10%-90% loss rate)
- Connection timeout scenarios
- Custom transport layer for HTTP chaos

#### Memory Pressure Testing
- Simulates high memory usage (up to 80% system memory)
- Concurrent batch operations under memory pressure
- Automatic memory cleanup and GC triggering
- Resource exhaustion scenarios

#### Data Corruption Testing
- Simulates corrupted JSON payloads
- Configurable corruption rates (0-100%)
- Validates system handling of malformed data
- Tests data integrity checks and error handling

Example chaos test configuration:
```go
type ChaosConfig struct {
    NetworkLatency  time.Duration
    PacketLossRate  float64
    MemoryPressure  bool
    CorruptDataRate float64
}
```

The chaos testing suite helps ensure:
- System resilience under network stress
- Proper handling of corrupted data
- Graceful degradation under memory pressure
- Recovery from various failure scenarios

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details. 
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
- Rate limiting per client
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

3. Start the server:
```bash
go run src/main.go
```

The server will start on `http://localhost:8080`

## Documentation

For detailed API documentation, including endpoints, request/response formats, and examples, see [API Documentation](docs/api.md).

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

## Monitoring & Observability

### Metrics (Prometheus)
- Request latency histograms
- Operation counters by type
- WebSocket connection metrics
- Event processing metrics
- Repository operation latency

### Structured Logging (Zap)
- Request/response logging
- Error tracking
- Performance metrics
- Audit trail

### Distributed Tracing (OpenTelemetry)
- End-to-end request tracing
- Cross-service correlation
- Performance bottleneck analysis
- Error propagation tracking
### Testing

1. Unit Tests
```bash
go test ./...
```

2. WebSocket Testing Interface
```bash
open test/websocket.html
```

3. Performance Testing
```bash
go test -bench=.
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details. 
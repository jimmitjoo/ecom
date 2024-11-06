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

API dokumentationen finns tillgänglig på:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`
- OpenAPI YAML: `http://localhost:8080/swagger/doc.yaml`

Dokumentationen genereras automatiskt från kodens kommentarer. För att uppdatera:

```bash
swag init -g src/main.go
```

Detta kommer att:
- Automatiskt generera OpenAPI/Swagger dokumentation från dina kodkommentarer
- Ge en interaktiv Swagger UI för att testa API:et
- Dokumentera både REST endpoints och WebSocket-anslutningar
- Uppdateras automatiskt när du kör `swag init`


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

### Development Tools

#### Hot Reloading
The project supports hot-reloading using Air, which automatically rebuilds and restarts the application when file changes are detected. This significantly improves the development experience.

To use hot-reloading:
1. Install Air: `go install github.com/air-verse/air@latest`
2. Run the server with Air: `air`

Air will monitor your source files and automatically rebuild when changes are detected.

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

#### Development Mode Logging
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
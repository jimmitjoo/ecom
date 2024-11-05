# E-commerce Product API

A robust and scalable API for product management in e-commerce systems, built with Go following clean architecture principles.

## Use Cases & Benefits

### Multi-Market Product Management
Perfect for businesses operating across multiple markets where:
- Products need different pricing per market
- Product information requires localization
- Stock levels vary by location
- Variants need market-specific attributes

### Real-time Inventory Management
Ideal for businesses requiring immediate stock updates:
- Live inventory tracking across warehouses
- Instant stock level notifications
- Real-time availability updates
- Automated stock threshold alerts

### Batch Operations for Scale
Designed for operations requiring bulk updates:
- Mass product imports
- Seasonal price updates
- Inventory reconciliation
- Bulk product modifications

### Event-Driven Architecture
Enables building reactive systems with:
- Real-time product updates
- Event-based integrations
- Audit trail capabilities
- Automated workflows

### Why Choose This Solution?

1. **Performance & Scalability**
   - Concurrent processing of batch operations
     * Configurable batch sizes (default max: 1000 items)
     * Parallel processing with worker pools
     * Progress tracking for large operations
   
   - Efficient in-memory data storage
     * Sub-millisecond read operations
     * O(1) lookup complexity
     * Easy migration path to persistent storage
   
   - Non-blocking event publishing
     * Async event emission
     * Buffered channels for event handling
     * Configurable event buffer sizes
   
   - Thread-safe implementations
     * Mutex-protected shared resources
     * Atomic operations where applicable
     * Lock-free algorithms for high concurrency

2. **Reliability & Consistency**
   - Guaranteed event ordering
   - Idempotent operations
   - Automatic retry mechanisms
   - Transaction-like batch operations

3. **Flexibility & Extensibility**
   - Clean architecture for easy modifications
   - Pluggable storage implementations
   - Event-driven for loose coupling
   - Modular design

4. **Developer Experience**
   - Clear API documentation
   - Consistent error handling
   - WebSocket testing interface
   - Comprehensive examples

## Features

- RESTful API for CRUD operations on products
- Batch operations for efficient bulk updates
- Real-time updates via WebSocket
- Product data validation
- In-memory data storage (easily replaceable with other data sources)
- CORS support
- Structured error handling

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

## Installation

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

## API Usage

### Single Operations

#### REST Endpoints
- `GET /products` - List all products
- `POST /products` - Create new product
- `GET /products/{id}` - Get specific product
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Delete product

#### Example Single Product Creation
```bash
curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{
        "sku": "TSHIRT-001",
        "base_title": "Basic T-shirt",
        "description": "Classic t-shirt in 100% cotton",
        "prices": [
            {
                "currency": "SEK",
                "amount": 299.00
            }
        ],
        "metadata": [
            {
                "market": "SE",
                "title": "Basic T-shirt",
                "description": "Classic t-shirt in 100% cotton"
            }
        ]
    }'
```

### Batch Operations

#### REST Endpoints
- `POST /products/batch` - Create multiple products
- `PUT /products/batch` - Update multiple products
- `DELETE /products/batch` - Delete multiple products

#### Example Batch Creation
```bash
curl -X POST http://localhost:8080/products/batch \
    -H "Content-Type: application/json" \
    -d '[
        {
            "sku": "TSHIRT-001",
            "base_title": "Basic T-shirt",
            "prices": [{
                "currency": "SEK",
                "amount": 299.00
            }],
            "metadata": [{
                "market": "SE",
                "title": "Basic T-shirt"
            }]
        },
        {
            "sku": "TSHIRT-002",
            "base_title": "Premium T-shirt",
            "prices": [{
                "currency": "SEK",
                "amount": 399.00
            }],
            "metadata": [{
                "market": "SE",
                "title": "Premium T-shirt"
            }]
        }
    ]'
```

#### Example Batch Deletion
```bash
curl -X DELETE http://localhost:8080/products/batch \
    -H "Content-Type: application/json" \
    -d '["prod_123", "prod_456"]'
```

#### Batch Operation Response Format
```json
[
    {
        "success": true,
        "id": "prod_123",
        "error": ""
    },
    {
        "success": false,
        "id": "prod_456",
        "error": "Product not found"
    }
]
```

### Real-time Updates (WebSocket)

- `ws://localhost:8080/ws` - WebSocket endpoint for real-time updates
- Receives events for all product operations (single and batch)
- Event types: product.created, product.updated, product.deleted

## Testing

### WebSocket Testing Interface
1. Open `test/websocket.html` in a browser
2. Use the interface to:
   - Create individual or batch products
   - View and select products
   - Delete individual or multiple products
   - See real-time updates for all operations

## Architecture

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

## Development

### Adding New Features

1. Start in the domain layer if new models are needed
2. Update or create new interfaces in the application layer
3. Implement in the infrastructure layer

### Code Standards

- Follow Go's official code standards
- Use gofmt for formatting
- Document public functions and types

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details. 
# E-commerce Product API

A robust and scalable API for product management in e-commerce systems, built with Go following clean architecture principles.

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
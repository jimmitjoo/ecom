# E-commerce Product API

A robust and scalable API for product management in e-commerce systems, built with Go following clean architecture principles.

## Features

- RESTful API for CRUD operations on products
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

### REST Endpoints

- `GET /products` - List all products
- `POST /products` - Create new product
- `GET /products/{id}` - Get specific product
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Delete product

### WebSocket

- `ws://localhost:8080/ws` - WebSocket endpoint for real-time updates

### Example Product Creation

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

## Testing WebSocket

1. Open `test/websocket.html` in a browser
2. Use the interface to create, update and delete products
3. See real-time updates in the browser

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
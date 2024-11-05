# Product API Documentation

## Overview
This API handles products for the e-commerce system. It supports CRUD operations (Create, Read, Update, Delete) for both individual products and batch operations, with real-time updates via WebSocket.

**Base URL:** `http://localhost:8080`

## Authentication
Currently, the API does not require authentication.

## Data Models

### Product Model
```json
{
    "id": "prod_123",
    "sku": "TSHIRT-001",
    "base_title": "Basic T-shirt",
    "description": "Classic t-shirt in 100% cotton",
    "prices": [
        {
            "currency": "SEK",
            "amount": 299.00
        }
    ],
    "variants": [
        {
            "id": "var_123",
            "sku": "TSHIRT-001-BLK-L",
            "attributes": {
                "color": "black",
                "size": "L"
            },
            "stock": [
                {
                    "location_id": "STH1",
                    "quantity": 100
                }
            ]
        }
    ],
    "metadata": [
        {
            "market": "SE",
            "title": "Basic T-shirt",
            "description": "Classic t-shirt in 100% cotton",
            "keywords": "t-shirt, basic, cotton"
        }
    ],
    "created_at": "2024-02-20T12:00:00Z",
    "updated_at": "2024-02-20T12:00:00Z"
}
```

### Validation Rules

#### Product
- `sku`: Required
- `base_title`: Required
- `prices`: At least one price must be specified

#### Price
- `currency`: Exactly 3 characters (e.g., "SEK", "USD")
- `amount`: Must be greater than or equal to 0

#### Variant
- `id`: Unique ID for the variant
- `sku`: Required
- `attributes`: At least one attribute must be specified

#### Stock
- `location_id`: Required
- `quantity`: Must be greater than or equal to 0

#### Metadata
- `market`: Required
- `title`: Required

## API Endpoints

### Single Product Operations

#### List Products
**GET** `/products`

Returns a list of all products.

**Response:** 200 OK
```json
[
    // Array of product objects
]
```

#### Get Product
**GET** `/products/{id}`

Returns a specific product by ID.

**Response:** 200 OK or 404 Not Found

#### Create Product
**POST** `/products`

Creates a new product.

**Request Body:** Product object (without id)
**Response:** 201 Created or 400 Bad Request

#### Update Product
**PUT** `/products/{id}`

Updates an existing product.

**Request Body:** Product object
**Response:** 200 OK, 400 Bad Request, or 404 Not Found

#### Delete Product
**DELETE** `/products/{id}`

Deletes a product.

**Response:** 204 No Content or 404 Not Found

### Batch Operations

#### Create Multiple Products
**POST** `/products/batch`

Creates multiple products in a single request.

**Request Body:**
```json
[
    // Array of product objects
]
```

**Response:** 201 Created
```json
[
    {
        "success": true,
        "id": "prod_abc123",
        "error": ""
    }
]
```

#### Update Multiple Products
**PUT** `/products/batch`

Updates multiple products in a single request.

#### Delete Multiple Products
**DELETE** `/products/batch`

Deletes multiple products in a single request.

**Request Body:**
```json
["prod_123", "prod_456"]
```

### Real-time Updates (WebSocket)

#### Connection Details
- Endpoint: `ws://localhost:8080/ws`
- Protocol: WebSocket (RFC 6455)
- Supported Operations: Text messages (JSON)

#### Reconnection Strategy
The client implements an exponential backoff strategy for reconnections:
- Initial retry delay: 2 seconds
- Maximum retries: 5
- Backoff multiplier: 2 (2s, 4s, 8s, 16s, 32s)
- Automatic reconnection when:
  - Connection is lost
  - Browser tab becomes active again
  - Network connectivity is restored

Example implementation:
```javascript
const MAX_RETRIES = 5;
const RETRY_DELAY_MS = 2000;

function connectWebSocket() {
    if (retryCount >= MAX_RETRIES) {
        console.error('Max retry attempts reached');
        return;
    }

    const delay = RETRY_DELAY_MS * Math.pow(2, retryCount);
    setTimeout(() => {
        // Attempt reconnection
        retryCount++;
    }, delay);
}
```

#### Handling Missed Events
The system handles missed events through several mechanisms:

1. **State Synchronization**
   - Client can request full state by fetching all products via REST API
   - Use GET /products endpoint after reconnection
   - Compare local state with server state

2. **Event Ordering**
   - Events include timestamps for ordering
   - Each event has a unique ID for deduplication
   - Events for the same product are processed in timestamp order

Example event with ordering information:
```json
{
    "id": "evt_123",
    "type": "product.updated",
    "data": {
        "product_id": "prod_123",
        "action": "updated",
        "product": {
            // Product data
        }
    },
    "timestamp": "2024-02-20T12:00:00Z"
}
```

#### Event Idempotency
Events are designed to be idempotent to handle duplicate processing:

1. **Event IDs**
   - Each event has a unique ID
   - Clients can track processed event IDs
   - Duplicate events can be safely ignored

2. **State-Based Updates**
   - Events contain full resource state
   - Processing same event multiple times is safe
   - Final state remains consistent

Example client-side idempotency handling:
```javascript
const processedEvents = new Set();

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    // Check if event was already processed
    if (processedEvents.has(data.id)) {
        console.log('Duplicate event ignored:', data.id);
        return;
    }
    
    // Process event
    handleEvent(data);
    
    // Mark event as processed
    processedEvents.add(data.id);
    
    // Cleanup old events (optional)
    if (processedEvents.size > 1000) {
        cleanupOldEvents();
    }
};
```

#### Best Practices
1. **Connection Management**
   - Implement heartbeat mechanism
   - Monitor connection state
   - Log reconnection attempts
   - Clear resources on disconnect

2. **Event Processing**
   - Maintain event order
   - Handle duplicates
   - Validate event data
   - Implement error handling

3. **State Management**
   - Keep local state
   - Implement state recovery
   - Handle partial updates
   - Validate state consistency

4. **Error Handling**
   - Log connection errors
   - Handle message parsing errors
   - Implement fallback mechanisms
   - Notify users of connection status

### Event Types and Format

The API uses the following event types for real-time updates:

#### Event Types
- `product.created` - Emitted when a new product is created
- `product.updated` - Emitted when an existing product is updated
- `product.deleted` - Emitted when a product is deleted

#### Event Structure
```json
{
    "id": "evt_123abc",        // Unique event ID
    "type": "product.created", // Event type
    "data": {
        "product_id": "prod_123", // ID of affected product
        "action": "created",      // Action performed
        "product": {              // Full product data
            "id": "prod_123",
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
            ],
            "created_at": "2024-02-20T12:00:00Z",
            "updated_at": "2024-02-20T12:00:00Z"
        }
    },
    "timestamp": "2024-02-20T12:00:00Z" // When the event occurred
}
```

#### Event Handling Notes
- All events include the full product data (except delete events)
- Events are ordered by timestamp
- Each event has a unique ID for deduplication
- Delete events include the product data before deletion
- Batch operations generate individual events for each affected product

## Error Handling

### Error Response Format
```json
{
    "error": "Descriptive error message"
}
```

### HTTP Status Codes
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Resource deleted
- `400 Bad Request` - Invalid request or validation error
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Code Examples

### cURL Examples

#### List Products
```bash
curl -X GET http://localhost:8080/products
```

#### Create Product
```bash
curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{
        "sku": "TSHIRT-001",
        "base_title": "Basic T-shirt",
        "prices": [{
            "currency": "SEK",
            "amount": 299.00
        }]
    }'
```

### JavaScript WebSocket Example
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received event:', data);
    
    switch(data.type) {
        case 'product.created':
            console.log('New product:', data.data.product);
            break;
        case 'product.updated':
            console.log('Updated product:', data.data.product);
            break;
        case 'product.deleted':
            console.log('Deleted product:', data.data.product_id);
            break;
    }
};
```

## Rate Limiting
Currently, no rate limiting is implemented.
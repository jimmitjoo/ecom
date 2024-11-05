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

#### WebSocket Connection
**WS** `/ws`

Establishes a WebSocket connection for real-time updates.

#### Event Types
- `product.created` - New product created
- `product.updated` - Product updated
- `product.deleted` - Product deleted

#### Event Format
```json
{
    "id": "event_123",
    "type": "product.created",
    "data": {
        "product_id": "prod_123",
        "action": "created",
        "product": {
            // Product data
        }
    },
    "timestamp": "2024-02-20T12:00:00Z"
}
```

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
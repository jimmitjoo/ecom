# Product API Documentation

## Overview
This API handles products for the e-commerce system. It supports basic CRUD operations (Create, Read, Update, Delete) for product management.

**Base URL:** `http://localhost:8080`

## Endpoints

### 1. List all products
**GET** `/products`

Retrieves a list of all products in the system.

**Response Code:** 200 OK

**Example Response:**
```json
[
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
]
```

### 2. Get specific product
**GET** `/products/{id}`

Retrieves detailed information about a specific product.

**Parameters:**
- `id` (path parameter) - Product's unique ID

**Response Codes:**
- 200 OK - Product found
- 404 Not Found - Product does not exist

**Example Response (200 OK):**
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
    ]
}
```

**Example Error Response (404 Not Found):**
```json
{
    "error": "Product with ID 'prod_123' not found"
}
```

### 3. Create new product
**POST** `/products`

Creates a new product in the system.

**Headers:**
- `Content-Type: application/json`

**Example Request:**
```json
{
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
            "description": "Classic t-shirt in 100% cotton",
            "keywords": "t-shirt, basic, cotton"
        }
    ]
}
```

**Response Codes:**
- 201 Created - Product created successfully
- 400 Bad Request - Validation error

### 4. Update product
**PUT** `/products/{id}`

Updates an existing product.

**Headers:**
- Content-Type: application/json

**Parameters:**
- `id` (path parameter) - Product's unique ID

**Request:** Same format as product creation

**Response Codes:**
- 200 OK - Product updated successfully
- 400 Bad Request - Validation error
- 404 Not Found - Product does not exist

### 5. Delete product
**DELETE** `/products/{id}`

Removes a product from the system.

**Parameters:**
- `id` (path parameter) - Product's unique ID

**Response Codes:**
- 204 No Content - Product deleted successfully
- 404 Not Found - Product does not exist

## Real-time Updates via WebSocket

The API supports real-time updates via WebSocket for receiving immediate updates when products are created, updated, or deleted.

### WebSocket Endpoint
**WS** `/ws`

### Event Format
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

### Event Types
- `product.created` - When a new product is created
- `product.updated` - When a product is updated
- `product.deleted` - When a product is deleted

### JavaScript Example
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received event:', data);
    
    switch(data.type) {
        case 'product.created':
            console.log('New product created:', data.data.product);
            break;
        case 'product.updated':
            console.log('Product updated:', data.data.product);
            break;
        case 'product.deleted':
            console.log('Product deleted:', data.data.product_id);
            break;
    }
};

ws.onerror = function(error) {
    console.error('WebSocket error:', error);
};

ws.onclose = function() {
    console.log('WebSocket connection closed');
};
```

## Validation Rules

### Product
- `sku`: Required
- `base_title`: Required
- `prices`: At least one price must be specified

### Price
- `currency`: Exactly 3 characters (e.g., "SEK", "USD")
- `amount`: Must be greater than or equal to 0

### Variant
- `id`: Unique ID for the variant
- `sku`: Required
- `attributes`: At least one attribute must be specified

### Stock
- `location_id`: Required
- `quantity`: Must be greater than or equal to 0

### Metadata
- `market`: Required
- `title`: Required

## Error Handling

The API returns structured error messages in the following format:
```json
{
    "error": "Descriptive error message"
}
```

### Common Error Codes
- `400 Bad Request` - Invalid request or validation error
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Examples with cURL

### List all products
```bash
curl -X GET http://localhost:8080/products
```

### Get specific product
```bash
curl -X GET http://localhost:8080/products/prod_123
```

### Create product
```bash
curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{
        "sku": "TSHIRT-001",
        "base_title": "Basic T-shirt",
        "prices": [
            {
                "currency": "SEK",
                "amount": 299.00
            }
        ]
    }'
```

### Update product
```bash
curl -X PUT http://localhost:8080/products/prod_123 \
    -H "Content-Type: application/json" \
    -d '{
        "sku": "TSHIRT-001",
        "base_title": "Updated T-shirt",
        "prices": [
            {
                "currency": "SEK",
                "amount": 399.00
            }
        ]
    }'
```

### Delete product
```bash
curl -X DELETE http://localhost:8080/products/prod_123
```

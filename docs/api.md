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

## Event Store and Versioning

### Event Structure
Events are stored with versioning information for conflict resolution and replay capabilities:

```json
{
    "id": "evt_123",
    "type": "product.updated",
    "entity_id": "prod_123",
    "version": 2,
    "sequence": 1234,
    "caused_by": "evt_122",
    "timestamp": "2024-02-20T12:00:00Z",
    "data": {
        "product_id": "prod_123",
        "action": "updated",
        "version": 2,
        "prev_hash": "abc123...", // Hash of previous state
        "changes": [
            {
                "field": "base_title",
                "old_value": "Basic T-shirt",
                "new_value": "Premium T-shirt"
            }
        ],
        "product": {
            // Full product state
        }
    }
}
```

### Event Versioning
Each event includes:
- `version`: Incremental version number for the entity
- `sequence`: Global sequence number for total ordering
- `prev_hash`: Hash of previous state for integrity verification
- `caused_by`: ID of the event that triggered this event (for causality tracking)

### Conflict Resolution
The system handles conflicts through:
1. **Optimistic Locking**
   - Version checking before updates
   - Automatic retry on version conflicts
   - Conflict resolution strategies

2. **Event Chain Validation**
   - Version sequence verification
   - Hash chain integrity checks
   - Causal ordering preservation

### Event Replay
Events can be replayed for:
- State reconstruction
- Audit purposes
- System recovery
- Data synchronization

Example replay request:
```bash
curl -X GET "http://localhost:8080/products/{id}/events?from_version=1"
```

Response:
```json
{
    "events": [
        {
            "id": "evt_123",
            "version": 1,
            // ... event data
        },
        {
            "id": "evt_124",
            "version": 2,
            // ... event data
        }
    ]
}
```

### Distributed Locks
The system uses distributed locks for concurrent operations:
- TTL-based locks
- Automatic lock cleanup
- Lock refresh mechanism
- Deadlock prevention

Example lock acquisition:
```go
acquired, err := lockManager.AcquireLock(ctx, "prod_123", 10*time.Second)
if err != nil {
    return err
}
if !acquired {
    return errors.New("resource locked")
}
defer lockManager.ReleaseLock("prod_123")
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

## Testing WebSocket Functionality

### Basic Connection Test
1. Open `test/websocket.html` in your browser
2. Check connection status in top-right corner
3. Open browser's developer tools (F12)
4. Monitor WebSocket connection in Network tab

### Testing Real-time Updates

#### Single Operations Test
1. Create a product:
   ```bash
   curl -X POST http://localhost:8080/products \
       -H "Content-Type: application/json" \
       -d '{
           "sku": "TEST-001",
           "base_title": "Test Product",
           "prices": [{
               "currency": "SEK",
               "amount": 299.00
           }],
           "metadata": [{
               "market": "SE",
               "title": "Test Product"
           }]
       }'
   ```
2. Observe WebSocket event in browser:
   ```json
   {
       "id": "evt_123",
       "type": "product.created",
       "data": {
           "product_id": "prod_abc",
           "action": "created",
           "product": {
               // Full product data
           }
       }
   }
   ```

#### Batch Operations Test
1. Create multiple products:
   ```javascript
   // In browser console
   batchCreateProducts(3);
   ```
2. Observe multiple events:
   - Each product generates separate event
   - Events maintain order
   - All products appear in list

#### Concurrent Operations Test
1. Open multiple browser tabs
2. Perform operations in different tabs
3. Verify all tabs receive updates
4. Check event ordering consistency

### Testing Error Scenarios

#### Connection Loss
1. Start server and connect
2. Stop server (`Ctrl+C`)
3. Observe reconnection attempts
4. Verify exponential backoff
5. Restart server
6. Verify automatic reconnection

#### Event Recovery
1. Create products while offline
2. Reconnect to WebSocket
3. Verify state synchronization:
   ```javascript
   // In browser console
   refreshProducts(); // Manual sync
   ```

#### Duplicate Events
1. Create product with same ID twice
2. Verify idempotent handling:
   ```javascript
   // In browser console
   const processedEvents = new Set();
   // Check processed events
   console.log(processedEvents);
   ```

### Testing Event Versioning

#### Version Conflicts
1. Update same product in two tabs:
   ```javascript
   // Tab 1
   const product = {
       "id": "prod_123",
       "version": 1,
       // ... other fields
   };
   
   // Tab 2 (same product, different version)
   const product = {
       "id": "prod_123",
       "version": 2,
       // ... other fields
   };
   ```
2. Observe conflict resolution

#### Event Chain Verification
1. Get event history:
   ```bash
   curl "http://localhost:8080/products/prod_123/events"
   ```
2. Verify chain integrity:
   ```json
   {
       "events": [
           {
               "id": "evt_1",
               "version": 1,
               "prev_hash": null
           },
           {
               "id": "evt_2",
               "version": 2,
               "prev_hash": "abc123..." // Hash of previous state
           }
       ]
   }
   ```

### Performance Testing

#### High-frequency Updates
1. Create rapid updates:
   ```javascript
   // In browser console
   async function rapidUpdates(count) {
       for (let i = 0; i < count; i++) {
           await batchCreateProducts(10);
           await new Promise(r => setTimeout(r, 100));
       }
   }
   rapidUpdates(5); // 50 products in 500ms
   ```
2. Monitor event processing
3. Check for dropped events

#### Large Batch Operations
1. Test with varying batch sizes:
   ```javascript
   batchCreateProducts(100);
   batchCreateProducts(500);
   batchCreateProducts(1000);
   ```
2. Monitor memory usage
3. Check event delivery timing

### Monitoring Tools

#### Browser DevTools
- Network tab: WebSocket frames
- Console: Event logging
- Memory: Resource usage

#### Server Logs
```bash
# Terminal 1 - Run server with detailed logging
LOG_LEVEL=debug go run src/main.go

# Terminal 2 - Follow logs
tail -f ecom.log
```

#### Metrics to Monitor
1. Connection Status:
   - Connected clients
   - Reconnection attempts
   - Connection duration

2. Event Processing:
   - Events per second
   - Processing latency
   - Queue length

3. Error Rates:
   - Failed deliveries
   - Version conflicts
   - Invalid events

### Common Issues and Solutions

1. **Event Order Inconsistency**
   - Symptom: Events appear out of order
   - Solution: Check sequence numbers
   ```javascript
   let lastSeq = 0;
   ws.onmessage = function(event) {
       const data = JSON.parse(event.data);
       if (data.sequence <= lastSeq) {
           console.warn('Out of order event detected');
       }
       lastSeq = data.sequence;
   };
   ```

2. **Memory Leaks**
   - Symptom: Growing memory usage
   - Solution: Clean up old events
   ```javascript
   // Implement cleanup
   function cleanupOldEvents() {
       const maxEvents = 1000;
       const events = Array.from(processedEvents);
       if (events.length > maxEvents) {
           events.slice(0, events.length - maxEvents)
               .forEach(id => processedEvents.delete(id));
       }
   }
   ```

3. **Lost Updates**
   - Symptom: Missing events after reconnection
   - Solution: Implement event replay
   ```javascript
   ws.onreconnect = async function() {
       const lastEventId = getLastProcessedEventId();
       const missed = await fetchEventsSince(lastEventId);
       processEvents(missed);
   };
   ```

### Best Practices

1. **Event Handling**
   - Always verify event order
   - Implement idempotency checks
   - Handle partial updates
   - Log unprocessable events

2. **Connection Management**
   - Implement heartbeat
   - Use exponential backoff
   - Clean up resources
   - Monitor connection health

3. **State Management**
   - Keep local state
   - Implement recovery
   - Validate consistency
   - Handle conflicts

4. **Error Handling**
   - Log all errors
   - Implement fallbacks
   - Notify users
   - Track error patterns
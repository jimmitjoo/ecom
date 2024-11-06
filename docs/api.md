# Product API Documentation

## Quick Start
1. Start the server: `go run src/main.go`
2. Create a product:
```bash
curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{
        "sku": "TSHIRT-001",
        "base_title": "Basic T-shirt",
        "prices": [{"currency": "SEK", "amount": 299.00}],
        "metadata": [{
            "market": "SE",
            "title": "Basic T-shirt"
        }]
    }'
```
3. Connect to WebSocket: `ws://localhost:8080/ws`

## Core Concepts

### Products
- Unique ID and SKU
- Multi-market support
- Variant handling
- Real-time inventory

### Events
- Versioned events
- Guaranteed ordering
- Idempotent processing
- Event replay

### Batch Operations
- Concurrent processing
- Atomic transactions
- Bulk updates
- Conflict resolution

## API Reference

### REST Endpoints
- `GET /products` - List all products
- `POST /products` - Create product
- `GET /products/{id}` - Get product
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Delete product

### Batch Endpoints
- `POST /products/batch` - Create multiple products
- `PUT /products/batch` - Update multiple products
- `DELETE /products/batch` - Delete multiple products

### WebSocket
- `ws://localhost:8080/ws` - Real-time updates
- Automatic reconnection
- Event deduplication
- State synchronization

## Technical Details

### Event Sourcing
Events are versioned and chained for data integrity:

```json
{
    "id": "evt_123",
    "type": "product.updated",
    "version": 2,
    "entity_id": "prod_123",
    "sequence": 1234,
    "caused_by": "evt_122",
    "timestamp": "2024-02-20T12:00:00Z",
    "data": {
        "product_id": "prod_123",
        "action": "updated",
        "version": 2,
        "prev_hash": "abc123...",
        "changes": [
            {
                "field": "base_title",
                "old_value": "Basic T-shirt",
                "new_value": "Premium T-shirt"
            }
        ]
    }
}
```

### Error Handling

All errors follow a consistent format:
```json
{
    "error": "Descriptive error message"
}
```

Common HTTP status codes:
- `400` - Invalid request data
- `404` - Resource not found
- `409` - Version conflict
- `429` - Rate limit exceeded
- `500` - Internal server error

### Performance Considerations

1. **Batch Operations**
   - Maximum batch size: 1000 items
   - Concurrent processing
   - Atomic transactions
   - Partial success handling

2. **WebSocket Connections**
   - Auto-reconnection with exponential backoff
   - Connection pooling
   - Message buffering
   - Heart-beat mechanism

3. **Event Processing**
   - Asynchronous event publishing
   - Event deduplication
   - Ordered delivery
   - Back-pressure handling

### Monitoring

1. **Metrics Available**
   ```
   # Request latency
   http_request_duration_seconds{handler="/products",method="POST"}
   
   # Active WebSocket connections
   websocket_connections_active
   
   # Event processing time
   event_processing_duration_seconds{event_type="product.created"}
   
   # Batch operation size
   batch_operation_size_bucket{le="100"}
   
   # Rate limiting
   rate_limit_exceeded_total{endpoint="/products"}
   rate_limit_remaining{ip="192.168.1.1"}
   ```

2. **Logging**
   ```json
   {
       "level": "info",
       "timestamp": "2024-02-20T12:00:00Z",
       "caller": "handlers/product_handler.go:42",
       "msg": "Product created",
       "product_id": "prod_123",
       "duration_ms": 45,
       "trace_id": "abc123"
   }
   ```

3. **Tracing**
   ```
   Trace ID: abc123
   ├── POST /products/batch
   │   ├── Validate products (10ms)
   │   ├── Process batch (150ms)
   │   │   ├── Create product 1 (45ms)
   │   │   ├── Create product 2 (48ms)
   │   │   └── Create product 3 (42ms)
   │   └── Publish events (25ms)
   ```

### Best Practices

1. **Error Handling**
   - Always validate input data
   - Return appropriate HTTP status codes
   - Include detailed error messages
   - Log errors with context

2. **Performance**
   - Use batch operations for bulk updates
   - Implement pagination for large datasets
   - Cache frequently accessed data
   - Monitor resource usage

3. **Security**
   - Validate all input
   - Use HTTPS in production
   - Implement rate limiting
   - Follow security headers best practices

4. **Monitoring**
   - Set up alerts for key metrics
   - Monitor error rates
   - Track response times
   - Watch resource usage

### Common Use Cases

1. **Product Import**
```bash
# Import products from CSV
curl -X POST http://localhost:8080/products/batch \
    -H "Content-Type: application/json" \
    -d @products.json
```

2. **Price Updates**
```bash
# Update prices for multiple products
curl -X PUT http://localhost:8080/products/batch \
    -H "Content-Type: application/json" \
    -d '[
        {
            "id": "prod_123",
            "prices": [{"currency": "SEK", "amount": 399.00}]
        }
    ]'
```

3. **Inventory Sync**
```javascript
// Connect to WebSocket for real-time updates
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'product.updated') {
        updateInventory(data.data.product);
    }
};
```

4. **Market-specific Updates**
```bash
# Update product metadata for specific market
curl -X PUT http://localhost:8080/products/prod_123 \
    -H "Content-Type: application/json" \
    -d '{
        "metadata": [{
            "market": "SE",
            "title": "Updated Title",
            "description": "New description"
        }]
    }'
```

### Troubleshooting

1. **Version Conflicts**
   - Get current version: `GET /products/{id}`
   - Fetch event history: `GET /products/{id}/events`
   - Resolve conflicts manually
   - Retry with correct version

2. **Connection Issues**
   - Check server status
   - Verify network connectivity
   - Review WebSocket logs
   - Monitor reconnection attempts

3. **Performance Problems**
   - Check batch sizes
   - Monitor system resources
   - Review slow queries
   - Analyze trace data

### Rate Limiting

Requests are rate limited per IP address. The API uses a sliding window algorithm to track request counts.

Response headers include rate limit information:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1614556800
```

When rate limit is exceeded, the API returns:
```json
{
    "error": "Rate limit exceeded",
    "retry_after": 45
}
```

Default limits:
- Regular endpoints: 100 requests per minute
- Batch operations: 20 requests per minute
- WebSocket connections: 5 concurrent connections per IP

Rate limits can be configured per endpoint and client if needed. Contact support for custom limits.


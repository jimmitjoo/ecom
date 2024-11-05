package models

import "time"

// EventType represents the type of event
type EventType string

const (
	// EventProductCreated is emitted when a product is created
	EventProductCreated EventType = "product.created"
	// EventProductUpdated is emitted when a product is updated
	EventProductUpdated EventType = "product.updated"
	// EventProductDeleted is emitted when a product is deleted
	EventProductDeleted EventType = "product.deleted"
)

// Event represents a system event
type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// ProductEvent represents product-specific event data
type ProductEvent struct {
	ProductID string   `json:"product_id"`
	Action    string   `json:"action"`
	Product   *Product `json:"product,omitempty"`
}

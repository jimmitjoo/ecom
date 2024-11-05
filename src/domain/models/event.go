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
	Version   int64       `json:"version"`   // Event version for ordering
	EntityID  string      `json:"entity_id"` // ID of the affected entity
	CausedBy  string      `json:"caused_by"` // ID of the event that caused this event
	Sequence  int64       `json:"sequence"`  // Global sequence number
}

// ProductEvent represents product-specific event data
type ProductEvent struct {
	ProductID string   `json:"product_id"`
	Action    string   `json:"action"`
	Product   *Product `json:"product,omitempty"`
	Version   int64    `json:"version"`   // Product version
	PrevHash  string   `json:"prev_hash"` // Hash of previous state
	Changes   []Change `json:"changes"`   // Specific changes made
}

type Change struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value"`
}

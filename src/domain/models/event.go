package models

import (
	"errors"
	"time"
)

// EventType defines the type of event
type EventType string

const (
	EventProductCreated EventType = "product.created"
	EventProductUpdated EventType = "product.updated"
	EventProductDeleted EventType = "product.deleted"
)

// Event represents a domain event
type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	EntityID  string      `json:"entity_id"`
	Version   int64       `json:"version"`
	Sequence  int64       `json:"sequence"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// ProductEvent contains product-specific event data
type ProductEvent struct {
	ProductID string   `json:"product_id"`
	Action    string   `json:"action"`
	Product   *Product `json:"product"`
	Version   int64    `json:"version"`
	PrevHash  string   `json:"prev_hash"`
	Changes   []Change `json:"changes,omitempty"`
}

// Change represents a field change in an update event
type Change struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// ValidateEvent validates an event
func ValidateEvent(event *Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}
	if event.ID == "" {
		return errors.New("event ID is required")
	}
	if event.Type == "" {
		return errors.New("event type is required")
	}
	if event.Data == nil {
		return errors.New("event data is required")
	}

	// Validera ProductEvent om det är en sådan
	if productEvent, ok := event.Data.(*ProductEvent); ok {
		if productEvent.ProductID == "" {
			return errors.New("product ID is required in product event")
		}
		if productEvent.Action == "" {
			return errors.New("action is required in product event")
		}
		if productEvent.Product == nil && event.Type != EventProductDeleted {
			return errors.New("product is required in product event except for delete events")
		}
	}

	return nil
}

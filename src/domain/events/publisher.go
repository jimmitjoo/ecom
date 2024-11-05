package events

import "github.com/jimmitjoo/ecom/src/domain/models"

// EventPublisher defines the interface for publishing and subscribing to events
type EventPublisher interface {
	// Publish sends an event to all subscribers
	Publish(event *models.Event) error

	// Subscribe registers a handler for a specific event type
	Subscribe(eventType models.EventType, handler func(*models.Event)) error

	// Unsubscribe removes a handler for a specific event type
	Unsubscribe(eventType models.EventType, handler func(*models.Event)) error
}

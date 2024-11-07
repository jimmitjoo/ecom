package memory

import (
	"reflect"
	"sync"

	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"
)

// MemoryEventPublisher implements an in-memory event publishing system
type MemoryEventPublisher struct {
	handlers map[models.EventType][]func(*models.Event)
	mu       sync.RWMutex
}

// NewMemoryEventPublisher creates a new in-memory event publisher
func NewMemoryEventPublisher() events.EventPublisher {
	return &MemoryEventPublisher{
		handlers: make(map[models.EventType][]func(*models.Event)),
	}
}

// Publish sends an event to all registered handlers for its type
func (p *MemoryEventPublisher) Publish(event *models.Event) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if handlers, exists := p.handlers[event.Type]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
	return nil
}

// Subscribe registers a new handler for a specific event type
func (p *MemoryEventPublisher) Subscribe(eventType models.EventType, handler func(*models.Event)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers[eventType] = append(p.handlers[eventType], handler)
	return nil
}

// Unsubscribe removes a handler for a specific event type
func (p *MemoryEventPublisher) Unsubscribe(eventType models.EventType, handler func(*models.Event)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if handlers, exists := p.handlers[eventType]; exists {
		// Create a new slice for handlers
		newHandlers := make([]func(*models.Event), 0)
		handlerValue := reflect.ValueOf(handler)

		// Copy all handlers except the one to be removed
		for _, h := range handlers {
			if reflect.ValueOf(h).Pointer() != handlerValue.Pointer() {
				newHandlers = append(newHandlers, h)
			}
		}

		// Update handlers for this event type
		if len(newHandlers) > 0 {
			p.handlers[eventType] = newHandlers
		} else {
			delete(p.handlers, eventType) // Remove the entire event type if there are no handlers left
		}
	}
	return nil
}

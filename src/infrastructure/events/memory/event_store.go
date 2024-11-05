package memory

import (
	"sync"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

// MemoryEventStore implements an in-memory event store
type MemoryEventStore struct {
	events []*models.Event
	mu     sync.RWMutex
}

// NewMemoryEventStore creates a new in-memory event store
func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		events: make([]*models.Event, 0),
	}
}

// StoreEvent stores an event in memory
func (s *MemoryEventStore) StoreEvent(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Skapa en djup kopia av eventet innan vi sparar det
	eventCopy := *event
	if productEvent, ok := event.Data.(*models.ProductEvent); ok {
		productEventCopy := *productEvent
		productEventCopy.Product = productEvent.Product.Clone()
		eventCopy.Data = &productEventCopy
	}

	s.events = append(s.events, &eventCopy)
	return nil
}

// GetEvents returns all events for a specific entity from a given version
func (s *MemoryEventStore) GetEvents(entityID string, fromVersion int64) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredEvents []*models.Event
	for _, event := range s.events {
		if event.EntityID == entityID && event.Version >= fromVersion {
			// Skapa en djup kopia av eventet
			eventCopy := *event
			if productEvent, ok := event.Data.(*models.ProductEvent); ok {
				productEventCopy := *productEvent
				productEventCopy.Product = productEvent.Product.Clone()
				eventCopy.Data = &productEventCopy
			}
			filteredEvents = append(filteredEvents, &eventCopy)
		}
	}

	return filteredEvents, nil
}

// GetSnapshot returns the latest snapshot for an entity
func (s *MemoryEventStore) GetSnapshot(entityID string) (*models.Product, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var latestEvent *models.Event
	var latestVersion int64

	// Hitta senaste eventet för denna entity
	for _, event := range s.events {
		if event.EntityID == entityID && event.Version > latestVersion {
			latestEvent = event
			latestVersion = event.Version
		}
	}

	if latestEvent == nil {
		return nil, 0, nil
	}

	// Konvertera event data till product
	if productEvent, ok := latestEvent.Data.(*models.ProductEvent); ok {
		return productEvent.Product, latestEvent.Version, nil
	}

	return nil, 0, nil
}

// CreateSnapshot stores a snapshot of the entity state
func (s *MemoryEventStore) CreateSnapshot(entityID string, product *models.Product, version int64) error {
	// I en minnesimplementation behöver vi inte spara snapshots
	// eftersom vi alltid har alla events tillgängliga
	return nil
}

package memory

import (
	"sort"
	"sync"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

type MemoryEventStore struct {
	events map[string][]*models.Event // entityID -> events
	mu     sync.RWMutex
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		events: make(map[string][]*models.Event),
	}
}

func (s *MemoryEventStore) StoreEvent(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.EntityID]; !exists {
		s.events[event.EntityID] = make([]*models.Event, 0)
	}
	s.events[event.EntityID] = append(s.events[event.EntityID], event)

	// Sort events by version to maintain order
	sort.Slice(s.events[event.EntityID], func(i, j int) bool {
		return s.events[event.EntityID][i].Version < s.events[event.EntityID][j].Version
	})

	return nil
}

func (s *MemoryEventStore) GetEvents(entityID string, fromVersion int64) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if events, exists := s.events[entityID]; exists {
		result := make([]*models.Event, 0)
		for _, event := range events {
			if event.Version >= fromVersion {
				result = append(result, event)
			}
		}
		return result, nil
	}
	return nil, nil
}

func (s *MemoryEventStore) GetSnapshot(entityID string) (*models.Product, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[entityID]
	if !exists || len(events) == 0 {
		return nil, 0, nil
	}

	// Get latest event
	latestEvent := events[len(events)-1]
	if productEvent, ok := latestEvent.Data.(*models.ProductEvent); ok {
		return productEvent.Product, latestEvent.Version, nil
	}

	return nil, 0, nil
}

func (s *MemoryEventStore) CreateSnapshot(entityID string, product *models.Product, version int64) error {
	// In a real implementation, this would store the snapshot in a persistent store
	// For memory implementation, we don't need to do anything as we always have
	// all events in memory
	return nil
}

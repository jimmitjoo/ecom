package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

func createTestEvent(entityID string, version int64, eventType models.EventType, prevHash string) *models.Event {
	product := &models.Product{
		ID:        entityID,
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
		Prices: []models.Price{
			{Amount: 100, Currency: "SEK"},
		},
		Metadata: []models.MarketMetadata{
			{Market: "SE", Title: "Test Product", Description: "Test"},
		},
		Version: version,
	}
	product.LastHash = product.CalculateHash()

	return &models.Event{
		ID:       fmt.Sprintf("evt_%s_%d", entityID, version),
		Type:     eventType,
		EntityID: entityID,
		Version:  version,
		Sequence: version,
		Data: &models.ProductEvent{
			ProductID: entityID,
			Action:    string(eventType),
			Product:   product,
			Version:   version,
			PrevHash:  prevHash,
		},
		Timestamp: time.Now(),
	}
}

func TestStoreAndGetEvents(t *testing.T) {
	store := NewMemoryEventStore()
	entityID := "test_product_1"

	// Create and store some events
	events := []*models.Event{
		createTestEvent(entityID, 1, models.EventProductCreated, ""),
		createTestEvent(entityID, 2, models.EventProductUpdated, "hash1"),
		createTestEvent(entityID, 3, models.EventProductUpdated, "hash2"),
	}

	// Store events
	for _, event := range events {
		err := store.StoreEvent(event)
		assert.NoError(t, err)
	}

	// Get all events
	retrieved, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verify that events are correctly stored and retrieved
	for i, event := range retrieved {
		assert.Equal(t, events[i].ID, event.ID)
		assert.Equal(t, events[i].Version, event.Version)
		assert.Equal(t, events[i].EntityID, event.EntityID)

		// Verify event data
		originalEventData := events[i].Data.(*models.ProductEvent)
		retrievedEventData := event.Data.(*models.ProductEvent)

		assert.Equal(t, originalEventData.ProductID, retrievedEventData.ProductID)
		assert.Equal(t, originalEventData.Version, retrievedEventData.Version)
		assert.Equal(t, originalEventData.PrevHash, retrievedEventData.PrevHash)
		assert.Equal(t, originalEventData.Product.LastHash, retrievedEventData.Product.LastHash)
	}
}

func TestGetEventsFromVersion(t *testing.T) {
	store := NewMemoryEventStore()
	entityID := "test_product_1"

	// Create and store events
	events := []*models.Event{
		createTestEvent(entityID, 1, models.EventProductCreated, ""),
		createTestEvent(entityID, 2, models.EventProductUpdated, "hash1"),
		createTestEvent(entityID, 3, models.EventProductUpdated, "hash2"),
		createTestEvent(entityID, 4, models.EventProductUpdated, "hash3"),
	}

	for _, event := range events {
		err := store.StoreEvent(event)
		assert.NoError(t, err)
	}

	// Test getting events from different versions
	testCases := []struct {
		fromVersion int64
		expected    int
	}{
		{1, 4}, // All events
		{2, 3}, // From version 2 and up
		{3, 2}, // From version 3 and up
		{4, 1}, // Only the last event
		{5, 0}, // No events
	}

	for _, tc := range testCases {
		retrieved, err := store.GetEvents(entityID, tc.fromVersion)
		assert.NoError(t, err)
		assert.Len(t, retrieved, tc.expected, "For fromVersion %d", tc.fromVersion)
	}
}

func TestEventDeepCopy(t *testing.T) {
	store := NewMemoryEventStore()
	entityID := "test_product_1"

	// Create and store an event
	original := createTestEvent(entityID, 1, models.EventProductCreated, "")
	err := store.StoreEvent(original)
	assert.NoError(t, err)

	// Get the event
	retrieved, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 1)

	// Modify the retrieved event
	retrievedEvent := retrieved[0]
	retrievedEventData := retrievedEvent.Data.(*models.ProductEvent)
	retrievedEventData.Product.BaseTitle = "Modified Title"

	// Get the event again and verify that the original hasn't changed
	retrievedAgain, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrievedAgain, 1)

	retrievedAgainData := retrievedAgain[0].Data.(*models.ProductEvent)
	assert.NotEqual(t, "Modified Title", retrievedAgainData.Product.BaseTitle)
	assert.Equal(t, "Test Product", retrievedAgainData.Product.BaseTitle)
}

func TestConcurrentAccess(t *testing.T) {
	store := NewMemoryEventStore()
	entityID := "test_product_1"

	// Create multiple goroutines that save and retrieve events concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(version int64) {
			event := createTestEvent(entityID, version, models.EventProductUpdated, "hash")
			err := store.StoreEvent(event)
			assert.NoError(t, err)
			done <- true
		}(int64(i + 1))
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that all events were stored correctly
	events, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, events, 10)
}

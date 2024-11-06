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

	// Skapa och spara några events
	events := []*models.Event{
		createTestEvent(entityID, 1, models.EventProductCreated, ""),
		createTestEvent(entityID, 2, models.EventProductUpdated, "hash1"),
		createTestEvent(entityID, 3, models.EventProductUpdated, "hash2"),
	}

	// Spara events
	for _, event := range events {
		err := store.StoreEvent(event)
		assert.NoError(t, err)
	}

	// Hämta alla events
	retrieved, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 3)

	// Verifiera att events är korrekt sparade och hämtade
	for i, event := range retrieved {
		assert.Equal(t, events[i].ID, event.ID)
		assert.Equal(t, events[i].Version, event.Version)
		assert.Equal(t, events[i].EntityID, event.EntityID)

		// Verifiera event data
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

	// Skapa och spara events
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

	// Testa att hämta events från olika versioner
	testCases := []struct {
		fromVersion int64
		expected    int
	}{
		{1, 4}, // Alla events
		{2, 3}, // Från version 2 och uppåt
		{3, 2}, // Från version 3 och uppåt
		{4, 1}, // Bara sista eventet
		{5, 0}, // Inga events
	}

	for _, tc := range testCases {
		retrieved, err := store.GetEvents(entityID, tc.fromVersion)
		assert.NoError(t, err)
		assert.Len(t, retrieved, tc.expected, "För fromVersion %d", tc.fromVersion)
	}
}

func TestEventDeepCopy(t *testing.T) {
	store := NewMemoryEventStore()
	entityID := "test_product_1"

	// Skapa och spara ett event
	original := createTestEvent(entityID, 1, models.EventProductCreated, "")
	err := store.StoreEvent(original)
	assert.NoError(t, err)

	// Hämta eventet
	retrieved, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 1)

	// Modifiera det hämtade eventet
	retrievedEvent := retrieved[0]
	retrievedEventData := retrievedEvent.Data.(*models.ProductEvent)
	retrievedEventData.Product.BaseTitle = "Modified Title"

	// Hämta eventet igen och verifiera att originalet inte ändrades
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

	// Skapa flera goroutines som sparar och hämtar events samtidigt
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(version int64) {
			event := createTestEvent(entityID, version, models.EventProductUpdated, "hash")
			err := store.StoreEvent(event)
			assert.NoError(t, err)
			done <- true
		}(int64(i + 1))
	}

	// Vänta på att alla goroutines är klara
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verifiera att alla events sparades korrekt
	events, err := store.GetEvents(entityID, 1)
	assert.NoError(t, err)
	assert.Len(t, events, 10)
}

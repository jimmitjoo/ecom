package memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

func createTestProduct() *models.Product {
	return &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
		Prices: []models.Price{
			{Currency: "SEK", Amount: 100},
		},
		Metadata: []models.MarketMetadata{
			{Market: "SE", Title: "Test Product", Description: "Test"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Test Create
	err := repo.Create(product)
	assert.NoError(t, err)

	// Test GetByID
	retrieved, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, retrieved.ID)
	assert.Equal(t, product.SKU, retrieved.SKU)
	assert.Equal(t, product.BaseTitle, retrieved.BaseTitle)
}

func TestGetNonExistentProduct(t *testing.T) {
	repo := NewProductRepository()

	_, err := repo.GetByID("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, models.ErrProductNotFound, err)
}

func TestUpdateProduct(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Create first
	err := repo.Create(product)
	assert.NoError(t, err)

	// Update
	product.BaseTitle = "Updated Title"
	product.Version++
	err = repo.Update(product)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.BaseTitle)
	assert.Equal(t, int64(2), updated.Version)
}

func TestUpdateNonExistentProduct(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	err := repo.Update(product)
	assert.Error(t, err)
	assert.Equal(t, models.ErrProductNotFound, err)
}

func TestDeleteProduct(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Create first
	err := repo.Create(product)
	assert.NoError(t, err)

	// Delete
	err = repo.Delete(product.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(product.ID)
	assert.Error(t, err)
	assert.Equal(t, models.ErrProductNotFound, err)
}

func TestDeleteNonExistentProduct(t *testing.T) {
	repo := NewProductRepository()

	err := repo.Delete("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, models.ErrProductNotFound, err)
}

func TestListProducts(t *testing.T) {
	repo := NewProductRepository()

	// Create multiple products
	products := []*models.Product{
		createTestProduct(),
		createTestProduct(),
		createTestProduct(),
	}

	for i, p := range products {
		p.ID = p.ID + string(rune('A'+i))
		p.SKU = p.SKU + string(rune('A'+i))
		err := repo.Create(p)
		assert.NoError(t, err)
	}

	// List all products
	listed, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, listed, len(products))
}

func TestEventStorage(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Create event
	createEvent := &models.Event{
		ID:       "evt_1",
		Type:     models.EventProductCreated,
		EntityID: product.ID,
		Version:  1,
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "created",
			Product:   product,
			Version:   1,
		},
		Timestamp: time.Now(),
	}

	// Store event
	err := repo.StoreEvent(createEvent)
	assert.NoError(t, err)

	// Update event
	product.BaseTitle = "Updated Title"
	product.Version = 2
	updateEvent := &models.Event{
		ID:       "evt_2",
		Type:     models.EventProductUpdated,
		EntityID: product.ID,
		Version:  2,
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "updated",
			Product:   product,
			Version:   2,
			PrevHash:  product.LastHash,
		},
		Timestamp: time.Now(),
	}

	err = repo.StoreEvent(updateEvent)
	assert.NoError(t, err)

	// Get events
	events, err := repo.GetEventsByProductID(product.ID, 1)
	assert.NoError(t, err)
	assert.Len(t, events, 2)

	// Verify event order and data
	assert.Equal(t, createEvent.ID, events[0].ID)
	assert.Equal(t, updateEvent.ID, events[1].ID)
	assert.Equal(t, int64(1), events[0].Version)
	assert.Equal(t, int64(2), events[1].Version)
}

func TestConcurrentAccess(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Create product
	err := repo.Create(product)
	assert.NoError(t, err)

	// Concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := repo.GetByID(product.ID)
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(version int64) {
			defer wg.Done()
			p := createTestProduct()
			p.ID = product.ID
			p.Version = version
			repo.Update(p)
		}(int64(i + 2))
	}

	wg.Wait()

	// Verify final state
	final, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.NotNil(t, final)
}

func TestEventVersionFiltering(t *testing.T) {
	repo := NewProductRepository()
	product := createTestProduct()

	// Create multiple events with different versions
	versions := []int64{1, 2, 3, 4, 5}
	for _, version := range versions {
		event := &models.Event{
			ID:       fmt.Sprintf("evt_%d", version),
			Type:     models.EventProductUpdated,
			EntityID: product.ID,
			Version:  version,
			Data: &models.ProductEvent{
				ProductID: product.ID,
				Action:    "updated",
				Product:   product,
				Version:   version,
			},
			Timestamp: time.Now(),
		}
		err := repo.StoreEvent(event)
		assert.NoError(t, err)
	}

	// Test different fromVersion values
	testCases := []struct {
		fromVersion int64
		expected    int
	}{
		{1, 5}, // All events
		{3, 3}, // Version 3 and up
		{5, 1}, // Only version 5
		{6, 0}, // No events
	}

	for _, tc := range testCases {
		events, err := repo.GetEventsByProductID(product.ID, tc.fromVersion)
		assert.NoError(t, err)
		assert.Len(t, events, tc.expected)

		for _, event := range events {
			assert.GreaterOrEqual(t, event.Version, tc.fromVersion)
		}
	}
}

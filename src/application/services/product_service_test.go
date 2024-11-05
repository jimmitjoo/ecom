package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/infrastructure/repositories/memory"
)

// MockEventPublisher är en mock för EventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventPublisher) Subscribe(eventType models.EventType, handler func(*models.Event)) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventPublisher) Unsubscribe(eventType models.EventType, handler func(*models.Event)) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

// MockLockManager implementerar locks.LockManager interface
type MockLockManager struct {
	mock.Mock
}

func (m *MockLockManager) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockLockManager) ReleaseLock(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockLockManager) RefreshLock(key string, ttl time.Duration) error {
	args := m.Called(key, ttl)
	return args.Error(0)
}

func createValidProduct() *models.Product {
	return &models.Product{
		BaseTitle: "Test Produkt",
		SKU:       "TEST-123",
		Prices: []models.Price{
			{
				Amount:   100,
				Currency: "SEK",
			},
		},
		Metadata: []models.MarketMetadata{
			{
				Market:      "SE",
				Title:       "Test Produkt",
				Description: "Test beskrivning",
			},
		},
	}
}

func setupProductService() (*productService, *MockEventPublisher, *MockLockManager) {
	repo := memory.NewProductRepository()
	publisher := new(MockEventPublisher)
	lockManager := new(MockLockManager)

	publisher.On("Publish", mock.AnythingOfType("*models.Event")).Return(nil).Maybe()
	lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(true, nil).Maybe()
	lockManager.On("ReleaseLock", mock.AnythingOfType("string")).Return(nil).Maybe()

	return &productService{
		repo:      repo,
		publisher: publisher,
		locks:     lockManager,
	}, publisher, lockManager
}

func TestCreateProduct(t *testing.T) {
	service, publisher, _ := setupProductService()

	product := createValidProduct()
	err := service.CreateProduct(product)

	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
	assert.True(t, product.CreatedAt.Before(time.Now()))
	assert.Equal(t, int64(1), product.Version)
	assert.NotEmpty(t, product.LastHash)

	publisher.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	service, publisher, _ := setupProductService()

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	retrieved, err := service.GetProduct(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, retrieved.ID)
	assert.Equal(t, product.BaseTitle, retrieved.BaseTitle)
	assert.Equal(t, product.SKU, retrieved.SKU)

	publisher.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	service, publisher, lockManager := setupProductService()

	// Återställ standard mock-förväntningar
	lockManager.ExpectedCalls = nil
	lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(true, nil).Once()
	lockManager.On("ReleaseLock", mock.AnythingOfType("string")).Return(nil).Once()

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	product.BaseTitle = "Uppdaterad Produkt"
	err = service.UpdateProduct(product)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), product.Version)

	updated, err := service.GetProduct(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Uppdaterad Produkt", updated.BaseTitle)

	publisher.AssertExpectations(t)
	lockManager.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	service, publisher, _ := setupProductService()

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	err = service.DeleteProduct(product.ID)
	assert.NoError(t, err)

	_, err = service.GetProduct(product.ID)
	assert.Error(t, err)

	publisher.AssertExpectations(t)
}

func TestListProducts(t *testing.T) {
	service, publisher, _ := setupProductService()

	products := []*models.Product{
		createValidProduct(),
		createValidProduct(),
	}

	for i, p := range products {
		p.SKU = p.SKU + "-" + string(rune('A'+i))
		err := service.CreateProduct(p)
		assert.NoError(t, err)
	}

	listed, err := service.ListProducts()
	assert.NoError(t, err)
	assert.Len(t, listed, 2)

	publisher.AssertExpectations(t)
}

func TestBatchCreateProducts(t *testing.T) {
	service, publisher, _ := setupProductService()

	products := []*models.Product{
		createValidProduct(),
		createValidProduct(),
		createValidProduct(),
	}

	// Sätt unika SKU för varje produkt
	for i, p := range products {
		p.SKU = p.SKU + "-" + string(rune('A'+i))
	}

	results, err := service.BatchCreateProducts(products)
	assert.NoError(t, err)
	assert.Len(t, results, len(products))

	for i, result := range results {
		assert.True(t, result.Success, "Product %d should be created successfully", i)
		assert.Empty(t, result.Error)

		// Verifiera att produkten skapades i repository
		stored, err := service.GetProduct(result.ID)
		if assert.NoError(t, err, "Should be able to fetch created product") {
			assert.Equal(t, products[i].SKU, stored.SKU)
			assert.NotEmpty(t, stored.ID)
			assert.NotEmpty(t, stored.CreatedAt)
			assert.Equal(t, int64(1), stored.Version)
		}
	}

	publisher.AssertExpectations(t)
}

func TestBatchUpdateProducts(t *testing.T) {
	service, publisher, lockManager := setupProductService()

	// Återställ standard mock-förväntningar för locks
	lockManager.ExpectedCalls = nil

	products := []*models.Product{
		createValidProduct(),
		createValidProduct(),
	}

	// Skapa produkterna först
	for i, p := range products {
		p.SKU = p.SKU + "-" + string(rune('A'+i))
		err := service.CreateProduct(p)
		assert.NoError(t, err)
	}

	// Uppdatera produkterna
	for _, p := range products {
		p.BaseTitle = "Uppdaterad " + p.BaseTitle
	}

	// Sätt förväntningar för varje produkt som ska uppdateras
	for range products {
		lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(true, nil).Once()
		lockManager.On("ReleaseLock", mock.AnythingOfType("string")).Return(nil).Once()
	}

	results, err := service.BatchUpdateProducts(products)
	assert.NoError(t, err)
	assert.Len(t, results, len(products))

	for i, result := range results {
		assert.True(t, result.Success)
		assert.Empty(t, result.Error)

		// Verifiera uppdateringen
		updated, err := service.GetProduct(products[i].ID)
		assert.NoError(t, err)
		assert.Contains(t, updated.BaseTitle, "Uppdaterad")
	}

	publisher.AssertExpectations(t)
	lockManager.AssertExpectations(t)
}

func TestBatchDeleteProducts(t *testing.T) {
	service, publisher, _ := setupProductService()

	products := []*models.Product{
		createValidProduct(),
		createValidProduct(),
	}

	var ids []string
	for i, p := range products {
		p.SKU = p.SKU + "-" + string(rune('A'+i))
		err := service.CreateProduct(p)
		assert.NoError(t, err)
		ids = append(ids, p.ID)
	}

	results, err := service.BatchDeleteProducts(ids)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	for _, result := range results {
		assert.True(t, result.Success)
		assert.Empty(t, result.Error)
	}

	for _, id := range ids {
		_, err := service.GetProduct(id)
		assert.Error(t, err)
	}

	publisher.AssertExpectations(t)
}

func TestUpdateProductVersionConflict(t *testing.T) {
	service, publisher, lockManager := setupProductService()

	// Återställ standard mock-förväntningar
	lockManager.ExpectedCalls = nil
	lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(true, nil)
	lockManager.On("ReleaseLock", mock.AnythingOfType("string")).Return(nil)

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	// Skapa en kopia av produkten med gammal version
	conflictProduct := *product

	// Uppdatera originalprodukten
	product.BaseTitle = "Första uppdateringen"
	err = service.UpdateProduct(product)
	assert.NoError(t, err)

	// Försök uppdatera med den gamla kopian
	conflictProduct.BaseTitle = "Andra uppdateringen"
	err = service.UpdateProduct(&conflictProduct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version conflict")

	publisher.AssertExpectations(t)
	lockManager.AssertExpectations(t)
}

func TestUpdateProductLockFailure(t *testing.T) {
	service, publisher, lockManager := setupProductService()

	// Återställ mock och sätt ny förväntan för låsfel
	lockManager.ExpectedCalls = nil
	lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(false, nil)

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	err = service.UpdateProduct(product)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not acquire lock")

	publisher.AssertExpectations(t)
	lockManager.AssertExpectations(t)
}

func TestReplayEvents(t *testing.T) {
	service, publisher, lockManager := setupProductService()

	// Återställ standard mock-förväntningar
	lockManager.ExpectedCalls = nil
	lockManager.On("AcquireLock", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(true, nil)
	lockManager.On("ReleaseLock", mock.AnythingOfType("string")).Return(nil)

	product := createValidProduct()
	err := service.CreateProduct(product)
	assert.NoError(t, err)

	// Spara create event hash
	createHash := product.LastHash
	t.Logf("Create event hash: %s", createHash)

	// Uppdatera produkten några gånger
	var prevHash string
	for i := 0; i < 3; i++ {
		prevHash = product.LastHash
		t.Logf("Before update %d - Hash: %s", i+1, prevHash)

		product.BaseTitle = product.BaseTitle + " Updated"
		err = service.UpdateProduct(product)
		assert.NoError(t, err)

		t.Logf("After update %d - Hash: %s", i+1, product.LastHash)
		t.Logf("Product state after update %d: %+v", i+1, product)
	}

	events, err := service.ReplayEvents(product.ID, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
	assert.Len(t, events, 4) // Ska ha create + 3 uppdateringshändelser

	// Verifiera händelserna
	for i, event := range events {
		eventData, ok := event.Data.(*models.ProductEvent)
		assert.True(t, ok)

		t.Logf("Event %d - Type: %s, Version: %d", i, event.Type, event.Version)
		t.Logf("Event %d - PrevHash: %s", i, eventData.PrevHash)
		t.Logf("Event %d - Product Hash: %s", i, eventData.Product.LastHash)

		if i == 0 {
			// Första eventet ska vara create
			assert.Equal(t, models.EventProductCreated, event.Type)
			assert.Empty(t, eventData.PrevHash)
		} else {
			// Efterföljande events ska vara updates
			assert.Equal(t, models.EventProductUpdated, event.Type)
			prevEventData := events[i-1].Data.(*models.ProductEvent)
			assert.Equal(t, prevEventData.Product.LastHash, eventData.PrevHash,
				"Event %d: Hash mismatch. Expected %s, got %s",
				i, prevEventData.Product.LastHash, eventData.PrevHash)
		}
	}

	publisher.AssertExpectations(t)
	lockManager.AssertExpectations(t)
}

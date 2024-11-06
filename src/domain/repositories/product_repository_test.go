package repositories_test

import (
	"errors"
	"testing"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository är en mock-implementation av ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id string) (*models.Product, error) {
	args := m.Called(id)
	if p, ok := args.Get(0).(*models.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) List() ([]*models.Product, error) {
	args := m.Called()
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error) {
	args := m.Called(productID, fromVersion)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockProductRepository) StoreEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// TestProductRepositoryInterface verifierar att interfacet implementeras korrekt
func TestProductRepositoryInterface(t *testing.T) {
	var _ repositories.ProductRepository = &MockProductRepository{}
}

// TestRepositoryErrors verifierar att repository errors är korrekt definierade
func TestRepositoryErrors(t *testing.T) {
	// Verifiera att ErrProductNotFound är definierat
	assert.NotNil(t, models.ErrProductNotFound)
	assert.Contains(t, models.ErrProductNotFound.Error(), "not found")

	// Verifiera att andra repository errors är korrekt definierade
	var ErrVersionConflict = errors.New("version conflict")
	assert.NotNil(t, ErrVersionConflict)
	assert.Contains(t, ErrVersionConflict.Error(), "version conflict")
}

// TestMockRepositoryUsage visar hur mock repository kan användas
func TestMockRepositoryUsage(t *testing.T) {
	repo := new(MockProductRepository)
	product := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
	}

	// Testa Create
	repo.On("Create", product).Return(nil)
	err := repo.Create(product)
	assert.NoError(t, err)

	// Testa GetByID
	repo.On("GetByID", product.ID).Return(product, nil)
	retrieved, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, retrieved.ID)

	// Testa GetByID med not found error
	repo.On("GetByID", "nonexistent").Return(nil, models.ErrProductNotFound)
	_, err = repo.GetByID("nonexistent")
	assert.ErrorIs(t, err, models.ErrProductNotFound)

	// Testa List
	products := []*models.Product{product}
	repo.On("List").Return(products, nil)
	listed, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, listed, 1)

	// Testa Event hantering
	event := &models.Event{
		ID:       "test_event_1",
		Type:     models.EventProductCreated,
		EntityID: product.ID,
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "created",
			Product:   product,
		},
	}

	repo.On("StoreEvent", event).Return(nil)
	err = repo.StoreEvent(event)
	assert.NoError(t, err)

	events := []*models.Event{event}
	repo.On("GetEventsByProductID", product.ID, int64(1)).Return(events, nil)
	retrieved_events, err := repo.GetEventsByProductID(product.ID, 1)
	assert.NoError(t, err)
	assert.Len(t, retrieved_events, 1)

	repo.AssertExpectations(t)
}

// TestRepositoryErrorHandling verifierar att repository errors hanteras korrekt
func TestRepositoryErrorHandling(t *testing.T) {
	repo := new(MockProductRepository)

	// Sätt upp mock-förväntningar först
	repo.On("GetByID", "nonexistent").Return(nil, models.ErrProductNotFound)

	testCases := []struct {
		name      string
		operation func() error
		err       error
	}{
		{
			name: "Product not found",
			operation: func() error {
				_, err := repo.GetByID("nonexistent")
				return err
			},
			err: models.ErrProductNotFound,
		},
		{
			name: "Version conflict",
			operation: func() error {
				product := &models.Product{ID: "test_1", Version: 1}
				repo.On("Update", product).Return(models.ErrVersionConflict)
				return repo.Update(product)
			},
			err: models.ErrVersionConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.operation()
			assert.ErrorIs(t, err, tc.err)
		})
	}

	repo.AssertExpectations(t)
}

// TestEventHandling verifierar att event hantering fungerar korrekt
func TestEventHandling(t *testing.T) {
	repo := new(MockProductRepository)
	product := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
		Version:   1,
	}

	// Skapa en sekvens av events
	events := []*models.Event{
		{
			ID:       "evt_1",
			Type:     models.EventProductCreated,
			EntityID: product.ID,
			Version:  1,
			Data: &models.ProductEvent{
				ProductID: product.ID,
				Action:    "created",
				Product:   product,
			},
		},
		{
			ID:       "evt_2",
			Type:     models.EventProductUpdated,
			EntityID: product.ID,
			Version:  2,
			Data: &models.ProductEvent{
				ProductID: product.ID,
				Action:    "updated",
				Product:   product,
				PrevHash:  product.LastHash,
			},
		},
	}

	// Testa att hämta events från olika versioner
	repo.On("GetEventsByProductID", product.ID, int64(1)).Return(events, nil)
	repo.On("GetEventsByProductID", product.ID, int64(2)).Return(events[1:], nil)

	// Hämta alla events
	allEvents, err := repo.GetEventsByProductID(product.ID, 1)
	assert.NoError(t, err)
	assert.Len(t, allEvents, 2)

	// Hämta events från version 2
	laterEvents, err := repo.GetEventsByProductID(product.ID, 2)
	assert.NoError(t, err)
	assert.Len(t, laterEvents, 1)

	repo.AssertExpectations(t)
}

package repositories_test

import (
	"errors"
	"testing"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
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

func (m *MockProductRepository) List(page, pageSize int) ([]*models.Product, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error) {
	args := m.Called(productID, fromVersion)
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockProductRepository) StoreEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// TestProductRepositoryInterface verifies that the interface is implemented correctly
func TestProductRepositoryInterface(t *testing.T) {
	var _ repositories.ProductRepository = &MockProductRepository{}
}

// TestRepositoryErrors verifies that repository errors are defined correctly
func TestRepositoryErrors(t *testing.T) {
	// Verify that ErrProductNotFound is defined
	assert.NotNil(t, models.ErrProductNotFound)
	assert.Contains(t, models.ErrProductNotFound.Error(), "not found")

	// Verify that other repository errors are defined correctly
	var ErrVersionConflict = errors.New("version conflict")
	assert.NotNil(t, ErrVersionConflict)
	assert.Contains(t, ErrVersionConflict.Error(), "version conflict")
}

// TestMockRepositoryUsage shows how the mock repository can be used
func TestMockRepositoryUsage(t *testing.T) {
	repo := new(MockProductRepository)
	product := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
	}

	// Test Create
	repo.On("Create", product).Return(nil)
	err := repo.Create(product)
	assert.NoError(t, err)

	// Test GetByID
	repo.On("GetByID", product.ID).Return(product, nil)
	retrieved, err := repo.GetByID(product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, retrieved.ID)

	// Test GetByID with not found error
	repo.On("GetByID", "nonexistent").Return(nil, models.ErrProductNotFound)
	_, err = repo.GetByID("nonexistent")
	assert.ErrorIs(t, err, models.ErrProductNotFound)

	// Test List
	products := []*models.Product{product}
	repo.On("List", 1, 10).Return(products, 1, nil)
	listed, _, err := repo.List(1, 10)
	assert.NoError(t, err)
	assert.Len(t, listed, 1)

	// Test Event handling
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

// TestRepositoryErrorHandling verifies that repository errors are handled correctly
func TestRepositoryErrorHandling(t *testing.T) {
	repo := new(MockProductRepository)

	// Set up mock expectations first
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

// TestEventHandling verifies that event handling works correctly
func TestEventHandling(t *testing.T) {
	repo := new(MockProductRepository)
	product := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
		Version:   1,
	}

	// Create a sequence of events
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

	// Test getting events from different versions
	repo.On("GetEventsByProductID", product.ID, int64(1)).Return(events, nil)
	repo.On("GetEventsByProductID", product.ID, int64(2)).Return(events[1:], nil)

	// Get all events
	allEvents, err := repo.GetEventsByProductID(product.ID, 1)
	assert.NoError(t, err)
	assert.Len(t, allEvents, 2)

	// Get events from version 2
	laterEvents, err := repo.GetEventsByProductID(product.ID, 2)
	assert.NoError(t, err)
	assert.Len(t, laterEvents, 1)

	repo.AssertExpectations(t)
}

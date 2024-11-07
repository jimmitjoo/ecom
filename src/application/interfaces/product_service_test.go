package interfaces_test

import (
	"testing"

	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/mock"
)

// TestBatchResultImplementation verifierar att BatchResult har alla nödvändiga fält
func TestBatchResultImplementation(t *testing.T) {
	result := &interfaces.BatchResult{
		ID:      "test_id",
		Success: true,
		Error:   "test error",
	}

	// Verifiera att alla fält är tillgängliga
	_ = result.ID
	_ = result.Success
	_ = result.Error
}

// MockProductService implementerar ProductService interface för test
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) ListProducts(page, pageSize int) ([]*models.Product, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *MockProductService) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) GetProduct(id string) (*models.Product, error) {
	args := m.Called(id)
	if p, ok := args.Get(0).(*models.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductService) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) BatchCreateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	args := m.Called(products)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

func (m *MockProductService) BatchUpdateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	args := m.Called(products)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

func (m *MockProductService) BatchDeleteProducts(ids []string) ([]*interfaces.BatchResult, error) {
	args := m.Called(ids)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

// TestProductServiceInterface verifierar att MockProductService implementerar interfacet
func TestProductServiceInterface(t *testing.T) {
	var _ interfaces.ProductService = &MockProductService{} // Kompileringstest
}

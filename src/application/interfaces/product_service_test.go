package interfaces_test

import (
	"testing"

	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/models"
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
	ListProductsFunc        func() ([]*models.Product, error)
	CreateProductFunc       func(*models.Product) error
	GetProductFunc          func(string) (*models.Product, error)
	UpdateProductFunc       func(*models.Product) error
	DeleteProductFunc       func(string) error
	BatchCreateProductsFunc func([]*models.Product) ([]*interfaces.BatchResult, error)
	BatchUpdateProductsFunc func([]*models.Product) ([]*interfaces.BatchResult, error)
	BatchDeleteProductsFunc func([]string) ([]*interfaces.BatchResult, error)
}

func (m *MockProductService) ListProducts() ([]*models.Product, error) {
	if m.ListProductsFunc != nil {
		return m.ListProductsFunc()
	}
	return nil, nil
}

func (m *MockProductService) CreateProduct(p *models.Product) error {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(p)
	}
	return nil
}

func (m *MockProductService) GetProduct(id string) (*models.Product, error) {
	if m.GetProductFunc != nil {
		return m.GetProductFunc(id)
	}
	return nil, nil
}

func (m *MockProductService) UpdateProduct(p *models.Product) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(p)
	}
	return nil
}

func (m *MockProductService) DeleteProduct(id string) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(id)
	}
	return nil
}

func (m *MockProductService) BatchCreateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	if m.BatchCreateProductsFunc != nil {
		return m.BatchCreateProductsFunc(products)
	}
	return nil, nil
}

func (m *MockProductService) BatchUpdateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	if m.BatchUpdateProductsFunc != nil {
		return m.BatchUpdateProductsFunc(products)
	}
	return nil, nil
}

func (m *MockProductService) BatchDeleteProducts(ids []string) ([]*interfaces.BatchResult, error) {
	if m.BatchDeleteProductsFunc != nil {
		return m.BatchDeleteProductsFunc(ids)
	}
	return nil, nil
}

// TestProductServiceInterface verifierar att MockProductService implementerar interfacet
func TestProductServiceInterface(t *testing.T) {
	var _ interfaces.ProductService = &MockProductService{} // Kompileringstest
}

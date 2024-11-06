package interfaces

import "github.com/jimmitjoo/ecom/src/domain/models"

// BatchResult represents the result of a batch operation
type BatchResult struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ProductService defines the interface for product operations
type ProductService interface {
	ListProducts() ([]*models.Product, error)
	CreateProduct(product *models.Product) error
	GetProduct(id string) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id string) error

	// Batch operations
	BatchCreateProducts(products []*models.Product) ([]*BatchResult, error)
	BatchUpdateProducts(products []*models.Product) ([]*BatchResult, error)
	BatchDeleteProducts(ids []string) ([]*BatchResult, error)
}

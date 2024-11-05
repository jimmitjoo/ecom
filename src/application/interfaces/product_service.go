package interfaces

import "github.com/jimmitjoo/ecom/src/domain/models"

// ProductService defines the interface for product operations
type ProductService interface {
	ListProducts() ([]*models.Product, error)
	CreateProduct(product *models.Product) error
	GetProduct(id string) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id string) error
}

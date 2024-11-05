package memory

import (
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
)

// ProductRepository implements an in-memory product repository
type ProductRepository struct {
	products map[string]*models.Product
}

// NewProductRepository creates a new in-memory product repository
func NewProductRepository() repositories.ProductRepository {
	return &ProductRepository{
		products: make(map[string]*models.Product),
	}
}

// Create stores a new product in memory
func (r *ProductRepository) Create(product *models.Product) error {
	r.products[product.ID] = product
	return nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, repositories.ErrProductNotFound
	}
	return product, nil
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	if _, exists := r.products[product.ID]; !exists {
		return repositories.ErrProductNotFound
	}
	r.products[product.ID] = product
	return nil
}

// Delete removes a product from storage
func (r *ProductRepository) Delete(id string) error {
	if _, exists := r.products[id]; !exists {
		return repositories.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}

// List returns all stored products
func (r *ProductRepository) List() ([]*models.Product, error) {
	products := make([]*models.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
)

// productService implements the ProductService interface
type productService struct {
	repo      repositories.ProductRepository
	publisher events.EventPublisher
}

// NewProductService creates a new product service instance
func NewProductService(repo repositories.ProductRepository, publisher events.EventPublisher) interfaces.ProductService {
	return &productService{
		repo:      repo,
		publisher: publisher,
	}
}

// ListProducts retrieves all products from the repository
func (s *productService) ListProducts() ([]*models.Product, error) {
	return s.repo.List()
}

// CreateProduct creates a new product and publishes a creation event
func (s *productService) CreateProduct(product *models.Product) error {
	// Generate unique ID and set timestamps
	product.ID = "prod_" + uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Store the product
	if err := s.repo.Create(product); err != nil {
		return err
	}

	// Create and publish product creation event
	event := &models.Event{
		ID:   uuid.New().String(),
		Type: models.EventProductCreated,
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "created",
			Product:   product,
		},
		Timestamp: time.Now(),
	}

	return s.publisher.Publish(event)
}

// GetProduct retrieves a specific product by ID
func (s *productService) GetProduct(id string) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// UpdateProduct updates an existing product and publishes an update event
func (s *productService) UpdateProduct(product *models.Product) error {
	// Update timestamp
	product.UpdatedAt = time.Now()

	// Update the product
	if err := s.repo.Update(product); err != nil {
		return err
	}

	// Create and publish product update event
	event := &models.Event{
		ID:   uuid.New().String(),
		Type: models.EventProductUpdated,
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "updated",
			Product:   product,
		},
		Timestamp: time.Now(),
	}

	return s.publisher.Publish(event)
}

// DeleteProduct removes a product and publishes a deletion event
func (s *productService) DeleteProduct(id string) error {
	// Delete the product
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Create and publish product deletion event
	event := &models.Event{
		ID:   uuid.New().String(),
		Type: models.EventProductDeleted,
		Data: &models.ProductEvent{
			ProductID: id,
			Action:    "deleted",
		},
		Timestamp: time.Now(),
	}

	return s.publisher.Publish(event)
}

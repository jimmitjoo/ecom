package repositories

import (
	"errors"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

var (
	// ErrProductNotFound is returned when a product cannot be found
	ErrProductNotFound = errors.New("product not found")
)

// ProductRepository defines the interface for product storage operations
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
	List() ([]*models.Product, error)
	GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error)
	StoreEvent(event *models.Event) error
}

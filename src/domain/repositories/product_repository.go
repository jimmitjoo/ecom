package repositories

import "github.com/jimmitjoo/ecom/src/domain/models"

// ProductRepository defines the interface for product storage
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
	List() ([]*models.Product, error)
	GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error)
	StoreEvent(event *models.Event) error
}

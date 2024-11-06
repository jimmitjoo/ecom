package interfaces

import (
	"errors"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
	List() ([]*models.Product, error)
}

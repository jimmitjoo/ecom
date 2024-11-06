package memory

import (
	"sync"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
	eventstore "github.com/jimmitjoo/ecom/src/infrastructure/events/memory"
)

// ProductRepository implements an in-memory product repository
type ProductRepository struct {
	products   map[string]*models.Product
	eventStore *eventstore.MemoryEventStore
	mu         sync.RWMutex // Mutex to protect map-operations
}

// NewProductRepository creates a new in-memory product repository
func NewProductRepository() repositories.ProductRepository {
	return &ProductRepository{
		products:   make(map[string]*models.Product),
		eventStore: eventstore.NewMemoryEventStore(),
	}
}

// Create stores a new product in memory
func (r *ProductRepository) Create(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[product.ID] = product
	return nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, exists := r.products[id]
	if !exists {
		return nil, models.ErrProductNotFound
	}
	return product, nil
}

// Update modifies an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[product.ID]; !exists {
		return models.ErrProductNotFound
	}
	r.products[product.ID] = product
	return nil
}

// Delete removes a product from storage
func (r *ProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[id]; !exists {
		return models.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}

// List returns all stored products
func (r *ProductRepository) List() ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	products := make([]*models.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

// GetEventsByProductID hämtar alla events för en produkt från en given version
func (r *ProductRepository) GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error) {
	return r.eventStore.GetEvents(productID, fromVersion)
}

func (r *ProductRepository) StoreEvent(event *models.Event) error {
	return r.eventStore.StoreEvent(event)
}

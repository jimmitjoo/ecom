package repositories

import (
	"sync"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

type MemoryProductRepository struct {
	products map[string]*models.Product
	events   map[string][]*models.Event
	mu       sync.RWMutex
}

func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[string]*models.Product),
		events:   make(map[string][]*models.Event),
	}
}

func (r *MemoryProductRepository) Create(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if product.CreatedAt.IsZero() {
		product.CreatedAt = time.Now()
	}
	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) GetByID(id string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if product, exists := r.products[id]; exists {
		return product, nil
	}
	return nil, models.ErrProductNotFound
}

func (r *MemoryProductRepository) Update(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return models.ErrProductNotFound
	}
	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return models.ErrProductNotFound
	}
	delete(r.products, id)
	return nil
}

func (r *MemoryProductRepository) List(page, pageSize int) ([]*models.Product, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allProducts := make([]*models.Product, 0, len(r.products))
	for _, product := range r.products {
		allProducts = append(allProducts, product)
	}

	total := len(allProducts)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []*models.Product{}, total, nil
	}

	if end > total {
		end = total
	}

	return allProducts[start:end], total, nil
}

func (r *MemoryProductRepository) StoreEvent(event *models.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.events[event.EntityID] == nil {
		r.events[event.EntityID] = make([]*models.Event, 0)
	}
	r.events[event.EntityID] = append(r.events[event.EntityID], event)
	return nil
}

func (r *MemoryProductRepository) GetEventsByProductID(productID string, fromVersion int64) ([]*models.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := r.events[productID]
	if events == nil {
		return []*models.Event{}, nil
	}

	result := make([]*models.Event, 0)
	for _, e := range events {
		if e.Version >= fromVersion {
			result = append(result, e)
		}
	}
	return result, nil
}

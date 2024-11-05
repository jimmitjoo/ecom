package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
	"github.com/jimmitjoo/ecom/src/infrastructure/locks"
)

// productService implements the ProductService interface
type productService struct {
	repo      repositories.ProductRepository
	publisher events.EventPublisher
	locks     locks.LockManager
	sequence  atomic.Int64
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

	// Initiate version and hash
	product.Version = 1
	product.LastHash = product.CalculateHash()

	if err := s.repo.Create(product); err != nil {
		return err
	}

	event := &models.Event{
		ID:       uuid.New().String(),
		Type:     models.EventProductCreated,
		EntityID: product.ID,
		Version:  product.Version,
		Sequence: s.getNextSequence(),
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "created",
			Product:   product,
			Version:   product.Version,
			PrevHash:  "", // No previous version for new products
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
	ctx := context.Background()

	// Acquire lock
	acquired, err := s.locks.AcquireLock(ctx, product.ID, 10*time.Second)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("could not acquire lock for update")
	}
	defer s.locks.ReleaseLock(product.ID)

	// Get current version
	current, err := s.repo.GetByID(product.ID)
	if err != nil {
		return err
	}

	// Check version conflicts
	if product.Version != current.Version {
		return errors.New("version conflict detected")
	}

	// Calculate changes
	changes := calculateChanges(current, product)

	// Store current hash before update
	prevHash := current.LastHash

	// Update version and timestamps
	product.UpdateVersion() // This increments version and updates hash
	product.UpdatedAt = time.Now()

	// Create event with changes
	event := &models.Event{
		ID:       uuid.New().String(),
		Type:     models.EventProductUpdated,
		EntityID: product.ID,
		Version:  product.Version,
		Sequence: s.getNextSequence(),
		Data: &models.ProductEvent{
			ProductID: product.ID,
			Action:    "updated",
			Product:   product,
			Version:   product.Version,
			PrevHash:  prevHash,
			Changes:   changes,
		},
		Timestamp: time.Now(),
	}

	// Update product and publish event
	if err := s.repo.Update(product); err != nil {
		return err
	}

	return s.publisher.Publish(event)
}

// DeleteProduct removes a product and publishes a deletion event
func (s *productService) DeleteProduct(id string) error {
	// Get product before deletion for event data
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

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
			Product:   product, // Include product data in event
		},
		Timestamp: time.Now(),
	}

	return s.publisher.Publish(event)
}

// BatchCreateProducts creates multiple products in parallel
func (s *productService) BatchCreateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	results := make([]*interfaces.BatchResult, len(products))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, product := range products {
		wg.Add(1)
		go func(index int, p *models.Product) {
			defer wg.Done()
			result := &interfaces.BatchResult{ID: p.ID}

			// Generate ID and timestamps
			p.ID = "prod_" + uuid.New().String()
			p.CreatedAt = time.Now()
			p.UpdatedAt = time.Now()

			if err := models.ValidateProduct(p); err != nil {
				result.Success = false
				result.Error = err.Error()
			} else if err := s.repo.Create(p); err != nil {
				result.Success = false
				result.Error = "Failed to create product"
			} else {
				result.Success = true
				// Publish event for each successfully created product
				event := &models.Event{
					ID:   uuid.New().String(),
					Type: models.EventProductCreated,
					Data: &models.ProductEvent{
						ProductID: p.ID,
						Action:    "created",
						Product:   p,
					},
					Timestamp: time.Now(),
				}
				s.publisher.Publish(event)
			}

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, product)
	}

	wg.Wait()
	return results, nil
}

// BatchUpdateProducts updates multiple products in parallel
func (s *productService) BatchUpdateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	results := make([]*interfaces.BatchResult, len(products))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, product := range products {
		wg.Add(1)
		go func(index int, p *models.Product) {
			defer wg.Done()
			result := &interfaces.BatchResult{ID: p.ID}

			p.UpdatedAt = time.Now()

			if err := models.ValidateProduct(p); err != nil {
				result.Success = false
				result.Error = err.Error()
			} else if err := s.repo.Update(p); err != nil {
				result.Success = false
				result.Error = "Failed to update product"
			} else {
				result.Success = true
				// Publish event for each successfully updated product
				event := &models.Event{
					ID:   uuid.New().String(),
					Type: models.EventProductUpdated,
					Data: &models.ProductEvent{
						ProductID: p.ID,
						Action:    "updated",
						Product:   p,
					},
					Timestamp: time.Now(),
				}
				s.publisher.Publish(event)
			}

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, product)
	}

	wg.Wait()
	return results, nil
}

// BatchDeleteProducts deletes multiple products in parallel
func (s *productService) BatchDeleteProducts(ids []string) ([]*interfaces.BatchResult, error) {
	results := make([]*interfaces.BatchResult, len(ids))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, id := range ids {
		wg.Add(1)
		go func(index int, productID string) {
			defer wg.Done()
			result := &interfaces.BatchResult{ID: productID}

			// Get product before deletion for event data
			product, err := s.repo.GetByID(productID)
			if err != nil {
				result.Success = false
				result.Error = "Failed to find product"
				mu.Lock()
				results[index] = result
				mu.Unlock()
				return
			}

			if err := s.repo.Delete(productID); err != nil {
				result.Success = false
				result.Error = "Failed to delete product"
			} else {
				result.Success = true
				// Publish event for each successfully deleted product
				event := &models.Event{
					ID:   uuid.New().String(),
					Type: models.EventProductDeleted,
					Data: &models.ProductEvent{
						ProductID: productID,
						Action:    "deleted",
						Product:   product, // Include product data in event
					},
					Timestamp: time.Now(),
				}
				s.publisher.Publish(event)
			}

			mu.Lock()
			results[index] = result
			mu.Unlock()
		}(i, id)
	}

	wg.Wait()
	return results, nil
}

// Helper function for publishing events
func (s *productService) publishEvent(eventType models.EventType, action string, product *models.Product) {
	var productID string
	if product != nil {
		productID = product.ID
	}

	event := &models.Event{
		ID:   uuid.New().String(),
		Type: eventType,
		Data: &models.ProductEvent{
			ProductID: productID,
			Action:    action,
			Product:   product,
		},
		Timestamp: time.Now(),
	}
	s.publisher.Publish(event)
}

func (s *productService) ReplayEvents(productID string, fromVersion int64) ([]*models.Event, error) {
	events, err := s.repo.GetEventsByProductID(productID, fromVersion)
	if err != nil {
		return nil, err
	}

	// Verify event chain
	for i := 1; i < len(events); i++ {
		curr := events[i]
		prev := events[i-1]

		if curr.Version != prev.Version+1 {
			return nil, errors.New("event chain broken")
		}

		// Verify event integrity
		productEvent := curr.Data.(*models.ProductEvent)
		if productEvent.PrevHash != calculateHash(prev.Data.(*models.ProductEvent).Product) {
			return nil, errors.New("event chain integrity violated")
		}
	}

	return events, nil
}

func calculateChanges(old, new *models.Product) []models.Change {
	changes := make([]models.Change, 0)

	// Compare fields and record changes
	if old.BaseTitle != new.BaseTitle {
		changes = append(changes, models.Change{
			Field:    "base_title",
			OldValue: old.BaseTitle,
			NewValue: new.BaseTitle,
		})
	}
	// Add more field comparisons...

	return changes
}

func calculateHash(product *models.Product) string {
	// Implement a hash function for the product state
	data, _ := json.Marshal(product)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *productService) getNextSequence() int64 {
	return s.sequence.Add(1)
}

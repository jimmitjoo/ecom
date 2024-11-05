package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

	// Create event first
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

	// Store event first
	if err := s.repo.StoreEvent(event); err != nil {
		return err
	}

	// Then create the product
	if err := s.repo.Create(product); err != nil {
		return err
	}

	// Finally publish the event
	return s.publisher.Publish(event)
}

// GetProduct retrieves a specific product by ID
func (s *productService) GetProduct(id string) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// UpdateProduct updates an existing product and publishes an update event
func (s *productService) UpdateProduct(product *models.Product) error {
	ctx := context.Background()

	acquired, err := s.locks.AcquireLock(ctx, product.ID, 10*time.Second)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("could not acquire lock for update")
	}
	defer s.locks.ReleaseLock(product.ID)

	current, err := s.repo.GetByID(product.ID)
	if err != nil {
		return err
	}

	if product.Version != current.Version {
		return errors.New("version conflict detected")
	}

	// Skapa en kopia av produkten innan vi uppdaterar den
	updatedProduct := product.Clone()
	updatedProduct.Version++
	updatedProduct.UpdatedAt = time.Now()
	updatedProduct.LastHash = updatedProduct.CalculateHash()

	// Skapa event med den uppdaterade produkten
	event := &models.Event{
		ID:       uuid.New().String(),
		Type:     models.EventProductUpdated,
		EntityID: updatedProduct.ID,
		Version:  updatedProduct.Version,
		Sequence: s.getNextSequence(),
		Data: &models.ProductEvent{
			ProductID: updatedProduct.ID,
			Action:    "updated",
			Product:   updatedProduct.Clone(), // Spara en kopia av den uppdaterade produkten
			Version:   updatedProduct.Version,
			PrevHash:  current.LastHash, // Använd current.LastHash som PrevHash
			Changes:   calculateChanges(current, updatedProduct),
		},
		Timestamp: time.Now(),
	}

	// Spara event först
	if err := s.repo.StoreEvent(event); err != nil {
		return err
	}

	// Sen uppdatera produkten
	if err := s.repo.Update(updatedProduct); err != nil {
		return err
	}

	// Kopiera tillbaka värdena till original-produkten
	*product = *updatedProduct

	return s.publisher.Publish(event)
}

// DeleteProduct removes a product and publishes a deletion event
func (s *productService) DeleteProduct(id string) error {
	// Get product before deletion for event data
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Create deletion event
	event := &models.Event{
		ID:       uuid.New().String(),
		Type:     models.EventProductDeleted,
		EntityID: id,
		Version:  product.Version + 1,
		Sequence: s.getNextSequence(),
		Data: &models.ProductEvent{
			ProductID: id,
			Action:    "deleted",
			Product:   product,
			Version:   product.Version + 1,
			PrevHash:  product.LastHash, // Använd current hash som prev hash
		},
		Timestamp: time.Now(),
	}

	// Store event first
	if err := s.repo.StoreEvent(event); err != nil {
		return err
	}

	// Then delete the product
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Finally publish the event
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

			// Skapa produkten först
			err := s.CreateProduct(p)

			mu.Lock()
			results[index] = &interfaces.BatchResult{
				ID:      p.ID, // Nu har produkten ett ID efter CreateProduct
				Success: err == nil,
			}
			if err != nil {
				results[index].Error = err.Error()
			}
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

			err := s.UpdateProduct(p)

			mu.Lock()
			results[index] = &interfaces.BatchResult{
				ID:      p.ID,
				Success: err == nil,
			}
			if err != nil {
				results[index].Error = err.Error()
			}
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

	if len(events) == 0 {
		return events, nil
	}

	// Sortera events efter version
	sort.Slice(events, func(i, j int) bool {
		return events[i].Version < events[j].Version
	})

	// Verifiera event chain
	for i := 0; i < len(events); i++ {
		curr := events[i]
		currEvent, ok := curr.Data.(*models.ProductEvent)
		if !ok {
			return nil, errors.New("invalid event data")
		}

		if i == 0 {
			// För första eventet i sekvensen
			if curr.Type == models.EventProductCreated {
				// Create-event ska inte ha någon PrevHash
				if currEvent.PrevHash != "" {
					return nil, errors.New("create event should not have prev hash")
				}
			} else {
				// Om första eventet inte är create, verifiera att det har en PrevHash
				if currEvent.PrevHash == "" {
					return nil, errors.New("non-create event must have prev hash")
				}
			}
			continue
		}

		// För efterföljande events
		prev := events[i-1]
		prevEvent, ok := prev.Data.(*models.ProductEvent)
		if !ok {
			return nil, errors.New("invalid event data")
		}

		// Kontrollera versionsnumren
		if curr.Version != prev.Version+1 {
			return nil, fmt.Errorf("event chain broken: curr version %d, prev version %d",
				curr.Version, prev.Version)
		}

		// Verifiera hash-kedjan
		if currEvent.PrevHash != prevEvent.Product.LastHash {
			return nil, fmt.Errorf("event chain integrity violated: expected hash %s, got %s",
				prevEvent.Product.LastHash, currEvent.PrevHash)
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

func (s *productService) getNextSequence() int64 {
	return s.sequence.Add(1)
}

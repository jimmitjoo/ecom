package chaos

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/application/services"
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/jimmitjoo/ecom/src/domain/repositories"
	"github.com/jimmitjoo/ecom/src/infrastructure/events/memory"
	memorylocks "github.com/jimmitjoo/ecom/src/infrastructure/locks/memory"
)

func setupTestService() (interfaces.ProductService, error) {
	repo := repositories.NewMemoryProductRepository()
	publisher := memory.NewMemoryEventPublisher()
	lockManager := memorylocks.NewMemoryLockManager()

	service := services.NewProductService(repo, publisher, lockManager)
	if service == nil {
		return nil, fmt.Errorf("failed to create product service")
	}

	return service, nil
}

func generateLargeProductBatch(count int) []*models.Product {
	products := make([]*models.Product, count)

	for i := 0; i < count; i++ {
		products[i] = generateValidProduct()
	}
	return products
}

func generateValidProduct() *models.Product {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &models.Product{
		ID:        fmt.Sprintf("prod_%d", r.Intn(10000)),
		SKU:       fmt.Sprintf("SKU-%d", r.Intn(10000)),
		BaseTitle: fmt.Sprintf("Product %d", r.Intn(10000)),
		Version:   1,
		Prices: []models.Price{
			{Amount: 100.0, Currency: "SEK"},
		},
		Metadata: []models.MarketMetadata{
			{
				Market:      "SE",
				Title:       "Test Title",
				Description: "Test Description",
			},
		},
	}
}

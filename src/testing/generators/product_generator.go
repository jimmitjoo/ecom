package generators

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

// ProductGenerator contains configuration for generating test products
type ProductGenerator struct {
	SkuPrefix  string
	Markets    []string
	Currencies []string
}

// NewProductGenerator creates a new generator with default values
func NewProductGenerator() *ProductGenerator {
	return &ProductGenerator{
		SkuPrefix:  "TEST",
		Markets:    []string{"SE", "NO", "DK", "FI"},
		Currencies: []string{"SEK", "NOK", "DKK", "EUR"},
	}
}

// GenerateProduct creates a single test product
func (g *ProductGenerator) GenerateProduct() *models.Product {
	rand.Seed(time.Now().UnixNano())

	sku := fmt.Sprintf("%s-%d", g.SkuPrefix, rand.Intn(100000))

	product := &models.Product{
		SKU:       sku,
		BaseTitle: fmt.Sprintf("Testprodukt %s", sku),
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add prices
	product.Prices = g.generatePrices()

	// Add metadata for each market
	product.Metadata = g.generateMetadata(product.BaseTitle)

	return product
}

// GenerateProducts creates multiple test products
func (g *ProductGenerator) GenerateProducts(count int) []*models.Product {
	products := make([]*models.Product, count)
	for i := 0; i < count; i++ {
		products[i] = g.GenerateProduct()
	}
	return products
}

// generatePrices creates test prices for different currencies
func (g *ProductGenerator) generatePrices() []models.Price {
	prices := make([]models.Price, len(g.Currencies))
	for i, currency := range g.Currencies {
		prices[i] = models.Price{
			Amount:   float64(100 + rand.Intn(9900)),
			Currency: currency,
		}
	}
	return prices
}

// generateMetadata creates test metadata for different markets
func (g *ProductGenerator) generateMetadata(baseTitle string) []models.MarketMetadata {
	metadata := make([]models.MarketMetadata, len(g.Markets))
	for i, market := range g.Markets {
		metadata[i] = models.MarketMetadata{
			Market:      market,
			Title:       fmt.Sprintf("%s (%s)", baseTitle, market),
			Description: fmt.Sprintf("Testbeskrivning för %s på marknad %s", baseTitle, market),
		}
	}
	return metadata
}

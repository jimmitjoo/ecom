package generators

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
)

// ProductGenerator innehåller konfiguration för att generera testprodukter
type ProductGenerator struct {
	SkuPrefix  string
	Markets    []string
	Currencies []string
}

// NewProductGenerator skapar en ny generator med standardvärden
func NewProductGenerator() *ProductGenerator {
	return &ProductGenerator{
		SkuPrefix:  "TEST",
		Markets:    []string{"SE", "NO", "DK", "FI"},
		Currencies: []string{"SEK", "NOK", "DKK", "EUR"},
	}
}

// GenerateProduct skapar en enskild testprodukt
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

	// Lägg till priser
	product.Prices = g.generatePrices()

	// Lägg till metadata för varje marknad
	product.Metadata = g.generateMetadata(product.BaseTitle)

	return product
}

// GenerateProducts skapar flera testprodukter
func (g *ProductGenerator) GenerateProducts(count int) []*models.Product {
	products := make([]*models.Product, count)
	for i := 0; i < count; i++ {
		products[i] = g.GenerateProduct()
	}
	return products
}

// generatePrices skapar testpriser för olika valutor
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

// generateMetadata skapar testmetadata för olika marknader
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

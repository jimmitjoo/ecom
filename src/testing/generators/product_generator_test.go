package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductGenerator(t *testing.T) {
	generator := NewProductGenerator()

	// Test generating a single product
	product := generator.GenerateProduct()
	assert.NotEmpty(t, product.SKU)
	assert.NotEmpty(t, product.BaseTitle)
	assert.Len(t, product.Prices, len(generator.Currencies))
	assert.Len(t, product.Metadata, len(generator.Markets))

	// Test generating multiple products
	products := generator.GenerateProducts(5)
	assert.Len(t, products, 5)

	// Verify that the products are unique
	skus := make(map[string]bool)
	for _, p := range products {
		assert.False(t, skus[p.SKU], "SKU should be unique")
		skus[p.SKU] = true
	}
}

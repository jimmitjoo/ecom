package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductGenerator(t *testing.T) {
	generator := NewProductGenerator()

	// Testa generering av enskild produkt
	product := generator.GenerateProduct()
	assert.NotEmpty(t, product.SKU)
	assert.NotEmpty(t, product.BaseTitle)
	assert.Len(t, product.Prices, len(generator.Currencies))
	assert.Len(t, product.Metadata, len(generator.Markets))

	// Testa generering av flera produkter
	products := generator.GenerateProducts(5)
	assert.Len(t, products, 5)

	// Verifiera att produkterna är unika
	skus := make(map[string]bool)
	for _, p := range products {
		assert.False(t, skus[p.SKU], "SKU ska vara unik")
		skus[p.SKU] = true
	}
}

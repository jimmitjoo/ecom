package models

import (
	"testing"
)

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr bool
	}{
		{
			name: "valid product",
			product: &Product{
				ID:        "prod_123",
				SKU:       "TEST-001",
				BaseTitle: "Test Product",
				Prices: []Price{
					{Currency: "SEK", Amount: 299.00},
				},
				Metadata: []MarketMetadata{
					{Market: "SE", Title: "Test Product"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			product: &Product{
				SKU: "TEST-001",
			},
			wantErr: true,
		},
		{
			name: "invalid price currency",
			product: &Product{
				ID:        "prod_123",
				SKU:       "TEST-001",
				BaseTitle: "Test Product",
				Prices: []Price{
					{Currency: "INVALID", Amount: 299.00},
				},
				Metadata: []MarketMetadata{
					{Market: "SE", Title: "Test Product"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProduct(tt.product)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduct_CalculateHash(t *testing.T) {
	product := &Product{
		ID:        "prod_123",
		SKU:       "TEST-001",
		BaseTitle: "Test Product",
		Version:   1,
	}

	hash1 := product.CalculateHash()
	if hash1 == "" {
		t.Error("CalculateHash() returned empty hash")
	}

	// Same product should generate same hash
	hash2 := product.CalculateHash()
	if hash1 != hash2 {
		t.Error("CalculateHash() not deterministic")
	}

	// Modified product should generate different hash
	product.BaseTitle = "Modified Product"
	hash3 := product.CalculateHash()
	if hash1 == hash3 {
		t.Error("CalculateHash() didn't change with modified product")
	}
}

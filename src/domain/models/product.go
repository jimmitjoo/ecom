package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
)

// Price represents a price for a specific market
type Price struct {
	Currency string  `json:"currency" validate:"required,len=3"`
	Amount   float64 `json:"amount" validate:"required,gte=0"`
}

// Stock represents inventory for a specific location
type Stock struct {
	LocationID string `json:"location_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"gte=0"`
}

// Variant represents product variants
type Variant struct {
	ID         string            `json:"id" validate:"required"`
	SKU        string            `json:"sku" validate:"required"`
	Attributes map[string]string `json:"attributes" validate:"required"` // e.g. {"size": "XL", "color": "blue"}
	Stock      []Stock           `json:"stock"`
}

// MarketMetadata contains market-specific information
type MarketMetadata struct {
	Market      string `json:"market" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}

// Product is the main product structure
type Product struct {
	ID          string           `json:"id" validate:"required"`
	SKU         string           `json:"sku" validate:"required"`
	BaseTitle   string           `json:"base_title" validate:"required"`
	Description string           `json:"description"`
	Prices      []Price          `json:"prices" validate:"required,dive"`
	Variants    []Variant        `json:"variants" validate:"dive"`
	Metadata    []MarketMetadata `json:"metadata" validate:"required,dive"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Version     int64            `json:"version"`   // Version number for optimistic locking
	LastHash    string           `json:"last_hash"` // Hash of last known state
}

func ValidateProduct(product *Product) error {
	validate := validator.New()
	return validate.Struct(product)
}

// CalculateHash generates a hash of the product's current state
func (p *Product) CalculateHash() string {
	// Använd den befintliga calculateHash-funktionen från service-lagret
	data, _ := json.Marshal(p)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// UpdateVersion increments the version and updates the hash
func (p *Product) UpdateVersion() {
	p.Version++
	p.LastHash = p.CalculateHash()
}

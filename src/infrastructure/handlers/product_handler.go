package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/models"

	"github.com/gorilla/mux"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	service interfaces.ProductService
}

// NewProductHandler creates a new product handler instance
func NewProductHandler(service interfaces.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// writeError is a helper function to write error responses
func (h *ProductHandler) writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.NewAPIError(message))
}

// ListProducts handles GET requests to retrieve all products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// CreateProduct handles POST requests to create a new product
// First validates the input, then creates the product, and finally validates the complete product
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Create the product first (which sets ID and timestamps)
	if err := h.service.CreateProduct(&product); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	// Validate the complete product after ID has been set
	if err := models.ValidateProduct(&product); err != nil {
		// If validation fails, clean up by deleting the product
		h.service.DeleteProduct(product.ID)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProduct handles GET requests to retrieve a specific product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.GetProduct(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Product with ID '%s' not found", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles PUT requests to update an existing product
// Validates the input, ensures the ID matches, and updates the product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Ensure the ID in the URL matches the product
	product.ID = id
	if err := models.ValidateProduct(&product); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateProduct(&product); err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Product with ID '%s' not found", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct handles DELETE requests to remove a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteProduct(id); err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Product with ID '%s' not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

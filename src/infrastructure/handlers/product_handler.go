package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// ListProducts godoc
// @Summary Lista alla produkter
// @Description Hämtar en lista över alla produkter
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} handlers.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// CreateProduct godoc
// @Summary Skapa en ny produkt
// @Description Skapar en ny produkt med angivna detaljer
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Produktdetaljer"
// @Success 201 {object} models.Product
// @Failure 400 {object} handlers.ErrorResponse
// @Router /products [post]
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

// GetProduct godoc
// @Summary Hämta en produkt
// @Description Hämtar en produkt med angivet ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Produkt ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} handlers.ErrorResponse
// @Router /products/{id} [get]
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

// UpdateProduct godoc
// @Summary Uppdatera en produkt
// @Description Uppdaterar en existerande produkt
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Produkt ID"
// @Param product body models.Product true "Uppdaterade produktdetaljer"
// @Success 200 {object} models.Product
// @Failure 400,404 {object} handlers.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Hämta först existerande produkt
	existingProduct, err := h.service.GetProduct(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, fmt.Sprintf("Product with ID '%s' not found", id))
		return
	}

	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Använd ID från URL:en, inte från request body
	updatedProduct.ID = id
	// Behåll version från existerande produkt
	updatedProduct.Version = existingProduct.Version
	// Behåll created_at från existerande produkt
	updatedProduct.CreatedAt = existingProduct.CreatedAt
	// Uppdatera updated_at till nu
	updatedProduct.UpdatedAt = time.Now()

	if err := h.service.UpdateProduct(&updatedProduct); err != nil {
		h.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update product: %v", err))
		return
	}

	h.sendSuccess(w, http.StatusOK, updatedProduct)
}

// DeleteProduct godoc
// @Summary Ta bort en produkt
// @Description Tar bort en produkt med angivet ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Produkt ID"
// @Success 204 "No Content"
// @Failure 404 {object} handlers.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteProduct(id); err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Product with ID '%s' not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BatchCreateProducts godoc
// @Summary Create multiple products in bulk
// @Description Creates multiple products simultaneously in a single request
// @Tags products
// @Accept json
// @Produce json
// @Param products body []models.Product true "Array of products to create"
// @Success 201 {array} models.Product "Array of created products"
// @Failure 400 {object} models.APIError "Invalid JSON data"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /products/batch [post]
func (h *ProductHandler) BatchCreateProducts(w http.ResponseWriter, r *http.Request) {
	var products []*models.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	results, err := h.service.BatchCreateProducts(products)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}

// BatchUpdateProducts godoc
// @Summary Batch update multiple products simultaneously
// @Description Updates multiple products in a single request. All products must exist and contain valid data.
// @Tags products
// @Accept json
// @Produce json
// @Param products body []models.Product true "Array of products to update with their IDs and new data"
// @Success 200 {array} models.Product "Array of updated products"
// @Failure 400 {object} models.APIError "Invalid JSON data or validation errors"
// @Failure 404 {object} models.APIError "One or more products not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /products/batch [put]
func (h *ProductHandler) BatchUpdateProducts(w http.ResponseWriter, r *http.Request) {
	var products []*models.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	results, err := h.service.BatchUpdateProducts(products)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// BatchDeleteProducts godoc
// @Summary Batch delete multiple products simultaneously
// @Description Deletes multiple products in a single request by their IDs. Returns results of deletion operations.
// @Tags products
// @Accept json
// @Produce json
// @Param productIDs body []string true "Array of product IDs to delete"
// @Success 200 {object} map[string]string "Map of product IDs to deletion status"
// @Failure 400 {object} models.APIError "Invalid JSON data"
// @Failure 404 {object} models.APIError "One or more products not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /products/batch [delete]
func (h *ProductHandler) BatchDeleteProducts(w http.ResponseWriter, r *http.Request) {
	var productIDs []string
	if err := json.NewDecoder(r.Body).Decode(&productIDs); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	results, err := h.service.BatchDeleteProducts(productIDs)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to delete products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *ProductHandler) sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func (h *ProductHandler) sendSuccess(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(SuccessResponse{
			Success: true,
			Data:    data,
		})
	}
}

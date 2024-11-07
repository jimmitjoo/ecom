package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gorilla/mux"
	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/domain/models"
)

// MockProductService is a mock for the ProductService interface
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) ListProducts(page, pageSize int) ([]*models.Product, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*models.Product), args.Int(1), args.Error(2)
}

func (m *MockProductService) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) GetProduct(id string) (*models.Product, error) {
	args := m.Called(id)
	if p, ok := args.Get(0).(*models.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductService) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) BatchCreateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	args := m.Called(products)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

func (m *MockProductService) BatchUpdateProducts(products []*models.Product) ([]*interfaces.BatchResult, error) {
	args := m.Called(products)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

func (m *MockProductService) BatchDeleteProducts(ids []string) ([]*interfaces.BatchResult, error) {
	args := m.Called(ids)
	return args.Get(0).([]*interfaces.BatchResult), args.Error(1)
}

func createTestProduct() *models.Product {
	return &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Test Product",
		Prices: []models.Price{
			{Currency: "SEK", Amount: 100},
		},
		Metadata: []models.MarketMetadata{
			{Market: "SE", Title: "Test Product", Description: "Test"},
		},
	}
}

func TestListProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []*models.Product{
		{ID: "1", BaseTitle: "Product 1"},
		{ID: "2", BaseTitle: "Product 2"},
	}
	totalItems := 10

	mockService.On("ListProducts", 1, 10).Return(products, totalItems, nil)

	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data       []*models.Product `json:"data"`
		Page       int               `json:"page"`
		PageSize   int               `json:"page_size"`
		TotalItems int               `json:"total_items"`
		TotalPages int               `json:"total_pages"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, products, response.Data)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, totalItems, response.TotalItems)
	assert.Equal(t, 1, response.TotalPages)

	mockService.AssertExpectations(t)
}

func TestListProductsWithCustomPagination(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []*models.Product{
		{ID: "1", BaseTitle: "Product 1"},
		{ID: "2", BaseTitle: "Product 2"},
	}
	totalItems := 20

	mockService.On("ListProducts", 2, 5).Return(products, totalItems, nil)

	req := httptest.NewRequest("GET", "/products?page=2&size=5", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data       []*models.Product `json:"data"`
		Page       int               `json:"page"`
		PageSize   int               `json:"page_size"`
		TotalItems int               `json:"total_items"`
		TotalPages int               `json:"total_pages"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, products, response.Data)
	assert.Equal(t, 2, response.Page)
	assert.Equal(t, 5, response.PageSize)
	assert.Equal(t, totalItems, response.TotalItems)
	assert.Equal(t, 4, response.TotalPages)

	mockService.AssertExpectations(t)
}

func TestCreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	product := createTestProduct()
	mockService.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(nil)

	body, _ := json.Marshal(product)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, product.SKU, response.SKU)

	mockService.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	product := createTestProduct()
	mockService.On("GetProduct", product.ID).Return(product, nil)

	req := httptest.NewRequest("GET", "/products/"+product.ID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": product.ID})
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, product.ID, response.ID)

	mockService.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	existingProduct := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Original Title",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedProduct := &models.Product{
		ID:        "test_prod_1",
		SKU:       "TEST-123",
		BaseTitle: "Updated Title",
		Version:   1,
		CreatedAt: existingProduct.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Mock GetProduct call
	mockService.On("GetProduct", "test_prod_1").Return(existingProduct, nil)

	// Mock UpdateProduct call
	mockService.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil)

	// Create request body
	body, _ := json.Marshal(updatedProduct)
	req, _ := http.NewRequest("PUT", "/products/test_prod_1", bytes.NewBuffer(body))

	// Set up Gorilla Mux router to handle URL parameters
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.UpdateProduct)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := "test_prod_1"
	mockService.On("DeleteProduct", productID).Return(nil)

	req := httptest.NewRequest("DELETE", "/products/"+productID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": productID})
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestBatchCreateProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []*models.Product{createTestProduct()}
	results := []*interfaces.BatchResult{{ID: products[0].ID, Success: true}}
	mockService.On("BatchCreateProducts", mock.AnythingOfType("[]*models.Product")).Return(results, nil)

	body, _ := json.Marshal(products)
	req := httptest.NewRequest("POST", "/products/batch", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.BatchCreateProducts(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response []*interfaces.BatchResult
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.True(t, response[0].Success)

	mockService.AssertExpectations(t)
}

func TestBatchUpdateProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []*models.Product{createTestProduct()}
	results := []*interfaces.BatchResult{{ID: products[0].ID, Success: true}}
	mockService.On("BatchUpdateProducts", mock.AnythingOfType("[]*models.Product")).Return(results, nil)

	body, _ := json.Marshal(products)
	req := httptest.NewRequest("PUT", "/products/batch", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.BatchUpdateProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*interfaces.BatchResult
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.True(t, response[0].Success)

	mockService.AssertExpectations(t)
}

func TestBatchDeleteProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	ids := []string{"test_prod_1"}
	results := []*interfaces.BatchResult{{ID: ids[0], Success: true}}
	mockService.On("BatchDeleteProducts", mock.AnythingOfType("[]string")).Return(results, nil)

	body, _ := json.Marshal(ids)
	req := httptest.NewRequest("DELETE", "/products/batch", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.BatchDeleteProducts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*interfaces.BatchResult
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.True(t, response[0].Success)

	mockService.AssertExpectations(t)
}

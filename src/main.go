package main

import (
	"log"
	"net/http"

	"github.com/jimmitjoo/ecom/src/application/services"
	"github.com/jimmitjoo/ecom/src/infrastructure/events/memory"
	"github.com/jimmitjoo/ecom/src/infrastructure/handlers"
	memoryRepo "github.com/jimmitjoo/ecom/src/infrastructure/repositories/memory"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Create repository instance
	repo := memoryRepo.NewProductRepository()

	// Create event publisher
	publisher := memory.NewMemoryEventPublisher()

	// Create service with repository and publisher
	service := services.NewProductService(repo, publisher)

	// Create handlers
	productHandler := handlers.NewProductHandler(service)
	wsHandler := handlers.NewWebSocketHandler(publisher)

	// Set up router
	r := mux.NewRouter()

	// Batch endpoints (must come before specific product endpoints)
	r.HandleFunc("/products/batch", productHandler.BatchCreateProducts).Methods("POST")
	r.HandleFunc("/products/batch", productHandler.BatchUpdateProducts).Methods("PUT")
	r.HandleFunc("/products/batch", productHandler.BatchDeleteProducts).Methods("DELETE")

	// REST endpoints for individual products
	r.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	r.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// WebSocket endpoint
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// CORS configuration
	corsMiddleware := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Use CORS middleware
	handler := corsMiddleware(r)

	log.Printf("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

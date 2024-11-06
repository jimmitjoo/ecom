package main

import (
	"log"
	"net/http"

	"github.com/jimmitjoo/ecom/src/application/services"
	"github.com/jimmitjoo/ecom/src/infrastructure/events/memory"
	"github.com/jimmitjoo/ecom/src/infrastructure/handlers"
	"github.com/jimmitjoo/ecom/src/infrastructure/locks"
	"github.com/jimmitjoo/ecom/src/infrastructure/middleware"
	"github.com/jimmitjoo/ecom/src/infrastructure/ratelimit"
	memoryRepo "github.com/jimmitjoo/ecom/src/infrastructure/repositories/memory"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jimmitjoo/ecom/docs" // Detta genereras av swag
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title E-commerce Product API
// @version 1.0
// @description A robust and scalable API for product management in e-commerce systems
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http ws

func main() {
	// Create repository instance
	repo := memoryRepo.NewProductRepository()

	// Create event publisher
	publisher := memory.NewMemoryEventPublisher()

	// Create lock manager
	lockManager := locks.NewMemoryLockManager()

	// Create product service
	productService := services.NewProductService(repo, publisher, lockManager)

	// Create handlers
	productHandler := handlers.NewProductHandler(productService)
	wsHandler := handlers.NewWebSocketHandler(publisher)

	// Set up router
	r := mux.NewRouter()

	// Set up rate limiter
	limiter := ratelimit.NewTokenBucketLimiter(10, 10) // 10 tokens/sec, max 10 tokens
	rateLimitMiddleware := middleware.RateLimitMiddleware(limiter)
	r.Use(rateLimitMiddleware)

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

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// CORS configuration
	corsMiddleware := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD",
		}),
		gorillaHandlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
			"Origin",
			"Accept",
		}),
		gorillaHandlers.ExposedHeaders([]string{
			"Content-Length",
			"Access-Control-Allow-Origin",
		}),
		gorillaHandlers.AllowCredentials(),
	)

	// Use CORS middleware
	handler := corsMiddleware(r)

	log.Printf("Repository initialized: %v", repo != nil)
	log.Printf("Publisher initialized: %v", publisher != nil)
	log.Printf("LockManager initialized: %v", lockManager != nil)
	log.Printf("ProductService initialized: %v", productService != nil)
	log.Printf("ProductHandler initialized: %v", productHandler != nil)

	log.Printf("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

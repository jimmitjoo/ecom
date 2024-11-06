package load

import (
	"testing"

	"github.com/jimmitjoo/ecom/src/application/interfaces"
	"github.com/jimmitjoo/ecom/src/application/services"
	eventmem "github.com/jimmitjoo/ecom/src/infrastructure/events/memory"
	"github.com/jimmitjoo/ecom/src/infrastructure/locks"
	"github.com/jimmitjoo/ecom/src/infrastructure/repositories/memory"
	"github.com/jimmitjoo/ecom/src/testing/generators"
)

func setupBenchmarkService() (interfaces.ProductService, error) {
	publisher := eventmem.NewMemoryEventPublisher()
	repository := memory.NewProductRepository()
	lockManager := locks.NewMemoryLockManager()
	return services.NewProductService(repository, publisher, lockManager), nil
}

func BenchmarkBatchOperations(b *testing.B) {
	scenarios := []struct {
		name      string
		batchSize int
	}{
		{"SmallBatch", 10},
		{"MediumBatch", 100},
		{"LargeBatch", 1000},
	}

	generator := generators.NewProductGenerator()
	service, err := setupBenchmarkService()
	if err != nil {
		b.Fatal(err)
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			products := generator.GenerateProducts(sc.batchSize)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				results, err := service.BatchCreateProducts(products)
				if err != nil {
					b.Fatal(err)
				}

				for _, result := range results {
					if !result.Success {
						b.Fatalf("Failed to create product: %s", result.Error)
					}
				}
			}
		})
	}
}

func BenchmarkParallelBatchOperations(b *testing.B) {
	scenarios := []struct {
		name      string
		batchSize int
	}{
		{"ParallelSmallBatch", 10},
		{"ParallelMediumBatch", 100},
		{"ParallelLargeBatch", 1000},
	}

	generator := generators.NewProductGenerator()
	service, err := setupBenchmarkService()
	if err != nil {
		b.Fatal(err)
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			products := generator.GenerateProducts(sc.batchSize)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					results, err := service.BatchCreateProducts(products)
					if err != nil {
						b.Fatal(err)
					}

					for _, result := range results {
						if !result.Success {
							b.Fatalf("Failed to create product: %s", result.Error)
						}
					}
				}
			})
		})
	}
}

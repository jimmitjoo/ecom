package chaos

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

// ChaosConfig contains configuration for chaos tests
type ChaosConfig struct {
	NetworkLatency  time.Duration
	PacketLossRate  float64
	MemoryPressure  bool
	CorruptDataRate float64
}

// simulateMemoryPressure simulates memory pressure
func simulateMemoryPressure(t *testing.T) {
	var memoryHog [][]byte
	defer func() {
		memoryHog = nil
		runtime.GC()
	}()

	// Allocate memory until we reach 80% usage
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		if float64(memStats.Alloc)/float64(memStats.Sys) > 0.8 {
			break
		}
		memoryHog = append(memoryHog, make([]byte, 1024*1024)) // 1MB at a time
	}
}

// DataCorruption simulates corrupt data
type DataCorruption struct {
	rate float64
}

func (dc *DataCorruption) CorruptData(data []byte) []byte {
	if dc.rate <= 0 {
		return data
	}

	corrupted := make([]byte, len(data))
	copy(corrupted, data)

	for i := range corrupted {
		if rand.Float64() < dc.rate {
			corrupted[i] = byte(rand.Intn(256))
		}
	}
	return corrupted
}

func TestMemoryPressure(t *testing.T) {
	// Set up product service with many concurrent calls
	service, _ := setupTestService()

	t.Run("Service Under Memory Pressure", func(t *testing.T) {
		var wg sync.WaitGroup

		// Start memory pressure in the background
		go simulateMemoryPressure(t)

		// Make many concurrent calls
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				products := generateLargeProductBatch(100)
				_, err := service.BatchCreateProducts(products)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}

func TestDataCorruption(t *testing.T) {
	// Increase corruption rate to ensure data is actually corrupted
	corruption := &DataCorruption{rate: 0.5}

	t.Run("Handle Corrupt Data", func(t *testing.T) {
		product := generateValidProduct()
		data, _ := json.Marshal(product)

		// Corrupt data multiple times if needed
		var corrupted models.Product
		var err error
		maxAttempts := 5

		for i := 0; i < maxAttempts; i++ {
			corruptData := corruption.CorruptData(data)
			err = json.Unmarshal(corruptData, &corrupted)

			// If we got an unmarshal error or data is corrupted, break
			if err != nil || corrupted.ID != product.ID {
				break
			}
		}

		// Verify that we either got an error or corrupted data
		if err != nil {
			assert.True(t,
				strings.Contains(err.Error(), "invalid character") ||
					strings.Contains(err.Error(), "json") ||
					strings.Contains(err.Error(), "unmarshal"),
				"Expected JSON parsing error, got: %v", err)
		} else {
			assert.NotEqual(t, product.ID, corrupted.ID,
				"Data should be corrupted but was still valid")
			assert.NotEqual(t, product.BaseTitle, corrupted.BaseTitle,
				"Expected corrupted data to be different")
		}
	})

	t.Run("Guaranteed Corruption", func(t *testing.T) {
		// Test with 100% corruption rate
		guaranteedCorruption := &DataCorruption{rate: 1.0}

		product := generateValidProduct()
		data, _ := json.Marshal(product)

		corruptData := guaranteedCorruption.CorruptData(data)
		var corrupted models.Product
		err := json.Unmarshal(corruptData, &corrupted)

		// With 100% corruption, we should always get an error
		assert.Error(t, err, "Expected error with 100% corruption rate")
	})
}

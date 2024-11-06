package chaos

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

// ChaosConfig innehåller konfiguration för chaos tester
type ChaosConfig struct {
	NetworkLatency  time.Duration
	PacketLossRate  float64
	MemoryPressure  bool
	CorruptDataRate float64
}

// NetworkChaos simulerar nätverksproblem
type NetworkChaos struct {
	latency  time.Duration
	lossRate float64
	mu       sync.Mutex
}

func NewNetworkChaos(latency time.Duration, lossRate float64) *NetworkChaos {
	return &NetworkChaos{
		latency:  latency,
		lossRate: lossRate,
	}
}

func (nc *NetworkChaos) WrapClient(client *http.Client) *http.Client {
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	client.Transport = &chaosTransport{
		base:  transport,
		chaos: nc,
	}
	return client
}

// MemoryPressure simulerar minnestryck
func simulateMemoryPressure(t *testing.T) {
	var memoryHog [][]byte
	defer func() {
		memoryHog = nil
		runtime.GC()
	}()

	// Allokera minne tills vi når 80% användning
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		if float64(memStats.Alloc)/float64(memStats.Sys) > 0.8 {
			break
		}
		memoryHog = append(memoryHog, make([]byte, 1024*1024)) // 1MB i taget
	}
}

// DataCorruption simulerar korrupt data
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

func TestNetworkChaos(t *testing.T) {
	chaos := NewNetworkChaos(100*time.Millisecond, 0.1)
	client := &http.Client{
		// Sätt en kort timeout för att inte hänga i testet
		Timeout: 1 * time.Second,
		// Använd ingen faktisk transport
		Transport: http.DefaultTransport,
	}
	chaosClient := chaos.WrapClient(client)

	t.Run("High Latency and Packet Loss", func(t *testing.T) {
		// Gör flera försök för att säkerställa att vi testar både latens och paketförlust
		for i := 0; i < 10; i++ {
			start := time.Now()
			resp, err := chaosClient.Get("http://example.com")
			duration := time.Since(start)

			if err != nil {
				assert.Contains(t, err.Error(), "simulerad paketförlust",
					"Fel som inte är paketförlust: %v", err)
			} else {
				assert.GreaterOrEqual(t, duration, 100*time.Millisecond,
					"Anropet tog mindre tid än förväntad latens")
				resp.Body.Close()
			}
		}
	})

	t.Run("Extreme Latency", func(t *testing.T) {
		extremeChaos := NewNetworkChaos(500*time.Millisecond, 0)
		extremeClient := extremeChaos.WrapClient(&http.Client{})

		start := time.Now()
		resp, err := extremeClient.Get("http://example.com")
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, duration, 500*time.Millisecond)
		resp.Body.Close()
	})

	t.Run("High Packet Loss", func(t *testing.T) {
		highLossChaos := NewNetworkChaos(0, 0.9)
		highLossClient := highLossChaos.WrapClient(&http.Client{})

		errorCount := 0
		attempts := 10

		for i := 0; i < attempts; i++ {
			_, err := highLossClient.Get("http://example.com")
			if err != nil {
				errorCount++
				assert.Contains(t, err.Error(), "simulerad paketförlust")
			}
		}

		assert.Greater(t, errorCount, attempts/2,
			"Förväntade mer än 50% paketförlust med 0.9 loss rate")
	})
}

func TestMemoryPressure(t *testing.T) {
	// Sätt upp produktservice med många samtidiga anrop
	service, _ := setupTestService()

	t.Run("Service Under Memory Pressure", func(t *testing.T) {
		var wg sync.WaitGroup

		// Starta minnestryck i bakgrunden
		go simulateMemoryPressure(t)

		// Gör många samtidiga anrop
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
	// Öka korruptionsfrekvensen för att säkerställa att data verkligen korrupteras
	corruption := &DataCorruption{rate: 0.5}

	t.Run("Handle Corrupt Data", func(t *testing.T) {
		product := generateValidProduct()
		data, _ := json.Marshal(product)

		// Korruptera data flera gånger om det behövs
		var corrupted models.Product
		var err error
		maxAttempts := 5

		for i := 0; i < maxAttempts; i++ {
			corruptData := corruption.CorruptData(data)
			err = json.Unmarshal(corruptData, &corrupted)

			// Om vi fick ett unmarshal-fel eller data är korrupt, break
			if err != nil || corrupted.ID != product.ID {
				break
			}
		}

		// Verifiera att vi antingen fick ett fel eller korrupt data
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
		// Testa med 100% korruptionsfrekvens
		guaranteedCorruption := &DataCorruption{rate: 1.0}

		product := generateValidProduct()
		data, _ := json.Marshal(product)

		corruptData := guaranteedCorruption.CorruptData(data)
		var corrupted models.Product
		err := json.Unmarshal(corruptData, &corrupted)

		// Med 100% korruption bör vi alltid få ett fel
		assert.Error(t, err, "Expected error with 100% corruption rate")
	})
}

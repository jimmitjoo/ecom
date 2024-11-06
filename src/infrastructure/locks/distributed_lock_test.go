package locks

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLockAcquisition(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resourceID := "test_resource"

	// Testa att låsa en resurs
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock on first attempt")

	// Testa att försöka låsa samma resurs igen
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.False(t, acquired, "Should not acquire lock when resource is already locked")

	// Vänta tills låset går ut
	time.Sleep(time.Second * 2)

	// Testa att låsa resursen igen efter att den gått ut
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock after previous lock expired")
}

func TestLockRelease(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resourceID := "test_resource"

	// Låsa resursen
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Släpp låset
	err = manager.ReleaseLock(resourceID)
	assert.NoError(t, err)

	// Verifiera att resursen kan låsas igen direkt
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock immediately after release")
}

func TestConcurrentLockAccess(t *testing.T) {
	manager := NewMemoryLockManager()
	resourceID := "test_resource"
	numGoroutines := 10
	successCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			acquired, err := manager.AcquireLock(context.Background(), resourceID, time.Second)
			assert.NoError(t, err)
			if acquired {
				mu.Lock()
				successCount++
				mu.Unlock()
				time.Sleep(time.Millisecond * 100) // Simulera arbete
				err := manager.ReleaseLock(resourceID)
				assert.NoError(t, err)
			}
		}()
	}

	wg.Wait()
	assert.Greater(t, successCount, 0, "At least one goroutine should acquire the lock")
	assert.Less(t, successCount, numGoroutines, "Not all goroutines should acquire the lock simultaneously")
}

func TestLockExpiration(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resourceID := "test_resource"
	shortDuration := time.Millisecond * 100

	// Låsa med kort TTL
	acquired, err := manager.AcquireLock(ctx, resourceID, shortDuration)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Vänta på att låset ska gå ut
	time.Sleep(shortDuration * 2)

	// Verifiera att en annan process kan låsa resursen
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock after expiration")
}

func TestMultipleResources(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resource1 := "resource1"
	resource2 := "resource2"

	// Låsa första resursen
	acquired1, err := manager.AcquireLock(ctx, resource1, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired1)

	// Låsa andra resursen
	acquired2, err := manager.AcquireLock(ctx, resource2, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired2, "Should acquire lock for different resource")

	// Släpp båda låsen
	err = manager.ReleaseLock(resource1)
	assert.NoError(t, err)
	err = manager.ReleaseLock(resource2)
	assert.NoError(t, err)
}

func TestCleanupExpiredLocks(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	shortDuration := time.Millisecond * 100

	// Skapa flera lås med kort TTL
	for i := 0; i < 5; i++ {
		resourceID := fmt.Sprintf("resource_%d", i)
		acquired, err := manager.AcquireLock(ctx, resourceID, shortDuration)
		assert.NoError(t, err)
		assert.True(t, acquired)
	}

	// Vänta på att låsen ska gå ut och cleanup ska köras
	time.Sleep(shortDuration * 3)

	// Verifiera att låsen har städats bort
	manager.mu.RLock()
	numLocks := len(manager.locks)
	manager.mu.RUnlock()
	assert.Equal(t, 0, numLocks, "All expired locks should be cleaned up")
}

func TestContextCancellation(t *testing.T) {
	manager := NewMemoryLockManager()
	resourceID := "test_resource"

	// Skapa en context med cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Låsa resursen
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Avbryt context
	cancel()

	// Försök låsa resursen med avbruten context
	_, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.Error(t, err, "Should fail when context is cancelled")
}

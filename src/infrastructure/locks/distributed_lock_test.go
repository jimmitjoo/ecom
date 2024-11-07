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

	// Test acquiring a lock on a resource
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock on first attempt")

	// Test trying to lock the same resource again
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.False(t, acquired, "Should not acquire lock when resource is already locked")

	// Wait until the lock expires
	time.Sleep(time.Second * 2)

	// Test acquiring the resource again after it expired
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock after previous lock expired")
}

func TestLockRelease(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resourceID := "test_resource"

	// Acquire the lock on the resource
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Release the lock
	err = manager.ReleaseLock(resourceID)
	assert.NoError(t, err)

	// Verify that the resource can be locked again immediately
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
				time.Sleep(time.Millisecond * 100) // Simulate work
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

	// Acquire lock with short TTL
	acquired, err := manager.AcquireLock(ctx, resourceID, shortDuration)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Wait until the lock expires
	time.Sleep(shortDuration * 2)

	// Verify that another process can acquire the lock
	acquired, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired, "Should acquire lock after expiration")
}

func TestMultipleResources(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	resource1 := "resource1"
	resource2 := "resource2"

	// Acquire lock on the first resource
	acquired1, err := manager.AcquireLock(ctx, resource1, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired1)

	// Acquire lock on the second resource
	acquired2, err := manager.AcquireLock(ctx, resource2, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired2, "Should acquire lock for different resource")

	// Release both locks
	err = manager.ReleaseLock(resource1)
	assert.NoError(t, err)
	err = manager.ReleaseLock(resource2)
	assert.NoError(t, err)
}

func TestCleanupExpiredLocks(t *testing.T) {
	manager := NewMemoryLockManager()
	ctx := context.Background()
	shortDuration := time.Millisecond * 100

	// Create multiple locks with short TTL
	for i := 0; i < 5; i++ {
		resourceID := fmt.Sprintf("resource_%d", i)
		acquired, err := manager.AcquireLock(ctx, resourceID, shortDuration)
		assert.NoError(t, err)
		assert.True(t, acquired)
	}

	// Wait until the locks expire and cleanup runs
	time.Sleep(shortDuration * 3)

	// Verify that the locks have been cleaned up
	manager.mu.RLock()
	numLocks := len(manager.locks)
	manager.mu.RUnlock()
	assert.Equal(t, 0, numLocks, "All expired locks should be cleaned up")
}

func TestContextCancellation(t *testing.T) {
	manager := NewMemoryLockManager()
	resourceID := "test_resource"

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Acquire the lock on the resource
	acquired, err := manager.AcquireLock(ctx, resourceID, time.Second)
	assert.NoError(t, err)
	assert.True(t, acquired)

	// Cancel the context
	cancel()

	// Try acquiring the resource with a cancelled context
	_, err = manager.AcquireLock(ctx, resourceID, time.Second)
	assert.Error(t, err, "Should fail when context is cancelled")
}

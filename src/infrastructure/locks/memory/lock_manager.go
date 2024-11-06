package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemoryLockManager struct {
	locks sync.Map
	mu    sync.Mutex
}

func NewMemoryLockManager() *MemoryLockManager {
	return &MemoryLockManager{}
}

func (m *MemoryLockManager) AcquireLock(ctx context.Context, resourceID string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.locks.Load(resourceID); exists {
		return false, nil
	}

	m.locks.Store(resourceID, time.Now().Add(ttl))
	return true, nil
}

func (m *MemoryLockManager) ReleaseLock(resourceID string) error {
	m.locks.Delete(resourceID)
	return nil
}

func (m *MemoryLockManager) RefreshLock(resourceID string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.locks.Load(resourceID); !exists {
		return fmt.Errorf("lock does not exist")
	}

	m.locks.Store(resourceID, time.Now().Add(ttl))
	return nil
}

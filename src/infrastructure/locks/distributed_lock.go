package locks

import (
	"context"
	"sync"
	"time"
)

type LockManager interface {
	AcquireLock(ctx context.Context, resourceID string, ttl time.Duration) (bool, error)
	ReleaseLock(resourceID string) error
	RefreshLock(resourceID string, ttl time.Duration) error
}

type MemoryLockManager struct {
	locks    map[string]*Lock
	mu       sync.RWMutex
	cleanupC chan struct{}
}

type Lock struct {
	owner     string
	expiresAt time.Time
	mu        sync.Mutex
}

func NewMemoryLockManager() *MemoryLockManager {
	m := &MemoryLockManager{
		locks:    make(map[string]*Lock),
		cleanupC: make(chan struct{}),
	}
	go m.cleanupExpiredLocks()
	return m
}

func (m *MemoryLockManager) AcquireLock(ctx context.Context, resourceID string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if lock, exists := m.locks[resourceID]; exists {
		if now.Before(lock.expiresAt) {
			return false, nil
		}
	}

	m.locks[resourceID] = &Lock{
		owner:     "owner", // In a distributed system, this would be a unique node ID
		expiresAt: now.Add(ttl),
	}
	return true, nil
}

func (m *MemoryLockManager) ReleaseLock(resourceID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.locks, resourceID)
	return nil
}

func (m *MemoryLockManager) cleanupExpiredLocks() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for id, lock := range m.locks {
				if now.After(lock.expiresAt) {
					delete(m.locks, id)
				}
			}
			m.mu.Unlock()
		case <-m.cleanupC:
			return
		}
	}
}

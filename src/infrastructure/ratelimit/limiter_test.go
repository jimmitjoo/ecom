package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketLimiter(t *testing.T) {
	limiter := NewTokenBucketLimiter(10, 10) // 10 tokens/sec, max 10 tokens
	key := "test-key"

	// Test initial allowance
	assert.True(t, limiter.Allow(key))

	// Test overloading
	for i := 0; i < 10; i++ {
		limiter.Allow(key)
	}
	assert.False(t, limiter.Allow(key), "Should deny after bucket is empty")

	// Test refill
	time.Sleep(time.Second)
	assert.True(t, limiter.Allow(key), "Should allow after refill")

	// Test reset
	limiter.Reset(key)
	assert.True(t, limiter.Allow(key), "Should allow after reset")
}

func TestSlidingWindowLimiter(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second) // 5 requests per second
	key := "test-key"

	// Test initial allowance
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(key), "Should allow initial requests")
	}

	// Test overloading
	assert.False(t, limiter.Allow(key), "Should deny after limit reached")

	// Wait until the window has moved
	time.Sleep(time.Second + 100*time.Millisecond) // Add some margin

	// Now all previous requests should have expired
	assert.True(t, limiter.Allow(key), "Should allow after window moved")

	// Test reset
	limiter.Reset(key)
	assert.True(t, limiter.Allow(key), "Should allow after reset")
}

func TestSlidingWindowGradualExpiry(t *testing.T) {
	limiter := NewSlidingWindowLimiter(3, 500*time.Millisecond)
	key := "test-key"

	// Fill the window
	assert.True(t, limiter.Allow(key))
	assert.True(t, limiter.Allow(key))
	assert.True(t, limiter.Allow(key))
	assert.False(t, limiter.Allow(key))

	// Wait until one request has expired
	time.Sleep(600 * time.Millisecond)

	// Now we should be able to make a new request
	assert.True(t, limiter.Allow(key), "Should allow after one request expired")
}

func TestConcurrentAccess(t *testing.T) {
	limiter := NewTokenBucketLimiter(100, 100)
	key := "concurrent-test"

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				limiter.Allow(key)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

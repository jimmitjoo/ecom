package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketLimiter(t *testing.T) {
	limiter := NewTokenBucketLimiter(10, 10) // 10 tokens/sec, max 10 tokens
	key := "test-key"

	// Testa initial tillåtelse
	assert.True(t, limiter.Allow(key))

	// Testa överbelastning
	for i := 0; i < 10; i++ {
		limiter.Allow(key)
	}
	assert.False(t, limiter.Allow(key), "Should deny after bucket is empty")

	// Testa påfyllning
	time.Sleep(time.Second)
	assert.True(t, limiter.Allow(key), "Should allow after refill")

	// Testa reset
	limiter.Reset(key)
	assert.True(t, limiter.Allow(key), "Should allow after reset")
}

func TestSlidingWindowLimiter(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second) // 5 requests per second
	key := "test-key"

	// Testa initial tillåtelse
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(key), "Should allow initial requests")
	}

	// Testa överbelastning
	assert.False(t, limiter.Allow(key), "Should deny after limit reached")

	// Vänta tills fönstret har flyttat sig
	time.Sleep(time.Second + 100*time.Millisecond) // Lägg till lite marginal

	// Nu borde alla tidigare requests ha förfallit
	assert.True(t, limiter.Allow(key), "Should allow after window moved")

	// Testa reset
	limiter.Reset(key)
	assert.True(t, limiter.Allow(key), "Should allow after reset")
}

func TestSlidingWindowGradualExpiry(t *testing.T) {
	limiter := NewSlidingWindowLimiter(3, 500*time.Millisecond)
	key := "test-key"

	// Fyll fönstret
	assert.True(t, limiter.Allow(key))
	assert.True(t, limiter.Allow(key))
	assert.True(t, limiter.Allow(key))
	assert.False(t, limiter.Allow(key))

	// Vänta tills en request bör ha förfallit
	time.Sleep(600 * time.Millisecond)

	// Nu borde vi kunna göra en ny request
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

	// Vänta på alla goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

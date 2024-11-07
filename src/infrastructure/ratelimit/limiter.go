package ratelimit

import (
	"sync"
	"time"
)

// RateLimiter interface defines methods for rate limiting
type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
}

// TokenBucketLimiter implements the token bucket algorithm
type TokenBucketLimiter struct {
	mu           sync.RWMutex
	tokens       map[string]float64
	lastRefill   map[string]time.Time
	rate         float64 // tokens per second
	capacity     float64 // max tokens
	refillPeriod time.Duration
}

func NewTokenBucketLimiter(rate float64, capacity float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		tokens:       make(map[string]float64),
		lastRefill:   make(map[string]time.Time),
		rate:         rate,
		capacity:     capacity,
		refillPeriod: time.Second,
	}
}

func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// Initialize key if it doesn't exist
	if _, exists := l.tokens[key]; !exists {
		l.tokens[key] = l.capacity
		l.lastRefill[key] = now
		return true
	}

	// Calculate tokens to refill
	elapsed := now.Sub(l.lastRefill[key])
	tokensToAdd := float64(elapsed) / float64(l.refillPeriod) * l.rate

	currentTokens := l.tokens[key] + tokensToAdd
	if currentTokens > l.capacity {
		currentTokens = l.capacity
	}

	// Check if we have enough tokens
	if currentTokens < 1 {
		return false
	}

	l.tokens[key] = currentTokens - 1
	l.lastRefill[key] = now
	return true
}

func (l *TokenBucketLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.tokens, key)
	delete(l.lastRefill, key)
}

// SlidingWindowLimiter implements the sliding window algorithm
type SlidingWindowLimiter struct {
	mu       sync.RWMutex
	windows  map[string][]time.Time
	limit    int
	duration time.Duration
}

func NewSlidingWindowLimiter(limit int, duration time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windows:  make(map[string][]time.Time),
		limit:    limit,
		duration: duration,
	}
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.duration)

	// Get or initialize the window
	window, exists := l.windows[key]
	if !exists {
		window = make([]time.Time, 0, l.limit)
	}

	// Filter out old timestamps
	validWindow := make([]time.Time, 0, len(window))
	for _, t := range window {
		if t.After(cutoff) {
			validWindow = append(validWindow, t)
		}
	}

	// Check if we are over the limit
	if len(validWindow) >= l.limit {
		l.windows[key] = validWindow // Save the updated window
		return false
	}

	// Add new timestamp and save
	validWindow = append(validWindow, now)
	l.windows[key] = validWindow

	return true
}

func (l *SlidingWindowLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.windows, key)
}

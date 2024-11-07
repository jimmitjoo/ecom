package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/infrastructure/ratelimit"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Create a mock handler that always responds with 200 OK
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		limit          int
		duration       time.Duration
		requests       int
		expectedStatus int
		description    string
	}{
		{
			name:           "Under the limit",
			limit:          5,
			duration:       time.Second,
			requests:       3,
			expectedStatus: http.StatusOK,
			description:    "Should allow requests under the limit",
		},
		{
			name:           "On the limit",
			limit:          5,
			duration:       time.Second,
			requests:       5,
			expectedStatus: http.StatusOK,
			description:    "Should allow requests up to the limit",
		},
		{
			name:           "Over the limit",
			limit:          5,
			duration:       time.Second,
			requests:       6,
			expectedStatus: http.StatusTooManyRequests,
			description:    "Should reject requests over the limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new limiter for each test
			limiter := ratelimit.NewSlidingWindowLimiter(tc.limit, tc.duration)
			middleware := RateLimitMiddleware(limiter)
			handler := middleware(nextHandler)

			// Simulate requests
			var lastStatus int
			for i := 0; i < tc.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:1234" // Simulate a client IP

				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				lastStatus = rec.Code
			}

			assert.Equal(t, tc.expectedStatus, lastStatus, tc.description)
		})
	}
}

func TestRateLimitMiddlewareWithDifferentIPs(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limiter := ratelimit.NewSlidingWindowLimiter(2, time.Second)
	middleware := RateLimitMiddleware(limiter)
	handler := middleware(nextHandler)

	// Test that different IP addresses have separate limits
	ips := []string{"192.168.1.1:1234", "192.168.1.2:1234"}

	for _, ip := range ips {
		// First request should succeed
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "First request for %s should succeed", ip)

		// Second request should succeed
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Second request for %s should succeed", ip)

		// Third request should be rejected
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTooManyRequests, rec.Code, "Third request for %s should be rejected", ip)
	}
}

func TestRateLimitMiddlewareReset(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limiter := ratelimit.NewSlidingWindowLimiter(2, 500*time.Millisecond)
	middleware := RateLimitMiddleware(limiter)
	handler := middleware(nextHandler)
	ip := "192.168.1.1:1234"

	// Make two requests (up to the limit)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Third request should be rejected
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Wait until the window has passed
	time.Sleep(600 * time.Millisecond)

	// Now a new request should succeed
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Request should succeed after the window has passed")
}

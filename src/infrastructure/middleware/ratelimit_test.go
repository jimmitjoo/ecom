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
	// Skapa en mock handler som alltid svarar med 200 OK
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
			name:           "Under gränsen",
			limit:          5,
			duration:       time.Second,
			requests:       3,
			expectedStatus: http.StatusOK,
			description:    "Bör tillåta requests under gränsen",
		},
		{
			name:           "På gränsen",
			limit:          5,
			duration:       time.Second,
			requests:       5,
			expectedStatus: http.StatusOK,
			description:    "Bör tillåta requests upp till gränsen",
		},
		{
			name:           "Över gränsen",
			limit:          5,
			duration:       time.Second,
			requests:       6,
			expectedStatus: http.StatusTooManyRequests,
			description:    "Bör neka requests över gränsen",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Skapa en ny limiter för varje test
			limiter := ratelimit.NewSlidingWindowLimiter(tc.limit, tc.duration)
			middleware := RateLimitMiddleware(limiter)
			handler := middleware(nextHandler)

			// Simulera requests
			var lastStatus int
			for i := 0; i < tc.requests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:1234" // Simulera en klient IP

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

	// Testa att olika IP-adresser har separata gränser
	ips := []string{"192.168.1.1:1234", "192.168.1.2:1234"}

	for _, ip := range ips {
		// Första requesten bör lyckas
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Första requesten för %s bör lyckas", ip)

		// Andra requesten bör lyckas
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Andra requesten för %s bör lyckas", ip)

		// Tredje requesten bör nekas
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTooManyRequests, rec.Code, "Tredje requesten för %s bör nekas", ip)
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

	// Gör två requests (upp till gränsen)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Tredje requesten bör nekas
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Vänta tills fönstret har passerat
	time.Sleep(600 * time.Millisecond)

	// Nu bör en ny request lyckas
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "Request bör lyckas efter att fönstret har passerat")
}

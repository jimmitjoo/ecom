package middleware

import (
	"net/http"

	"github.com/jimmitjoo/ecom/src/infrastructure/ratelimit"
)

func RateLimitMiddleware(limiter ratelimit.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use IP as the key for rate limiting
			key := r.RemoteAddr

			if !limiter.Allow(key) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter creates a new rate limiter middleware.
// For now, we use a simple global limiter for demonstration.
// In a real production environment, this should ideally be keyed by IP or API token.
func RateLimit(next http.Handler) http.Handler {
	// Allow 20 requests per second with a burst of 50.
	limiter := rate.NewLimiter(rate.Every(time.Second/20), 50)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

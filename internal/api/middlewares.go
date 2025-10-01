package api

import (
	"net/http"

	"golang.org/x/time/rate"
)

// rateLimitMiddleware: a simple global rate limiter for all incoming requests.
func rateLimitMiddleware(next http.Handler, limiter *rate.Limiter) http.Handler {

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})

	return httpHandler
}

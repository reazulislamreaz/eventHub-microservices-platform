package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimit applies per-IP token bucket rate limiting.
func RateLimit(rps float64, burst int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := make(map[string]*rate.Limiter)

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		if lim, ok := visitors[ip]; ok {
			return lim
		}
		lim := rate.NewLimiter(rate.Limit(rps), burst)
		visitors[ip] = lim
		return lim
	}

	go func() {
		for range time.Tick(time.Minute * 5) {
			mu.Lock()
			visitors = make(map[string]*rate.Limiter)
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = fwd
			}
			if !getLimiter(ip).Allow() {
				http.Error(w, `{"error":"Too Many Requests","code":429}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

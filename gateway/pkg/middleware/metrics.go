package middleware

import (
	"net/http"

	"github.com/eventhub/gateway/internal/handler/rest"
)

// Metrics counts requests for the Prometheus endpoint.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rest.IncRequestCount()
		next.ServeHTTP(w, r)
	})
}

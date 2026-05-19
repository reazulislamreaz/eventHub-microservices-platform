package rest

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var (
	requestsTotal uint64
)

// IncRequestCount increments the request counter (called from middleware).
func IncRequestCount() {
	atomic.AddUint64(&requestsTotal, 1)
}

// Metrics exposes Prometheus-compatible metrics for observability stacks.
func Metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	total := atomic.LoadUint64(&requestsTotal)
	_, _ = fmt.Fprintf(w, "# HELP eventhub_http_requests_total Total HTTP requests through gateway\n")
	_, _ = fmt.Fprintf(w, "# TYPE eventhub_http_requests_total counter\n")
	_, _ = fmt.Fprintf(w, "eventhub_http_requests_total %d\n", total)
	_, _ = fmt.Fprintf(w, "# HELP eventhub_up Gateway is running\n")
	_, _ = fmt.Fprintf(w, "# TYPE eventhub_up gauge\n")
	_, _ = fmt.Fprintf(w, "eventhub_up 1\n")
}

package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eventhub/pkg/grpcutil"
	"go.uber.org/zap"
)

// ReadyCheck reports dependency readiness for orchestrators.
type ReadyCheck struct {
	Status       string            `json:"status" example:"ready"`
	Service      string            `json:"service" example:"eventhub-gateway"`
	Dependencies map[string]string `json:"dependencies"`
	Timestamp    string            `json:"timestamp"`
}

// ReadyWithDeps checks downstream gRPC services before reporting ready.
func ReadyWithDeps(w http.ResponseWriter, r *http.Request, addrs map[string]string, log *zap.Logger) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := map[string]string{}
	allOK := true

	checks := []struct {
		name    string
		addr    string
		service string
	}{
		{"user-service", addrs["user"], "user.v1.UserService"},
		{"event-service", addrs["event"], "event.v1.EventService"},
		{"ticket-service", addrs["ticket"], "ticket.v1.TicketService"},
	}

	for _, c := range checks {
		if err := grpcutil.WaitForService(ctx, c.addr, c.service, nil, 1); err != nil {
			status[c.name] = "unavailable"
			allOK = false
			if log != nil {
				log.Warn("dependency not ready", zap.String("service", c.name), zap.Error(err))
			}
		} else {
			status[c.name] = "ok"
		}
	}

	resp := ReadyCheck{
		Service:      "eventhub-gateway",
		Dependencies: status,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if !allOK {
		resp.Status = "not_ready"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		resp.Status = "ready"
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(resp)
}

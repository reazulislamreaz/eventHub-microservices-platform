package rest

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents service health status.
// @Description Health check response
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Service   string `json:"service" example:"eventhub-gateway"`
	Timestamp string `json:"timestamp" example:"2026-05-19T12:00:00Z"`
}

// Health godoc
// @Summary      Health check
// @Description  Returns gateway health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	resp := HealthResponse{
		Status:    "ok",
		Service:   "eventhub-gateway",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Ready godoc
// @Summary      Readiness check
// @Description  Returns gateway readiness status
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /ready [get]
func Ready(w http.ResponseWriter, _ *http.Request) {
	Health(w, nil)
}

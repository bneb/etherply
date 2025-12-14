// Package server provides HTTP and WebSocket handlers for the EtherPly sync engine.
package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// HealthChecker provides health and readiness check functionality.
type HealthChecker struct {
	store     store.Store
	startTime time.Time
}

// NewHealthChecker creates a new HealthChecker with the given dependencies.
func NewHealthChecker(s store.Store) *HealthChecker {
	return &HealthChecker{
		store:     s,
		startTime: time.Now(),
	}
}

// HealthResponse represents the JSON response for health endpoints.
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// HandleHealthz handles liveness probe requests.
// This endpoint returns 200 OK if the process is running.
// Kubernetes uses this to determine if the container should be restarted.
//
// Path: /healthz
func (h *HealthChecker) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).Round(time.Second).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleReadyz handles readiness probe requests.
// This endpoint returns 200 OK only if all dependencies are healthy.
// Kubernetes uses this to determine if the pod should receive traffic.
//
// Path: /readyz
func (h *HealthChecker) HandleReadyz(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	allHealthy := true

	// Check store connectivity
	if err := h.store.Ping(); err != nil {
		checks["store"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["store"] = "ok"
	}

	response := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	if allHealthy {
		response.Status = "ok"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		response.Status = "unhealthy"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

package handler

import (
	"encoding/json"
	"marketfuck/internal/application/port/in"
	"net/http"
)

type HealthHandler struct {
	healthService in.HealthService
}

func NewHealthHandler(healthService in.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.healthService == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status": "unhealthy",
			"error":  "health service not initialized",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	health, err := h.healthService.HealthCheck(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusOK
	if health.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

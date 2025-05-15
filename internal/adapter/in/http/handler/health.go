package handler

import (
	"net/http"

	"marketfuck/internal/application/port/in"
)

type HealthHandler struct {
	healthService in.HealthService
}

func NewHealthHandler(healthService in.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// обрабатывает запрос на проверку работоспособности системы
func (h *HealthHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	
}

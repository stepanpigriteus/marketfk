package handler

import (
	"fmt"
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
	"net/http"
	"strings"
)

type PriceHandler struct {
	priceService in.PriceService
	logger       port.Logger
}

func NewPriceHandler(priceService in.PriceService, logger port.Logger) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
		logger:       logger,
	}
}


func (h *PriceHandler) HandleGetLatestPrice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetLatestPrice handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	pairName := parts[2]
	fmt.Println(">>>>", pairName)
	latestPrice, err := h.priceService.GetLatestPrice(ctx, pairName)
	if err != nil {
		h.logger.Error("Incorrect GetAveragePrice result")
	}
	fmt.Println(">>>>", latestPrice)
}

// обрабатывает запрос на получение последней цены с конкретной биржи
func (h *PriceHandler) HandleGetLatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
}

// обрабатывает запрос на получение наивысшей цены за период
func (h *PriceHandler) HandleGetHighestPrice(w http.ResponseWriter, r *http.Request) {
}

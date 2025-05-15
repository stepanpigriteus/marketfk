package handler

import (
	"net/http"

	"marketfuck/internal/application/port/in"
)

type PriceHandler struct {
	priceService in.PriceService
}

func NewPriceHandler(priceService in.PriceService) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

// обрабатывает запрос на получение последней цены
func (h *PriceHandler) HandleGetLatestPrice(w http.ResponseWriter, r *http.Request) {
}

// обрабатывает запрос на получение последней цены с конкретной биржи
func (h *PriceHandler) HandleGetLatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
}

// обрабатывает запрос на получение наивысшей цены за период
func (h *PriceHandler) HandleGetHighestPrice(w http.ResponseWriter, r *http.Request) {
}

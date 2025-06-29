package handler

import (
	"encoding/json"
	"fmt"
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
	"net/http"
	"strings"
	"time"
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

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	pairName := parts[3]
	aggregatedPrice, err := h.priceService.GetLatestPrice(ctx, pairName)
	if err != nil {
		h.logger.Error("Incorrect GetLatestPrice result")
		http.Error(w, `{"error":"failed to fetch latest price"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(aggregatedPrice)
}

// обрабатывает запрос на получение последней цены с конкретной биржи
func (h *PriceHandler) HandleGetLatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetLatestPriceByExange handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	pairName := parts[4]
	exchangeID := parts[3]

	if exchangeID == "" || pairName == "" {
		http.Error(w, "invalid exchange or symbol", http.StatusBadRequest)
		return
	}

	latestPriceByEx, err := h.priceService.GetLatestPriceByExchange(ctx, exchangeID, pairName)
	if err != nil {
		h.logger.Error("Incorrect GetAveragePrice result")
	}
	jsonResponse, err := json.Marshal(latestPriceByEx)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *PriceHandler) HandleGetHighestPriceInPeriod(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetHighestPriceInPeriod handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	symbol := parts[3]
	if symbol == "" {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")

	if periodStr == "" {
		http.Error(w, `{"error":"missing period query parameter"}`, http.StatusBadRequest)
		return
	}

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
		return
	}

	highestPriceInPeriod, err := h.priceService.GetHighestPriceInPeriod(ctx, symbol, period)
	if err != nil {
		h.logger.Error("Incorrect GetAveragePrice result")
	}

	jsonResponse, err := json.Marshal(highestPriceInPeriod)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

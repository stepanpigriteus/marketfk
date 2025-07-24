package handler

import (
	"encoding/json"
	"fmt"
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
	"marketfuck/internal/domain/model"
	"marketfuck/pkg/utils"
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

func (h *PriceHandler) HandleGetHighestPrice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info((">>> GetHighestPrice handler called"))
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	var highestPrice model.AggregatedPrice

	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	pairName := parts[3]
	highestPrice, err := h.priceService.GetHighestPrice(ctx, pairName)
	if err != nil {
		h.logger.Error("Incorrect GetAveragePrice result")
	}
	jsonResponse, err := json.Marshal(highestPrice)
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
	exchange := ""
	w.Header().Set("Content-Type", "application/json")
	var highestPriceInPeriod model.AggregatedPrice
	flag := "max"
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
	formSymb, err := utils.PairNameValidFormatter(symbol)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")
	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)
	if periodStr == "" {
		highestPriceInPeriod, err = h.priceService.GetHighestPrice(ctx, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, period), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			highestPriceInPeriod, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchange, flag)
			fmt.Println(">>>>>cashe", highestPriceInPeriod)
		} else {
			highestPriceInPeriod, err = h.priceService.GetHighestPriceInPeriod(ctx, formSymb, period)
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}

	jsonResponse, err := json.Marshal(highestPriceInPeriod)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *PriceHandler) HandleGetHighestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetHighestPriceByExangeInPeriod handler called")
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
	formSymb, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}

	if exchangeID == "" || pairName == "" {
		http.Error(w, "invalid exchange or symbol", http.StatusBadRequest)
		return
	}
	var latestPriceByEx model.AggregatedPrice
	flag := "max"
	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")
	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)

	if periodStr == "" {
		latestPriceByEx, err = h.priceService.GetHighestPriceByExchange(ctx, exchangeID, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid GetHighestPriceByExchange result: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, period), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			latestPriceByEx, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchangeID, flag)
			fmt.Println(">>>>>cashe", latestPriceByEx)
		} else {
			latestPriceByEx, err = h.priceService.GetHighestPriceByExchangeInPeriod(ctx, exchangeID, formSymb, period)
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}

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

func (h *PriceHandler) HandleGetLowestPrice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info((">>> GetHighestPrice handler called"))
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	var lowestPrice model.AggregatedPrice
	flag := "min"
	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	symbol := parts[3]
	formSymb, err := utils.PairNameValidFormatter(symbol)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}
	exchange := ""
	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")

	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)
	if periodStr == "" {
		lowestPrice, err = h.priceService.GetLowestPrice(ctx, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid GetLowestPrice result: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, period), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, err), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			lowestPrice, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchange, flag)
			fmt.Println(">>>>>cashe", lowestPrice)
		} else {
			lowestPrice, err = h.priceService.GetLowestPriceInPeriod(ctx, formSymb, period)
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}
	// lovestPrice, err := h.priceService.GetLowestPrice(ctx, pairName)
	// if err != nil {
	// 	h.logger.Error("Incorrect GetAveragePrice result")
	// }

	jsonResponse, err := json.Marshal(lowestPrice)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *PriceHandler) HandleGetLowestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetHighestPriceByExangeInPeriod handler called")
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
	formSymb, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}

	if exchangeID == "" || pairName == "" {
		http.Error(w, "invalid exchange or symbol", http.StatusBadRequest)
		return
	}
	var latestPriceByEx model.AggregatedPrice
	flag := "min"
	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")
	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)

	if periodStr == "" {
		latestPriceByEx, err = h.priceService.GetLowestPriceByExchange(ctx, exchangeID, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid GetHighestPriceByExchange result: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, period), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			latestPriceByEx, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchangeID, flag)
			fmt.Println(">>>>>cashe", latestPriceByEx)
		} else {
			latestPriceByEx, err = h.priceService.GetLowestPriceByExchangeInPeriod(ctx, exchangeID, formSymb, period)
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}

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

func (h *PriceHandler) HandleGetAveragePrice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info((">>> GetAveragePrice handler called"))
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	var averagePrice model.AggregatedPrice
	flag := "min"
	if ctx == nil {
		http.Error(w, `{"error":"invalid context"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	symbol := parts[3]
	formSymb, err := utils.PairNameValidFormatter(symbol)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}
	exchange := ""
	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")

	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)
	if periodStr == "" {
		averagePrice, err = h.priceService.GetAveragePrice(ctx, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid GetLowestPrice result: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, period), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, err), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			averagePrice, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchange, flag)
			fmt.Println(">>>>>cashe", averagePrice)
		} else {
			averagePrice, err = h.priceService.GetAveragePriceInPeriod(ctx, formSymb, period) // sdfsdfsdf
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}

	jsonResponse, err := json.Marshal(averagePrice)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *PriceHandler) HandleGetAveragePriceByExchange(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetHighestPriceByExangeInPeriod handler called")
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
	formSymb, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		http.Error(w, "invalid symbol in path", http.StatusBadRequest)
		return
	}

	if exchangeID == "" || pairName == "" {
		http.Error(w, "invalid exchange or symbol", http.StatusBadRequest)
		return
	}
	var averagePriceByEx model.AggregatedPrice
	flag := "min"
	queryParams := r.URL.Query()
	periodStr := queryParams.Get("period")
	h.logger.Info("Processing request", "symbol", formSymb, "period", periodStr)

	if periodStr == "" {
		averagePriceByEx, err = h.priceService.GetAveragePriceByExchange(ctx, exchangeID, formSymb)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid GetHighestPriceByExchange result: %s"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		period, err := time.ParseDuration(periodStr)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid period: %s"}`, err), http.StatusBadRequest)
			return
		}
		if period <= 0 {
			http.Error(w, fmt.Sprintf(`{"error":"negative period in query %s"}`, period), http.StatusBadRequest)
			return
		}

		if period < time.Minute {
			averagePriceByEx, err = h.priceService.GetHighestPriceFromCache(ctx, formSymb, period, exchangeID, flag)
			fmt.Println(">>>>>cashe", averagePriceByEx)
		} else {
			averagePriceByEx, err = h.priceService.GetAveragePriceByExchangeInPeriod(ctx, exchangeID, formSymb, period)
			fmt.Println(">>>>>base")
		}

		if err != nil {
			h.logger.Error("Incorrect GetHightPriceFromPeriod result")
			return
		}
	}

	if err != nil {
		h.logger.Error("Incorrect GetAveragePrice result")
	}
	jsonResponse, err := json.Marshal(averagePriceByEx)
	if err != nil {
		http.Error(w, `{"error": "failed to serialize response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

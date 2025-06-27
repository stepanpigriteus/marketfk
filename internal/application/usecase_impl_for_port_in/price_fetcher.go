package usecase

import (
	"context"
	"fmt"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port/in"
	"marketfuck/internal/domain/model"
	"strings"
	"time"
)

type priceService struct {
	priceRepo postgres.PriceRepository
}

func NewPriceService(repo postgres.PriceRepository) in.PriceService {
	return &priceService{priceRepo: repo}
}

func (s *priceService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	validPairs := map[string]bool{
		"BTCUSDT":  true,
		"DOGEUSDT": true,
		"TONUSDT":  true,
		"SOLUSDT":  true,
		"ETHUSDT":  true,
	}
	pairName = strings.ToUpper(pairName)
	if len(pairName) == 0 || !validPairs[pairName] {
		return model.AggregatedPrice{}, fmt.Errorf("incorrect PairName")
	}

	price, err := s.priceRepo.GetLatestPrice(ctx, pairName)
	if err != nil {
		fmt.Printf("Error in GetLatestPrice: %v\n", err) // Логируем ошибку
		return model.AggregatedPrice{}, err
	}

	// Возвращаем результат
	return price, nil
}

func (s *priceService) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	validPairs := map[string]bool{
		"BTCUSDT":  true,
		"DOGEUSDT": true,
		"TONUSDT":  true,
		"SOLUSDT":  true,
		"ETHUSDT":  true,
	}
	pairName = strings.ToUpper(pairName)
	if len(pairName) == 0 || !validPairs[pairName] {
		return model.AggregatedPrice{}, fmt.Errorf("incorrect PairName")
	}

	price, err := s.priceRepo.GetLatestPriceByExchange(ctx, exchangeID, pairName)
	if err != nil {
		fmt.Printf("Error in GetLatestPrice: %v\n", err) // Логируем ошибк у
		return model.AggregatedPrice{}, err
	}

	return price, nil
}

func (s *priceService) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	validPairs := map[string]bool{
		"BTCUSDT":  true,
		"DOGEUSDT": true,
		"TONUSDT":  true,
		"SOLUSDT":  true,
		"ETHUSDT":  true,
	}
	pairName = strings.ToUpper(pairName)
	fmt.Println(pairName,period, ctx)
	if len(pairName) == 0 || !validPairs[pairName] {
		return model.AggregatedPrice{}, fmt.Errorf("incorrect PairName")
	}
	price, err := s.priceRepo.GetHighestPriceInPeriod(ctx, pairName, period)
	if err != nil {
		fmt.Printf("Error in GetHighestPriceInPeriod: %v\n", err) // Логируем ошибк у
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetHighestPrice(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	return price, nil
}

func (s *priceService) GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetLowestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetAveragePrice(ctx context.Context, pairName string, period time.Duration) (float64, error) {
	var avg float64
	return avg, nil
}

func (s *priceService) GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (float64, error) {
	var avg float64
	return avg, nil
}

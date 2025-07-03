package service

import (
	"context"
	"errors"
	"marketfuck/internal/application/port/out"
	"marketfuck/internal/domain/model"
	"time"
)

type MarketService struct {
	aggregatedRepo out.PriceRepository
}

func NewMarketService(repo out.PriceRepository) *MarketService {
	return &MarketService{aggregatedRepo: repo}
}

func (s *MarketService) SavePrice(ctx context.Context, price []model.AggregatedPrice) error {
	if len(price) == 0 {
		return errors.New("price lenght equal zero")
	}
	err := s.aggregatedRepo.SavePrice(ctx, price)
	if err != nil {
		return err
	}

	return nil
}

func (s *MarketService) GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return s.aggregatedRepo.GetPricesInPeriodByExchange(ctx, exchangeID, pairName, startTime, endTime)
}

func (s *MarketService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLatestPrice(ctx, pairName)
}

func (s *MarketService) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetHighestPriceInPeriod(ctx, pairName, period)
}

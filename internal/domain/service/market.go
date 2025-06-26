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

func (s *MarketService) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error) {
	return s.aggregatedRepo.GetLatestPriceByExchange(ctx, exchangeID, pairName)
}

func (s *MarketService) GetPricesInPeriod(ctx context.Context, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return s.aggregatedRepo.GetPricesInPeriod(ctx, pairName, startTime, endTime)
}

func (s *MarketService) GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return s.aggregatedRepo.GetPricesInPeriodByExchange(ctx, exchangeID, pairName, startTime, endTime)
}

func (s *MarketService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLatestPrice(ctx, pairName)
}

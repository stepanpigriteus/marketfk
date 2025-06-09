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

func (s *MarketService) SaveAggregatedPricesBatch(ctx context.Context, prices []model.AggregatedPrice) error {
	// добавить расчет средних и т.д.
	// сериализовать в структуру
	if batchRepo, ok := s.aggregatedRepo.(interface {
		SaveAggregatedPricesBatch(context.Context, []model.AggregatedPrice) error
	}); ok {
		return batchRepo.SaveAggregatedPricesBatch(ctx, prices)
	}
	return errors.New("repository does not support batch insert")
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

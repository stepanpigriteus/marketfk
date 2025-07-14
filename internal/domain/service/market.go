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
		return errors.New("price length equal zero")
	}
	err := s.aggregatedRepo.SavePrice(ctx, price)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarketService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLatestPrice(ctx, pairName)
}

func (s *MarketService) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLatestPriceByExchange(ctx, exchangeID, pairName)
}

func (s *MarketService) GetHighestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetHighestPrice(ctx, pairName)
}

func (s *MarketService) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetHighestPriceInPeriod(ctx, pairName, period)
}

func (s *MarketService) GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetHighestPriceByExchange(ctx, exchangeID, pairName)
}

func (s *MarketService) GetHighestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetHighestPriceByExchangeInPeriod(ctx, exchangeID, pairName, period)
}

func (s *MarketService) GetLowestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLovestPrice(ctx, pairName)
}

func (s *MarketService) GetLowestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLowestPriceInPeriod(ctx, pairName, period)
}

func (s *MarketService) GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLowestPriceByExchange(ctx, exchangeID, pairName)
}

func (s *MarketService) GetLowestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetLowestPriceByExchangeInPeriod(ctx, exchangeID, pairName, period)
}

func (s *MarketService) GetAveragePrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetAveragePrice(ctx, pairName)
}

func (s *MarketService) GetAveragePriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetAveragePriceInPeriod(ctx, pairName, period)
}

func (s *MarketService) GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetAveragePriceByExchange(ctx, exchangeID, pairName)
}

func (s *MarketService) GetAveragePriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	return s.aggregatedRepo.GetAveragePriceByExchangeInPeriod(ctx, exchangeID, pairName, period)
}

func (s *MarketService) GetLastRecordTime(ctx context.Context, pairName string) (time.Time, error) {
	return s.aggregatedRepo.GetLastRecordTime(ctx, pairName)
}

// func (s *MarketService) GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
// 	return s.aggregatedRepo.GetPricesInPeriodByExchange(ctx, exchangeID, pairName, startTime, endTime)
// }

// func (s *MarketService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
// 	return s.aggregatedRepo.GetLatestPrice(ctx, pairName)
// }

// func (s *MarketService) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
// 	return s.aggregatedRepo.GetHighestPriceInPeriod(ctx, pairName, period)
// }

// func (s *MarketService) GetHighestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
// 	return s.aggregatedRepo.GetHighestPrice(ctx, pairName)
// }

package usecase

import (
	"context"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port/in"
	"marketfuck/internal/domain/model"
	"time"
)

type priceService struct {
	priceRepo postgres.PriceRepository
}

func NewPriceService(repo postgres.PriceRepository) in.PriceService {
	return &priceService{priceRepo: repo}
}

func (s *priceService) GetLatestPrice(ctx context.Context, pairName string) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetHighestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
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

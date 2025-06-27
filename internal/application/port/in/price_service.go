package in

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

type PriceService interface {
	GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error)
	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error)
	GetHighestPrice(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)
	GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)
	GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error)
	GetLowestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error)
	GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error)
	GetAveragePrice(ctx context.Context, pairName string, period time.Duration) (float64, error)
	GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (float64, error)
}

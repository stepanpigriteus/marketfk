package in

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

type PriceService interface {
	// next
	GetHighestPriceFromCache(ctx context.Context, pairName string, period time.Duration, exchange string) (model.AggregatedPrice, error)

	//+
	GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) // GET /prices/latest/{symbol}

	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error)
	GetHighestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error)
	GetHighestPriceByExchange(ctx context.Context, exchangeID string, pairName string) (model.AggregatedPrice, error)
	GetHighestPriceByExchangeInPeriod(ctx context.Context, exchangeID string, pairName string, period time.Duration) (model.AggregatedPrice, error)
	GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)
	GetLowestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error)
	//
	//-

}

// type PriceService interface {
// 	GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error)

// 	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error)

// 	GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)

// 	GetHighestPriceInPeriodByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error)

// 	GetLowestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)

// 	GetLowestPriceInPeriodByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error)

// 	GetAveragePriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error)

// 	GetAveragePriceInPeriodByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error)
// }

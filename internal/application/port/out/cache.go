package out

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

type Cache interface {
	SavePrice(ctx context.Context, price model.Price) error
	GetLatestPrice(ctx context.Context, pairName string) (model.Price, error)
	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error)
	GetPricesForPeriod(ctx context.Context, pairName, exchangeID string, period time.Duration) ([]model.Price, error)
	ClearOldPrices(ctx context.Context, pairName, exchangeID string, before time.Time) error
	Health(ctx context.Context) error
}

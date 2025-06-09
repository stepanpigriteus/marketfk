package out

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

type PriceRepository interface {
	SavePrice(ctx context.Context, price []model.AggregatedPrice) error
	// изменить модель
	GetLatestPrice(ctx context.Context, pairName string) (model.Price, error)
	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error)
	GetPricesInPeriod(ctx context.Context, pairName string, startTime, endTime time.Time) ([]model.Price, error)
	GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error)
}

type ModeRepository interface {
	SetMode(ctx context.Context, mode string) error
	GetMode(ctx context.Context) (string, error)
}

type HealthRepository interface {
	CheckConnection(ctx context.Context) (bool, error)
}

package in

import (
	"context"
	"time"
	"marketfuck/internal/domain/model"
)

type PriceService interface {
	GetLatestPrice(ctx context.Context, pairName string) (model.Price, error)
	GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error)
	GetHighestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error)
	GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error)
	GetLowestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error)
	GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error)
	GetAveragePrice(ctx context.Context, pairName string, period time.Duration) (float64, error)
	GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (float64, error)
}

type ModeService interface {
	SwitchToTestMode(ctx context.Context) error
	SwitchToLiveMode(ctx context.Context) error
}


type HealthService interface {
	HealthCheck(ctx context.Context) (SystemHealth, error)
}


type SystemHealth struct {
	Status         string            `json:"status"`
	Connections    map[string]string `json:"connections"` 
	RedisActive    bool              `json:"redis_active"`
	PostgresActive bool              `json:"postgres_active"`
}

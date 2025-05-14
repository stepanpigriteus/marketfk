package out

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

type Storage interface {
	SaveAggregatedData(ctx context.Context, data AggregatedData) error
	GetPriceHistory(ctx context.Context, pairName, exchangeID string, period time.Duration) ([]AggregatedData, error)
	Health(ctx context.Context) error
}

type AggregatedData struct {
	PairName     string         `json:"pair_name"`
	Exchange     model.Exchange `json:"exchange"`
	Timestamp    time.Time      `json:"timestamp"`
	AveragePrice float64        `json:"average_price"`
	MinPrice     float64        `json:"min_price"`
	MaxPrice     float64        `json:"max_price"`
}

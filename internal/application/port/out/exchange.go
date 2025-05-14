package out

import (
	"context"
	"marketfuck/internal/domain/model"
)

type ExchangeClient interface {
	Connect(ctx context.Context) error
	Subscribe(ctx context.Context, pairs []model.Pair) (<-chan model.Price, error)
	Health(ctx context.Context) error
	Close(ctx context.Context) error
}

package in

import (
	"context"
	"marketfuck/internal/domain/model"
)


type ExchangeDataService interface {
	Subscribe(ctx context.Context, pairs []model.Pair) (<-chan model.Price, error)
	Health(ctx context.Context) error
}

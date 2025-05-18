package out

import (
	"context"
	"time"

	"marketfuck/internal/domain/model"
)

type CacheClient interface {
	// сохраняет цену в кеш
	SetPrice(ctx context.Context, key string, price model.Price, expiration time.Duration) error
	// получает цену из кеша
	GetPrice(ctx context.Context, key string) (model.Price, bool, error)
	// сохраняет режим в кеш
	SetMode(ctx context.Context, mode string) error
	// получает режим из кеша
	GetMode(ctx context.Context) (string, bool, error)
	// проверяет соединение с кешем
	CheckConnection(ctx context.Context) (bool, error)

	Close() error
}

package out

import (
	"context"
	"marketfuck/internal/domain/model"
	"time"
)

// Базовый интерфейс кеша
type CacheClient interface {
	SetPrice(ctx context.Context, key string, price model.Price, expiration time.Duration) error
	GetPrice(ctx context.Context, key string) (model.Price, bool, error)
	SetMode(ctx context.Context, mode string) error
	GetMode(ctx context.Context) (string, bool, error)
	CheckConnection(ctx context.Context) (bool, error)
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Close() error
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}



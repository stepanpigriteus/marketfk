package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"marketfuck/internal/domain/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	return &RedisCache{client: rdb}, nil
}

func (r *RedisCache) SetPrice(ctx context.Context, key string, price model.Price, expiration time.Duration) error {
	data, err := json.Marshal(price)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) GetPrice(ctx context.Context, key string) (model.Price, bool, error) {
	var price model.Price
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return price, false, nil
	}
	if err != nil {
		return price, false, err
	}

	err = json.Unmarshal([]byte(data), &price)
	return price, true, err
}

func (r *RedisCache) SetMode(ctx context.Context, mode string) error {
	return r.client.Set(ctx, "mode", mode, 0).Err()
}

func (r *RedisCache) GetMode(ctx context.Context) (string, bool, error) {
	mode, err := r.client.Get(ctx, "mode").Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return mode, err == nil, err
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.Scan(ctx, cursor, match, count).Result()
}

func (r *RedisCache) CheckConnection(ctx context.Context) (bool, error) {
	err := r.client.Ping(ctx).Err()
	return err == nil, err
}

func (r *RedisCache) Clear(ctx context.Context) error {
	keys, err := r.Keys(ctx, "*")
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := r.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

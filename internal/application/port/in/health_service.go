package in

import "context"

type HealthService interface {
	HealthCheck(ctx context.Context) (SystemHealth, error)
}

type SystemHealth struct {
	Status         string            `json:"status"`
	Connections    map[string]string `json:"connections"`
	RedisActive    bool              `json:"redis_active"`
	PostgresActive bool              `json:"postgres_active"`
}

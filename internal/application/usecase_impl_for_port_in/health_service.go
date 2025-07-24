package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"marketfuck/internal/application/port/in"
	"marketfuck/internal/application/port/out"
	"net"
	"sync"
	"time"
)

type HealthServiceImpl struct {
	cacheClient out.CacheClient
	db          *sql.DB
	exchanges   []ExchangeConfig // конфигурация бирж
}

type ExchangeConfig struct {
	Name string
	Host string
	Port string
}

func NewHealthService(cacheClient out.CacheClient, db *sql.DB, exchanges []ExchangeConfig) in.HealthService {
	return &HealthServiceImpl{
		cacheClient: cacheClient,
		db:          db,
		exchanges:   exchanges,
	}
}

func (h *HealthServiceImpl) HealthCheck(ctx context.Context) (in.SystemHealth, error) {
	health := in.SystemHealth{
		Status:      "healthy",
		Connections: make(map[string]string),
	}

	// Используем WaitGroup для параллельной проверки всех компонентов
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Проверка Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		isActive, err := h.checkRedis(ctx)
		mu.Lock()
		defer mu.Unlock()

		health.RedisActive = isActive
		if err != nil {
			health.Connections["redis"] = fmt.Sprintf("error: %v", err)
			health.Status = "unhealthy"
		} else {
			health.Connections["redis"] = "ok"
		}
	}()

	// Проверка PostgreSQL
	wg.Add(1)
	go func() {
		defer wg.Done()
		isActive, err := h.checkPostgres(ctx)
		mu.Lock()
		defer mu.Unlock()

		health.PostgresActive = isActive
		if err != nil {
			health.Connections["postgres"] = fmt.Sprintf("error: %v", err)
			health.Status = "unhealthy"
		} else {
			health.Connections["postgres"] = "ok"
		}
	}()

	// Проверка бирж
	for i, exchange := range h.exchanges {
		wg.Add(1)
		go func(idx int, ex ExchangeConfig) {
			defer wg.Done()
			err := h.checkExchange(ctx, ex)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				health.Connections[ex.Name] = fmt.Sprintf("error: %v", err)
				health.Status = "unhealthy"
			} else {
				health.Connections[ex.Name] = "ok"
			}
		}(i, exchange)
	}

	wg.Wait()

	return health, nil
}

func (h *HealthServiceImpl) checkRedis(ctx context.Context) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return h.cacheClient.CheckConnection(ctxTimeout)
}

func (h *HealthServiceImpl) checkPostgres(ctx context.Context) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := h.db.PingContext(ctxTimeout)
	if err != nil {
		return false, err
	}

	// Дополнительная проверка - выполним простой запрос
	var result int
	err = h.db.QueryRowContext(ctxTimeout, "SELECT 1").Scan(&result)
	if err != nil {
		return false, fmt.Errorf("query test failed: %w", err)
	}

	return true, nil
}

func (h *HealthServiceImpl) checkExchange(ctx context.Context, exchange ExchangeConfig) error {
	address := fmt.Sprintf("%s:%s", exchange.Host, exchange.Port)

	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("TCP connection failed: %w", err)
	}
	defer conn.Close()

	return nil
}

package runner

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/internal/domain/model"
	"marketfuck/internal/domain/service"
	"marketfuck/pkg/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetupSignalHandler() chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	return sigCh
}

func InitDependencies(cfg *config.Config, logger port.Logger) (*sql.DB, *redis.RedisCache, *service.MarketService) {
	logger.Info("[2/4] Connecting to the DB")
	time.Sleep(4 * time.Second)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)

	db, err := postgres.ConnectDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	logger.Info("[3/4] Connecting to Redis")
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	redisClient, err := redis.NewRedisCache(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Невозможно подключиться к Redis: %v", err)
	}

	repo := postgres.NewPriceRepository(db)
	marketService := service.NewMarketService(repo)

	return db, redisClient, marketService
}

func RunPriceSaver(ctx context.Context, redisClient *redis.RedisCache, marketService *service.MarketService, logger port.Logger) {
	var delay int64 = 10
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()
	var defRedisData []model.AggregatedPrice

	for {
		select {
		case <-ctx.Done():
			logger.Info("Выход из цикла сохранения цен по сигналу завершения.")
			return
		case <-ticker.C:
			redisData, threshold := usecase.GetAllPrices(redisClient, delay)
			// fmt.Println(redisData)
			if len(defRedisData) == 0 {
				if err := marketService.SavePrice(ctx, redisData); err != nil {
					defRedisData = append(defRedisData, redisData...)
					logger.Error("Ошибка при сохранении данных: ", err)
				} else {
					logger.Info("Данные успешно сохранены в базу.")
				}
			} else {
				defRedisData = append(defRedisData, redisData...)
				if err := marketService.SavePrice(ctx, defRedisData); err != nil {
					logger.Error("Ошибка при сохранении данных: ", err)
				} else {
					logger.Info("Данные успешно сохранены в базу.")
					defRedisData = nil
				}
			}

			if err := usecase.CleanupOldPrices(redisClient, threshold, delay); err != nil {
				logger.Error("Ошибка при удалении старых данных: ", err)
			} else {
				logger.Info("Старые данные успешно удалены из Redis.")
			}
		}
	}
}

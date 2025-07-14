package main

import (
	"context"
	"fmt"
	"log"
	"marketfuck/internal/adapter/in/http"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/internal/domain/model"
	"marketfuck/internal/domain/service"
	"marketfuck/pkg/concurrency"
	"marketfuck/pkg/config"
	"marketfuck/pkg/logger"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	var counter atomic.Uint64
	logger := logger.NewSlogAdapter()
	logger.Info("[1/4] Reading configurations")
	cfg := config.LoadConfig()
	ctx := context.Background()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	logger.Info("[2/4] An attempt to connect to the DB")
	time.Sleep(4 * time.Second)
	db, err := postgres.ConnectDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	logger.Info("[3/4] An attempt to connect to the Redis")
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	redisClient, err := redis.NewRedisCache(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Невозможно подключиться к Redis: %v", err)
	}
	defer redisClient.Close()

	repo := postgres.NewPriceRepository(db)
	marketService := service.NewMarketService(repo)


	
	server := http.NewServer("8081", db, logger, redisClient)
	logger.Info("[4/4] Time to run server!")
	go server.RunServer()

	fanIn := concurrency.GenAggr(&counter, *redisClient)
	go usecase.PriceAggregator(redisClient, fanIn)
	var delay int64 = 60
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()
	var defRedisData []model.AggregatedPrice

	for {
		select {
		case <-ticker.C:
			redisData, threshold := usecase.GetAllPrices(redisClient, delay)
			// fmt.Println(redisData)

			if len(defRedisData) == 0 {
				if err := marketService.SavePrice(ctx, redisData); err != nil {
					defRedisData = append(defRedisData, redisData...)
					logger.Error("Ошибка при сохранении данных: ", err)
				} else {
					fmt.Println("Данные успешно сохранены в базу.")
				}
			} else {
				defRedisData = append(defRedisData, redisData...)
				if err := marketService.SavePrice(ctx, defRedisData); err != nil {
					logger.Error("Ошибка при сохранении данных: ", err)
				} else {
					fmt.Println("Данные успешно сохранены в базу.")
					defRedisData = nil
				}
			}

			// После сохранения данных очищаем старые записи из Redis
			if err := usecase.CleanupOldPrices(redisClient, threshold, delay); err != nil {
				logger.Error("Ошибка при удалении старых данных: ", err)
			} else {
				fmt.Println("Старые данные из редис успешно удалены.")
			}
		}
	}

	// Закрываем канал цен - это приведет к завершению FanOut

	logger.Info("Все операции завершены")
}

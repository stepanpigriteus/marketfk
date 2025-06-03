package main

import (
	"fmt"
	"log"
	"marketfuck/internal/adapter/in/http"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
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

	server := http.NewServer("8081", db, logger)
	logger.Info("[4/4] Time to run server!")
	go server.RunServer()

	// нужно возвращать данные из агрегации и передавать в функцию управления редисом и базой
	concurrency.GenAggr(&counter, *redisClient)

	// Закрываем канал цен - это приведет к завершению FanOut

	logger.Info("Все операции завершены")
}

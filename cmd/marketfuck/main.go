package main

import (
	"fmt"
	"log"
	"marketfuck/internal/adapter/in/exchange/live"
	"marketfuck/internal/adapter/in/http"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/domain/model"
	"marketfuck/pkg/concurrency"
	"marketfuck/pkg/config"
	"marketfuck/pkg/logger"
	"sync"
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

	/// переписать под мапу
	ports := []string{"40101", "40102", "40103"}
	var wg sync.WaitGroup
	priceCh1 := make(chan model.Price, 100)
	priceCh2 := make(chan model.Price, 100)
	priceCh3 := make(chan model.Price, 100)

	priceChannels := [3]chan model.Price{priceCh1, priceCh2, priceCh3}
	const workerCount int = 50
	outChannels := [workerCount]chan model.Price{}
	for i := 0; i < workerCount; i++ {
		outChannels[i] = make(chan model.Price)
	}
	// Каналы для каждого воркера
	workerChans := make([]chan model.Price, workerCount)
	for i := 0; i < workerCount; i++ {
		workerChans[i] = make(chan model.Price, 100)
	}

	// Запуск воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		worker := concurrency.NewWorker(i, workerChans[i], &wg, &counter, outChannels[i])
		go worker.Run()
	}

	// Запуск Fan-Out для распределения данных по воркерам
	for _, chann := range priceChannels {
		go concurrency.FanOut(chann, workerChans)
	}

	// Подключение к биржам и получение данных
	for i, port := range ports {
		wg.Add(1)
		go live.GenConnectAndRead(port, &wg, priceChannels[i])
	}

	fanInChan := make(chan model.Price, workerCount*300)

	concurrency.FanIn(workerCount, fanInChan, outChannels[:], &wg)

	go func() {
		wg.Wait()
		close(fanInChan)
	}()

	for price := range fanInChan {
		fmt.Println("Получена цена:", price, counter.Load())
	}

	wg.Wait()

	// Закрываем канал цен - это приведет к завершению FanOut
	close(priceCh1)
	close(priceCh2)
	close(priceCh3)

	logger.Info("Все операции завершены")
}

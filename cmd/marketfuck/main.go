package main

import (
	"context"
	"marketfuck/internal/adapter/in/http"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/pkg/concurrency"
	"marketfuck/pkg/config"
	"marketfuck/pkg/logger"
	"marketfuck/pkg/runner"
	"sync"
	"sync/atomic"

	_ "github.com/lib/pq"
)

func main() {
	var counter atomic.Uint64
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logger.NewSlogAdapter()
	logger.Info("[1/4] Reading configurations")
	cfg := config.LoadConfig()
	ready := make(chan struct{})
	sigCh := runner.SetupSignalHandler()

	db, redisClient, marketService := runner.InitDependencies(cfg, logger)
	defer db.Close()
	defer redisClient.Close()

	server := http.NewServer("8081", db, logger, redisClient)
	wg := &sync.WaitGroup{}

	logger.Info("[4/4] Starting HTTP server")
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.RunServer(); err != nil {
			logger.Error("Ошибка запуска сервера: ", err)
			cancel()
		}
	}()

	go func() {
		<-ready
		<-sigCh
		logger.Info("Получен сигнал завершения. Остановка...")
		cancel()
		server.GracefulShutdown()
	}()

	fanIn := concurrency.GenAggr(ctx, &counter, *redisClient)
	close(ready)
	wg.Add(1)
	go func() {
		defer wg.Done()
		usecase.PriceAggregator(redisClient, fanIn)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runner.RunPriceSaver(ctx, redisClient, marketService, logger)
	}()

	<-sigCh
	logger.Info("Получен сигнал завершения. Остановка...")
	cancel()
	server.GracefulShutdown()
	wg.Wait()
	logger.Info("Все операции завершены.")
}

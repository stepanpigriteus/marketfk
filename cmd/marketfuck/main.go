package main

import (
	"context"
	"marketfuck/cmd/testgen"
	"marketfuck/internal/adapter/in/http"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/internal/domain/model"
	"marketfuck/internal/domain/service"
	"marketfuck/pkg/concurrency"
	"marketfuck/pkg/config"
	"marketfuck/pkg/logger"
	"marketfuck/pkg/runner"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	var counter atomic.Uint64
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logger.NewSlogAdapter()
	logger.Info("[1/4] Reading configurations")
	cfg := config.LoadConfig()

	go func() {
		testgen.StartFakeExchangeWithCtx(ctx, "50201", "exchangeT1")
	}()
	go func() {
		testgen.StartFakeExchangeWithCtx(ctx, "50202", "exchangeT2")
	}()
	go func() {
		testgen.StartFakeExchangeWithCtx(ctx, "50203", "exchangeT3")
	}()
	time.Sleep(2 * time.Second)
	sigCh := runner.SetupSignalHandler()

	db, redisClient, marketService := runner.InitDependencies(cfg, logger)
	defer db.Close()
	defer redisClient.Close()

	exchanges := []usecase.ExchangeConfig{
		{Name: "exchangeT1", Host: "localhost", Port: "50201"},
		{Name: "exchangeT2", Host: "localhost", Port: "50202"},
		{Name: "exchangeT3", Host: "localhost", Port: "50203"},
		{Name: "exchange1", Host: "localhost", Port: "40201"},
		{Name: "exchange2", Host: "localhost", Port: "40202"},
		{Name: "exchange3", Host: "localhost", Port: "40203"},
	}

	healthService := usecase.NewHealthService(redisClient, db, exchanges)

	aggregator := func(ctx context.Context, counter *atomic.Uint64, redis redis.RedisCache, ports []string, host string) <-chan model.Price {
		return concurrency.GenAggr(ctx, counter, redis, ports, host)
	}

	modeService := service.NewModeService(redisClient, &counter, aggregator)

	if err := modeService.SwitchToLiveMode(ctx); err != nil && err.Error() != "already in this mode" {
		logger.Error("Не удалось запустить в live режиме", "error", err)
		cancel()
		return
	}
	logger.Info("Successfully switched to LiveMode")

	server := http.NewServer("8081", db, logger, redisClient, modeService, healthService)

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

package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"marketfuck/internal/adapter/in/http/handler"
	"marketfuck/internal/adapter/in/http/router"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/pkg/errors"
	"net/http"
	"os"
	"time"
)

type server struct {
	port        string
	db          *sql.DB
	logger      port.Logger
	services    *in.AllServices
	redisClient *redis.RedisCache
	srv         *http.Server
}

func NewServer(port string, db *sql.DB, logger port.Logger, redisClient *redis.RedisCache) *server {
	priceRepo := postgres.NewPriceRepository(db)
	priceService := usecase.NewPriceService(*priceRepo, redisClient)
	services := &in.AllServices{
		PriceService: priceService,
	}

	return &server{
		port:        port,
		db:          db,
		logger:      logger,
		services:    services,
		redisClient: redisClient,
	}
}

func (s *server) RunServer() error {
	if s.port == "" {
		s.logger.Error("Port is not set")
		os.Exit(1)
	}
	mux := http.NewServeMux()
	handlers := handler.NewAllHandlers(
		s.services.HealthService,
		s.services.ModeService,
		s.services.PriceService,
		s.logger,
	)

	router.RegisterRoutes(mux, handlers)
	mux.Handle("/", &handleDef{})

	srv := &http.Server{
		Addr:         "0.0.0.0:" + s.port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	s.srv = srv

	s.logger.Info("Starting server", "port", s.port)

	// Не вызываем os.Exit — возвращаем ошибку выше
	return s.srv.ListenAndServe()
}

type handleDef struct{}

func (h *handleDef) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusInternalServerError
	if r.Method == "OPTIONS" {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)

	response := errors.Error{
		Message: "Undefined Error, please check your method or endpoint correctness",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

func (s *server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error("Ошибка при остановке сервера: ", err)
	} else {
		s.logger.Info("Сервер завершён корректно.")
	}
}

package http

import (
	"database/sql"
	"encoding/json"
	"marketfuck/internal/adapter/in/http/handler"
	"marketfuck/internal/adapter/in/http/router"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port"
	"marketfuck/internal/application/port/in"
	usecase "marketfuck/internal/application/usecase_impl_for_port_in"
	"marketfuck/pkg/errors"
	"net/http"
	"os"
)

type server struct {
	port     string
	db       *sql.DB
	logger   port.Logger
	services *in.AllServices
}

func NewServer(port string, db *sql.DB, logger port.Logger) *server {
	priceRepo := postgres.NewPriceRepository(db)
	priceService := usecase.NewPriceService(*priceRepo)
	services := &in.AllServices{
		PriceService: priceService,
	}

	return &server{
		port:     port,
		db:       db,
		logger:   logger,
		services: services,
	}
}

func (s *server) RunServer() {
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

	s.logger.Info("Starting server", "port", s.port)

	err := http.ListenAndServe("0.0.0.0:"+s.port, mux)
	if err != nil {
		s.logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
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

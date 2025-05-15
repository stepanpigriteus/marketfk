package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"marketfuck/internal/application/port"
	"marketfuck/pkg/errors"
)

type server struct {
	port   string
	db     *sql.DB
	logger port.Logger
	// services *service.AllServices
}

func NewServer(port string, db *sql.DB, logger port.Logger) *server {
	return &server{
		port:   port,
		db:     db,
		logger: logger,
		// services: services,
	}
}

type handleDef struct{}

func RunServer(s *server) {
	if s.port == "" {
		s.logger.Error("Port is not set")
		os.Exit(1)
	}
	mux := http.NewServeMux()
	fmt.Println(mux)
}

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

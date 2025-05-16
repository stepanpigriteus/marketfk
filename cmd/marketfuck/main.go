package main

import (
	"fmt"
	"log"
	"marketfuck/internal/adapter/in/http"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/pkg/config"
	"marketfuck/pkg/logger"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.NewSlogAdapter()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
	
	time.Sleep(6 * time.Second)
	db, err := postgres.ConnectDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	server := http.NewServer("8081", db, logger)
	server.RunServer()
}

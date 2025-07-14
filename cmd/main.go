package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/izymalhaw/go-crud/yishakterefe/docs"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/api/handlers"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/config"
	customlogger "github.com/izymalhaw/go-crud/yishakterefe/internal/core/logger"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/repository"
	person_service "github.com/izymalhaw/go-crud/yishakterefe/internal/services/person"
)

const (
	version = "1.0.0"
)

func main() {
	// Load from .env file before reading ENV vars
	_ = godotenv.Load()

	// Load config with proper error logging
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Initialize logger
	logger := customlogger.NewLogger(cfg.Env, cfg.LogLevel, version)

	// Initialize PostgreSQL database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify database connection
	if err = db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected successfully")

	// Initialize repository and service
	store := repository.NewPostgresUserRepo(db)
	personService := person_service.NewPersonSvc(store)

	// Start server
	webSrv := handlers.NewApp(cfg.Port, personService, logger)
	logger.Info("server running")
	webSrv.Run()
}

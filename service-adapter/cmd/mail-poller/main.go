package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimustaev/service-adapter/internal/adapter/mailhog"
	"github.com/aimustaev/service-adapter/internal/config"
	"github.com/aimustaev/service-adapter/internal/db"
	"github.com/aimustaev/service-adapter/internal/service"
	"github.com/aimustaev/service-adapter/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Infof("Starting mail-poller with configuration: %+v", cfg)

	// Initialize database
	database, err := db.NewDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Mailhog adapter
	adapter := mailhog.NewMailhogAdapter(cfg.MailhogHost, cfg.MailhogPort, log, database)

	// Create service
	svc := service.NewService(adapter, log)

	// Create context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Start service
	if err := svc.Start(ctx); err != nil {
		log.Fatalf("Service error: %v", err)
	}
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimustaev/service-communications/internal/adapter/mailhog"
	"github.com/aimustaev/service-communications/internal/adapter/telegram"
	"github.com/aimustaev/service-communications/internal/config"
	"github.com/aimustaev/service-communications/internal/db"
	"github.com/aimustaev/service-communications/internal/gateway"
	"github.com/aimustaev/service-communications/internal/kafka"
	"github.com/aimustaev/service-communications/internal/service"
	"github.com/aimustaev/service-communications/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Infof("Starting service-communications with configuration: %+v", cfg)

	// Initialize database
	database, err := db.NewDB(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize gateway client
	gatewayClient, err := gateway.NewClient(cfg.GetGatewayAddress())
	if err != nil {
		log.Fatalf("Failed to create gateway client: %v", err)
	}
	defer gatewayClient.Close()

	// Initialize Kafka producer
	producer, err := kafka.NewProducer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic, log)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize adapters
	mailhogAdapter := mailhog.NewMailhogAdapter(cfg.MailhogHost, cfg.MailhogPort, log, database)
	telegramAdapter := telegram.NewTelegramAdapter(cfg.TelegramBotToken, cfg.TelegramChatID, log, database)

	// Create services for each adapter
	mailhogService := service.NewService(mailhogAdapter, log, gatewayClient, producer)
	telegramService := service.NewService(telegramAdapter, log, gatewayClient, producer)

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

	// Start services in separate goroutines
	go func() {
		if err := mailhogService.Start(ctx); err != nil {
			log.Errorf("Mailhog service error: %v", err)
		}
	}()

	go func() {
		if err := telegramService.Start(ctx); err != nil {
			log.Errorf("Telegram service error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
}

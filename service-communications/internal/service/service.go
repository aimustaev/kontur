package service

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aimustaev/service-communications/internal/adapter"
	"github.com/aimustaev/service-communications/internal/gateway"
	"github.com/aimustaev/service-communications/internal/kafka"
)

// Service represents the main service for processing messages
type Service struct {
	adapter  adapter.Adapter
	logger   *logrus.Logger
	server   *http.Server
	gateway  *gateway.Client
	producer *kafka.Producer
}

// NewService creates a new service instance
func NewService(adapter adapter.Adapter, logger *logrus.Logger, gatewayClient *gateway.Client, producer *kafka.Producer) *Service {
	return &Service{
		adapter:  adapter,
		logger:   logger,
		gateway:  gatewayClient,
		producer: producer,
	}
}

// Start begins processing messages
func (s *Service) Start(ctx context.Context) error {
	if err := s.adapter.Connect(ctx); err != nil {
		return err
	}
	defer s.adapter.Disconnect(ctx)

	// Start HTTP server for health checks
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("HTTP server error: %v", err)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Shutdown HTTP server
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.server.Shutdown(shutdownCtx); err != nil {
				s.logger.Errorf("HTTP server shutdown error: %v", err)
			}
			return nil
		case <-ticker.C:
			if err := s.processMessages(ctx); err != nil {
				s.logger.Errorf("Error processing messages: %v", err)
			}
		}
	}
}

// processMessages retrieves and processes messages
func (s *Service) processMessages(ctx context.Context) error {
	messages, err := s.adapter.GetMessages(ctx)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		s.logger.Infof("Processing message: %s", msg.ID)
		s.logger.Infof("From: %s", msg.From)
		s.logger.Infof("To: %s", msg.To)
		s.logger.Infof("Subject: %s", msg.Subject)
		s.logger.Infof("Body: %s", msg.Body)
		s.logger.Info("---")

		// Send message to Kafka
		if err := s.producer.SendMessage(ctx, msg); err != nil {
			s.logger.Errorf("Failed to send message to Kafka: %v", err)
			continue
		}

		//// Отправляем сообщение в gateway сервис
		//if err := s.gateway.SendMessage(ctx, msg.ID, msg.Body, msg.From); err != nil {
		//	s.logger.Errorf("Failed to send message to gateway: %v", err)
		//	continue
		//}

		if err := s.adapter.MarkAsProcessed(ctx, msg.ID); err != nil {
			s.logger.Errorf("Failed to mark message as processed: %v", err)
		}
	}

	return nil
}

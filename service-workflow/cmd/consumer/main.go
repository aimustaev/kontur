package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"

	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/kafka"
	"github.com/aimustaev/service-workflow/internal/temporal"
	"github.com/aimustaev/service-workflow/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Temporal client configuration
	temporalConfig := temporal.DefaultConfig()
	temporalConfig.HostPort = cfg.GetTemporalAddr()
	temporalConfig.Namespace = cfg.Temporal.Namespace

	// Create Temporal client
	temporalClient, err := temporal.NewClient(temporalConfig)
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer temporalClient.Close()
	log.Println("Temporal client initialized successfully")

	// Create usecases
	startWorkflowUseCase := usecase.NewStartWorkflowUseCase(temporalClient.GetClient())
	startV2WorkflowUseCase := usecase.NewStartV2WorkflowUseCase(temporalClient.GetClient())

	// Create message handler
	messageHandler := usecase.NewMessageHandler(startWorkflowUseCase, startV2WorkflowUseCase)

	// Create Kafka client
	kafkaClient, err := kafka.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()

	// Create context that will be canceled on termination signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages
	go func() {
		err = kafkaClient.Consume(ctx, func(msg *sarama.ConsumerMessage) {
			// Parse the message into Message struct
			var message usecase.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return
			}

			// Handle the message
			if err := messageHandler.HandleMessage(ctx, message); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	log.Println("Waiting for messages...")
	// Wait for termination signal
	<-signals
	log.Println("Shutting down...")
}

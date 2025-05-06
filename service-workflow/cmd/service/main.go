package main

import (
	"log"
	"net/http"

	"github.com/aimustaev/service-workflow/internal/api"
	"github.com/aimustaev/service-workflow/internal/config"
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

	// Create handlers
	startHandler := api.NewStartWorkflowHandler(startWorkflowUseCase)

	// Create HTTP server
	http.HandleFunc("/start", startHandler.Handle)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start HTTP server
	log.Printf("Starting HTTP server on %s\n", cfg.GetHTTPAddr())
	if err := http.ListenAndServe(cfg.GetHTTPAddr(), nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

package main

import (
	"log"

	"go.temporal.io/sdk/worker"

	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/temporal"
	"github.com/aimustaev/service-workflow/internal/ticket"
	"github.com/aimustaev/service-workflow/internal/workflow"
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

	// Create ticket service client
	ticketClient, err := ticket.NewClient(cfg)
	if err != nil {
		log.Fatalln("Unable to create ticket service client", err)
	}
	defer ticketClient.Close()

	// Create worker
	log.Println("Creating worker...")
	w := worker.New(temporalClient.GetClient(), "workflow-ticket", worker.Options{})

	// Register workflows
	log.Println("Registering workflows...")
	workflow.RegisterWorkflows(w, ticketClient.GetClient(), temporalClient.GetClient())

	// Start worker
	log.Println("Starting worker...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

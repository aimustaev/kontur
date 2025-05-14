package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/worker"

	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/temporal"
	"github.com/aimustaev/service-workflow/internal/ticket"
	"github.com/aimustaev/service-workflow/internal/workflow"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize PostgreSQL connection
	db, err := sqlx.Connect("postgres", cfg.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("PostgreSQL connection established successfully")

	// Initialize config repository
	configRepo := config.NewPostgresConfigRepository(db)
	log.Println("Config repository initialized successfully")

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
	workflow.RegisterWorkflows(w, ticketClient.GetClient(), temporalClient.GetClient(), configRepo)

	// Start worker
	log.Println("Starting worker...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

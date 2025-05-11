package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/aimustaev/service-tickets/internal/config"
	"github.com/aimustaev/service-tickets/internal/repository"
	"github.com/aimustaev/service-tickets/internal/service"
	"github.com/aimustaev/service-tickets/proto"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize repository
	repo, err := repository.NewPostgresTicketRepository(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	ticketService := service.NewTicketService(repo)
	proto.RegisterTicketServiceServer(s, ticketService)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

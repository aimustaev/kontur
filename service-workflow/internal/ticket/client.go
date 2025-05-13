package ticket

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/aimustaev/service-workflow/internal/config"
	"github.com/aimustaev/service-workflow/internal/generated/proto"
)

// Client represents a ticket service client
type Client struct {
	client proto.TicketServiceClient
	conn   *grpc.ClientConn
}

// NewClient creates a new ticket service client
func NewClient(cfg *config.Config) (*Client, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(
		cfg.GetTicketServiceAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ticket service: %w", err)
	}

	// Create client
	client := proto.NewTicketServiceClient(conn)

	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateTicket creates a new ticket
func (c *Client) CreateTicket(ctx context.Context, request *proto.CreateTicketRequest) (*proto.TicketResponse, error) {
	log.Printf("Creating ticket with request: %+v", request)
	response, err := c.client.CreateTicket(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	return response, nil
}

// CreateTicket creates a new ticket
func (c *Client) UpdateTicket(ctx context.Context, request *proto.UpdateTicketRequest) (*proto.TicketResponse, error) {
	log.Printf("Update ticket with request: %+v", request)
	response, err := c.client.UpdateTicket(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	return response, nil
}

func (c *Client) GetClient() proto.TicketServiceClient {
	return c.client
}

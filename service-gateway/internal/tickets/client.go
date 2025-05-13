package tickets

import (
	"context"
	"fmt"
	"os"

	pb "github.com/aimustaev/service-gateway/internal/generated/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.TicketServiceClient
}

func NewClient() (*Client, error) {
	host := os.Getenv("TICKET_SERVICE_HOST")
	port := os.Getenv("TICKET_SERVICE_PORT")

	if host == "" || port == "" {
		return nil, fmt.Errorf("TICKET_SERVICE_HOST and TICKET_SERVICE_PORT must be set")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tickets service: %w", err)
	}

	return &Client{
		client: pb.NewTicketServiceClient(conn),
	}, nil
}

func (c *Client) GetAllTickets(ctx context.Context) (*pb.GetAllTicketsResponse, error) {
	return c.client.GetAllTickets(ctx, &pb.GetAllTicketsRequest{})
}

func (c *Client) GetTicketMessages(ctx context.Context, ticketID string) (*pb.GetTicketMessagesResponse, error) {
	return c.client.GetTicketMessages(ctx, &pb.GetTicketMessagesRequest{
		TicketId: ticketID,
	})
}

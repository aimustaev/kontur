package gateway

import (
	"context"
	"fmt"
	"time"

	pb "github.com/aimustaev/service-communications/internal/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.GatewayServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewGatewayServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendMessage(ctx context.Context, messageID, content, sender string) error {
	req := &pb.NewMessageRequest{
		MessageId: messageID,
		Content:   content,
		Sender:    sender,
		Timestamp: time.Now().Unix(),
	}

	stream, err := c.client.HandleNewMessage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Читаем ответы из стрима
	for {
		resp, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("error receiving response: %w", err)
		}

		// Обрабатываем статус
		switch resp.Status {
		case "success":
			// Продолжаем читать ответы
			continue
		case "processing":
			// Продолжаем читать ответы
			continue
		case "completed":
			// Успешное завершение
			return nil
		default:
			// Неизвестный статус считаем ошибкой
			return fmt.Errorf("gateway returned unknown status: %s", resp.Status)
		}
	}
}

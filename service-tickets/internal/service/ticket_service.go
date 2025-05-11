package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aimustaev/service-tickets/internal/model"
	"github.com/aimustaev/service-tickets/internal/repository"
	"github.com/aimustaev/service-tickets/proto"
)

// TicketService implements the gRPC ticket service
type TicketService struct {
	proto.UnimplementedTicketServiceServer
	repo repository.TicketRepository
}

// NewTicketService creates a new instance of TicketService
func NewTicketService(repo repository.TicketRepository) *TicketService {
	return &TicketService{
		repo: repo,
	}
}

// CreateTicket implements the CreateTicket RPC method
func (s *TicketService) CreateTicket(ctx context.Context, req *proto.CreateTicketRequest) (*proto.TicketResponse, error) {
	ticket := &model.Ticket{
		ID:         generateID(), // You'll need to implement this function
		VerticalID: req.VerticalId,
		UserID:     req.UserId,
		Assign:     req.Assign,
		SkillID:    req.SkillId,
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, status.Error(codes.Internal, "failed to create ticket")
	}

	return convertToProtoTicket(ticket), nil
}

// GetTicket implements the GetTicket RPC method
func (s *TicketService) GetTicket(ctx context.Context, req *proto.GetTicketRequest) (*proto.TicketResponse, error) {
	ticket, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to get ticket")
	}

	return convertToProtoTicket(ticket), nil
}

// UpdateTicket implements the UpdateTicket RPC method
func (s *TicketService) UpdateTicket(ctx context.Context, req *proto.UpdateTicketRequest) (*proto.TicketResponse, error) {
	ticket := &model.Ticket{
		ID:         req.Id,
		VerticalID: req.VerticalId,
		UserID:     req.UserId,
		Assign:     req.Assign,
		SkillID:    req.SkillId,
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to update ticket")
	}

	// Get updated ticket
	updatedTicket, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated ticket")
	}

	return convertToProtoTicket(updatedTicket), nil
}

// DeleteTicket implements the DeleteTicket RPC method
func (s *TicketService) DeleteTicket(ctx context.Context, req *proto.DeleteTicketRequest) (*proto.DeleteTicketResponse, error) {
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete ticket")
	}

	return &proto.DeleteTicketResponse{Success: true}, nil
}

// Helper function to convert model.Ticket to proto.TicketResponse
func convertToProtoTicket(ticket *model.Ticket) *proto.TicketResponse {
	return &proto.TicketResponse{
		Id:         ticket.ID,
		VerticalId: ticket.VerticalID,
		UserId:     ticket.UserID,
		Assign:     ticket.Assign,
		SkillId:    ticket.SkillID,
		CreatedAt:  ticket.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  ticket.UpdatedAt.Format(time.RFC3339),
	}
}

// Helper function to generate a unique ID
func generateID() string {
	// You can implement your own ID generation logic here
	// For example, using UUID or any other unique identifier
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// Helper function to generate a random string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

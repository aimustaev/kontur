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
	ticketRepo  repository.TicketRepository
	messageRepo repository.MessageRepository
}

// NewTicketService creates a new instance of TicketService
func NewTicketService(ticketRepo repository.TicketRepository, messageRepo repository.MessageRepository) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		messageRepo: messageRepo,
	}
}

// CreateTicket implements the CreateTicket RPC method
func (s *TicketService) CreateTicket(ctx context.Context, req *proto.CreateTicketRequest) (*proto.TicketResponse, error) {
	ticket := &model.Ticket{
		ID:          generateID(),
		Status:      req.Status,
		User:        req.User,
		Agent:       stringPtr(req.Agent),
		ProblemID:   int64Ptr(req.ProblemId),
		VerticalID:  int64Ptr(req.VerticalId),
		SkillID:     int64Ptr(req.SkillId),
		UserGroupID: int64Ptr(req.UserGroupId),
		Channel:     req.Channel,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, status.Error(codes.Internal, "failed to create ticket")
	}

	return convertToProtoTicket(ticket), nil
}

// GetTicket implements the GetTicket RPC method
func (s *TicketService) GetTicket(ctx context.Context, req *proto.GetTicketRequest) (*proto.TicketResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, req.Id)
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
		ID:          req.Id,
		Status:      req.Status,
		User:        req.User,
		Agent:       stringPtr(req.Agent),
		ProblemID:   int64Ptr(req.ProblemId),
		VerticalID:  int64Ptr(req.VerticalId),
		SkillID:     int64Ptr(req.SkillId),
		UserGroupID: int64Ptr(req.UserGroupId),
		Channel:     req.Channel,
		UpdatedAt:   time.Now(),
	}

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to update ticket")
	}

	// Get updated ticket
	updatedTicket, err := s.ticketRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated ticket")
	}

	return convertToProtoTicket(updatedTicket), nil
}

// DeleteTicket implements the DeleteTicket RPC method
func (s *TicketService) DeleteTicket(ctx context.Context, req *proto.DeleteTicketRequest) (*proto.DeleteTicketResponse, error) {
	if err := s.ticketRepo.Delete(ctx, req.Id); err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete ticket")
	}

	return &proto.DeleteTicketResponse{Success: true}, nil
}

// GetActiveTicketsByUser implements the GetActiveTicketsByUser RPC method
func (s *TicketService) GetActiveTicketsByUser(ctx context.Context, req *proto.GetActiveTicketsByUserRequest) (*proto.GetActiveTicketsByUserResponse, error) {
	tickets, err := s.ticketRepo.GetActiveByUser(ctx, req.User)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get active tickets")
	}

	response := &proto.GetActiveTicketsByUserResponse{
		Tickets: make([]*proto.TicketResponse, len(tickets)),
	}

	for i, ticket := range tickets {
		response.Tickets[i] = convertToProtoTicket(ticket)
	}

	return response, nil
}

// AddMessageToTicket implements the AddMessageToTicket RPC method
func (s *TicketService) AddMessageToTicket(ctx context.Context, req *proto.AddMessageToTicketRequest) (*proto.MessageResponse, error) {
	// Check if ticket exists
	_, err := s.ticketRepo.GetByID(ctx, req.TicketId)
	if err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to get ticket")
	}

	message := &model.Message{
		ID:          generateID(),
		TicketID:    req.TicketId,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Subject:     req.Subject,
		Body:        req.Body,
		Channel:     req.Channel,
		CreatedAt:   time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, status.Error(codes.Internal, "failed to create message")
	}

	return &proto.MessageResponse{
		Message: convertToProtoMessage(message),
	}, nil
}

// AddMessageToActiveTicket implements the AddMessageToActiveTicket RPC method
func (s *TicketService) AddMessageToActiveTicket(ctx context.Context, req *proto.AddMessageToActiveTicketRequest) (*proto.MessageResponse, error) {
	// Get active tickets for user
	tickets, err := s.ticketRepo.GetActiveByUser(ctx, req.User)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get active tickets")
	}

	if len(tickets) == 0 {
		return nil, status.Error(codes.NotFound, "no active tickets found for user")
	}

	// Use the most recent ticket
	ticket := tickets[0]

	message := &model.Message{
		ID:          generateID(),
		TicketID:    ticket.ID,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Subject:     req.Subject,
		Body:        req.Body,
		Channel:     req.Channel,
		CreatedAt:   time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, status.Error(codes.Internal, "failed to create message")
	}

	return &proto.MessageResponse{
		Message: convertToProtoMessage(message),
	}, nil
}

// GetTicketMessages implements the GetTicketMessages RPC method
func (s *TicketService) GetTicketMessages(ctx context.Context, req *proto.GetTicketMessagesRequest) (*proto.GetTicketMessagesResponse, error) {
	// Check if ticket exists
	_, err := s.ticketRepo.GetByID(ctx, req.TicketId)
	if err != nil {
		if err == repository.ErrTicketNotFound {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to get ticket")
	}

	messages, err := s.messageRepo.GetByTicketID(ctx, req.TicketId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get messages")
	}

	response := &proto.GetTicketMessagesResponse{
		Messages: make([]*proto.Message, len(messages)),
	}

	for i, message := range messages {
		response.Messages[i] = convertToProtoMessage(message)
	}

	return response, nil
}

// GetAllTickets implements the GetAllTickets RPC method
func (s *TicketService) GetAllTickets(ctx context.Context, req *proto.GetAllTicketsRequest) (*proto.GetAllTicketsResponse, error) {
	tickets, err := s.ticketRepo.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get tickets")
	}

	response := &proto.GetAllTicketsResponse{
		Tickets: make([]*proto.TicketResponse, len(tickets)),
	}

	for i, ticket := range tickets {
		response.Tickets[i] = convertToProtoTicket(ticket)
	}

	return response, nil
}

// Helper function to convert model.Ticket to proto.TicketResponse
func convertToProtoTicket(ticket *model.Ticket) *proto.TicketResponse {
	resp := &proto.TicketResponse{
		Id:        ticket.ID,
		Status:    ticket.Status,
		User:      ticket.User,
		Channel:   ticket.Channel,
		CreatedAt: ticket.CreatedAt.Format(time.RFC3339),
		UpdatedAt: ticket.UpdatedAt.Format(time.RFC3339),
	}

	if ticket.Agent != nil {
		resp.Agent = *ticket.Agent
	}
	if ticket.ProblemID != nil {
		resp.ProblemId = *ticket.ProblemID
	}
	if ticket.VerticalID != nil {
		resp.VerticalId = *ticket.VerticalID
	}
	if ticket.SkillID != nil {
		resp.SkillId = *ticket.SkillID
	}
	if ticket.UserGroupID != nil {
		resp.UserGroupId = *ticket.UserGroupID
	}

	return resp
}

// Helper function to convert model.Message to proto.Message
func convertToProtoMessage(message *model.Message) *proto.Message {
	return &proto.Message{
		Id:          message.ID,
		TicketId:    message.TicketID,
		FromAddress: message.FromAddress,
		ToAddress:   message.ToAddress,
		Subject:     message.Subject,
		Body:        message.Body,
		Channel:     message.Channel,
		CreatedAt:   message.CreatedAt.Format(time.RFC3339),
	}
}

// Helper functions for pointer conversions
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// Helper function to generate a unique ID
func generateID() string {
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

package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	generated "github.com/aimustaev/service-gateway/internal/generated/proto"
)

type GatewayServer struct {
	generated.UnimplementedGatewayServiceServer
}

func (s *GatewayServer) HandleNewMessage(req *generated.NewMessageRequest, stream generated.GatewayService_HandleNewMessageServer) error {
	log.Printf("Received message: %s from %s", req.Content, req.Sender)

	// Send acknowledgment response
	response := &generated.NewMessageResponse{
		Status:  "success",
		Message: "Message received successfully",
	}

	if err := stream.Send(response); err != nil {
		log.Printf("Error sending response: %v", err)
		return err
	}

	// Send processing status
	time.Sleep(time.Second)
	response = &generated.NewMessageResponse{
		Status:  "processing",
		Message: "Message is being processed",
	}

	if err := stream.Send(response); err != nil {
		log.Printf("Error sending response: %v", err)
		return err
	}

	// Send final status
	time.Sleep(time.Second)
	response = &generated.NewMessageResponse{
		Status:  "completed",
		Message: "Message processing completed",
	}

	if err := stream.Send(response); err != nil {
		log.Printf("Error sending response: %v", err)
		return err
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	generated.RegisterGatewayServiceServer(s, &GatewayServer{})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

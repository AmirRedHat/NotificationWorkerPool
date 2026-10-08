package main

import (
	"context"
	"log"
	"net"
	"task01/notification/services"

	"google.golang.org/grpc"

	pb "task01/notification/proto"
)

type Server struct {
	pb.UnimplementedNotificationRPCServiceServer
}

func (server *Server) TransportMessage(ctx context.Context, in *pb.NotificationMessageRequest) (*pb.NotificationMessageResponse, error) {
	err := services.EnqueueNotificationJob(ctx, in.Message, in.Kind, in.Receiver)
	if err != nil {
		return nil, err
	}
	return &pb.NotificationMessageResponse{Success: true, Message: in.Message, Kind: in.Kind}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNotificationRPCServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package grpc

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "task01/bff/proto"
)

func SendNotificationViaGRPC(ctx context.Context, kind string, message string, receiver string) error {
	notificationAddress := os.Getenv("GRPC_NOTIFICATION_ADDRESS")
	conn, err := grpc.NewClient(notificationAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return err
	}
	defer conn.Close()
	c := pb.NewNotificationRPCServiceClient(conn)
	response, err := c.TransportMessage(ctx, &pb.NotificationMessageRequest{Message: message, Kind: kind, Receiver: receiver})
	if err != nil {
		log.Printf("could not send: %v\n", err)
		return err
	}
	log.Printf("response: %v\n", response)
	return nil
}

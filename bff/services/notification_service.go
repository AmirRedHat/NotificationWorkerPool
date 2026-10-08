package services

import (
	"context"
	"task01/bff/config"
	"task01/bff/grpc"
)

func SendNotificationService(ctx context.Context, kind string, message string, receiver string, config config.Config) error {
	// send to notification service via RPC
	// you can validate or get specific configuration from config data
	err := grpc.SendNotificationViaGRPC(ctx, kind, message, receiver)
	if err != nil {
		return err
	}
	return nil
}

package services

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func EnqueuePushNotification(ctx context.Context, receiver string, message string, conn *redis.Client) error {
	log.Printf("EnqueuePushNotification: %v", message)
	return nil
}

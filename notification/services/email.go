package services

import (
	"context"
	"fmt"
	"log"
	commonServices "task01/common/services"

	"github.com/redis/go-redis/v9"
)

func EnqueueEmailNotification(ctx context.Context, receiver string, message string, conn *redis.Client) error {
	log.Printf("EnqueueEmailNotification: %v", message)
	value := fmt.Sprintf("%v:%v", receiver, message)
	log.Printf("Value: %v", value)
	return commonServices.AddRedisMessage(ctx, conn, "email", value)
}

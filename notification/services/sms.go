package services

import (
	"context"
	"fmt"
	"log"
	commonServices "task01/common/services"

	"github.com/redis/go-redis/v9"
)

func EnqueueSMSNotification(ctx context.Context, receiver string, message string, conn *redis.Client) error {
	log.Printf("EnqueueSMSNotification: %v", message)
	value := fmt.Sprintf("%v:%v", receiver, message)
	log.Printf("Value: %v", value)
	return commonServices.AddRedisMessage(ctx, conn, "sms", value)
}

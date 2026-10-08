package services

import (
	"context"
	"fmt"
	"task01/common/services"
)

func EnqueueNotificationJob(ctx context.Context, message string, kind string, receiver string) error {

	// define redis connection for more operations with one connection
	redisConnection := services.NewRedisClient()
	defer redisConnection.Close()

	// use factory for finding
	switch kind {
	case "SMS":
		return EnqueueSMSNotification(ctx, receiver, message, redisConnection)
	case "Email":
		return EnqueueEmailNotification(ctx, receiver, message, redisConnection)
	case "Push":
		return EnqueuePushNotification(ctx, receiver, message, redisConnection)
	default:
		return fmt.Errorf("invalid notification kind")
	}
}

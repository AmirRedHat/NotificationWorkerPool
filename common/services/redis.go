package services

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})
}

func AddRedisMessage(ctx context.Context, conn *redis.Client, key string, message string) error {
	err := conn.LPush(ctx, key, message).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetNotificationViaRedis(ctx context.Context, conn *redis.Client, key string) (string, string, error) {
	values := conn.BLPop(ctx, 1*time.Second, key).Val()

	if len(values) == 0 {
		return "", "", nil
	}
	value := values[1]
	notificationValues := strings.Split(value, ":")
	return notificationValues[0], notificationValues[1], nil
}

package services

import (
	"context"
	"fmt"
	"log"
	"task01/worker/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewEmailSender(workerNumber int, attemptNumber int, redisConn *redis.Client) *EmailSenderProvider {
	if attemptNumber <= 0 {
		attemptNumber = 5
	}
	return &EmailSenderProvider{
		Label:         "Email",
		WorkerNumber:  workerNumber,
		RedisConn:     redisConn,
		RedisKey:      "email",
		AttemptNumber: attemptNumber,
	}
}

type EmailSenderProvider struct {
	Label         string
	WorkerNumber  int
	RedisConn     *redis.Client
	RedisKey      string
	AttemptNumber int
}

func (p *EmailSenderProvider) Send(ctx context.Context, receiver string, message string) error {
	log.Printf("Email processing started = receiver: %s |  message: %s", receiver, message)

	// block receiver = 123456 | for test
	if receiver == "09151234567" {
		return fmt.Errorf("%s is not allowed", receiver)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Second * 7):
		log.Printf("Email processing finished = receiver: %s |  message: %s", receiver, message)
		return nil
	}
}

func (p *EmailSenderProvider) GetWorkerNumber() int {
	var c config.Config
	config.LoadConfig(&c)
	return c.Workers.Email.Number
}

func (p *EmailSenderProvider) GetSenderProviderLabel() string {
	return p.Label
}

func (p *EmailSenderProvider) GetAttemptNumber() int {
	return p.AttemptNumber
}

func (p *EmailSenderProvider) GetRedisConn() *redis.Client {
	return p.RedisConn
}

func (p *EmailSenderProvider) GetRedisKey() string {
	return p.RedisKey
}

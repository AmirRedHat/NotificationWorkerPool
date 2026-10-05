package services

import (
	"context"
	"fmt"
	"log"
	"task01/worker/config"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewSMSSender(workerNumber int, attemptNumber int, redisConn *redis.Client) *SMSSenderProvider {
	if attemptNumber <= 0 {
		attemptNumber = 5
	}
	return &SMSSenderProvider{
		Label:         "SMS",
		WorkerNumber:  workerNumber,
		RedisConn:     redisConn,
		RedisKey:      "sms",
		AttemptNumber: attemptNumber,
	}
}

type SMSSenderProvider struct {
	Label         string
	WorkerNumber  int
	RedisConn     *redis.Client
	RedisKey      string
	AttemptNumber int
}

func (p *SMSSenderProvider) Send(ctx context.Context, receiver string, message string) error {
	log.Printf("SMS processing started = receiver: %s |  message: %s", receiver, message)

	// block receiver = 123456 | for test
	if receiver == "09151234567" {
		return fmt.Errorf("%s is not allowed", receiver)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Second * 3):
		log.Printf("SMS processing finished = receiver: %s |  message: %s", receiver, message)
		return nil
	}
}

func (p *SMSSenderProvider) GetWorkerNumber() int {
	var c config.Config
	config.LoadConfig(&c)
	return c.Workers.Sms.Number
}

func (p *SMSSenderProvider) GetSenderProviderLabel() string {
	return p.Label
}

func (p *SMSSenderProvider) GetAttemptNumber() int {
	return p.AttemptNumber
}

func (p *SMSSenderProvider) GetRedisConn() *redis.Client {
	return p.RedisConn
}

func (p *SMSSenderProvider) GetRedisKey() string {
	return p.RedisKey
}

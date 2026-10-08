package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	commonServices "task01/common/services"
	"task01/worker/config"
	"task01/worker/services"
	"time"

	"github.com/redis/go-redis/v9"
)

func workerLookup(ctx context.Context, redisConn *redis.Client, config *config.WorkersConfig, errChan chan error) {
	DefaultAttempt := 3
	wg := &sync.WaitGroup{}
	defer close(errChan)

	// run SMS worker pool
	if config.Sms.Number > 0 {
		wg.Add(1)
		go func() {
			smsSigChan := make(chan os.Signal, 3)
			signal.Notify(smsSigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
			smsSender := services.NewSMSSender(config.Sms.Number, DefaultAttempt, redisConn)
			worker := services.NewWorker(smsSender, DefaultAttempt, 4*time.Second)
			if err := worker.Run(ctx, smsSigChan, wg); err != nil {
				errChan <- err
			}
		}()
	}
	// run Email worker pool - not implemented
	if config.Email.Number > 0 {
		wg.Add(1)
		go func() {
			emailSigChan := make(chan os.Signal, 3)
			signal.Notify(emailSigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
			emailSender := services.NewEmailSender(config.Email.Number, DefaultAttempt, redisConn)
			worker := services.NewWorker(emailSender, DefaultAttempt, 7*time.Second)
			if err := worker.Run(ctx, emailSigChan, wg); err != nil {
				errChan <- err
			}
		}()
	}
	// run Push worker pool - not implemented

	wg.Wait()
}

func main() {

	conn := commonServices.NewRedisClient()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run workers with config
	var c config.Config
	config.LoadConfig(&c)

	errChan := make(chan error)
	workerLookup(ctx, conn, &c.Workers, errChan)
	close(errChan)
}

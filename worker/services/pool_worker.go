package services

import (
	"context"
	"log"
	"os"
	"sync"
	"syscall"
	commonServices "task01/common/services"
	"task01/worker/schema"
	"time"

	"github.com/redis/go-redis/v9"
)

type SenderProvider interface {
	Send(ctx context.Context, receiver string, message string) error
	GetWorkerNumber() int
	GetSenderProviderLabel() string
	GetAttemptNumber() int
	GetRedisConn() *redis.Client
	GetRedisKey() string
}

type Worker struct {
	Sender        SenderProvider
	SenderLabel   string
	WorkerNumber  int
	RedisConn     *redis.Client
	RedisKey      string
	DLQ           chan schema.DeadNotification
	AttemptNumber int
	Timeout       time.Duration
}

func NewWorker(sender SenderProvider, attemptNumber int, timeout time.Duration) *Worker {
	return &Worker{
		sender,
		sender.GetSenderProviderLabel(),
		sender.GetWorkerNumber(),
		sender.GetRedisConn(),
		sender.GetRedisKey(),
		make(chan schema.DeadNotification),
		attemptNumber,
		timeout,
	}
}

func (worker *Worker) reload() {
	log.Printf("%s processing reloaded = worker: %d", worker.SenderLabel, worker.WorkerNumber)
	worker.WorkerNumber = worker.Sender.GetWorkerNumber()
	worker.SenderLabel = worker.Sender.GetSenderProviderLabel()
	worker.AttemptNumber = worker.Sender.GetAttemptNumber()
	worker.RedisConn = worker.Sender.GetRedisConn()
	worker.RedisKey = worker.Sender.GetRedisKey()
	log.Printf("%s processing reloaded = worker: %d", worker.SenderLabel, worker.WorkerNumber)
}

func (worker *Worker) start(wg *sync.WaitGroup, errChan chan error) {
	for range worker.WorkerNumber {
		wg.Add(1)
		go func() {
			err := worker.doTask(wg)
			if err != nil {
				errChan <- err
			}
		}()
	}
}

func (worker *Worker) Run(ctx context.Context, sigChan chan os.Signal, parentWG *sync.WaitGroup) error {

	errChan := make(chan error)
	defer parentWG.Done()

	log.Printf("%s starting worker: %d", worker.SenderLabel, worker.WorkerNumber)
	for {
		wg := &sync.WaitGroup{}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-sigChan:
			switch sig {
			case os.Interrupt:
				return worker.Shutdown(ctx)
			case syscall.SIGQUIT:
				worker.reload()
			default:
				log.Printf("Signal: %v", sig.String())
			}
		case err := <-worker.DLQ:
			log.Printf("DeadLock Queue : %s", err)
		default:
			worker.start(wg, errChan)
		}
		wg.Wait()
	}
}

func (worker *Worker) Shutdown(ctx context.Context) error {
	log.Printf("shutting down %s Pool Worker", worker.SenderLabel)
	if err := worker.RedisConn.Ping(ctx).Err(); err != nil {
		return nil
	}
	return worker.RedisConn.Close()
}

func (worker *Worker) doRetry(ctx context.Context, receiver string, message string, taskChan chan error) {
	retryAttemptCount := 1
	for range worker.AttemptNumber {
		if ctx.Err() != nil {
			return
		}
		taskRes := worker.Sender.Send(ctx, receiver, message)
		if taskRes != nil && ctx.Err() == nil {
			if retryAttemptCount >= worker.AttemptNumber {
				worker.DLQ <- schema.DeadNotification{Receiver: receiver, Message: message, Error: taskRes}
			}
			log.Printf("[%s] Retrying task #%d | receiver: %v | message: %v | DELAY: %v", worker.SenderLabel, retryAttemptCount, receiver, message, (1<<retryAttemptCount)*time.Second)
			time.Sleep((1 << retryAttemptCount) * time.Second)
			retryAttemptCount++
			continue
		}
		taskChan <- taskRes
		break
	}
}

func (worker *Worker) doTask(wg *sync.WaitGroup) error {
	defer wg.Done()

	getNotifCtx, cancelNotif := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelNotif()
	receiver, message, err := commonServices.GetNotificationViaRedis(getNotifCtx, worker.RedisConn, worker.RedisKey)
	if err != nil {
		return err
	}
	if receiver == "" || message == "" {
		return nil
	}

	taskCtx, cancelTask := context.WithTimeout(context.Background(), worker.Timeout)
	defer cancelTask()

	taskChan := make(chan error, 1)
	go worker.doRetry(taskCtx, receiver, message, taskChan)

	select {
	case <-taskCtx.Done():
		log.Println("task cancelled")
		return taskCtx.Err()
	case taskErr := <-taskChan:
		return taskErr
	}

}

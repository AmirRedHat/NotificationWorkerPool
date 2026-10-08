package services

import "context"

type BaseNotification interface {
	Send(ctx context.Context) error
}

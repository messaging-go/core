package core

import (
	"context"

	"github.com/messaging-go/core/internal/middleware"
)

type Middleware[T any] interface {
	Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error
}

type core[MessageType any] struct {
	chain           middleware.Processor[*MessageType, error]
	shouldContinue  bool
	resultsObserver func(error)
}
type Processor[MessageType any] func(ctx context.Context, item *MessageType) error

type MessageProcessor[MessageType any] interface {
	AddMiddleware(middleware Middleware[*MessageType]) MessageProcessor[MessageType]
	Stop()
	Run(processor Processor[MessageType])
}

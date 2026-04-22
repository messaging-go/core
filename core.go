package core

import (
	"context"

	"github.com/messaging-go/core/internal/middleware"
)

func (c *core[MessageType]) Run(processor Processor[MessageType]) {
	c.chain.AddMiddleware(middleware.FinalMiddleware(func(ctx context.Context, item *MessageType) error {
		return processor(ctx, item)
	}))

	for c.shouldContinue.Load() {
		ctx := context.Background()
		c.resultsObserver(c.chain.Process(ctx, nil))
	}
}

func (c *core[MessageType]) AddMiddleware(middleware Middleware[*MessageType]) MessageProcessor[MessageType] {
	c.chain.AddMiddleware(middleware)

	return c
}

func (c *core[MessageType]) Stop() {
	c.shouldContinue.Store(false)
}

func New[MessageType any](resultsObserver func(error)) MessageProcessor[MessageType] {
	instance := &core[MessageType]{
		chain:           middleware.New[*MessageType, error](),
		resultsObserver: resultsObserver,
	}
	instance.shouldContinue.Store(true)

	return instance
}

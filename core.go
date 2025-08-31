package core

import (
	"context"

	"github.com/messaging-go/core/internal/middleware"
)

func (c *core[MessageType]) Run(processor Processor[MessageType]) {
	c.chain.AddMiddleware(middleware.FinalMiddleware(func(ctx context.Context, item *MessageType) error {
		return processor(ctx, item)
	}))

	for c.shouldContinue {
		ctx := context.Background()
		c.resultsObserver(c.chain.Process(ctx, nil))
	}
}

func (c *core[MessageType]) AddMiddleware(middleware Middleware[*MessageType]) MessageProcessor[MessageType] {
	c.chain.AddMiddleware(middleware)

	return c
}

func (c *core[MessageType]) Stop() {
	c.shouldContinue = false
}

func New[MessageType any](resultsObserver func(error)) MessageProcessor[MessageType] {
	return &core[MessageType]{
		chain:           middleware.New[*MessageType, error](),
		shouldContinue:  true,
		resultsObserver: resultsObserver,
	}
}

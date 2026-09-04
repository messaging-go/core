package core

import (
	"context"

	"github.com/messaging-go/core/internal/middleware"
)

func (c *core[MessageType]) Run(processor Processor[MessageType]) {
	c.chain.AddMiddleware(middleware.FinalMiddleware(func(ctx context.Context, item *MessageType) error {
		return processor(ctx, item)
	}))

	for {
		select {
		case <-c.closeRequest:
			return
		default:
			c.resultsObserver(c.chain.Process(context.Background(), nil))
		}
	}
}

func (c *core[MessageType]) AddMiddleware(middleware Middleware[*MessageType]) MessageProcessor[MessageType] {
	c.chain.AddMiddleware(middleware)

	return c
}

func (c *core[MessageType]) Stop() {
	close(c.closeRequest)
}

func New[MessageType any](resultsObserver func(error)) MessageProcessor[MessageType] {
	return &core[MessageType]{
		chain:           middleware.New[*MessageType, error](),
		resultsObserver: resultsObserver,
		closeRequest:    make(chan struct{}),
	}
}

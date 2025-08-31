package core_test

import (
	"context"
	"fmt"
	"time"

	"github.com/messaging-go/core"
)

type consumerMiddleware[T any] struct {
	data []T
}

func (c *consumerMiddleware[T]) Process(ctx context.Context, _ T, next func(ctx context.Context, item T) error) error {
	// we're simulating messages being consumed from somewhere, it can be kafka, sqs, rabbitmq, etc.
	// basically the item for this middleware will be empty, and we need to inject a new message and call next.
	if len(c.data) == 0 {
		return nil // nothing left to process
	}

	// Take the first element
	item := c.data[0]
	// Remove it from the slice
	c.data = c.data[1:]

	// Call next with the item
	return next(ctx, item)
}

func newMockKafkaConsumer() core.Middleware[*int] {
	return &consumerMiddleware[*int]{
		data: []*int{ptr(1), ptr(2), ptr(3), ptr(4), ptr(5)},
	}
}

func ExampleNew() {
	processor := core.New[int](func(err error) {
		if err != nil {
			panic(err)
		}
	})
	processor.AddMiddleware(newMockKafkaConsumer())

	go func() {
		time.Sleep(100 * time.Millisecond)
		processor.Stop()
	}()

	processor.Run(func(ctx context.Context, item *int) error {
		if item == nil {
			return nil
		}

		fmt.Println(*item)

		return nil
	})
	// Output: 1
	// 2
	// 3
	// 4
	// 5
}

func ptr[T any](val T) *T {
	return &val
}

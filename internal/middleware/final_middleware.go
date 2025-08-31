package middleware

import "context"

type finalMiddleware[IN, OUT any] struct {
	processor func(ctx context.Context, item IN) OUT
}

func (m finalMiddleware[IN, OUT]) Process(ctx context.Context, item IN, _ func(ctx context.Context, item IN) OUT) OUT {
	return m.processor(ctx, item)
}

func FinalMiddleware[IN, OUT any](fn func(ctx context.Context, item IN) OUT) Middleware[IN, OUT] {
	return finalMiddleware[IN, OUT]{
		processor: fn,
	}
}

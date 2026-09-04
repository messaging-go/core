package middleware

import (
	"context"
	"fmt"
)

type Middleware[IN any, OUT any] interface {
	Process(ctx context.Context, item IN, next func(ctx context.Context, item IN) OUT) OUT
}

type Processor[IN any, OUT any] interface {
	AddMiddleware(middleware Middleware[IN, OUT])
	Process(ctx context.Context, input IN) OUT
}

type middlewareHandler[IN, OUT any] func(ctx context.Context, item IN) OUT

type middlewareChain[IN, OUT any] func(index int) middlewareHandler[IN, OUT]

type stack[IN any, OUT any] struct {
	middlewares []Middleware[IN, OUT]
}

func (r *stack[IN, OUT]) AddMiddleware(mw Middleware[IN, OUT]) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *stack[IN, OUT]) Process(ctx context.Context, options IN) OUT {
	// chain(index) returns the entry point for the chain starting at index.
	// The position is carried by the argument, not by shared mutable state, so
	// a middleware may call its `next` zero, one, or many times (retries), and
	// from several goroutines at once (fan-out), without corrupting the chain.
	var chain middlewareChain[IN, OUT]
	chain = func(index int) middlewareHandler[IN, OUT] {
		return func(c context.Context, item IN) OUT {
			if index >= len(r.middlewares) {
				panic(fmt.Sprintf(
					"middleware: next() called past the end of the chain (%d middleware(s)); "+
						"the last middleware must be a FinalMiddleware and must not call next",
					len(r.middlewares),
				))
			}

			return r.middlewares[index].Process(c, item, chain(index+1))
		}
	}

	return chain(0)(ctx, options)
}

func New[IN, OUT any]() Processor[IN, OUT] {
	return &stack[IN, OUT]{}
}

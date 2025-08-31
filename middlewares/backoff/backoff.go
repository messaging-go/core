package backoff

import (
	"context"
	"fmt"
	"time"

	"github.com/messaging-go/core"
	"github.com/messaging-go/core/middlewares/backoff/policies"
)

type backoff[T any] struct {
	policy           policies.Policy
	errorCounter     int
	maxErrorTracking int
}

func (b *backoff[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	delay := b.policy(b.errorCounter)
	if delay == 0 {
		return b.trackError(next(ctx, item))
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("backoff aborted due to context cancellation: %w", ctx.Err())
	case <-timer.C:
	}

	return b.trackError(next(ctx, item))
}

func (b *backoff[T]) trackError(err error) error {
	if err == nil {
		b.errorCounter = max(0, b.errorCounter-1)

		return nil
	}

	b.errorCounter = min(b.maxErrorTracking, b.errorCounter+1)

	return err
}

func Middleware[T any](policy policies.Policy, maxError int) core.Middleware[T] {
	return &backoff[T]{policy: policy, maxErrorTracking: maxError, errorCounter: 0}
}

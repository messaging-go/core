package backoff_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/messaging-go/core/middlewares/backoff"
	"github.com/messaging-go/core/middlewares/backoff/policies"
	"github.com/stretchr/testify/assert"
)

func TestBackoffProcess(t *testing.T) {
	t.Parallel()
	t.Run("no error", func(t *testing.T) {
		t.Parallel()
		mw := backoff.Middleware[int](policies.Constant(0), 5)

		err := mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		})

		assert.NoError(t, err)
	})
	t.Run("it returns when there's error, and it tracks the number of errors", func(t *testing.T) {
		t.Parallel()
		latestErrorCountSeen := 0
		mw := backoff.Middleware[int](func(count int) time.Duration {
			latestErrorCountSeen = count

			return 0
		}, 3)
		assert.Equal(t, 0, latestErrorCountSeen)
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		assert.Equal(t, 0, latestErrorCountSeen) // we see the past error
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		assert.Equal(t, 1, latestErrorCountSeen)
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		assert.Equal(t, 2, latestErrorCountSeen)
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		assert.Equal(t, 3, latestErrorCountSeen)
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		assert.Equal(t, 3, latestErrorCountSeen)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		assert.Equal(t, 3, latestErrorCountSeen)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		assert.Equal(t, 2, latestErrorCountSeen)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		assert.Equal(t, 1, latestErrorCountSeen)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		assert.Equal(t, 0, latestErrorCountSeen)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		assert.Equal(t, 0, latestErrorCountSeen)
	})
	t.Run("it delays processing when there's an error", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			mw := backoff.Middleware[int](policies.Constant(time.Second*5), 5)
			start := time.Now()

			_ = mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			})
			_ = mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			})
			_ = mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			})
			assert.Equal(t, time.Second*5*2, time.Since(start))
		})
	})
	t.Run("can cancel the wait with context", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			mw := backoff.Middleware[int](policies.Constant(time.Second*5), 5)
			start := time.Now()

			_ = mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			})
			ctx, _ := context.WithDeadline(t.Context(), time.Now().Add(time.Second))
			_ = mw.Process(ctx, 0, func(ctx context.Context, item int) error {
				return assert.AnError
			})
			assert.Equal(t, time.Second, time.Since(start))
		})
	})
}

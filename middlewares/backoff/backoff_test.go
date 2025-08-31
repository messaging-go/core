package backoff_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/messaging-go/core/middlewares/backoff"
	"github.com/messaging-go/core/middlewares/backoff/policies"
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
		middleware := backoff.Middleware[int](func(count int) time.Duration {
			latestErrorCountSeen = count

			return 0
		}, 3)

		require.Equal(t, 0, latestErrorCountSeen)
		require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		require.Equal(t, 0, latestErrorCountSeen) // we see the past error
		require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		require.Equal(t, 1, latestErrorCountSeen)
		require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		require.Equal(t, 2, latestErrorCountSeen)
		require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		require.Equal(t, 3, latestErrorCountSeen)
		require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
		require.Equal(t, 3, latestErrorCountSeen)
		require.NoError(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		require.Equal(t, 3, latestErrorCountSeen)
		require.NoError(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		require.Equal(t, 2, latestErrorCountSeen)
		require.NoError(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		require.Equal(t, 1, latestErrorCountSeen)
		require.NoError(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		require.Equal(t, 0, latestErrorCountSeen)
		require.NoError(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
		require.Equal(t, 0, latestErrorCountSeen)
	})
	t.Run("it delays processing when there's an error", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			t.Helper()

			middleware := backoff.Middleware[int](policies.Constant(time.Second*5), 5)
			start := time.Now()

			require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			}))
			require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			}))
			require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			}))
			require.Equal(t, time.Second*5*2, time.Since(start))
		})
	})
	t.Run("can cancel the wait with context", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			t.Helper()

			middleware := backoff.Middleware[int](policies.Constant(time.Second*5), 5)
			start := time.Now()

			require.Error(t, middleware.Process(t.Context(), 0, func(ctx context.Context, item int) error {
				return assert.AnError
			}))

			ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(time.Second))
			defer cancel()

			require.Error(t, middleware.Process(ctx, 0, func(ctx context.Context, item int) error {
				return assert.AnError
			}))

			require.Equal(t, time.Second, time.Since(start))
		})
	})
}

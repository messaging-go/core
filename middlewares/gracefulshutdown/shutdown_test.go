package gracefulshutdown_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/messaging-go/core/middlewares/gracefulshutdown"
	"github.com/messaging-go/core/test/mocks"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()
	t.Run("returns data returned by next middleware", func(t *testing.T) {
		t.Parallel()

		mw := gracefulshutdown.Middleware[int](nil)
		assert.NoError(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return nil
		}))
	})
	t.Run("returns error when next middleware returns error", func(t *testing.T) {
		t.Parallel()

		mw := gracefulshutdown.Middleware[int](nil)
		assert.Error(t, mw.Process(t.Context(), 0, func(ctx context.Context, item int) error {
			return assert.AnError
		}))
	})
	t.Run("returns error when context is canceled", func(t *testing.T) {
		t.Parallel()

		signalCh := make(chan os.Signal, 1)
		mockCore := mocks.NewMockMessageProcessor[int](gomock.NewController(t))
		mockCore.EXPECT().Stop()

		middleware := gracefulshutdown.MiddlewareWithStopSignalChannel[int](signalCh, mockCore)

		called := false

		// Start Process in a goroutine so we can send the "signal"
		go func() {
			err := middleware.Process(t.Context(), 1, func(ctx context.Context, item int) error {
				<-ctx.Done() // wait for cancellation

				called = true

				return ctx.Err()
			})
			assert.ErrorIs(t, err, context.Canceled)
		}()

		// simulate a signal
		signalCh <- nil // content doesn't matter, just triggers shutdown

		time.Sleep(50 * time.Millisecond) // give goroutine a moment to react

		assert.True(t, called)
	})
}

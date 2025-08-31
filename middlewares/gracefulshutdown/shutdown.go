package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/messaging-go/core"
)

type shutdown[T any] struct {
	coreInstance core.MessageProcessor[T]
	sigCh        <-chan os.Signal
}

// Process spawns a short-lived goroutine to cancel ctx if shutdown arrives.
func (s *shutdown[T]) Process(ctx context.Context, item T, next func(ctx context.Context, item T) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan struct{})
	defer close(done)

	// Goroutine to cancel context if a signal arrives
	go func() {
		select {
		case <-s.sigCh:
			s.coreInstance.Stop()
			cancel()
		case <-done:
			// Message finished, exit goroutine
		}
	}()

	return next(ctx, item)
}

// Middleware creates the shutdown middleware, sets up signal.Notify once.
func Middleware[T any](c core.MessageProcessor[T]) core.Middleware[T] {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	return &shutdown[T]{
		coreInstance: c,
		sigCh:        sigCh,
	}
}

func MiddlewareWithStopSignalChannel[T any](
	signalChannel chan os.Signal,
	c core.MessageProcessor[T],
) core.Middleware[T] {
	return &shutdown[T]{
		coreInstance: c,
		sigCh:        signalChannel,
	}
}

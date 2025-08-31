package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/messaging-go/core"
	"github.com/messaging-go/core/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("can process batch of items", func(t *testing.T) {
		t.Parallel()
		mockMessageInjector := mocks.NewMockMiddleware[*int](gomock.NewController(t))
		mockMessageInjector.
			EXPECT().
			Process(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, i *int, next func(context.Context, *int) error) error {
				assert.Nil(t, i)
				mockMessage := 1

				return next(ctx, &mockMessage)
			})
		mockMessageInjector.
			EXPECT().
			Process(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, i *int, next func(context.Context, *int) error) error {
				assert.Nil(t, i)
				mockMessage := 2

				return next(ctx, &mockMessage)
			})
		mockMessageInjector.
			EXPECT().
			Process(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, i *int, next func(context.Context, *int) error) error {
				assert.Nil(t, i)
				mockMessage := 3

				return next(ctx, &mockMessage)
			})
		mockMessageInjector.
			EXPECT().
			Process(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, i *int, next func(context.Context, *int) error) error {
				assert.Nil(t, i)

				return next(ctx, nil)
			}).AnyTimes()
		processor := core.New[int]()
		processor.AddMiddleware(mockMessageInjector)
		go func() {
			time.Sleep(10 * time.Millisecond)
			processor.Stop()
		}()
		var receivedMessages []int
		processor.Run(func(ctx context.Context, item *int) error {
			if item == nil {
				return nil
			}
			receivedMessages = append(receivedMessages, *item)

			return nil
		})
		assert.Equal(t, []int{1, 2, 3}, receivedMessages)
	})
}

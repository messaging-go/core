package policies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/messaging-go/core/middlewares/backoff/policies"
)

func TestJitter(t *testing.T) {
	t.Parallel()
	t.Run("when count is 1", func(t *testing.T) {
		t.Parallel()

		base := 1 * time.Second
		p := policies.Jitter(policies.Constant(base), 0.5)

		for i := 0; i < 10; i++ {
			d := p(1)
			assert.GreaterOrEqual(t, d, 500*time.Millisecond) // 50% lower bound
			assert.LessOrEqual(t, d, 1500*time.Millisecond)   // 50% upper bound
		}
	})
	t.Run("when count is 0", func(t *testing.T) {
		t.Parallel()

		base := 1 * time.Second
		p := policies.Jitter(policies.Constant(base), 0.5)

		for i := 0; i < 10; i++ {
			d := p(0)
			assert.Equal(t, time.Duration(0), d)
		}
	})
}

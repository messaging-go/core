package policies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/messaging-go/core/middlewares/backoff/policies"
)

func TestConstant(t *testing.T) {
	t.Parallel()

	constantPolicy := policies.Constant(2 * time.Second)
	assert.Equal(t, 0*time.Second, constantPolicy(0))
	assert.Equal(t, 2*time.Second, constantPolicy(1))
	assert.Equal(t, 2*time.Second, constantPolicy(2))
	assert.Equal(t, 2*time.Second, constantPolicy(5))
}

func TestLinear(t *testing.T) {
	t.Parallel()

	linearPolicy := policies.Linear(500 * time.Millisecond)
	assert.Equal(t, 0*time.Second, linearPolicy(0))
	assert.Equal(t, 500*time.Millisecond, linearPolicy(1))
	assert.Equal(t, 1500*time.Millisecond, linearPolicy(3))
}

func TestExponential(t *testing.T) {
	t.Parallel()

	exponentialPolicy := policies.Exponential(2, 200*time.Millisecond, 5*time.Second)
	assert.Equal(t, 0*time.Second, exponentialPolicy(0))
	assert.Equal(t, 200*time.Millisecond, exponentialPolicy(1))  // 200ms * 2^0
	assert.Equal(t, 400*time.Millisecond, exponentialPolicy(2))  // 200ms * 2^1
	assert.Equal(t, 1600*time.Millisecond, exponentialPolicy(4)) // 200ms * 2^3
	assert.Equal(t, 5*time.Second, exponentialPolicy(10))        // capped at max
}

func TestFibonacci(t *testing.T) {
	t.Parallel()

	fibonacciPolicy := policies.Fibonacci(100*time.Millisecond, 2*time.Second)
	assert.Equal(t, 0*time.Second, fibonacciPolicy(0))        // 0
	assert.Equal(t, 100*time.Millisecond, fibonacciPolicy(1)) // Fib(1) = 1
	assert.Equal(t, 100*time.Millisecond, fibonacciPolicy(2)) // Fib(2) = 1
	assert.Equal(t, 200*time.Millisecond, fibonacciPolicy(3)) // Fib(3) = 2
	assert.Equal(t, 300*time.Millisecond, fibonacciPolicy(4)) // Fib(4) = 3
	assert.Equal(t, 500*time.Millisecond, fibonacciPolicy(5)) // Fib(5) = 5
	assert.Equal(t, 2*time.Second, fibonacciPolicy(20))       // capped at max
}

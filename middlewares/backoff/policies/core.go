package policies

import (
	"math"
	"time"
)

// Linear increases linearly with each attempt: base * count.
func Linear(base time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}

		return base * time.Duration(count)
	}
}

// Constant returns the same duration every time.
func Constant(duration time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}

		return duration
	}
}

// Exponential doubles each time: base * 2^(count-1).
func Exponential(exponent float64, base time.Duration, maxDuration time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}

		duration := float64(base) * math.Pow(exponent, float64(count-1))
		if maxDuration <= 0 {
			return time.Duration(duration)
		}

		return min(time.Duration(duration), maxDuration)
	}
}

// Fibonacci backoff: base * Fib(count).
func Fibonacci(base time.Duration, maxDuration time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}

		first, second := 0, 1
		for i := 0; i < count; i++ {
			first, second = second, first+second
		}

		duration := time.Duration(first) * base
		if maxDuration <= 0 {
			return duration
		}

		return min(duration, maxDuration)
	}
}

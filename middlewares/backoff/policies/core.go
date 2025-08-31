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
func Constant(d time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}

		return d
	}
}

// Exponential doubles each time: base * 2^(count-1).
func Exponential(exponent float64, base time.Duration, max time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}
		d := float64(base) * math.Pow(exponent, float64(count-1))
		if max > 0 && time.Duration(d) > max {
			return max
		}

		return time.Duration(d)
	}
}

// Fibonacci backoff: base * Fib(count).
func Fibonacci(base time.Duration, max time.Duration) Policy {
	return func(count int) time.Duration {
		if count <= 0 {
			return 0
		}
		a, b := 0, 1
		for i := 0; i < count; i++ {
			a, b = b, a+b
		}
		d := time.Duration(a) * base
		if max > 0 && d > max {
			return max
		}

		return d
	}
}

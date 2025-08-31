package policies

import (
	"math/rand/v2"
	"time"
)

// Jitter wraps another policy and randomizes it by factor.
// factor = 0.5 means +/- 50%.
func Jitter(policy Policy, factor float64) Policy {
	return func(count int) time.Duration {
		duration := policy(count)
		if duration <= 0 {
			return duration
		}
		// random in [1-factor, 1+factor]
		multiplier := 1 + (rand.Float64()*2-1)*factor //nolint:gosec // this is not used for cryptography.

		return time.Duration(float64(duration) * multiplier)
	}
}

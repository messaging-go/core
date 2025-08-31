package policies

import (
	"math/rand/v2"
	"time"
)

// Jitter wraps another policy and randomizes it by factor.
// factor = 0.5 means +/- 50%.
func Jitter(p Policy, factor float64) Policy {
	return func(count int) time.Duration {
		d := p(count)
		if d <= 0 {
			return d
		}
		// random in [1-factor, 1+factor]
		mult := 1 + (rand.Float64()*2-1)*factor

		return time.Duration(float64(d) * mult)
	}
}

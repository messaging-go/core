package policies

import "time"

type Policy = func(count int) time.Duration
